package encutil

import (
	"bytes"
	"sync"
)

const maxSize = 1 << 16

var bufPool = &sync.Pool{
	New: func() any {
		return bytes.NewBuffer(make([]byte, 500))
	},
}

func putBuf(b *bytes.Buffer) {
	if b.Cap() > maxSize {
		return
	}

	bufPool.Put(b)
}

func Join(sep string, elems ...string) string {
	switch len(elems) {
	case 0:
		return ""
	case 1:
		return elems[0]
	}

	n := len(sep) * (len(elems) - 1)
	for i := 0; i < len(elems); i++ {
		n += len(elems[i])
	}

	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	buf.WriteString(elems[0])

	for _, s := range elems[1:] {
		buf.WriteString(sep)
		buf.WriteString(s)
	}

	txt := buf.String()
	putBuf(buf)

	return txt
}
