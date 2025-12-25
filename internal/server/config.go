package server

import (
	"fmt"
	"sync"

	"github.com/AdityaTaggar05/godis/internal/protocol"
)

type ServerConfig struct {
	mu sync.RWMutex

	Port       int
	Dir        string
	DBFilename string
	MaxMemory  int64
}

func (s *ServerConfig) String() string {
	return fmt.Sprintf(
		"Server Configuration\n(PORT): %d\n(DIR): %s\n(DBFilename): %s\n(MaxMemory): %d",
		s.Port, s.Dir, s.DBFilename, s.MaxMemory,
	)
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
