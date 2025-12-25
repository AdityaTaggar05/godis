package types

import (
	"strconv"
	"time"

	"github.com/AdityaTaggar05/godis/internal/protocol"
	"github.com/AdityaTaggar05/godis/internal/store"
)

func GetString(key string) (string, bool) {
	data, ok := store.Load(key)
	if !ok {
		return "", false
	}

	if data.Typ != store.StringType {
		panicRedisWrongType()
	}

	return data.Value.(string), true
}

func GET(cmd protocol.Cmd) string {
	if len(cmd.Args) != 1 {
		panicRedisWrongNumArgs("GET")
	}
	str, _ := GetString(cmd.Args[0])
	return str
}

func SET(cmd protocol.Cmd) {
	if len(cmd.Args) < 2 {
		panicRedisWrongNumArgs("SET")
	}

	data := store.Data{
		Typ:   store.StringType,
		Value: cmd.Args[1],
	}

	if len(cmd.Args) > 2 {
		n, _ := strconv.Atoi(cmd.Args[3])

		switch cmd.Args[2] {
		case "EX":
			data.ExpiresAt = time.Now().Add(time.Duration(n) * time.Second)
		case "PX":
			data.ExpiresAt = time.Now().Add(time.Duration(n) * time.Millisecond)
		}
	}

	store.Store(cmd.Args[0], data)
}
