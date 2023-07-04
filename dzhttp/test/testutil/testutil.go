package testutil

import "net"

func RandomTCPPort() (int, error) {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, err
	}

	port := l.Addr().(*net.TCPAddr).Port
	_ = l.Close()

	return port, nil
}
