package gw

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
)

// CommandsReader parses a NATS connection input stream into commands
type CommandsReader struct {
	io.Reader
	br *bufio.Reader
}

// NewCommandsReader creates a CommandsReader
func NewCommandsReader(src io.Reader) CommandsReader {
	return CommandsReader{
		Reader: src,
		br:     bufio.NewReader(src),
	}
}

// NextCommand returns the next command in the input stream
func (cr CommandsReader) NextCommand() ([]byte, error) {
	return cr.nextCommand()
}

func (cr CommandsReader) nextCommand() ([]byte, error) {
	var msg []byte

	line, err := cr.br.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	for bytes.Equal(line, []byte("\r\n")) {
		line, err = cr.br.ReadBytes('\n')
		if err != nil {
			return nil, err
		}
	}

	if len(line) == 0 {
		return nil, fmt.Errorf("Unexpected empty line")
	}
	if len(line) < 3 {
		return nil, fmt.Errorf("Invalid command: %v", line)
	}
	op := line[0:3]
	if bytes.Equal(op, []byte("MSG")) || bytes.Equal(op, []byte("PUB")) {
		msg = line[:]
		splitted := bytes.Split(line, []byte(" "))
		sizeStr := splitted[len(splitted)-1]
		sizeStr = sizeStr[:len(sizeStr)-2]
		size, err := strconv.Atoi(string(sizeStr))
		if err != nil {
			return nil, fmt.Errorf("Error reading %s size: %s", op, err)
		}
		// the '-2' is to account for the trailing \r\n which is after the payload
		for size > -2 {
			chunk, err := cr.br.ReadBytes('\n')
			if err != nil {
				return nil, fmt.Errorf("Error reading %s payload: %s", op, err)
			}
			size -= len(chunk)
			msg = append(msg, chunk...)
		}
		if size != -2 {
			return nil, fmt.Errorf(
				"Error reading %s payload. Got %d extra bytes", op, -size-2)
		}
	} else {
		msg = line
	}
	return msg, nil
}
