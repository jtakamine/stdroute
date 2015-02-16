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
	"strings"
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

func (rpc *TestRPC) Write(m string, success *bool) (err error) {
	fmt.Printf("received \"%s\"\n", m)
	*success = false
	if strings.Contains(m, "true") {
		*success = true
	}

	return nil
}
func listen(t *testing.T, port int) {
	server := rpc.NewServer()
	err := server.RegisterName("Log", &TestRPC{})
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
	c1 := exec.Command("echo", "line 1\nline 2\nline 3\nline 4(return true)\nline 5\nline 6(return true)\nline 7\nline 8\nline 9(return true)")
	c2 := exec.Command("go", "run", "main.go", "-dest", "localhost:"+strconv.Itoa(port), "-method", "Log.Write")

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
