package protocol

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "ERROR message\r\n", (&Error{Err: errors.New("message")}).Error())
}

func TestClientError(t *testing.T) {
	t.Parallel()

	assert.Equal(
		t,
		"CLIENT_ERROR message\r\n",
		(&ClientError{Err: errors.New("message")}).Error(),
	)
}

func TestServerError(t *testing.T) {
	t.Parallel()

	assert.Equal(
		t,
		"SERVER_ERROR message\r\n",
		(&ServerError{Err: errors.New("message")}).Error(),
	)
}

func TestVersionRequest(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "version\r\n", (&VersionRequest{}).String())
}

func TestVersionResponse(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "VERSION version\r\n", (&VersionResponse{Version: "version"}).String())
}

func TestShutdownRequest(t *testing.T) {
	t.Parallel()

	cases := []struct {
		request    *ShutdownRequest
		serialized string
	}{
		{
			request:    &ShutdownRequest{},
			serialized: "shutdown\r\n",
		},
		{
			request:    &ShutdownRequest{Type: "graceful"},
			serialized: "shutdown graceful\r\n",
		},
	}

	for _, testCase := range cases {
		assert.Equal(t, testCase.serialized, testCase.request.String())
	}
}

func TestShutdownResponse(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "OK\r\n", (&ShutdownResponse{}).String())
}

func TestFlushRequest(t *testing.T) {
	t.Parallel()

	cases := []struct {
		request    *FlushRequest
		serialized string
	}{
		{
			request:    &FlushRequest{},
			serialized: "flush_all\r\n",
		},
		{
			request:    &FlushRequest{Delay: 5 * time.Second},
			serialized: "flush_all 5\r\n",
		},
	}

	for _, testCase := range cases {
		assert.Equal(t, testCase.serialized, testCase.request.String())
	}
}

func TestFlushResponse(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "OK\r\n", (&FlushResponse{}).String())
}

func TestQuitRequest(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "quit\r\n", (&QuitRequest{}).String())
}

func TestQuitResponse(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "OK\r\n", (&QuitResponse{}).String())
}

func TestStatsRequest(t *testing.T) {
	t.Parallel()

	cases := []struct {
		request    *StatsRequest
		serialized string
	}{
		{
			request:    &StatsRequest{},
			serialized: "stats\r\n",
		},
		{
			request:    &StatsRequest{Type: "settings"},
			serialized: "stats settings\r\n",
		},
		{
			request:    &StatsRequest{Type: "detail"},
			serialized: "stats detail\r\n",
		},
	}

	for _, testCase := range cases {
		assert.Equal(t, testCase.serialized, testCase.request.String())
	}
}

func TestStatsResponse(t *testing.T) {
	t.Parallel()

	cases := []struct {
		response   *StatsResponse
		serialized string
	}{
		{
			response:   &StatsResponse{},
			serialized: "END\r\n",
		},
		{
			response: &StatsResponse{
				Stats: []StatsItem{
					{Name: "foo", Value: "bar"},
				},
			},
			serialized: "STAT foo bar\r\nEND\r\n",
		},
		{
			response: &StatsResponse{
				Stats: []StatsItem{
					{Name: "foo", Value: "bar"},
					{Name: "baz", Value: "qux"},
				},
			},
			serialized: "STAT foo bar\r\nSTAT baz qux\r\nEND\r\n",
		},
	}

	for _, testCase := range cases {
		assert.Equal(t, testCase.serialized, testCase.response.String())
	}
}

func TestWatchRequest(t *testing.T) {
	t.Parallel()

	cases := []struct {
		request    *WatchRequest
		serialized string
	}{
		{
			request:    &WatchRequest{},
			serialized: "watch\r\n",
		},
		{
			request:    &WatchRequest{Loggers: []string{"fetchers"}},
			serialized: "watch fetchers\r\n",
		},
		{
			request:    &WatchRequest{Loggers: []string{"fetchers", "mutations"}},
			serialized: "watch fetchers mutations\r\n",
		},
	}

	for _, testCase := range cases {
		assert.Equal(t, testCase.serialized, testCase.request.String())
	}
}

func TestWatchResponse(t *testing.T) {
	t.Parallel()

	cases := []struct {
		response   *WatchResponse
		serialized string
	}{
		{
			response:   &WatchResponse{},
			serialized: "",
		},
		{
			response: &WatchResponse{
				Logs: []LogEntry{
					{
						Attributes: []LogAttribute{
							{Key: "gid", Value: "foo"},
						},
					},
				},
			},
			serialized: "gid=foo\r\n",
		},
		{
			response: &WatchResponse{
				Logs: []LogEntry{
					{
						Attributes: []LogAttribute{
							{Key: "gid", Value: "foo"},
							{Key: "key", Value: "bar"},
						},
					},
				},
			},
			serialized: "gid=foo key=bar\r\n",
		},
		{
			response: &WatchResponse{
				Logs: []LogEntry{
					{
						Attributes: []LogAttribute{
							{Key: "gid", Value: "foo"},
							{Key: "key", Value: "bar"},
						},
					},
					{
						Attributes: []LogAttribute{
							{Key: "gid", Value: "baz"},
						},
					},
				},
			},
			serialized: "gid=foo key=bar\r\ngid=baz\r\n",
		},
	}

	for _, testCase := range cases {
		assert.Equal(t, testCase.serialized, testCase.response.String())
	}
}

func TestTouchRequest(t *testing.T) {
	t.Parallel()

	cases := []struct {
		request    *TouchRequest
		serialized string
	}{
		{
			request: &TouchRequest{
				Key:        "key",
				Expiration: 5 * time.Second,
				NoReply:    false,
			},
			serialized: "touch key 5\r\n",
		},
		{
			request: &TouchRequest{
				Key:        "key",
				Expiration: 5 * time.Second,
				NoReply:    true,
			},
			serialized: "touch key 5 noreply\r\n",
		},
	}

	for _, testCase := range cases {
		assert.Equal(t, testCase.serialized, testCase.request.String())
	}
}

func TestTouchResponse(t *testing.T) {
	t.Parallel()

	cases := []struct {
		response   *TouchResponse
		serialized string
	}{
		{
			response:   &TouchResponse{Found: false},
			serialized: "NOT_FOUND\r\n",
		},
		{
			response:   &TouchResponse{Found: true},
			serialized: "TOUCHED\r\n",
		},
	}

	for _, testCase := range cases {
		assert.Equal(t, testCase.serialized, testCase.response.String())
	}
}

func TestDeleteRequest(t *testing.T) {
	t.Parallel()

	cases := []struct {
		request    *DeleteRequest
		serialized string
	}{
		{
			request: &DeleteRequest{
				Key:     "key",
				NoReply: false,
			},
			serialized: "delete key\r\n",
		},
		{
			request: &DeleteRequest{
				Key:     "key",
				NoReply: true,
			},
			serialized: "delete key noreply\r\n",
		},
	}

	for _, testCase := range cases {
		assert.Equal(t, testCase.serialized, testCase.request.String())
	}
}

func TestDeleteResponse(t *testing.T) {
	t.Parallel()

	cases := []struct {
		response   *DeleteResponse
		serialized string
	}{
		{
			response:   &DeleteResponse{Found: false},
			serialized: "NOT_FOUND\r\n",
		},
		{
			response:   &DeleteResponse{Found: true},
			serialized: "DELETED\r\n",
		},
	}

	for _, testCase := range cases {
		assert.Equal(t, testCase.serialized, testCase.response.String())
	}
}

func TestIncrRequest(t *testing.T) {
	t.Parallel()

	cases := []struct {
		request    *IncrRequest
		serialized string
	}{
		{
			request: &IncrRequest{
				Key:     "key",
				Delta:   5,
				NoReply: false,
			},
			serialized: "incr key 5\r\n",
		},
		{
			request: &IncrRequest{
				Key:     "key",
				Delta:   5,
				NoReply: true,
			},
			serialized: "incr key 5 noreply\r\n",
		},
	}

	for _, testCase := range cases {
		assert.Equal(t, testCase.serialized, testCase.request.String())
	}
}

func TestIncrResponse(t *testing.T) {
	t.Parallel()

	cases := []struct {
		response   *IncrResponse
		serialized string
	}{
		{
			response:   &IncrResponse{Found: false},
			serialized: "NOT_FOUND\r\n",
		},
		{
			response:   &IncrResponse{Found: true, Value: 6},
			serialized: "6\r\n",
		},
	}

	for _, testCase := range cases {
		assert.Equal(t, testCase.serialized, testCase.response.String())
	}
}

func TestDecrRequest(t *testing.T) {
	t.Parallel()

	cases := []struct {
		request    *DecrRequest
		serialized string
	}{
		{
			request: &DecrRequest{
				Key:     "key",
				Delta:   5,
				NoReply: false,
			},
			serialized: "decr key 5\r\n",
		},
		{
			request: &DecrRequest{
				Key:     "key",
				Delta:   5,
				NoReply: true,
			},
			serialized: "decr key 5 noreply\r\n",
		},
	}

	for _, testCase := range cases {
		assert.Equal(t, testCase.serialized, testCase.request.String())
	}
}

func TestDecrResponse(t *testing.T) {
	t.Parallel()

	cases := []struct {
		response   *DecrResponse
		serialized string
	}{
		{
			response:   &DecrResponse{Found: false},
			serialized: "NOT_FOUND\r\n",
		},
		{
			response:   &DecrResponse{Found: true, Value: 6},
			serialized: "6\r\n",
		},
	}

	for _, testCase := range cases {
		assert.Equal(t, testCase.serialized, testCase.response.String())
	}
}

func TestGetRequest(t *testing.T) {
	t.Parallel()

	cases := []struct {
		request    *GetRequest
		serialized string
	}{
		{
			request:    &GetRequest{Keys: []string{"foo"}},
			serialized: "get foo\r\n",
		},
		{
			request:    &GetRequest{Keys: []string{"foo", "bar"}},
			serialized: "get foo bar\r\n",
		},
	}

	for _, testCase := range cases {
		assert.Equal(t, testCase.serialized, testCase.request.String())
	}
}

func TestGetResponse(t *testing.T) {
	t.Parallel()

	cases := []struct {
		response   *GetResponse
		serialized string
	}{
		{
			response: &GetResponse{
				Values: []*Retrieval{
					{
						Key:   "foo",
						Flags: 5,
						Size:  3,
						Data:  []byte("bar"),
					},
				},
			},
			serialized: "VALUE foo 5 3\r\nbar\r\nEND\r\n",
		},
		{
			response: &GetResponse{
				Values: []*Retrieval{
					{
						Key:   "foo",
						Flags: 5,
						Size:  3,
						Data:  []byte("bar"),
					},
					{
						Key:   "baz",
						Flags: 0,
						Size:  3,
						Data:  []byte("qux"),
					},
				},
			},
			serialized: "VALUE foo 5 3\r\nbar\r\nVALUE baz 0 3\r\nqux\r\nEND\r\n",
		},
	}

	for _, testCase := range cases {
		assert.Equal(t, testCase.serialized, testCase.response.String())
	}
}

func TestGetsRequest(t *testing.T) {
	t.Parallel()

	cases := []struct {
		request    *GetsRequest
		serialized string
	}{
		{
			request:    &GetsRequest{Keys: []string{"foo"}},
			serialized: "gets foo\r\n",
		},
		{
			request:    &GetsRequest{Keys: []string{"foo", "bar"}},
			serialized: "gets foo bar\r\n",
		},
	}

	for _, testCase := range cases {
		assert.Equal(t, testCase.serialized, testCase.request.String())
	}
}

func TestGetsResponse(t *testing.T) {
	t.Parallel()

	cases := []struct {
		response   *GetsResponse
		serialized string
	}{
		{
			response: &GetsResponse{
				Values: []*Retrieval{
					{
						Key:       "foo",
						Flags:     5,
						Size:      3,
						CasUnique: 4,
						Data:      []byte("bar"),
					},
				},
			},
			serialized: "VALUE foo 5 3 4\r\nbar\r\nEND\r\n",
		},
		{
			response: &GetsResponse{
				Values: []*Retrieval{
					{
						Key:       "foo",
						Flags:     5,
						Size:      3,
						CasUnique: 4,
						Data:      []byte("bar"),
					},
					{
						Key:       "baz",
						Flags:     0,
						Size:      3,
						CasUnique: 8,
						Data:      []byte("qux"),
					},
				},
			},
			serialized: "VALUE foo 5 3 4\r\nbar\r\nVALUE baz 0 3 8\r\nqux\r\nEND\r\n",
		},
	}

	for _, testCase := range cases {
		assert.Equal(t, testCase.serialized, testCase.response.String())
	}
}

func TestGatRequest(t *testing.T) {
	t.Parallel()

	cases := []struct {
		request    *GatRequest
		serialized string
	}{
		{
			request: &GatRequest{
				Expiration: 5 * time.Second,
				Keys:       []string{"foo"},
			},
			serialized: "gat 5 foo\r\n",
		},
		{
			request: &GatRequest{
				Expiration: 5 * time.Second,
				Keys:       []string{"foo", "bar"},
			},
			serialized: "gat 5 foo bar\r\n",
		},
	}

	for _, testCase := range cases {
		assert.Equal(t, testCase.serialized, testCase.request.String())
	}
}

func TestGatResponse(t *testing.T) {
	t.Parallel()

	cases := []struct {
		response   *GatResponse
		serialized string
	}{
		{
			response: &GatResponse{
				Values: []*Retrieval{
					{
						Key:   "foo",
						Flags: 5,
						Size:  3,
						Data:  []byte("bar"),
					},
				},
			},
			serialized: "VALUE foo 5 3\r\nbar\r\nEND\r\n",
		},
		{
			response: &GatResponse{
				Values: []*Retrieval{
					{
						Key:   "foo",
						Flags: 5,
						Size:  3,
						Data:  []byte("bar"),
					},
					{
						Key:   "baz",
						Flags: 0,
						Size:  3,
						Data:  []byte("qux"),
					},
				},
			},
			serialized: "VALUE foo 5 3\r\nbar\r\nVALUE baz 0 3\r\nqux\r\nEND\r\n",
		},
	}

	for _, testCase := range cases {
		assert.Equal(t, testCase.serialized, testCase.response.String())
	}
}

func TestGatsRequest(t *testing.T) {
	t.Parallel()

	cases := []struct {
		request    *GatsRequest
		serialized string
	}{
		{
			request: &GatsRequest{
				Expiration: 5 * time.Second,
				Keys:       []string{"foo"},
			},
			serialized: "gats 5 foo\r\n",
		},
		{
			request: &GatsRequest{
				Expiration: 5 * time.Second,
				Keys:       []string{"foo", "bar"},
			},
			serialized: "gats 5 foo bar\r\n",
		},
	}

	for _, testCase := range cases {
		assert.Equal(t, testCase.serialized, testCase.request.String())
	}
}

func TestGatsResponse(t *testing.T) {
	t.Parallel()

	cases := []struct {
		response   *GatsResponse
		serialized string
	}{
		{
			response: &GatsResponse{
				Values: []*Retrieval{
					{
						Key:       "foo",
						Flags:     5,
						Size:      3,
						CasUnique: 4,
						Data:      []byte("bar"),
					},
				},
			},
			serialized: "VALUE foo 5 3 4\r\nbar\r\nEND\r\n",
		},
		{
			response: &GatsResponse{
				Values: []*Retrieval{
					{
						Key:       "foo",
						Flags:     5,
						Size:      3,
						CasUnique: 4,
						Data:      []byte("bar"),
					},
					{
						Key:       "baz",
						Flags:     0,
						Size:      3,
						CasUnique: 8,
						Data:      []byte("qux"),
					},
				},
			},
			serialized: "VALUE foo 5 3 4\r\nbar\r\nVALUE baz 0 3 8\r\nqux\r\nEND\r\n",
		},
	}

	for _, testCase := range cases {
		assert.Equal(t, testCase.serialized, testCase.response.String())
	}
}

func TestSetRequest(t *testing.T) {
	t.Parallel()

	cases := []struct {
		request    *SetRequest
		serialized string
	}{
		{
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
			serialized: "set key 1 2 3\r\nfoo\r\n",
		},
		{
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
			serialized: "set key 1 2 3 noreply\r\nfoo\r\n",
		},
	}

	for _, testCase := range cases {
		assert.Equal(t, testCase.serialized, testCase.request.String())
	}
}

func TestSetResponse(t *testing.T) {
	t.Parallel()

	cases := []struct {
		response   *SetResponse
		serialized string
	}{
		{
			response:   &SetResponse{Status: Stored},
			serialized: "STORED\r\n",
		},
		{
			response:   &SetResponse{Status: NotStored},
			serialized: "NOT_STORED\r\n",
		},
		{
			response:   &SetResponse{Status: Exists},
			serialized: "EXISTS\r\n",
		},
		{
			response:   &SetResponse{Status: NotFound},
			serialized: "NOT_FOUND\r\n",
		},
	}

	for _, testCase := range cases {
		assert.Equal(t, testCase.serialized, testCase.response.String())
	}
}

func TestAddRequest(t *testing.T) {
	t.Parallel()

	cases := []struct {
		request    *AddRequest
		serialized string
	}{
		{
			request: &AddRequest{
				Payload: &Storage{
					Key:        "key",
					Flags:      1,
					Expiration: 2 * time.Second,
					Size:       3,
					Data:       []byte("foo"),
					NoReply:    false,
				},
			},
			serialized: "add key 1 2 3\r\nfoo\r\n",
		},
		{
			request: &AddRequest{
				Payload: &Storage{
					Key:        "key",
					Flags:      1,
					Expiration: 2 * time.Second,
					Size:       3,
					Data:       []byte("foo"),
					NoReply:    true,
				},
			},
			serialized: "add key 1 2 3 noreply\r\nfoo\r\n",
		},
	}

	for _, testCase := range cases {
		assert.Equal(t, testCase.serialized, testCase.request.String())
	}
}

func TestAddResponse(t *testing.T) {
	t.Parallel()

	cases := []struct {
		response   *AddResponse
		serialized string
	}{
		{
			response:   &AddResponse{Status: Stored},
			serialized: "STORED\r\n",
		},
		{
			response:   &AddResponse{Status: NotStored},
			serialized: "NOT_STORED\r\n",
		},
		{
			response:   &AddResponse{Status: Exists},
			serialized: "EXISTS\r\n",
		},
		{
			response:   &AddResponse{Status: NotFound},
			serialized: "NOT_FOUND\r\n",
		},
	}

	for _, testCase := range cases {
		assert.Equal(t, testCase.serialized, testCase.response.String())
	}
}

func TestReplaceRequest(t *testing.T) {
	t.Parallel()

	cases := []struct {
		request    *ReplaceRequest
		serialized string
	}{
		{
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
			serialized: "replace key 1 2 3\r\nfoo\r\n",
		},
		{
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
			serialized: "replace key 1 2 3 noreply\r\nfoo\r\n",
		},
	}

	for _, testCase := range cases {
		assert.Equal(t, testCase.serialized, testCase.request.String())
	}
}

func TestReplaceResponse(t *testing.T) {
	t.Parallel()

	cases := []struct {
		response   *ReplaceResponse
		serialized string
	}{
		{
			response:   &ReplaceResponse{Status: Stored},
			serialized: "STORED\r\n",
		},
		{
			response:   &ReplaceResponse{Status: NotStored},
			serialized: "NOT_STORED\r\n",
		},
		{
			response:   &ReplaceResponse{Status: Exists},
			serialized: "EXISTS\r\n",
		},
		{
			response:   &ReplaceResponse{Status: NotFound},
			serialized: "NOT_FOUND\r\n",
		},
	}

	for _, testCase := range cases {
		assert.Equal(t, testCase.serialized, testCase.response.String())
	}
}

func TestAppendRequest(t *testing.T) {
	t.Parallel()

	cases := []struct {
		request    *AppendRequest
		serialized string
	}{
		{
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
			serialized: "append key 1 2 3\r\nfoo\r\n",
		},
		{
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
			serialized: "append key 1 2 3 noreply\r\nfoo\r\n",
		},
	}

	for _, testCase := range cases {
		assert.Equal(t, testCase.serialized, testCase.request.String())
	}
}

func TestAppendResponse(t *testing.T) {
	t.Parallel()

	cases := []struct {
		response   *AppendResponse
		serialized string
	}{
		{
			response:   &AppendResponse{Status: Stored},
			serialized: "STORED\r\n",
		},
		{
			response:   &AppendResponse{Status: NotStored},
			serialized: "NOT_STORED\r\n",
		},
		{
			response:   &AppendResponse{Status: Exists},
			serialized: "EXISTS\r\n",
		},
		{
			response:   &AppendResponse{Status: NotFound},
			serialized: "NOT_FOUND\r\n",
		},
	}

	for _, testCase := range cases {
		assert.Equal(t, testCase.serialized, testCase.response.String())
	}
}

func TestPrependRequest(t *testing.T) {
	t.Parallel()

	cases := []struct {
		request    *PrependRequest
		serialized string
	}{
		{
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
			serialized: "prepend key 1 2 3\r\nfoo\r\n",
		},
		{
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
			serialized: "prepend key 1 2 3 noreply\r\nfoo\r\n",
		},
	}

	for _, testCase := range cases {
		assert.Equal(t, testCase.serialized, testCase.request.String())
	}
}

func TestPrependResponse(t *testing.T) {
	t.Parallel()

	cases := []struct {
		response   *PrependResponse
		serialized string
	}{
		{
			response:   &PrependResponse{Status: Stored},
			serialized: "STORED\r\n",
		},
		{
			response:   &PrependResponse{Status: NotStored},
			serialized: "NOT_STORED\r\n",
		},
		{
			response:   &PrependResponse{Status: Exists},
			serialized: "EXISTS\r\n",
		},
		{
			response:   &PrependResponse{Status: NotFound},
			serialized: "NOT_FOUND\r\n",
		},
	}

	for _, testCase := range cases {
		assert.Equal(t, testCase.serialized, testCase.response.String())
	}
}

func TestCasRequest(t *testing.T) {
	t.Parallel()

	cases := []struct {
		request    *CasRequest
		serialized string
	}{
		{
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
			serialized: "cas key 1 2 3 4\r\nfoo\r\n",
		},
		{
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
			serialized: "cas key 1 2 3 4 noreply\r\nfoo\r\n",
		},
	}

	for _, testCase := range cases {
		assert.Equal(t, testCase.serialized, testCase.request.String())
	}
}

func TestCasResponse(t *testing.T) {
	t.Parallel()

	cases := []struct {
		response   *CasResponse
		serialized string
	}{
		{
			response:   &CasResponse{Status: Stored},
			serialized: "STORED\r\n",
		},
		{
			response:   &CasResponse{Status: NotStored},
			serialized: "NOT_STORED\r\n",
		},
		{
			response:   &CasResponse{Status: Exists},
			serialized: "EXISTS\r\n",
		},
		{
			response:   &CasResponse{Status: NotFound},
			serialized: "NOT_FOUND\r\n",
		},
	}

	for _, testCase := range cases {
		assert.Equal(t, testCase.serialized, testCase.response.String())
	}
}
