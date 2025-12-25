package server

import (
	"fmt"

	"github.com/AdityaTaggar05/godis/internal/protocol"
	"github.com/AdityaTaggar05/godis/internal/types"
)

func (h *ConnHandler) dispatch(cmd protocol.Cmd) {
	switch cmd.Command {
	case "PING":
		h.writeOutput(protocol.EncodeSimple("PONG"))
	case "ECHO":
		h.writeOutput(protocol.EncodeBulk(cmd.Args[0]))
	case "SET":
		types.SET(cmd)
		h.writeOutput(protocol.EncodeSimple("OK"))
	case "GET":
		h.writeOutput(protocol.EncodeBulk(types.GET(cmd)))
	case "RPUSH":
		h.writeOutput(protocol.EncodeInteger(types.RPUSH(cmd)))
	case "LPUSH":
		h.writeOutput(protocol.EncodeInteger(types.LPUSH(cmd)))
	case "LLEN":
		h.writeOutput(protocol.EncodeInteger(types.LLEN(cmd)))
	case "LPOP":
		removed := types.LPOP(cmd)
		var resp []byte

		if len(removed) == 0 {
			resp = protocol.EncodeBulk("")
		} else if len(removed) == 1 {
			resp = protocol.EncodeBulk(removed[0])
		} else {
			resp = protocol.EncodeMultiBulk(removed)
		}

		h.writeOutput(resp)
	case "BLPOP":
		h.writeOutput(protocol.EncodeMultiBulk(types.BLPOP(cmd)))
	case "LRANGE":
		h.writeOutput(protocol.EncodeMultiBulk(types.LRANGE(cmd)))
	case "XADD":
		h.writeOutput(protocol.EncodeBulk(types.XADD(cmd)))
	case "XRANGE":
		entries := types.XRANGE(cmd)
		data := make([]any, 0)

		for _, e := range entries {
			data = append(data, e.ToArray())
		}
		h.writeOutput(protocol.EncodeArray(data))
	case "XREAD":
		stEntries := types.XREAD(cmd)
		data := make([]any, 0)

		for k, st := range stEntries {
			stEntry := make([]any, 0)
			stEntry = append(stEntry, k)

			entries := make([]any, 0)

			for _, e := range st {
				entries = append(entries, e.ToArray())
			}

			stEntry = append(stEntry, entries)
			data = append(data, stEntry)
		}
		h.writeOutput(protocol.EncodeArray(data))
	case "CONFIG":
		h.writeOutput(CONFIG(cmd))
	case "QUIT":
		h.writeOutput(protocol.EncodeSimple("OK"))
	case "TYPE":
		h.writeOutput(protocol.EncodeSimple(types.TYPE(cmd)))
	default:
		fmt.Printf("[DEBUG] Unknown command: %v\r\n", cmd)
		h.writeOutput(protocol.EncodeError(fmt.Sprintf("ERR unknown command: %v", cmd)))
	}
}
