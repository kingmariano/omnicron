package groq

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
)

const minBufferSize = 4096
const maxBufferSize = 1 << 16

var (
	headerID    = []byte("id:")
	headerData  = []byte("data:")
	headerEvent = []byte("event:")
	headerRetry = []byte("retry:")
)

type StreamReader struct {
	scanner *bufio.Scanner
}

type SSEvent struct {
	ID    []byte
	Data  []byte
	Event []byte
	Retry []byte
}

func (sse *SSEvent) String() string {
	return fmt.Sprintf(
		"SSEvent{ID: %s, Data: %s, Event: %s, Retry: %s}",
		sse.ID, sse.Data, sse.Event, sse.Retry,
	)
}

func NewStreamReader(eventStream io.Reader) *StreamReader {
	scanner := bufio.NewScanner(eventStream)
	scanner.Buffer(make([]byte, minBufferSize), maxBufferSize)

	split := func(data []byte, atEOF bool) (int, []byte, error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}

		if i, nlen := splitHasDoubleNewline(data); i >= 0 {
			return i + nlen, data[0:i], nil
		}
		if atEOF {
			return len(data), data, nil
		}
		return 0, nil, nil
	}
	scanner.Split(split)

	return &StreamReader{
		scanner: scanner,
	}
}

func SSEventFromBytes(event []byte) (*SSEvent, error) {
	sse := &SSEvent{}
	for _, line := range bytes.FieldsFunc(event, func(r rune) bool { return r == '\n' || r == '\r' }) {
		switch {
		case bytes.HasPrefix(line, headerID):
			sse.ID = append([]byte(nil), trimHeader(len(headerID), line)...)
		case bytes.HasPrefix(line, headerData):
			sse.Data = append(sse.Data[:], append(trimHeader(len(headerData), line), byte('\n'))...)
		case bytes.Equal(line, bytes.TrimSuffix(headerData, []byte(":"))):
			sse.Data = append(sse.Data, byte('\n'))
		case bytes.HasPrefix(line, headerEvent):
			sse.Event = append([]byte(nil), trimHeader(len(headerEvent), line)...)
		case bytes.HasPrefix(line, headerRetry):
			sse.Retry = append([]byte(nil), trimHeader(len(headerRetry), line)...)
		default:
			// TODO return error?
		}
	}

	return sse, nil
}

func (e *StreamReader) Next() (*SSEvent, error) {
	if e.scanner.Scan() {
		event := e.scanner.Bytes()
		return SSEventFromBytes(event)
	}
	if err := e.scanner.Err(); err != nil {
		if err == context.Canceled {
			return nil, io.EOF
		}
		return nil, err
	}
	return nil, io.EOF
}

func splitHasDoubleNewline(data []byte) (int, int) {
	patterns := []struct {
		pattern []byte
		length  int
	}{
		{[]byte("\n\n"), 2},
		{[]byte("\r\r"), 2},
		{[]byte("\n\r\n"), 3},
		{[]byte("\r\n\n"), 3},
		{[]byte("\r\n\r\n"), 4},
	}

	// Initialize minPos to -1 (no match found yet)
	minPos := -1

	// Loop over each pattern
	for _, p := range patterns {
		pos := bytes.Index(data, p.pattern)
		if pos != -1 && (minPos == -1 || pos < minPos) {
			minPos = pos
			// If one of shortest patterns is found, return immediately with its length
			if p.length == 2 {
				return pos, 2
			}
			return pos, p.length
		}
	}

	// Return -1 (no match found)
	return -1, 0
}

func trimHeader(size int, data []byte) []byte {
	if data == nil || len(data) < size {
		return data
	}

	data = data[size:]
	// Remove optional leading whitespace
	if len(data) > 0 && data[0] == 32 {
		data = data[1:]
	}
	// Remove trailing new line
	if len(data) > 0 && data[len(data)-1] == 10 {
		data = data[:len(data)-1]
	}
	return data
}
