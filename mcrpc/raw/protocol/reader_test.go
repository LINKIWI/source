package protocol

import (
	"io"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestReaderReadRequestError(t *testing.T) {
	t.Parallel()

	streams := []io.Reader{
		strings.NewReader(""),
		strings.NewReader("none"),
	}

	for _, stream := range streams {
		cmd, err := NewReader(stream).ReadASCIICommand()

		assert.Error(t, err)
		assert.Nil(t, cmd)
	}
}

func TestReaderReadRequestSingleDelimiter(t *testing.T) {
	t.Parallel()

	requests := []Request{
		&VersionRequest{},
		&ShutdownRequest{Type: "graceful"},
		&FlushRequest{Delay: 5 * time.Second},
		&StatsRequest{},
		&StatsRequest{Type: "settings"},
		&WatchRequest{Loggers: []string{"fetchers", "mutations"}},
		&TouchRequest{Key: "key", Expiration: 5 * time.Second},
		&DeleteRequest{Key: "key"},
		&IncrRequest{Key: "key", Delta: 5},
		&DecrRequest{Key: "key", Delta: 5},
		&GetRequest{Keys: []string{"k"}},
		&GetRequest{Keys: []string{"foo", "bar"}},
		&GetsRequest{Keys: []string{"foo", "bar"}},
		&GatRequest{Keys: []string{"k"}, Expiration: 5 * time.Second},
		&GatsRequest{Keys: []string{"k"}, Expiration: 5 * time.Second},
	}

	for _, req := range requests {
		reader := NewReader(strings.NewReader(req.String()))
		cmd, err := reader.ReadASCIICommand()

		assert.NoError(t, err)
		assert.Equal(t, req.String(), string(cmd))
	}
}

func TestReaderReadRequestDoubleDelimiter(t *testing.T) {
	t.Parallel()

	payload := &Storage{
		Key:        "key",
		Flags:      1,
		Expiration: 2 * time.Second,
		Size:       3,
		CasUnique:  4,
		Data:       []byte("data"),
		NoReply:    true,
	}

	requests := []Request{
		&SetRequest{Payload: payload},
		&AddRequest{Payload: payload},
		&ReplaceRequest{Payload: payload},
		&AppendRequest{Payload: payload},
		&PrependRequest{Payload: payload},
		&CasRequest{Payload: payload},
	}

	for _, req := range requests {
		reader := NewReader(strings.NewReader(req.String()))
		cmd, err := reader.ReadASCIICommand()

		assert.NoError(t, err)
		assert.Equal(t, req.String(), string(cmd))
	}
}
