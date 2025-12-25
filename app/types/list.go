package types

import (
	"strconv"
	"sync"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/protocol"
	"github.com/codecrafters-io/redis-starter-go/app/store"
)

type waiter struct {
	ch chan string
}

func removeWaiter(waiters []*waiter, w *waiter) []*waiter {
	for i, x := range waiters {
		if x == w {
			return append(waiters[:i], waiters[i+1:]...)
		}
	}
	return waiters
}

type List struct {
	items   []string
	waiters []*waiter
	mu      sync.Mutex
}

func newList() *List {
	return &List{
		items:   make([]string, 0),
		waiters: make([]*waiter, 0),
	}
}

func GetList(key string) (*List, bool) {
	data, ok := store.Load(key)
	if !ok {
		return nil, false
	}

	if data.Typ != store.ListType {
		panicRedisWrongType()
	}

	return data.Value.(*List), true
}

func LPOP(cmd protocol.Cmd) []string {
	if len(cmd.Args) < 1 {
		return []string{}
	}

	l, found := GetList(cmd.Args[0])
	count := 1

	if !found {
		return []string{}
	}

	if len(cmd.Args) == 2 {
		count, _ = strconv.Atoi(cmd.Args[1])
	}

	var removed []string

	l.mu.Lock()
	defer l.mu.Unlock()

	if len(l.items) == 0 {
		return l.items
	}

	if count > len(l.items) {
		count = len(l.items)
	}

	removed = l.items[:count]
	l.items = l.items[count:]
	return removed
}

func RPUSH(cmd protocol.Cmd) int {
	if len(cmd.Args) < 2 {
		panicRedisWrongNumArgs("RPUSH")
	}

	_, _ = store.LoadOrStore(cmd.Args[0], store.Data{Typ: store.ListType, Value: newList()})

	l, _ := GetList(cmd.Args[0])
	l.mu.Lock()
	defer l.mu.Unlock()

	n := 0

	for _, val := range cmd.Args[1:] {
		if len(l.waiters) > 0 {
			w := l.waiters[0]
			l.waiters = l.waiters[1:]
			w.ch <- val
			n++
		} else {
			l.items = append(l.items, val)
		}
	}

	return len(l.items) + n
}

func LPUSH(cmd protocol.Cmd) int {
	if len(cmd.Args) < 2 {
		panicRedisWrongNumArgs("LPUSH")
	}

	_, _ = store.LoadOrStore(cmd.Args[0], store.Data{Typ: store.ListType, Value: newList()})
	l, _ := GetList(cmd.Args[0])
	l.mu.Lock()
	defer l.mu.Unlock()

	n := 0

	for _, val := range cmd.Args[1:] {
		if len(l.waiters) > 0 {
			w := l.waiters[0]
			l.waiters = l.waiters[1:]
			w.ch <- val
			n++
		} else {
			l.items = append([]string{val}, l.items...)
		}
	}

	return len(l.items) + n
}

func LLEN(cmd protocol.Cmd) int {
	list, found := GetList(cmd.Args[0])
	if found {
		return len(list.items)
	}
	return 0
}

func LRANGE(cmd protocol.Cmd) []string {
	if len(cmd.Args) < 3 {
		panicRedisWrongNumArgs("LRANGE")
	}

	start, _ := strconv.Atoi(cmd.Args[1])
	stop, _ := strconv.Atoi(cmd.Args[2])

	l, found := GetList(cmd.Args[0])

	if !found {
		return []string{}
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	n := len(l.items)

	if start < 0 {
		start = max(n+start, 0)
	}
	if stop < 0 {
		stop = max(n+stop, 0)
	}
	if start >= n || start > stop {
		return []string{}
	}
	if stop >= n {
		stop = n - 1
	}

	result := make([]string, stop-start+1)
	copy(result, l.items[start:stop+1])
	return result
}

func BLPOP(cmd protocol.Cmd) []string {
	if len(cmd.Args) < 2 {
		panicRedisWrongNumArgs("BLPOP")
	}

	key := cmd.Args[0]

	_, _ = store.LoadOrStore(key, store.Data{Typ: store.ListType, Value: newList()})
	l, _ := GetList(key)
	l.mu.Lock()

	if len(l.items) > 0 {
		elem := l.items[0]
		l.items = l.items[1:]
		l.mu.Unlock()
		return []string{key, elem}
	}

	w := &waiter{ch: make(chan string, 1)}
	l.waiters = append(l.waiters, w)
	l.mu.Unlock()

	timelimit, _ := strconv.ParseFloat(cmd.Args[1], 32)
	if timelimit == 0 {
		elem := <-w.ch
		return []string{key, elem}
	}

	timer := time.NewTimer(time.Duration(timelimit*1000) * time.Millisecond)
	defer timer.Stop()

	select {
	case elem := <-w.ch:
		return []string{key, elem}
	case <-timer.C:
		l.mu.Lock()
		l.waiters = removeWaiter(l.waiters, w)
		l.mu.Unlock()
		return nil
	}
}
