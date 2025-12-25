package types

import (
	"github.com/codecrafters-io/redis-starter-go/app/protocol"
	"github.com/codecrafters-io/redis-starter-go/app/store"
)

func TYPE(cmd protocol.Cmd) string {
	if len(cmd.Args) != 1 {
		panicRedisWrongNumArgs("TYPE")
	}

	data, ok := store.Load(cmd.Args[0])
	if !ok {
		return "none"
	}

	switch data.Typ {
	case store.StringType:
		return "string"
	case store.ListType:
		return "list"
	case store.StreamType:
		return "stream"
	default:
		return "none"
	}
}
