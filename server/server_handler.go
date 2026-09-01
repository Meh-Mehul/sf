package server

// To be run on worker
//
// Code for the worker side (im also calling it server)
// Workflow:
// connects by sending WORKER to Central server
// then handles readyness checks
// and shell cmd executions
import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/Meh-Mehul/sf/shell"
)

type ServerOpts struct {
	Address    string        // ip address (static of the main server)
	MaxSession int           // currently only 1
	RetryDelay time.Duration // how much to retry until removal (only at initial time)
}

type Server struct {
	CurrSessions   []string
	MaxSession     int
	RequestingAddr string
	sh             *shell.Shell
}

func StartServer(opts *ServerOpts) {
	if opts.RetryDelay == 0 {
		opts.RetryDelay = 3 * time.Second
	}

	for {
		log.Printf("Dialing central hub at %s...\n", opts.Address)
		conn, err := net.Dial("tcp", opts.Address)
		if err != nil {
			log.Printf("Dial failed: %v, retrying in %v...\n", err, opts.RetryDelay)
			time.Sleep(opts.RetryDelay)
			continue
		}
		// set up shell
		sh, err := shell.NewShell()
		if err != nil {
			log.Printf("Failed to spawn shell: %v\n", err)
			conn.Close()
			time.Sleep(opts.RetryDelay)
			continue
		}

		srv := &Server{
			CurrSessions:   []string{"default"},
			MaxSession:     opts.MaxSession,
			RequestingAddr: opts.Address,
			sh:             sh,
		}

		err = srv.handleConn(conn)
		sh.Close()
		conn.Close()
		if err != nil {
			log.Printf("Connection dropped (%v), reconnecting...\n", err)
		}
		time.Sleep(opts.RetryDelay)
	}
}

// first send WORKER, then loop infinitely for incoming
// tasks
func (s *Server) handleConn(conn net.Conn) error {
	r := bufio.NewReader(conn)
	if _, err := conn.Write([]byte("WORKER\n")); err != nil {
		return err
	}
	ack, err := r.ReadString('\n')
	if err != nil || strings.TrimSpace(ack) != "OK" {
		return fmt.Errorf("bad handshake: %v, ack: %s", err, ack)
	}
	// just get->run->return loop
	for {
		lenStr, err := r.ReadString('\n')
		if err != nil {
			return err
		}
		sz, err := strconv.Atoi(strings.TrimSpace(lenStr))
		if err != nil {
			return fmt.Errorf("bad payload size: %v", err)
		}
		cmdBuf := make([]byte, sz)
		if _, err := io.ReadFull(r, cmdBuf); err != nil {
			return err
		}
		cmdStr := string(cmdBuf)
		// no run the cmd
		outStr, runErr := s.sh.Run(cmdStr)
		exitCode := 0
		if runErr != nil {
			exitCode = 1
			if strings.Contains(runErr.Error(), "status ") {
				parts := strings.Split(runErr.Error(), "status ")
				if len(parts) == 2 {
					if parsed, pe := strconv.Atoi(strings.TrimSpace(parts[1])); pe == nil {
						exitCode = parsed
					}
				}
			}
		}

		resHeader := fmt.Sprintf("%d\n%d\n", exitCode, len(outStr))
		if _, err := conn.Write([]byte(resHeader)); err != nil {
			return err
		}
		if _, err := conn.Write([]byte(outStr)); err != nil {
			return err
		}
	}
}
