package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"strings"
	"time"
)

var client *rpc.Client

func main() {
	dest := parseArgs()

	conn, err := net.DialTimeout("tcp", dest, time.Second*5)
	if err != nil {
		panic(err)
	}

	client = jsonrpc.NewClient(conn)

	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		line = strings.TrimSuffix(line, "\n")
		if err != nil {
			if err == io.EOF {
				fmt.Println("Stdin closed. Exiting.")
				return
			}
			panic(err)
		}

		err = route(line, dest)
		if err != nil {
			fmt.Printf("%v", err)
			panic(err)
		}
	}
}

var parseArgs = func() (dest string) {
	flag.Parse()
	args := flag.Args()

	dest = ""
	switch len(args) {
	case 0:
		dest = os.Getenv("STDROUTE_DEST")
	case 1:
		dest = args[0]
	}

	if dest == "" {
		fmt.Println("Example Usage:\n\n\t$ someProgram | stdroute www.example.com:8000/jsonrpc\n\n\tor\n\n\t$ export STDROUTE_DEST=www.example.com:8000/jsonrpc\n\t$ someProgram | stdroute\n")
		os.Exit(1)
	}

	return dest
}

func route(m string, dest string) (err error) {
	var success bool

	err = client.Call("Stdin.Write", m, &success)
	if err != nil {
		return err
	}

	return nil
}
