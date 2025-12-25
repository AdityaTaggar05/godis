package server

import (
	"sync"

	"github.com/codecrafters-io/redis-starter-go/app/protocol"
)

type ServerConfig struct {
	mu sync.RWMutex

	Port       int
	Dir        string
	DBFilename string
	MaxMemory  int64
}

var getters = map[string]func() string{
	"dir": func() string {
		return Server.Dir
	},
	"dbfilename": func() string {
		return Server.DBFilename
	},
}

// var setters = map[string]func(string) error{
// 	"dir": func(v string) error {
// 		Server.Dir = v
// 		return nil
// 	},
// 	"dbfilename": func(v string) error {
// 		Server.DBFilename = v
// 		return nil
// 	},
// }

func CONFIG(cmd protocol.Cmd) []byte {
	buf := make([]byte, 0)

	switch cmd.Args[0] {
	case "GET":
		data := make([]string, 0)
		for _, key := range cmd.Args[1:] {
			data = append(data, key)
			data = append(data, getters[key]())
		}
		buf = append(buf, protocol.EncodeMultiBulk(data)...)
		// case "SET":
	}

	return buf
}
