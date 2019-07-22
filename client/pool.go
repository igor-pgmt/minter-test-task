package client

type pool chan []byte

func newPool(maxCap, poolCap uint64) *pool {
	p := make(pool, poolCap)
	for i := uint64(0); i < poolCap; i++ {
		b := make([]byte, maxCap)
		(&p).putBytes(b)
	}
	return &p
}

func (p *pool) getBytes() (b []byte) {
	bt, ok := <-*p
	if ok {
		return bt
	}
	panic("the pool channel is closed")
}

func (p *pool) putBytes(b []byte) {
	b = b[:0]
	select {
	case *p <- b:
	default:
	}

	return
}
