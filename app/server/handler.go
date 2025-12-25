package server

import (
	"bufio"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/protocol"
)

type ConnHandler struct {
	conn net.Conn
	in   chan protocol.Cmd
}

func (h *ConnHandler) readLoop() {
	reader := bufio.NewReader(h.conn)

	for {
		cmd, err := protocol.ReadCommand(reader)

		if err != nil {
			close(h.in)
			return
		}

		h.in <- *cmd
	}
}
