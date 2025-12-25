# godis ⚡  
`godis` is a lightweight Redis clone implemented from scratch in **Go**, designed to explore **TCP servers, concurrency safety, correct Redis semantics and clean internal architecture**.

It implements a meaningful subset of Redis commands while preserving correct semantics such as **FIFO ordering**, **blocking behaviour**, and **safe concurrent access**.

---

## 🚀 Features

### Core
- Built as a **TCP server**
- Redis-compatible **RESP protocol**
- Concurrent client handling
- Safe access using `mutexes`, `condition variables`, and `sync.Map`
- Redis-style command dispatch

---

## 📦 Supported Data Types & Commands

### 🔑 Strings
- `SET`
  - Supports expiration:
    - `EX` (seconds)
    - `PX` (milliseconds)
- `GET`

---

### 📃 Lists
- `LPUSH`
- `RPUSH`
- `LPOP`
- `BLPOP` *(blocking with timeout, FIFO-safe)*
- `LLEN`
- `LRANGE`

---

### 🌊 Streams
- `XADD`
- `XRANGE`
- `XREAD`

Streams are implemented with:
- Blocking reads where applicable
- Auto-generated IDs (*)
- Strict ID ordering enforcement

---

### ⚙️ Server / Utility Commands
- `TYPE`
- `CONFIG GET / SET`
- `PING`, `ECHO`, `QUIT`

---

## 🧠 Concurrency Model

`godis` is built to handle **multiple concurrent clients safely and efficiently**.

### Techniques Used:
- `sync.Map` for the global keyspace (read-heavy workload)
- Fine-grained mutexes per data structure:
  - Lists & Streams have their own locks
- Blocking commands (BLPOP, XREAD) use:
  - Waiter queues
  - Condition variables
- FIFO fairness guaranteed for blocked clients

---

## 📜 Demo

This project is fully compatible with the official Redis CLI, allowing it to be demonstrated exactly like a real Redis server.
First, ensure that redis-cli is installed on your system. For Windows, unfortunately, there is no direct way to use it other than WSL or a VM.

Once that's taken care of, start the server with
```
go run ./cmd/server/
```

Then, in another terminal window, connect using Redis CLI and use the server:
```
redis-cli -p 6379
```

Or run the demo script using
```
./demo.sh
```
