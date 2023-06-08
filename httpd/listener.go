package httpd

import (
	"bytes"
	"net"
	"time"

	"github.com/vitorsalgado/mocha/v3/lib"
)

type MListener struct {
	wrapped net.Listener
	rPipes  []lib.Piping
	wPipes  []lib.Piping
}

func NewListener(wrapped net.Listener, rPipes []lib.Piping, wPipes []lib.Piping) net.Listener {
	return &MListener{
		wrapped: wrapped,
		rPipes:  rPipes,
		wPipes:  wPipes,
	}
}

func (l *MListener) Accept() (net.Conn, error) {
	conn, err := l.wrapped.Accept()
	if err != nil {
		return nil, err
	}

	return newConn(conn, l.rPipes, l.wPipes), nil
}

func (l *MListener) Close() error {
	return l.wrapped.Close()
}

func (l *MListener) Addr() net.Addr {
	return l.wrapped.Addr()
}

type MConn struct {
	wrapped net.Conn
	rPipes  []lib.Piping
	wPipes  []lib.Piping
}

func newConn(c net.Conn, rPipes []lib.Piping, wPipes []lib.Piping) *MConn {
	return &MConn{
		wrapped: c,
		rPipes:  rPipes,
		wPipes:  wPipes,
	}
}

func (c *MConn) Read(b []byte) (n int, err error) {
	n, err = c.wrapped.Read(b)
	if err != nil {
		return n, err
	}

	connector := lib.NewConnector(c.rPipes)
	bb := make([]byte, 0)
	buf := bytes.NewBuffer(bb)

	_, err = connector.Connect(b[:n], buf)
	if err != nil {
		return n, err
	}

	copy(b[:n], buf.Bytes())

	return n, nil
}

func (c *MConn) Write(b []byte) (n int, err error) {
	connector := lib.NewConnector(c.wPipes)
	return connector.Connect(b, c.wrapped)
}

func (c *MConn) Close() error {
	return c.wrapped.Close()
}

func (c *MConn) LocalAddr() net.Addr {
	return c.wrapped.LocalAddr()
}

func (c *MConn) RemoteAddr() net.Addr {
	return c.wrapped.RemoteAddr()
}

func (c *MConn) SetDeadline(t time.Time) error {
	return c.wrapped.SetDeadline(t)
}

func (c *MConn) SetReadDeadline(t time.Time) error {
	return c.wrapped.SetReadDeadline(t)
}

func (c *MConn) SetWriteDeadline(t time.Time) error {
	return c.wrapped.SetWriteDeadline(t)
}
