package server

import (
	"bufio"
	"net"

	"github.com/AdityaTaggar05/godis/internal/protocol"
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
