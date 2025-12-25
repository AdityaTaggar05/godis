package types

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type StreamID struct {
	ms  int64
	seq int64
}

func (id StreamID) toString() string {
	return fmt.Sprintf("%d-%d", id.ms, id.seq)
}

func parseStreamID(st *Stream, id string) (StreamID, error) {
	t := strings.Split(id, "-")
	ms, _ := strconv.ParseInt(t[0], 10, 64)

	if len(t) == 1 {
		return StreamID{ms: ms}, fmt.Errorf("no seq provided")
	}

	var seq int64
	if t[1] == "*" {
		if ms == st.lastID.ms {
			seq = st.lastID.seq + 1
		} else {
			if ms == 0 {
				seq = 1
			} else {
				seq = 0
			}
		}
	} else {
		seq, _ = strconv.ParseInt(t[1], 10, 64)
	}

	return StreamID{ms: ms, seq: seq}, nil
}

func nextStreamID(st *Stream) StreamID {
	now := time.Now().UnixMilli()

	if now == st.lastID.ms {
		st.lastID.seq++
	} else {
		st.lastID.ms = now
		st.lastID.seq = 0
	}
	return st.lastID
}

func (id1 StreamID) compare(id2 StreamID) int {
	if id1.ms < id2.ms || (id1.ms == id2.ms && id1.seq < id2.seq) {
		return -1
	} else if id1.ms == id2.ms && id1.seq == id2.seq {
		return 0
	}
	return 1
}
