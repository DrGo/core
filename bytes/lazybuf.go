package bytes

// source: modified from go and upspin code

// A Lazybuf is a lazily constructed path buffer.
// It supports append, reading previously appended bytes,
// and retrieving the final string. It does not allocate a buffer
// to hold the output until that output diverges from s.
type Lazybuf struct {
	s   string
	buf []byte
	w   int
}

func (b *Lazybuf) index(i int) byte {
	if b.buf != nil {
		return b.buf[i]
	}
	return b.s[i]
}

func (b *Lazybuf) append(c byte) {
	if b.buf == nil {
		if b.w < len(b.s) && b.s[b.w] == c {
			b.w++
			return
		}
		b.buf = make([]byte, len(b.s))
		copy(b.buf, b.s[:b.w])
	}
	b.buf[b.w] = c
	b.w++
}

func (b *Lazybuf) string() string {
	if b.buf == nil {
		return b.s[:b.w]
	}
	return string(b.buf[:b.w])
}

// LazyBuffer is a []byte that is lazily (re-)allocated when its
// Bytes method is called.
type LazyBuffer []byte

// Bytes returns a []byte that has length n. It re-uses the underlying
// LazyBuffer []byte if it is at least n bytes in length.
func (b *LazyBuffer) Bytes(n int) []byte {
	if *b == nil || len(*b) < n {
		*b = make([]byte, n)
	}
	return (*b)[:n]
}
