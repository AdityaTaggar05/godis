package server

import (
	"flag"
	"fmt"
	"sync"
)

var DB sync.Map
var Server *ServerConfig

func Start() {
	DB = sync.Map{}

	port := flag.Int("port", 6379, "port to listen on")
	rdbDir := flag.String("dir", "", "directory of rdb file")
	dbfilename := flag.String("dbfilename", "", "file name of rdb")
	flag.Parse()

	Server = &ServerConfig{
		Port:       *port,
		Dir:        *rdbDir,
		DBFilename: *dbfilename,
		MaxMemory:  0,
	}

	fmt.Println(Server)
}
