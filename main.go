package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"net/rpc/jsonrpc"
	"os"
	"time"
)

func main() {
	dest := parseArgs()

	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("Stdin closed. Exiting.")
				return
			}
			panic(err)
		}

		route(line, dest)
	}
}

var parseArgs = func() (dest string) {
	flag.Parse()
	args := flag.Args()

	dest = ""
	switch len(args) {
	case 1:
		dest = args[0]
	case 0:
		dest = os.Getenv("STDROUTE_DEST")
	}

	if dest == "" {
		fmt.Println("Example Usage:\n\n\t$ someProgram | stdroute www.example.com:8000/jsonrpc\n\n\tor\n\n\t$ export STDROUTE_DEST=www.example.com:8000/jsonrpc\n\t$ someProgram | stdroute\n")
		os.Exit(1)
	}

	return dest
}

func route(msg string, dest string) (err error) {
	conn, err := net.DialTimeout("tcp", dest, time.Second*5)
	if err != nil {
		return err
	}
	c := jsonrpc.NewClient(conn)

	var success bool

	err = c.Call("Stdin.Write", msg, &success)
	if err != nil {
		return err
	}

	return nil
}
