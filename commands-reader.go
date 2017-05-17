package gw

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
)

type CommandsReader struct {
	io.Reader
	br *bufio.Reader
}

func NewCommandsReader(src io.Reader) CommandsReader {
	return CommandsReader{
		Reader: src,
		br:     bufio.NewReader(src),
	}
}

func (cr CommandsReader) nextCommand() ([]byte, error) {
	var msg []byte
	line, err := cr.br.ReadBytes('\n')
	if err != nil {
		return nil, err
	}

	if bytes.Equal(line[0:3], []byte("MSG")) {
		msg = line[:]
		splitted := bytes.Split(line, []byte(" "))
		sizeStr := splitted[len(splitted)-1]
		sizeStr = sizeStr[:len(sizeStr)-2]
		size, err := strconv.Atoi(string(sizeStr))
		if err != nil {
			return nil, fmt.Errorf("Error reading MSG size: %s", err)
		}
		// the '-2' is to account for the trailing \r\n which is after the payload
		for size > -2 {
			chunk, err := cr.br.ReadBytes('\n')
			if err != nil {
				return nil, fmt.Errorf("Error reading MSG payload: %s", err)
			}
			size -= len(chunk)
			msg = append(msg, chunk...)
		}
		if size != -2 {
			return nil, fmt.Errorf(
				"Error reading MSG payload. Got %d extra bytes", -size-2)
		}
	} else {
		msg = line
	}
	return msg, nil
}
