package main

import (
	"flag"
	"io"
	"log"
	"net"
	"os"
)

var args struct {
	address  string
	filepath string
}

func init() {
	flag.StringVar(&args.address, "address", "", "remote address to connect to")
	flag.StringVar(&args.filepath, "filepath", "", "file to save")
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

	file, err := os.OpenFile(args.filepath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	conn, err := net.Dial("tcp", args.address)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	_, err = io.Copy(file, conn)
	if err != nil {
		log.Fatalf("transfer error: %s", err)
	}
}
