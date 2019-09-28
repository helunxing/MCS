package tcp

import (
	"bufio"
	"io"
)

func (s *Server) readKey(r *bufio.Reader) (string, error) {
	// 读取一个字符串并转成整形
	klen, e := readLen(r)
	if e != nil {
		return "", e
	}
	// 将键内容复制
	k := make([]byte, klen)
	_, e = io.ReadFull(r, k)
	if e != nil {
		return "", e
	}
	return string(k), nil
}

func (s *Server) readKeyAndValue(r *bufio.Reader) (string, []byte, error) {
	// 读取键长度
	klen, e := readLen(r)
	if e != nil {
		return "", nil, e
	}
	// 读取值长度
	vlen, e := readLen(r)
	if e != nil {
		return "", nil, e
	}
	// 读取键
	k := make([]byte, klen)
	_, e = io.ReadFull(r, k)
	if e != nil {
		return "", nil, e
	}
	// 读取值
	v := make([]byte, vlen)
	_, e = io.ReadFull(r, v)
	if e != nil {
		return "", nil, e
	}
	return string(k), v, nil
}
