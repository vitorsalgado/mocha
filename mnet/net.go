package mnet

import (
	"fmt"
	"net"
)

type MNet struct {
	listener net.Listener
}

func New(addr string) *MNet {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	return &MNet{listener: l}
}

func (m *MNet) Listen() {
	go func() {
		for {
			conn, err := m.listener.Accept()
			if err != nil {
				panic(err)
			}

			// _, err = net.Dial("tcp", "127.0.0.1:6379")
			// if err != nil {
			// 	conn.Close()
			// 	fmt.Println(err.Error())
			// 	continue
			// }

			// bb, _ := io.ReadAll(conn)
			// reader := io.NopCloser(bytes.NewBuffer(bb))

			//			go func() {

			// p := make([]byte, 4)
			// for {
			// 	n, err := conn.Read(p)
			// 	if err == io.EOF {
			// 		break
			// 	}
			// 	// fmt.Print(string(p[:n]))
			// 	conn.Write(p[:n])
			// }

			// scanner := bufio.NewScanner(conn)
			// for scanner.Scan() {
			// 	b := scanner.Bytes()

			// 	fmt.Println(string(b))
			// 	conn2.Write(b)
			// }
			// }()

			buf := make([]byte, 32*1024)

			go func() {
				for {
					n, _ := conn.Read(buf)

					if n > 0 {
						bb := buf[0:n]

						fmt.Println(string(bb))

						conn.Write(bb)
					}
				}
			}()

			// go io.Copy(conn2, conn)

			// go io.Copy(conn, conn)
		}
	}()
}
