package server

import (
	"flag"
	"fmt"
	"net"
	"sync"

	"github.com/codecrafters-io/redis-starter-go/app/protocol"
)

var DB sync.Map
var Server *ServerConfig

func Start() {
	DB = sync.Map{}

	port := flag.Int("port", 6379, "port to listen on")
	rdb_dir := flag.String("dir", "", "directory of rdb file")
	dbfilename := flag.String("dbfilename", "", "file name of rdb")
	flag.Parse()

	Server = &ServerConfig{
		Port:       *port,
		Dir:        *rdb_dir,
		DBFilename: *dbfilename,
		MaxMemory:  0,
	}

	fmt.Println(Server)
}

func NewConnHandler(conn net.Conn) *ConnHandler {
	return &ConnHandler{
		conn: conn,
		in:   make(chan protocol.Cmd),
	}
}

func (h *ConnHandler) Handle() {
	defer h.conn.Close()

	go h.readLoop()

	for cmd := range h.in {
		fmt.Printf("[DEBUG] Received command: %s\r\n", cmd)

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
