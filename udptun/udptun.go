package main

import (
	"flag"
	"log"
	"net"

	"github.com/pkg/errors"
)

var pool map[int]*net.UDPConn

// ./udptun.exe -l ":6000" -r "127.0.0.1:4000"
func main() {
	var laddr string
	var raddr string
	flag.StringVar(&laddr, "l", ":6000", "local addr")
	flag.StringVar(&raddr, "r", "127.0.0.1:40000", "remote addr")
	flag.Parse()

	pool = make(map[int]*net.UDPConn)

	udpaddr, err := net.ResolveUDPAddr("udp", laddr)
	if err != nil {
		log.Println(errors.WithStack(err))
	}
	conn, err := net.ListenUDP("udp", udpaddr)
	if err != nil {
		log.Println(errors.WithStack(err))
	}

	rep, err := net.ResolveUDPAddr("udp", raddr)
	if err != nil {
		log.Println(errors.WithStack(err))
	}
	// network := "udp4"
	// if udpaddr.IP.To4() == nil {
	// 	network = "udp"
	// }

	listen(conn, rep)
}

func listen(lc *net.UDPConn, rep *net.UDPAddr) {
	// buf := make([]byte, 1500)
	for {
		buf := make([]byte, 1500)
		n, addr, err := lc.ReadFrom(buf)
		if err != nil {
			log.Println(errors.WithStack(err))
		}

		// debug
		// fmt.Println("From ", addr, " len=", n)
		// fmt.Println(buf[:n])
		// fmt.Println(addr.(*net.UDPAddr).Port)

		// 如果没有则新建，按照port来
		if pool[addr.(*net.UDPAddr).Port] == nil {
			udpaddr := &net.UDPAddr{
				IP:   net.IPv4(127, 27, 7, 1),
				Port: addr.(*net.UDPAddr).Port,
			}
			nc, err := net.ListenUDP("udp", udpaddr)
			if err != nil {
				log.Println(errors.WithStack(err))
			}
			pool[addr.(*net.UDPAddr).Port] = nc
			go newConnDeamon(lc, nc, addr) // 监听的conn，新的conn，新建时收到的地址
		}

		// 通过port索引并发送
		n, err = pool[addr.(*net.UDPAddr).Port].WriteToUDP(buf[:n], rep)
		if err != nil {
			log.Println(errors.WithStack(err))
		}
		// debug
		// fmt.Println("From ", pool[addr.(*net.UDPAddr).Port].LocalAddr().String(), "To ", rep, " len=", n)

	}
}

// 监听的conn，新的conn，新建时收到的地址
func newConnDeamon(lc, conn *net.UDPConn, la net.Addr) {
	// buf := make([]byte, 1500)
	for {
		buf := make([]byte, 1500)
		n, _, err := conn.ReadFrom(buf)
		if err != nil {
			log.Println(errors.WithStack(err))
		}
		// debug
		// fmt.Println("回传：From ", addr, " len=", n)
		// fmt.Println(buf[:n])
		// fmt.Println(addr.(*net.UDPAddr).Port)

		// 回传
		n, err = lc.WriteTo(buf[:n], la)
		if err != nil {
			log.Println(errors.WithStack(err))
		}
		// fmt.Println("回传：From ", lc.LocalAddr(), "To ", la, " len=", n)

	}
}
