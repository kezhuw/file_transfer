package main

import (
	"flag"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"syscall"
)

var args struct {
	address  string
	filepath string
}

func init() {
	flag.StringVar(&args.address, "address", "", "local address to listen on")
	flag.StringVar(&args.filepath, "filepath", "", "file to transfer")
}

func parse_args() {
	flag.Parse()

	if args.address == "" || args.filepath == "" {
		flag.PrintDefaults()
		os.Exit(0)
	}
}

func main() {
	parse_args()

	file, err := os.Open(args.filepath)
	if err != nil {
		log.Fatal(err)
	}

	listener, err := net.Listen("tcp", args.address)
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			break
		}

		fd, err := syscall.Dup(int(file.Fd()))
		if err != nil {
			log.Println(err)
			conn.Close()
			break
		}
		sourceFile := os.NewFile(uintptr(fd), args.filepath)

		wg.Add(1)
		go func(conn net.Conn, file *os.File, wg *sync.WaitGroup) {
			defer conn.Close()
			defer wg.Done()
			_, err := io.Copy(conn, file)
			if err != nil {
				log.Printf("transfer to %s error: %s", conn.RemoteAddr(), err)
			}
		}(conn, sourceFile, &wg)
	}
	wg.Wait()
}
