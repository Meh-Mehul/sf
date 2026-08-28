package shell

import (
	"bufio"
	"fmt"
	"github.com/creack/pty"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// implementation of a persistent session shell in golang

type Shell struct {
	ptmx *os.File
	cmd  *exec.Cmd
	mu   sync.Mutex // so that only one can write to the shell
}

func NewShell() (*Shell, error) {
	cmd := exec.Command("bash", "--noprofile", "--norc")

	ptmx, err := pty.Start()
	if err != nil {
		return nil, err
	}
	return &Shell{
		ptmx: ptmx,
		cmd:  cmd,
	}, nil
}

func (s *Shell) Close() error {
	if s.cmd.Process != nil {
		_ = s.cmd.Process.Kill()
	}
	return s.ptmx.Close()
}

func (s *Shell) Run(command string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	marker := fmt.Sprintf("__COMMAND_DONE_%d__", time.Now().UnixNano())
	wrapped := fmt.Sprintf(
		"%s\nprintf '\\n%s:%%s\\n' \"$?\"\n",
		command,
		marker,
	)
	if _, err := io.WriteString(s.ptmx, wrapped); err != nil {
		return "", err
	}
	reader := bufio.NewReader(s.ptmx)
	var output strings.Builder
	for {
		line, err := reader.ReadString('\n')
		if len(line) > 0 {
			if strings.HasPrefix(line, marker+":") {
				exitCode := strings.TrimSpace(
					strings.TrimPrefix(line, marker+":"),
				)
				if exitCode != "0" {
					return output.String(), fmt.Errorf(
						"command exited with status %s",
						exitCode,
					)
				}
				return output.String(), nil
			}
			output.WriteString(line)
		}
		if err != nil {
			return output.String(), err
		}
	}
}
