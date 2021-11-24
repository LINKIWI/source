package protocol

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseVersion(t *testing.T) {
	t.Parallel()

	cases := []struct {
		command string
		request *VersionRequest
	}{
		{
			command: "version",
			request: nil,
		},
		{
			command: "version\n",
			request: nil,
		},
		{
			command: "version extra\r\n",
			request: nil,
		},
		{
			command: "version\r\n",
			request: &VersionRequest{},
		},
	}

	for _, testCase := range cases {
		parsed, err := NewASCIIParser().Parse([]byte(testCase.command))

		if testCase.request != nil {
			assert.NoError(t, err)
			assert.Equal(t, parsed, testCase.request)
		} else {
			assert.Error(t, err)
		}
	}
}

func TestParseShutdown(t *testing.T) {
	t.Parallel()

	cases := []struct {
		command string
		request *ShutdownRequest
	}{
		{
			command: "shutdown",
			request: nil,
		},
		{
			command: "shutdown foo bar\r\n",
			request: nil,
		},
		{
			command: "shutdown\r\n",
			request: &ShutdownRequest{},
		},
		{
			command: "shutdown graceful\r\n",
			request: &ShutdownRequest{Type: "graceful"},
		},
	}

	for _, testCase := range cases {
		parsed, err := NewASCIIParser().Parse([]byte(testCase.command))

		if testCase.request != nil {
			assert.NoError(t, err)
			assert.Equal(t, parsed, testCase.request)
		} else {
			assert.Error(t, err)
		}
	}
}

func TestParseFlush(t *testing.T) {
	t.Parallel()

	cases := []struct {
		command string
		request *FlushRequest
	}{
		{
			command: "flush_all",
			request: nil,
		},
		{
			command: "flush_all \r\n",
			request: nil,
		},
		{
			command: "flush_all foo\r\n",
			request: nil,
		},
		{
			command: "flush_all\r\n",
			request: &FlushRequest{},
		},
		{
			command: "flush_all 5\r\n",
			request: &FlushRequest{Delay: 5 * time.Second},
		},
	}

	for _, testCase := range cases {
		parsed, err := NewASCIIParser().Parse([]byte(testCase.command))

		if testCase.request != nil {
			assert.NoError(t, err)
			assert.Equal(t, parsed, testCase.request)
		} else {
			assert.Error(t, err)
		}
	}
}

func TestParseQuit(t *testing.T) {
	t.Parallel()

	cases := []struct {
		command string
		request *QuitRequest
	}{
		{
			command: "quit",
			request: nil,
		},
		{
			command: "quit \r\n",
			request: nil,
		},
		{
			command: "quit foo\r\n",
			request: nil,
		},
		{
			command: "quit\r\n",
			request: &QuitRequest{},
		},
	}

	for _, testCase := range cases {
		parsed, err := NewASCIIParser().Parse([]byte(testCase.command))

		if testCase.request != nil {
			assert.NoError(t, err)
			assert.Equal(t, parsed, testCase.request)
		} else {
			assert.Error(t, err)
		}
	}
}

func TestParseStats(t *testing.T) {
	t.Parallel()

	cases := []struct {
		command string
		request *StatsRequest
	}{
		{
			command: "stats",
			request: nil,
		},
		{
			command: "stats \r\n",
			request: nil,
		},
		{
			command: "stats 1\r\n",
			request: nil,
		},
		{
			command: "stats foo bar\r\n",
			request: nil,
		},
		{
			command: "stats\r\n",
			request: &StatsRequest{},
		},
		{
			command: "stats settings\r\n",
			request: &StatsRequest{Type: "settings"},
		},
		{
			command: "stats   settings\r\n",
			request: &StatsRequest{Type: "settings"},
		},
	}

	for _, testCase := range cases {
		parsed, err := NewASCIIParser().Parse([]byte(testCase.command))

		if testCase.request != nil {
			assert.NoError(t, err)
			assert.Equal(t, parsed, testCase.request)
		} else {
			assert.Error(t, err)
		}
	}
}

func TestParseWatch(t *testing.T) {
	t.Parallel()

	cases := []struct {
		command string
		request *WatchRequest
	}{
		{
			command: "watch",
			request: nil,
		},
		{
			command: "watch \r\n",
			request: nil,
		},
		{
			command: "watch\r\n",
			request: &WatchRequest{},
		},
		{
			command: "watch foo\r\n",
			request: &WatchRequest{Loggers: []string{"foo"}},
		},
		{
			command: "watch foo bar\r\n",
			request: &WatchRequest{Loggers: []string{"foo", "bar"}},
		},
	}

	for _, testCase := range cases {
		parsed, err := NewASCIIParser().Parse([]byte(testCase.command))

		if testCase.request != nil {
			assert.NoError(t, err)
			assert.Equal(t, parsed, testCase.request)
		} else {
			assert.Error(t, err)
		}
	}
}

func TestParseTouch(t *testing.T) {
	t.Parallel()

	cases := []struct {
		command string
		request *TouchRequest
	}{
		{
			command: "touch",
			request: nil,
		},
		{
			command: "touch\r\n",
			request: nil,
		},
		{
			command: "touch key\r\n",
			request: nil,
		},
		{
			command: "touch key foo\r\n",
			request: nil,
		},
		{
			command: "touch key 5\r\n",
			request: &TouchRequest{
				Key:        "key",
				Expiration: 5 * time.Second,
				NoReply:    false,
			},
		},
		{
			command: "touch key 5 noreply\r\n",
			request: &TouchRequest{
				Key:        "key",
				Expiration: 5 * time.Second,
				NoReply:    true,
			},
		},
	}

	for _, testCase := range cases {
		parsed, err := NewASCIIParser().Parse([]byte(testCase.command))

		if testCase.request != nil {
			assert.NoError(t, err)
			assert.Equal(t, parsed, testCase.request)
		} else {
			assert.Error(t, err)
		}
	}
}

func TestParseDelete(t *testing.T) {
	t.Parallel()

	cases := []struct {
		command string
		request *DeleteRequest
	}{
		{
			command: "delete",
			request: nil,
		},
		{
			command: "delete\r\n",
			request: nil,
		},
		{
			command: "delete key foo\r\n",
			request: nil,
		},
		{
			command: "delete key\r\n",
			request: &DeleteRequest{Key: "key", NoReply: false},
		},
		{
			command: "delete key noreply\r\n",
			request: &DeleteRequest{Key: "key", NoReply: true},
		},
	}

	for _, testCase := range cases {
		parsed, err := NewASCIIParser().Parse([]byte(testCase.command))

		if testCase.request != nil {
			assert.NoError(t, err)
			assert.Equal(t, parsed, testCase.request)
		} else {
			assert.Error(t, err)
		}
	}
}

func TestParseIncr(t *testing.T) {
	t.Parallel()

	cases := []struct {
		command string
		request *IncrRequest
	}{
		{
			command: "incr",
			request: nil,
		},
		{
			command: "incr key foo\r\n",
			request: nil,
		},
		{
			command: "incr key 5\r\n",
			request: &IncrRequest{
				Key:     "key",
				Delta:   5,
				NoReply: false,
			},
		},
		{
			command: "incr key 5 noreply\r\n",
			request: &IncrRequest{
				Key:     "key",
				Delta:   5,
				NoReply: true,
			},
		},
	}

	for _, testCase := range cases {
		parsed, err := NewASCIIParser().Parse([]byte(testCase.command))

		if testCase.request != nil {
			assert.NoError(t, err)
			assert.Equal(t, parsed, testCase.request)
		} else {
			assert.Error(t, err)
		}
	}
}

func TestParseDecr(t *testing.T) {
	t.Parallel()

	cases := []struct {
		command string
		request *DecrRequest
	}{
		{
			command: "decr",
			request: nil,
		},
		{
			command: "decr key foo\r\n",
			request: nil,
		},
		{
			command: "decr key 5\r\n",
			request: &DecrRequest{
				Key:     "key",
				Delta:   5,
				NoReply: false,
			},
		},
		{
			command: "decr key 5 noreply\r\n",
			request: &DecrRequest{
				Key:     "key",
				Delta:   5,
				NoReply: true,
			},
		},
	}

	for _, testCase := range cases {
		parsed, err := NewASCIIParser().Parse([]byte(testCase.command))

		if testCase.request != nil {
			assert.NoError(t, err)
			assert.Equal(t, parsed, testCase.request)
		} else {
			assert.Error(t, err)
		}
	}
}

func TestParseGet(t *testing.T) {
	t.Parallel()

	cases := []struct {
		command string
		request *GetRequest
	}{
		{
			command: "get",
			request: nil,
		},
		{
			command: "get\r\n",
			request: nil,
		},
		{
			command: "get \r\n",
			request: nil,
		},
		{
			command: "get foo",
			request: nil,
		},
		{
			command: "get foo\r\n",
			request: &GetRequest{Keys: []string{"foo"}},
		},
		{
			command: "get foo bar\r\n",
			request: &GetRequest{Keys: []string{"foo", "bar"}},
		},
		{
			command: "get foo bar baz\r\n",
			request: &GetRequest{Keys: []string{"foo", "bar", "baz"}},
		},
	}

	for _, testCase := range cases {
		parsed, err := NewASCIIParser().Parse([]byte(testCase.command))

		if testCase.request != nil {
			assert.NoError(t, err)
			assert.Equal(t, parsed, testCase.request)
		} else {
			assert.Error(t, err)
		}
	}
}

func TestParseGets(t *testing.T) {
	t.Parallel()

	cases := []struct {
		command string
		request *GetsRequest
	}{
		{
			command: "gets",
			request: nil,
		},
		{
			command: "gets\r\n",
			request: nil,
		},
		{
			command: "gets \r\n",
			request: nil,
		},
		{
			command: "gets foo",
			request: nil,
		},
		{
			command: "gets foo\r\n",
			request: &GetsRequest{Keys: []string{"foo"}},
		},
		{
			command: "gets foo bar\r\n",
			request: &GetsRequest{Keys: []string{"foo", "bar"}},
		},
		{
			command: "gets foo bar baz\r\n",
			request: &GetsRequest{Keys: []string{"foo", "bar", "baz"}},
		},
	}

	for _, testCase := range cases {
		parsed, err := NewASCIIParser().Parse([]byte(testCase.command))

		if testCase.request != nil {
			assert.NoError(t, err)
			assert.Equal(t, parsed, testCase.request)
		} else {
			assert.Error(t, err)
		}
	}
}

func TestParseGat(t *testing.T) {
	t.Parallel()

	cases := []struct {
		command string
		request *GatRequest
	}{
		{
			command: "gat",
			request: nil,
		},
		{
			command: "gat\r\n",
			request: nil,
		},
		{
			command: "gat \r\n",
			request: nil,
		},
		{
			command: "gat foo",
			request: nil,
		},
		{
			command: "gat foo\r\n",
			request: nil,
		},
		{
			command: "gat 5\r\n",
			request: nil,
		},
		{
			command: "gat 5 foo\r\n",
			request: &GatRequest{
				Expiration: 5 * time.Second,
				Keys:       []string{"foo"},
			},
		},
		{
			command: "gat 5 foo bar\r\n",
			request: &GatRequest{
				Expiration: 5 * time.Second,
				Keys:       []string{"foo", "bar"},
			},
		},
		{
			command: "gat 5 foo bar baz\r\n",
			request: &GatRequest{
				Expiration: 5 * time.Second,
				Keys:       []string{"foo", "bar", "baz"},
			},
		},
	}

	for _, testCase := range cases {
		parsed, err := NewASCIIParser().Parse([]byte(testCase.command))

		if testCase.request != nil {
			assert.NoError(t, err)
			assert.Equal(t, parsed, testCase.request)
		} else {
			assert.Error(t, err)
		}
	}
}

func TestParseGats(t *testing.T) {
	t.Parallel()

	cases := []struct {
		command string
		request *GatsRequest
	}{
		{
			command: "gats",
			request: nil,
		},
		{
			command: "gats\r\n",
			request: nil,
		},
		{
			command: "gats \r\n",
			request: nil,
		},
		{
			command: "gats foo",
			request: nil,
		},
		{
			command: "gats foo\r\n",
			request: nil,
		},
		{
			command: "gats 5\r\n",
			request: nil,
		},
		{
			command: "gats 5 foo\r\n",
			request: &GatsRequest{
				Expiration: 5 * time.Second,
				Keys:       []string{"foo"},
			},
		},
		{
			command: "gats 5 foo bar\r\n",
			request: &GatsRequest{
				Expiration: 5 * time.Second,
				Keys:       []string{"foo", "bar"},
			},
		},
		{
			command: "gats 5 foo bar baz\r\n",
			request: &GatsRequest{
				Expiration: 5 * time.Second,
				Keys:       []string{"foo", "bar", "baz"},
			},
		},
	}

	for _, testCase := range cases {
		parsed, err := NewASCIIParser().Parse([]byte(testCase.command))

		if testCase.request != nil {
			assert.NoError(t, err)
			assert.Equal(t, parsed, testCase.request)
		} else {
			assert.Error(t, err)
		}
	}
}

func TestParseSet(t *testing.T) {
	t.Parallel()

	cases := []struct {
		command string
		request *SetRequest
	}{
		{
			command: "set",
			request: nil,
		},
		{
			command: "set key flags 2 3\r\nfoo\r\n",
			request: nil,
		},
		{
			command: "set key 1 exptime 3\r\nfoo\r\n",
			request: nil,
		},
		{
			command: "set key 1 2 size\r\nfoo\r\n",
			request: nil,
		},
		{
			command: "set key 1 2 3\nfoo\r\n",
			request: nil,
		},
		{
			command: "set key 1 2 3\r\nfoo",
			request: nil,
		},
		{
			command: "set key 1 2 3\r\n",
			request: nil,
		},
		{
			command: "set k ey 1 2 3\r\nfoo\r\n",
			request: nil,
		},
		{
			command: "set 1 2 3\r\nfoo\r\n",
			request: nil,
		},
		{
			command: "set key 1 2 3\r\nfoo\r\n",
			request: &SetRequest{
				Payload: &Storage{
					Key:        "key",
					Flags:      1,
					Expiration: 2 * time.Second,
					Size:       3,
					Data:       []byte("foo"),
					NoReply:    false,
				},
			},
		},
		{
			command: "set key 1 2 3 noreply\r\nfoo\r\n",
			request: &SetRequest{
				Payload: &Storage{
					Key:        "key",
					Flags:      1,
					Expiration: 2 * time.Second,
					Size:       3,
					Data:       []byte("foo"),
					NoReply:    true,
				},
			},
		},
	}

	for _, testCase := range cases {
		parsed, err := NewASCIIParser().Parse([]byte(testCase.command))

		if testCase.request != nil {
			assert.NoError(t, err)
			assert.Equal(t, parsed, testCase.request)
		} else {
			assert.Error(t, err)
		}
	}
}

func TestParseReplace(t *testing.T) {
	t.Parallel()

	cases := []struct {
		command string
		request *ReplaceRequest
	}{
		{
			command: "replace",
			request: nil,
		},
		{
			command: "replace key flags 2 3\r\nfoo\r\n",
			request: nil,
		},
		{
			command: "replace key 1 exptime 3\r\nfoo\r\n",
			request: nil,
		},
		{
			command: "replace key 1 2 size\r\nfoo\r\n",
			request: nil,
		},
		{
			command: "replace key 1 2 3\nfoo\r\n",
			request: nil,
		},
		{
			command: "replace key 1 2 3\r\nfoo",
			request: nil,
		},
		{
			command: "replace key 1 2 3\r\n",
			request: nil,
		},
		{
			command: "replace k ey 1 2 3\r\nfoo\r\n",
			request: nil,
		},
		{
			command: "replace 1 2 3\r\nfoo\r\n",
			request: nil,
		},
		{
			command: "replace key 1 2 3\r\nfoo\r\n",
			request: &ReplaceRequest{
				Payload: &Storage{
					Key:        "key",
					Flags:      1,
					Expiration: 2 * time.Second,
					Size:       3,
					Data:       []byte("foo"),
					NoReply:    false,
				},
			},
		},
		{
			command: "replace key 1 2 3 noreply\r\nfoo\r\n",
			request: &ReplaceRequest{
				Payload: &Storage{
					Key:        "key",
					Flags:      1,
					Expiration: 2 * time.Second,
					Size:       3,
					Data:       []byte("foo"),
					NoReply:    true,
				},
			},
		},
	}

	for _, testCase := range cases {
		parsed, err := NewASCIIParser().Parse([]byte(testCase.command))

		if testCase.request != nil {
			assert.NoError(t, err)
			assert.Equal(t, parsed, testCase.request)
		} else {
			assert.Error(t, err)
		}
	}
}

func TestParseAppend(t *testing.T) {
	t.Parallel()

	cases := []struct {
		command string
		request *AppendRequest
	}{
		{
			command: "append",
			request: nil,
		},
		{
			command: "append key flags 2 3\r\nfoo\r\n",
			request: nil,
		},
		{
			command: "append key 1 exptime 3\r\nfoo\r\n",
			request: nil,
		},
		{
			command: "append key 1 2 size\r\nfoo\r\n",
			request: nil,
		},
		{
			command: "append key 1 2 3\nfoo\r\n",
			request: nil,
		},
		{
			command: "append key 1 2 3\r\nfoo",
			request: nil,
		},
		{
			command: "append key 1 2 3\r\n",
			request: nil,
		},
		{
			command: "append k ey 1 2 3\r\nfoo\r\n",
			request: nil,
		},
		{
			command: "append 1 2 3\r\nfoo\r\n",
			request: nil,
		},
		{
			command: "append key 1 2 3\r\nfoo\r\n",
			request: &AppendRequest{
				Payload: &Storage{
					Key:        "key",
					Flags:      1,
					Expiration: 2 * time.Second,
					Size:       3,
					Data:       []byte("foo"),
					NoReply:    false,
				},
			},
		},
		{
			command: "append key 1 2 3 noreply\r\nfoo\r\n",
			request: &AppendRequest{
				Payload: &Storage{
					Key:        "key",
					Flags:      1,
					Expiration: 2 * time.Second,
					Size:       3,
					Data:       []byte("foo"),
					NoReply:    true,
				},
			},
		},
	}

	for _, testCase := range cases {
		parsed, err := NewASCIIParser().Parse([]byte(testCase.command))

		if testCase.request != nil {
			assert.NoError(t, err)
			assert.Equal(t, parsed, testCase.request)
		} else {
			assert.Error(t, err)
		}
	}
}

func TestParsePrepend(t *testing.T) {
	t.Parallel()

	cases := []struct {
		command string
		request *PrependRequest
	}{
		{
			command: "prepend",
			request: nil,
		},
		{
			command: "prepend key flags 2 3\r\nfoo\r\n",
			request: nil,
		},
		{
			command: "prepend key 1 exptime 3\r\nfoo\r\n",
			request: nil,
		},
		{
			command: "prepend key 1 2 size\r\nfoo\r\n",
			request: nil,
		},
		{
			command: "prepend key 1 2 3\nfoo\r\n",
			request: nil,
		},
		{
			command: "prepend key 1 2 3\r\nfoo",
			request: nil,
		},
		{
			command: "prepend key 1 2 3\r\n",
			request: nil,
		},
		{
			command: "prepend k ey 1 2 3\r\nfoo\r\n",
			request: nil,
		},
		{
			command: "prepend 1 2 3\r\nfoo\r\n",
			request: nil,
		},
		{
			command: "prepend key 1 2 3\r\nfoo\r\n",
			request: &PrependRequest{
				Payload: &Storage{
					Key:        "key",
					Flags:      1,
					Expiration: 2 * time.Second,
					Size:       3,
					Data:       []byte("foo"),
					NoReply:    false,
				},
			},
		},
		{
			command: "prepend key 1 2 3 noreply\r\nfoo\r\n",
			request: &PrependRequest{
				Payload: &Storage{
					Key:        "key",
					Flags:      1,
					Expiration: 2 * time.Second,
					Size:       3,
					Data:       []byte("foo"),
					NoReply:    true,
				},
			},
		},
	}

	for _, testCase := range cases {
		parsed, err := NewASCIIParser().Parse([]byte(testCase.command))

		if testCase.request != nil {
			assert.NoError(t, err)
			assert.Equal(t, parsed, testCase.request)
		} else {
			assert.Error(t, err)
		}
	}
}

func TestParseCas(t *testing.T) {
	t.Parallel()

	cases := []struct {
		command string
		request *CasRequest
	}{
		{
			command: "cas",
			request: nil,
		},
		{
			command: "cas key flags 2 3\r\nfoo\r\n",
			request: nil,
		},
		{
			command: "cas key 1 exptime 3\r\nfoo\r\n",
			request: nil,
		},
		{
			command: "cas key 1 2 size\r\nfoo\r\n",
			request: nil,
		},
		{
			command: "cas key 1 2 3\nfoo\r\n",
			request: nil,
		},
		{
			command: "cas key 1 2 3\r\nfoo",
			request: nil,
		},
		{
			command: "cas key 1 2 3\r\n",
			request: nil,
		},
		{
			command: "cas k ey 1 2 3\r\nfoo\r\n",
			request: nil,
		},
		{
			command: "cas 1 2 3\r\nfoo\r\n",
			request: nil,
		},
		{
			command: "cas key 1 2 3\r\nfoo\r\n",
			request: nil,
		},
		{
			command: "cas key 1 2 3 4\r\nfoo\r\n",
			request: &CasRequest{
				Payload: &Storage{
					Key:        "key",
					Flags:      1,
					Expiration: 2 * time.Second,
					Size:       3,
					CasUnique:  4,
					Data:       []byte("foo"),
					NoReply:    false,
				},
			},
		},
		{
			command: "cas key 1 2 3 4 noreply\r\nfoo\r\n",
			request: &CasRequest{
				Payload: &Storage{
					Key:        "key",
					Flags:      1,
					Expiration: 2 * time.Second,
					Size:       3,
					CasUnique:  4,
					Data:       []byte("foo"),
					NoReply:    true,
				},
			},
		},
	}

	for _, testCase := range cases {
		parsed, err := NewASCIIParser().Parse([]byte(testCase.command))

		if testCase.request != nil {
			assert.NoError(t, err)
			assert.Equal(t, parsed, testCase.request)
		} else {
			assert.Error(t, err)
		}
	}
}
