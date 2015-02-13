package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"os/exec"
	"strconv"
	"testing"
	"time"
)

func TestMain(t *testing.T) {
	port := 9090
	go listen(t, port)
	time.Sleep(1 * time.Second)
	testGoRun(t, port)
}

type TestRPC struct{}

func (rpc *TestRPC) Write(msg string, success *bool) (err error) {
	fmt.Printf("received: %s\n", msg)
	_true := true
	success = &_true

	return nil
}

func listen(t *testing.T, port int) {
	server := rpc.NewServer()
	err := server.RegisterName("Stdin", &TestRPC{})
	if err != nil {
		t.Error(err)
		return
	}

	l, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		t.Error(err)
		return
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			t.Error(err)
			return
		}

		go server.ServeCodec(jsonrpc.NewServerCodec(conn))
	}

	return
}

func testGoRun(t *testing.T, port int) {
	c1 := exec.Command("echo", "line 1\nline 2\nline 3")
	c2 := exec.Command("go", "run", "main.go", "localhost:"+strconv.Itoa(port))

	r, w := io.Pipe()
	c1.Stdout = w
	c2.Stdin = r

	var b2 bytes.Buffer
	c2.Stdout = &b2

	c1.Start()
	c2.Start()
	c1.Wait()
	w.Close()
	c2.Wait()
	io.Copy(os.Stdout, &b2)
}
