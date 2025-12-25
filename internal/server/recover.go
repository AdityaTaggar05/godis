package server

import (
	"github.com/AdityaTaggar05/godis/internal/protocol"
	"github.com/AdityaTaggar05/godis/internal/types"
)

func (h *ConnHandler) handleRecover(r any) {
	switch r.(type) {
	case types.RedisWrongTypeError:
		h.writeOutput(protocol.EncodeError("WRONGTYPE Operation against a key holding the wrong kind of value"))
	case types.RedisInvalidStreamIDError:
		h.writeOutput(protocol.EncodeError("ERR The ID specified in XADD is equal or smaller than the target stream top item"))
	case types.RedisMinimumStreamIDError:
		h.writeOutput(protocol.EncodeError("ERR The ID specified in XADD must be greater than 0-0"))
	default:
		panic(r)
	}
}
