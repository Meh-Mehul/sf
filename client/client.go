package client

// to be ran in client server
// Just gets->pushes and waits

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

type ClientOpts struct {
	Address string
}

func PushJob(opts *ClientOpts, cmd string) (string, int, error) {
	conn, err := net.Dial("tcp", opts.Address)
	if err != nil {
		return "", 1, fmt.Errorf("dial failed: %v", err)
	}
	defer conn.Close()
	r := bufio.NewReader(conn)
	if _, err := conn.Write([]byte("CLIENT\n")); err != nil {
		return "", 1, err
	}
	ack, err := r.ReadString('\n')
	if err != nil || strings.TrimSpace(ack) != "OK" {
		return "", 1, fmt.Errorf("bad handshake: %v, ack: %s", err, ack)
	}
	payload := fmt.Sprintf("%d\n%s", len(cmd), cmd)
	if _, err := conn.Write([]byte(payload)); err != nil {
		return "", 1, err
	}
	codeStr, err := r.ReadString('\n')
	if err != nil {
		return "", 1, fmt.Errorf("read exit code failed: %v", err)
	}
	exitCode, pe := strconv.Atoi(strings.TrimSpace(codeStr))
	if pe != nil {
		exitCode = 1
	}
	lenStr, err := r.ReadString('\n')
	if err != nil {
		return "", exitCode, fmt.Errorf("read output length failed: %v", err)
	}
	sz, pe := strconv.Atoi(strings.TrimSpace(lenStr))
	if pe != nil {
		sz = 0
	}
	outBuf := make([]byte, sz)
	if _, err := io.ReadFull(r, outBuf); err != nil {
		return "", exitCode, fmt.Errorf("read output payload failed: %v", err)
	}

	return string(outBuf), exitCode, nil
}
