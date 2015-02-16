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
	dest, method := parseArgs()

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

		err = route(line, dest, method)
		if err != nil {
			fmt.Printf("%v", err)
			panic(err)
		}
	}
}

var parseArgs = func() (dest string, method string) {
	flag.StringVar(&dest, "dest", "", "The destination endpoint to which this program's Stdin will be forwarded")
	flag.StringVar(&method, "method", "", "The RPC method to call when forwarding a message from Stdin")

	flag.Parse()
	if dest == "" || method == "" {
		usg := "Example Usage:\n\n"
		usg += "\t$ someProgram | stdroute -dest www.example.com:9090/jsonrpc -method Log.Write"

		fmt.Println(usg)
		os.Exit(1)
	}

	return dest, method
}

func route(m string, dest string, method string) (err error) {
	var success bool

	err = client.Call(method, m, &success)
	if err != nil {
		return err
	}

	return nil
}
