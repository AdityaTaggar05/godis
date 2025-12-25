package types

import (
	"sync"

	"github.com/AdityaTaggar05/godis/internal/protocol"
	"github.com/AdityaTaggar05/godis/internal/store"
)

type StreamEntry struct {
	id     StreamID
	fields map[string]string
}

type Stream struct {
	entries []StreamEntry
	lastID  StreamID
	mu      sync.Mutex
	cond    *sync.Cond
}

func newStream() *Stream {
	st := Stream{
		entries: make([]StreamEntry, 0),
		lastID:  StreamID{ms: 0, seq: 0},
	}
	st.cond = sync.NewCond(&st.mu)
	return &st
}

func GetStream(key string) (*Stream, bool) {
	data, ok := store.Load(key)
	if !ok {
		return nil, false
	}

	if data.Typ != store.StreamType {
		panicRedisWrongType()
	}

	return data.Value.(*Stream), true
}

func (entry StreamEntry) ToArray() []any {
	arr := make([]string, 0)

	for key, value := range entry.fields {
		arr = append(arr, key)
		arr = append(arr, value)
	}
	return []any{entry.id.toString(), arr}
}

func XADD(cmd protocol.Cmd) string {
	if len(cmd.Args)%2 != 0 {
		panicRedisWrongNumArgs("XADD")
	}

	key := cmd.Args[0]

	_, _ = store.LoadOrStore(key, store.Data{Typ: store.StreamType, Value: newStream()})
	st, _ := GetStream(key)

	var id StreamID

	if cmd.Args[1] == "*" {
		id = nextStreamID(st)
	} else {
		id, _ = parseStreamID(st, cmd.Args[1])

		if id.ms == 0 && id.seq == 0 {
			panicRedisMinimumStreamID()
		}

		if id.ms < st.lastID.ms {
			panicRedisInvalidStreamID()
		}

		if id.ms == st.lastID.ms && id.seq <= st.lastID.seq {
			panicRedisInvalidStreamID()
		}
	}

	entry := StreamEntry{
		id:     id,
		fields: make(map[string]string),
	}

	for i := 2; i < len(cmd.Args); i += 2 {
		entry.fields[cmd.Args[i]] = cmd.Args[i+1]
	}

	st.entries = append(st.entries, entry)
	st.lastID = id

	return id.toString()
}

func XRANGE(cmd protocol.Cmd) []StreamEntry {
	if len(cmd.Args) < 3 {
		panicRedisWrongNumArgs("XRANGE")
	}

	st, _ := GetStream(cmd.Args[0])

	start, _ := parseStreamID(st, cmd.Args[1])
	stop, err := parseStreamID(st, cmd.Args[2])
	tillEnd := err != nil

	entries := make([]StreamEntry, 0)

	for _, entry := range st.entries {
		after := entry.id.compare(start) >= 0
		before := entry.id.compare(stop) <= 0 || (entry.id.ms == stop.ms && tillEnd)

		if after && before {
			entries = append(entries, entry)
		}

		if after && !before {
			break
		}
	}

	return entries
}

func XREAD(cmd protocol.Cmd) map[string][]StreamEntry {
	if len(cmd.Args) < 3 {
		panicRedisWrongNumArgs("XREAD")
	}

	data := make(map[string][]StreamEntry)

	query := cmd.Args[1:]
	n := len(query) / 2

	for i := 0; i < n; i += 1 {
		entries := make([]StreamEntry, 0)

		st, _ := GetStream(query[i])
		for _, entry := range st.entries {
			id, _ := parseStreamID(st, query[n+i])

			if entry.id.compare(id) >= 0 {
				entries = append(entries, entry)
			}
		}

		data[query[i]] = entries
	}

	return data
}
