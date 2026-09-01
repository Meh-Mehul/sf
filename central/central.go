package central

// rough arch
// CLIENT -------------> CENTRAL HUB (Job queue) <-------------> WORKER (maintains connections and readiness checks)

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
)

// i call the central server the 'central' (?)
type CentralOpts struct {
	Address string
}

// a typical job is just a command
type JobResult struct {
	ExitCode int
	Output   string
	Err      error
}

// channel is just to pipe result back
type Job struct {
	Cmd     string
	ResChan chan *JobResult
}

// To be run in the static server
type CentralHub struct {
	mu      sync.Mutex
	worker  net.Conn      // conn fd
	wreader *bufio.Reader // just a redundant reader
	jobs    chan *Job
	wReady  chan struct{} // to check if worker ready or not?wReady
}

func NewCentralHub() *CentralHub {
	return &CentralHub{
		jobs:   make(chan *Job, 100),
		wReady: make(chan struct{}, 1),
	}
}

// just starts the Hub
func StartHub(opts *CentralOpts) error {
	hub := NewCentralHub()
	l, err := net.Listen("tcp", opts.Address)
	if err != nil {
		return err
	}
	defer l.Close()
	log.Printf("Central hub listening on %s\n", opts.Address)

	go hub.workerDispatcher() // set-up worker job dispatching

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("Accept err: %v\n", err)
			continue
		}
		go hub.handleInbound(conn) // for both client and worker
	}
}

func (h *CentralHub) handleInbound(conn net.Conn) {
	r := bufio.NewReader(conn)
	greeting, err := r.ReadString('\n')
	if err != nil {
		conn.Close()
		return
	}
	role := strings.TrimSpace(greeting)

	switch role {
	case "WORKER":
		h.registerWorker(conn, r)
	case "CLIENT":
		h.handleClient(conn, r)
	default:
		conn.Write([]byte("ERR unknown role\n"))
		conn.Close()
	}
}

func (h *CentralHub) registerWorker(conn net.Conn, r *bufio.Reader) {
	h.mu.Lock()
	if h.worker != nil {
		h.worker.Close()
	}
	h.worker = conn
	h.wreader = r
	h.mu.Unlock()

	conn.Write([]byte("OK\n"))
	log.Printf("Worker connected from %s\n", conn.RemoteAddr())

	select {
	case h.wReady <- struct{}{}:
	default:
	}
}

// erase the old worker data (if any error occured)
func (h *CentralHub) eraseWorker(chk bool) {
	h.mu.Lock()
	if chk {
		h.worker.Close()
		h.worker = nil
		h.wreader = nil
	}
	h.mu.Unlock()
}

// takes a job at head, and puts it to worker
// worker's side code in its file
func (h *CentralHub) workerDispatcher() {
	for job := range h.jobs {
		for {
			h.mu.Lock()
			w := h.worker
			wr := h.wreader
			h.mu.Unlock()

			if w == nil {
				<-h.wReady // simple semaphore-ish blocking but
				// this blocks till a worker is registered
				continue
			}
			// payload format :
			// sz
			// <cmd>

			// response format:
			// exitcode
			// len_out
			// <output>
			payload := fmt.Sprintf("%d\n%s", len(job.Cmd), job.Cmd)
			if _, err := w.Write([]byte(payload)); err != nil {
				h.eraseWorker(h.worker == w)
				continue
			}
			codeStr, err := wr.ReadString('\n')
			if err != nil {
				h.eraseWorker(h.worker == w)
				continue
			}
			exitCode, pe := strconv.Atoi(strings.TrimSpace(codeStr))
			if pe != nil {
				exitCode = 1
			}
			lenStr, err := wr.ReadString('\n')
			if err != nil {
				h.eraseWorker(h.worker == w)
				continue
			}

			sz, pe := strconv.Atoi(strings.TrimSpace(lenStr))
			if pe != nil {
				sz = 0
			}

			outBuf := make([]byte, sz)
			if _, err := io.ReadFull(wr, outBuf); err != nil {
				h.eraseWorker(h.worker == w)
				continue
			}
			// send back job's result to its pipe once it has been executed at worker
			job.ResChan <- &JobResult{
				ExitCode: exitCode,
				Output:   string(outBuf),
				Err:      nil,
			}
			break
		}
	}
}

// handle when new client is attached (pushes to job chan)
func (h *CentralHub) handleClient(conn net.Conn, r *bufio.Reader) {
	defer conn.Close()
	if _, err := conn.Write([]byte("OK\n")); err != nil {
		return
	}

	// req format:
	// sz
	// <cmd>

	lenStr, err := r.ReadString('\n')
	if err != nil {
		return
	}

	sz, err := strconv.Atoi(strings.TrimSpace(lenStr))
	if err != nil {
		return
	}

	cmdBuf := make([]byte, sz)
	if _, err := io.ReadFull(r, cmdBuf); err != nil {
		return
	}

	job := &Job{
		Cmd:     string(cmdBuf),
		ResChan: make(chan *JobResult, 1),
	}
	// create a job and push to channel
	h.jobs <- job
	// we block here till the job is fully executed on worker-side and result is sent back
	res := <-job.ResChan

	// res format:
	// exitcode
	// sz
	// <output>
	resHeader := fmt.Sprintf("%d\n%d\n", res.ExitCode, len(res.Output))
	conn.Write([]byte(resHeader))
	conn.Write([]byte(res.Output))
}
