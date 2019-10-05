package tcp

import (
	"bufio"
	"errors"
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
	// 检查其是否由本节点处理
	key := string(k)
	addr, ok := s.ShouldProcess(key)
	if !ok {
		return "", errors.New("redirect " + addr)
	}
	return key, nil
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
	// 检查其是否由本节点处理
	key := string(k)
	addr, ok := s.ShouldProcess(key)
	if !ok {
		return "", nil, errors.New("redirect " + addr)
	}
	// 读取值
	v := make([]byte, vlen)
	_, e = io.ReadFull(r, v)
	if e != nil {
		return "", nil, e
	}
	return key, v, nil
}
