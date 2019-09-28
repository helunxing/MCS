package cache

type Stat struct {
	Count     int64
	KeySize   int64
	ValueSize int64
}

// 添加时更改状态
func (s *Stat) add(k string, v []byte) {
	s.Count += 1
	s.KeySize += int64(len(k))
	s.ValueSize += int64(len(v))
}

// 删除时更改状态
func (s *Stat) del(k string, v []byte) {
	s.Count -= 1
	s.KeySize -= int64(len(k))
	s.ValueSize -= int64(len(v))
}
