package colorlog

import "sync"

type buffer struct {
	Text     []byte
	inverted bool
}

var bufPool = sync.Pool{
	New: func() any {
		return &buffer{
			Text:     make([]byte, 0, 1024),
			inverted: false,
		}
	},
}

func newBuffer() *buffer {
	return bufPool.Get().(*buffer)
}

func (b *buffer) Free() {
	// To reduce peak allocation, return only smaller buffers to the pool.
	const maxBufferSize = 16 << 10
	if cap(b.Text) <= maxBufferSize {
		b.Text = (b.Text)[:0]
		b.inverted = false
		bufPool.Put(b)
	}
}
func (b *buffer) Write(bytes []byte) (int, error) {
	b.Text = append(b.Text, bytes...)
	return len(bytes), nil
}

func (b *buffer) WriteByte(char byte) error {
	b.Text = append(b.Text, char)
	return nil
}

func (b *buffer) WriteString(str string) (int, error) {
	b.Text = append(b.Text, str...)
	return len(str), nil
}

func (b *buffer) WriteStringIf(ok bool, str string) (int, error) {
	if !ok {
		return 0, nil
	}
	return b.WriteString(str)
}

func (b *buffer) Inverse() bool {
	b.inverted = !b.inverted
	return !b.inverted
}

func (b *buffer) IsInverted() bool {
	return b.inverted
}
