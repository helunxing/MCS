package cache

// #include "rocksdb/c.h"
// #cgo CFLAGS: -I${SRCDIR}/../../../rocksdb/include
// #cgo LDFLAGS: -L${SRCDIR}/../../../rocksdb -lrocksdb -lz -lpthread -lsnappy -lstdc++ -lm -O3
import "C"
import "runtime"

func newRocksdbCache() *rocksdbCache {
	// 创建对象指针
	options := C.rocksdb_options_create()
	// 设置并发线程数
	C.rocksdb_options_increase_parallelism(options, C.int(runtime.NumCPU()))
	// 不存在则创建新目录
	C.rocksdb_options_set_create_if_missing(options, 1)
	// 打开目录
	var e *C.char //这个变量本应该手动回收
	db := C.rocksdb_open(options, C.CString("/mnt/rocksdb"), &e)
	if e != nil {
		panic(C.GoString(e))
	}
	// 删除指针
	C.rocksdb_options_destroy(options)

	c := make(chan *pair, 5000)
	wo := C.rocksdb_writeoptions_create()
	// 启动一个新的协程接收c中的数据
	go write_func(db, c, wo)
	// 返回rocksdbcache结构体指针
	return &rocksdbCache{db, C.rocksdb_readoptions_create(), wo, e, c}
}
