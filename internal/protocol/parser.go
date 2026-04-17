package protocol

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Cmd struct {
	Command string
	Args    []string
}

type value struct {
	Type  byte
	Str   string
	Int   int64
	Array []value
}

func ReadCommand(r *bufio.Reader) (Cmd, error) {
	v, err := ReadResp(r)

	if err != nil {
		return Cmd{}, err
	}

	if v.Type != '*' || len(v.Array) == 0 {
		return Cmd{}, fmt.Errorf("expected array command")
	}

	cmd := Cmd{
		Command: strings.ToUpper(v.Array[0].Str),
		Args:    make([]string, len(v.Array)-1),
	}

	for i := 1; i < len(v.Array); i++ {
		cmd.Args[i-1] = v.Array[i].Str
	}

	return cmd, nil
}

func ReadResp(r *bufio.Reader) (value, error) {
	t, err := r.ReadByte()
	if err != nil {
		return value{}, err
	}

	switch t {
	case '$':
		return readBulkString(r)
	case '*':
		return readArray(r)
	case '+':
		return readSimpleString(r)
	case '-':
		return readError(r)
	case ':':
		return readInteger(r)
	default:
		return value{}, fmt.Errorf("unknown RESP type: %q", t)
	}
}

func readLine(r *bufio.Reader) (string, error) {
	line, err := r.ReadBytes('\n')
	if err != nil {
		return "", err
	}

	length := len(line)

	if length < 2 || line[length-2] != '\r' {
		return "", fmt.Errorf("RESP Err: invalid line ending")
	}

	return string(line[:length-2]), nil
}

func readCount(r *bufio.Reader) (int64, error) {
	line, err := readLine(r)
	if err != nil {
		return -1, err
	}

	count, err := strconv.ParseInt(line, 10, 64)
	if err != nil {
		return -1, err
	}

	return count, nil
}

func readError(r *bufio.Reader) (value, error) {
	line, err := readLine(r)
	if err != nil {
		return value{}, err
	}

	return value{
		Type: '-',
		Str:  line,
	}, nil
}

func readInteger(r *bufio.Reader) (value, error) {
	line, err := readCount(r)
	if err != nil {
		return value{}, err
	}

	return value{
		Type: ':',
		Int:  line,
	}, nil
}

func readSimpleString(r *bufio.Reader) (value, error) {
	line, err := readLine(r)
	if err != nil {
		return value{}, err
	}

	return value{
		Type: '+',
		Str:  line,
	}, nil
}

func readBulkString(r *bufio.Reader) (value, error) {
	count, err := readCount(r)
	if err != nil {
		return value{}, err
	}

	buf := make([]byte, count)
	io.ReadFull(r, buf)
	r.ReadBytes('\n')

	return value{
		Type: '$',
		Str:  string(buf),
	}, nil
}

func readArray(r *bufio.Reader) (value, error) {
	count, err := readCount(r)
	if err != nil {
		return value{}, err
	}

	val := value{
		Type: '*',
	}

	for range count {
		item, err := ReadResp(r)
		if err != nil {
			return value{}, err
		}

		val.Array = append(val.Array, item)
	}

	return val, nil
}
