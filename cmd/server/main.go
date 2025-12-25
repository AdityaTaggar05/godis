package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/AdityaTaggar05/godis/internal/server"
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	fmt.Println("Logs from your program will appear here!")

	server.Start()

	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", server.Server.Port))
	if err != nil {
		fmt.Printf("Failed to bind to port %d\n", server.Server.Port)
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection")
			os.Exit(1)
		}
		handler := server.NewConnHandler(conn)

		go handler.Handle()
	}
}
