//go:generate ragel -G2 -Z memcache.rl -o memcache.go
//go:generate gofmt -s -w memcache.go

package protocol

import (
	"fmt"
	"strings"
)

func (e *Error) String() string {
	return fmt.Sprintf("ERROR %v\r\n", e.Err)
}

func (e *Error) Error() string {
	return e.String()
}

func (e *Error) Unwrap() error {
	return e.Err
}

func (e *ClientError) String() string {
	return fmt.Sprintf("CLIENT_ERROR %v\r\n", e.Err)
}

func (e *ClientError) Error() string {
	return e.String()
}

func (e *ClientError) Unwrap() error {
	return e.Err
}

func (e *ServerError) String() string {
	return fmt.Sprintf("SERVER_ERROR %v\r\n", e.Err)
}

func (e *ServerError) Error() string {
	return e.String()
}

func (e *ServerError) Unwrap() error {
	return e.Err
}

// IsNoReply is always false.
func (v *VersionRequest) IsNoReply() bool {
	return false
}

func (v *VersionRequest) String() string {
	return "version\r\n"
}

func (v *VersionResponse) String() string {
	return fmt.Sprintf("VERSION %s\r\n", v.Version)
}

// IsNoReply is always false.
func (s *ShutdownRequest) IsNoReply() bool {
	return false
}

func (s *ShutdownRequest) String() string {
	if s.Type == "" {
		return "shutdown\r\n"
	}

	return fmt.Sprintf("shutdown %s\r\n", s.Type)
}

func (s *ShutdownResponse) String() string {
	return "OK\r\n"
}

// IsNoReply is always false.
func (f *FlushRequest) IsNoReply() bool {
	return false
}

func (f *FlushRequest) String() string {
	if f.Delay == 0 {
		return "flush_all\r\n"
	}

	return fmt.Sprintf("flush_all %s\r\n", fmt.Sprintf("%.f", f.Delay.Seconds()))
}

func (f *FlushResponse) String() string {
	return "OK\r\n"
}

// IsNoReply is always false.
func (s *StatsRequest) IsNoReply() bool {
	return false
}

func (s *StatsRequest) String() string {
	if s.Type == "" {
		return "stats\r\n"
	}

	return fmt.Sprintf("stats %s\r\n", s.Type)
}

func (s *StatsResponse) String() string {
	sb := strings.Builder{}

	for _, stat := range s.Stats {
		sb.WriteString(stat.String())
	}

	sb.WriteString("END\r\n")

	return sb.String()
}

func (s *StatsItem) String() string {
	return fmt.Sprintf("STAT %s %s\r\n", s.Name, s.Value)
}

// IsNoReply is always false.
func (w *WatchRequest) IsNoReply() bool {
	return false
}

func (w *WatchRequest) String() string {
	if len(w.Loggers) == 0 {
		return "watch\r\n"
	}

	return fmt.Sprintf("watch %s\r\n", strings.Join(w.Loggers, " "))
}

func (w *WatchResponse) String() string {
	sb := strings.Builder{}

	for _, entry := range w.Logs {
		sb.WriteString(entry.String())
	}

	return sb.String()
}

func (l *LogAttribute) String() string {
	return fmt.Sprintf("%s=%s", l.Key, l.Value)
}

func (l *LogEntry) String() string {
	sb := strings.Builder{}

	for idx, attribute := range l.Attributes {
		sb.WriteString(attribute.String())

		if idx < len(l.Attributes)-1 {
			sb.WriteString(" ")
		}
	}

	sb.WriteString("\r\n")

	return sb.String()
}

// IsNoReply returns the NoReply attribute of the request.
func (t *TouchRequest) IsNoReply() bool {
	return t.NoReply
}

func (t *TouchRequest) String() string {
	exptime := fmt.Sprintf("%.f", t.Expiration.Seconds())

	if t.NoReply {
		return fmt.Sprintf("touch %s %s noreply\r\n", t.Key, exptime)
	}

	return fmt.Sprintf("touch %s %s\r\n", t.Key, exptime)
}

func (t *TouchResponse) String() string {
	if t.Found {
		return "TOUCHED\r\n"
	}

	return "NOT_FOUND\r\n"
}

// IsNoReply returns the NoReply attribute of the request.
func (d *DeleteRequest) IsNoReply() bool {
	return d.NoReply
}

func (d *DeleteRequest) String() string {
	if d.NoReply {
		return fmt.Sprintf("delete %s noreply\r\n", d.Key)
	}

	return fmt.Sprintf("delete %s\r\n", d.Key)
}

func (d *DeleteResponse) String() string {
	if d.Found {
		return "DELETED\r\n"
	}

	return "NOT_FOUND\r\n"
}

// IsNoReply returns the NoReply attribute of the request.
func (i *IncrRequest) IsNoReply() bool {
	return i.NoReply
}

func (i *IncrRequest) String() string {
	if i.NoReply {
		return fmt.Sprintf("incr %s %d noreply\r\n", i.Key, i.Delta)
	}

	return fmt.Sprintf("incr %s %d\r\n", i.Key, i.Delta)
}

func (i *IncrResponse) String() string {
	if i.Found {
		return fmt.Sprintf("%d\r\n", i.Value)
	}

	return "NOT_FOUND\r\n"
}

// IsNoReply returns the NoReply attribute of the request.
func (d *DecrRequest) IsNoReply() bool {
	return d.NoReply
}

func (d *DecrRequest) String() string {
	if d.NoReply {
		return fmt.Sprintf("decr %s %d noreply\r\n", d.Key, d.Delta)
	}

	return fmt.Sprintf("decr %s %d\r\n", d.Key, d.Delta)
}

func (d *DecrResponse) String() string {
	if d.Found {
		return fmt.Sprintf("%d\r\n", d.Value)
	}

	return "NOT_FOUND\r\n"
}

// IsNoReply is always false.
func (g *GetRequest) IsNoReply() bool {
	return false
}

func (g *GetRequest) String() string {
	return fmt.Sprintf("get %s\r\n", strings.Join(g.Keys, " "))
}

func (g *GetResponse) String() string {
	sb := strings.Builder{}

	for _, value := range g.Values {
		sb.WriteString(value.String())
	}

	sb.WriteString("END\r\n")

	return sb.String()
}

// IsNoReply is always false.
func (g *GetsRequest) IsNoReply() bool {
	return false
}

func (g *GetsRequest) String() string {
	return fmt.Sprintf("gets %s\r\n", strings.Join(g.Keys, " "))
}

func (g *GetsResponse) String() string {
	return (&GetResponse{Values: g.Values}).String()
}

// IsNoReply is always false.
func (g *GatRequest) IsNoReply() bool {
	return false
}

func (g *GatRequest) String() string {
	return fmt.Sprintf(
		"gat %s %s\r\n",
		fmt.Sprintf("%.f", g.Expiration.Seconds()),
		strings.Join(g.Keys, " "),
	)
}

func (g *GatResponse) String() string {
	return (&GetResponse{Values: g.Values}).String()
}

// IsNoReply is always false.
func (g *GatsRequest) IsNoReply() bool {
	return false
}

func (g *GatsRequest) String() string {
	return fmt.Sprintf(
		"gats %s %s\r\n",
		fmt.Sprintf("%.f", g.Expiration.Seconds()),
		strings.Join(g.Keys, " "),
	)
}

func (g *GatsResponse) String() string {
	return (&GetResponse{Values: g.Values}).String()
}

func (r *Retrieval) String() string {
	if r.CasUnique == 0 {
		return fmt.Sprintf(
			"VALUE %s %d %d\r\n%s\r\n",
			r.Key,
			r.Flags,
			r.Size,
			string(r.Data),
		)
	}

	return fmt.Sprintf(
		"VALUE %s %d %d %d\r\n%s\r\n",
		r.Key,
		r.Flags,
		r.Size,
		r.CasUnique,
		string(r.Data),
	)
}

// IsNoReply returns the NoReply attribute of the request.
func (s *SetRequest) IsNoReply() bool {
	return s.Payload.NoReply
}

func (s *SetRequest) String() string {
	return fmt.Sprintf("set %s", s.Payload.String())
}

func (s *SetResponse) String() string {
	switch s.Status {
	case Stored:
		return "STORED\r\n"
	case NotStored:
		return "NOT_STORED\r\n"
	case Exists:
		return "EXISTS\r\n"
	case NotFound:
		return "NOT_FOUND\r\n"
	default:
		return "OK\r\n"
	}
}

// IsNoReply returns the NoReply attribute of the request.
func (a *AddRequest) IsNoReply() bool {
	return a.Payload.NoReply
}

func (a *AddRequest) String() string {
	return fmt.Sprintf("add %s", a.Payload.String())
}

func (a *AddResponse) String() string {
	return (&SetResponse{Status: a.Status}).String()
}

// IsNoReply returns the NoReply attribute of the request.
func (r *ReplaceRequest) IsNoReply() bool {
	return r.Payload.NoReply
}

func (r *ReplaceRequest) String() string {
	return fmt.Sprintf("replace %s", r.Payload.String())
}

func (r *ReplaceResponse) String() string {
	return (&SetResponse{Status: r.Status}).String()
}

// IsNoReply returns the NoReply attribute of the request.
func (a *AppendRequest) IsNoReply() bool {
	return a.Payload.NoReply
}

func (a *AppendRequest) String() string {
	return fmt.Sprintf("append %s", a.Payload.String())
}

func (a *AppendResponse) String() string {
	return (&SetResponse{Status: a.Status}).String()
}

// IsNoReply returns the NoReply attribute of the request.
func (p *PrependRequest) IsNoReply() bool {
	return p.Payload.NoReply
}

func (p *PrependRequest) String() string {
	return fmt.Sprintf("prepend %s", p.Payload.String())
}

func (p *PrependResponse) String() string {
	return (&SetResponse{Status: p.Status}).String()
}

// IsNoReply returns the NoReply attribute of the request.
func (c *CasRequest) IsNoReply() bool {
	return c.Payload.NoReply
}

func (c *CasRequest) String() string {
	return fmt.Sprintf("cas %s", c.Payload.String())
}

func (c *CasResponse) String() string {
	return (&SetResponse{Status: c.Status}).String()
}

func (s *Storage) String() string {
	exptime := fmt.Sprintf("%.f", s.Expiration.Seconds())

	switch {
	case s.CasUnique != 0 && s.NoReply:
		return fmt.Sprintf(
			"%s %d %s %d %d noreply\r\n%s\r\n",
			s.Key,
			s.Flags,
			exptime,
			s.Size,
			s.CasUnique,
			string(s.Data),
		)
	case s.CasUnique != 0 && !s.NoReply:
		return fmt.Sprintf(
			"%s %d %s %d %d\r\n%s\r\n",
			s.Key,
			s.Flags,
			exptime,
			s.Size,
			s.CasUnique,
			string(s.Data),
		)
	case s.CasUnique == 0 && s.NoReply:
		return fmt.Sprintf(
			"%s %d %s %d noreply\r\n%s\r\n",
			s.Key,
			s.Flags,
			exptime,
			s.Size,
			string(s.Data),
		)
	default:
		return fmt.Sprintf(
			"%s %d %s %d\r\n%s\r\n",
			s.Key,
			s.Flags,
			exptime,
			s.Size,
			string(s.Data),
		)
	}
}
