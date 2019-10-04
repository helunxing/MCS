package cache

// #include <stdlib.h>
// #include "rocksdb/c.h"
// #cgo CFLAGS: -I${SRCDIR}/../../../rocksdb/include
// #cgo LDFLAGS: -L${SRCDIR}/../../../rocksdb -lrocksdb -lz -lpthread -lsnappy -lstdc++ -lm -O3
import "C"
import (
	"time"
	"unsafe"
)

const BATCH_SIZE = 100

func flush_batch(db *C.rocksdb_t, b *C.rocksdb_writebatch_t, o *C.rocksdb_writeoptions_t) {
	var e *C.char
	C.rocksdb_write(db, o, b, &e)
	if e != nil {
		panic(C.GoString(e))
	}
	C.rocksdb_writebatch_clear(b)
}

func write_func(db *C.rocksdb_t, c chan *pair, o *C.rocksdb_writeoptions_t) {
	// 当前写入键值对数量计数
	count := 0
	// 计数器
	t := time.NewTimer(time.Second)
	// 用于批量写入的rocksdb结构体
	b := C.rocksdb_writebatch_create()
	// 等待timer和channel中的事件
	for {
		select {
		case p := <-c:
			count++
			key := C.CString(p.k)
			value := C.CBytes(p.v)
			C.rocksdb_writebatch_put(b, key, C.size_t(len(p.k)), (*C.char)(value), C.size_t(len(p.v)))
			C.free(unsafe.Pointer(key))
			C.free(value)
			if count == BATCH_SIZE {
				flush_batch(db, b, o)
				count = 0
			}
			// 若计时器已经触发，则要取走其首个Time结构体，避免干扰
			if !t.Stop() {
				<-t.C
			}
			t.Reset(time.Second)
		case <-t.C:
			if count != 0 {
				flush_batch(db, b, o)
				count = 0
			}
			t.Reset(time.Second)
		}
	}
}

func (c *rocksdbCache) Set(key string, value []byte) error {
	// 只需将键值传入ch中即可
	c.ch <- &pair{key, value}
	return nil
}
