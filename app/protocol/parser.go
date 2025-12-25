package protocol

import (
	"bufio"
	"bytes"
	"fmt"
	"strconv"
)

type Cmd struct {
	Command string
	Args    []string
}

func ReadCommand(reader *bufio.Reader) (*Cmd, error) {
	line, err := reader.ReadBytes('\n')
	if err != nil {
		return nil, err
	}

	length := len(line)

	if length < 2 || line[length-2] != '\r' {
		return nil, fmt.Errorf("invalid command")
	}

	line = bytes.TrimSuffix(line, []byte{'\r', '\n'})

	switch line[0] {
	case '+':
		return parseInlineCommand(string(line[1:])), nil
	case '*':
		return parseMultiBulkCommand(reader, line)
	default:
		args := bytes.Split(line, []byte{' '})
		return &Cmd{Command: string(args[0]), Args: bytesToStrings(args[1:])}, nil
	}
}

func parseInlineCommand(content string) *Cmd {
	return &Cmd{Command: content}
}

func parseBulkString(reader *bufio.Reader, line []byte) (string, error) {
	length, err := parseCount(line)
	if err != nil {
		return "", err
	}

	if length == 0 {
		return "", nil
	}

	content := make([]byte, length+2)
	_, err = reader.Read(content)
	if err != nil {
		return "", err
	}

	if content[length] != '\r' || content[length+1] != '\n' {
		return "", fmt.Errorf("invalid command")
	}

	return string(content[:length]), nil
}

func parseMultiBulkCommand(reader *bufio.Reader, line []byte) (*Cmd, error) {
	count, err := parseCount(line)
	if err != nil {
		return nil, err
	}

	cmd := Cmd{}
	cmd.Args = make([]string, count)

	for i := 0; i < count; i++ {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			return nil, err
		}

		length := len(line)

		if length < 2 || line[length-2] != '\r' {
			return nil, fmt.Errorf("invalid command")
		}

		line = bytes.TrimSuffix(line, []byte{'\r', '\n'})

		switch line[0] {
		case '$':
			content, err := parseBulkString(reader, line)
			if err != nil {
				return nil, err
			}
			cmd.Args[i] = content
		default:
			return nil, fmt.Errorf("invalid command")
		}
	}

	cmd.Command = cmd.Args[0]
	cmd.Args = cmd.Args[1:]

	return &cmd, nil
}

func parseCount(line []byte) (int, error) {
	count, err := parseInteger(line)
	if err != nil {
		return 0, err
	}

	if count < 0 {
		return 0, fmt.Errorf("invalid command")
	}

	return count, nil
}

func parseInteger(line []byte) (int, error) {
	return strconv.Atoi(string(line[1:]))
}

func bytesToStrings(bytes [][]byte) []string {
	strings := make([]string, len(bytes))
	for i, b := range bytes {
		strings[i] = string(b)
	}
	return strings
}
