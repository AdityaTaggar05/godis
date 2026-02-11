package server

import (
	"bufio"
	"fmt"
	"io"
	"net"

	"github.com/AdityaTaggar05/godis/internal/protocol"
)

type ConnHandler struct {
	conn net.Conn
}

func NewConnHandler(conn net.Conn) *ConnHandler {
	return &ConnHandler{
		conn: conn,
	}
}

func (h *ConnHandler) Handle() {
	defer h.conn.Close()

	reader := bufio.NewReader(h.conn)

	for {
		cmd, err := protocol.ReadCommand(reader)
		if err != nil {
			if err == io.EOF {
				return
			}
			fmt.Printf("[ERR] invalid sequence of characters: %s\n", err.Error())

			return
		}

		fmt.Printf("[DEBUG] Received command: %s\n", cmd)

		func() {
			defer func() {
				if r := recover(); r != nil {
					h.handleRecover(r)
				}
			}()

			h.dispatch(cmd)
		}()
	}
}

func (h *ConnHandler) writeOutput(b []byte) error {
	_, err := h.conn.Write(b)

	if err != nil {
		fmt.Println("Error writing data to connection: ", err.Error())
		return err
	}

	return nil
}
