package protocol

import (
    "errors"
    "strconv"
    "time"
)

%%{
    machine memcache_ascii;

    action mark {
        m.mark()
    }

    # Common

    sp = ' '+;
    crlf = '\r\n';
    key = (any - ' '+){1,250} >mark %{ keys = append(keys, m.text()) };
    key_list = (key sp)* key;
    subcommand = /[a-z]/+ >mark %{ subcommands = append(subcommands, m.text()) };
    subcommand_list = (subcommand sp)* subcommand;
    noreply = 'noreply' %{ noreply = true };
    exptime = /[0-9]/+ >mark %{
        if parsed, err := strconv.ParseUint(m.text(), 10, 64); err == nil {
            exptime = time.Duration(parsed) * time.Second
        }
    };

    # Get commands

    get = 'get' sp key_list crlf %{ cmd = Get };
    gets = 'gets' sp key_list crlf %{ cmd = Gets };
    gat = 'gat' sp exptime sp key_list crlf %{ cmd = Gat };
    gats = 'gats' sp exptime sp key_list crlf %{ cmd = Gats };

    # Storage commands

    flags = /[0-9]/+ >mark %{
        if parsed, err := strconv.ParseUint(m.text(), 10, 16); err == nil {
            flags = uint16(parsed)
        }
    };
    size = /[0-9]/+ >mark %{
        if parsed, err := strconv.Atoi(m.text()); err == nil {
            size = parsed
        }
    };
    cas_unique = /[0-9]/+ >mark %{
        if parsed, err := strconv.ParseUint(m.text(), 10, 64); err == nil {
            casUnique = parsed
        }
    };
    storage_data = (any - crlf)+ >mark %{ data = m.bytes() };
    storage_header = key sp flags sp exptime sp size;

    set = 'set' sp storage_header (sp noreply)? crlf storage_data crlf %{ cmd = Set };
    add = 'add' sp storage_header (sp noreply)? crlf storage_data crlf %{ cmd = Add };
    replace = 'replace' sp storage_header (sp noreply)? crlf storage_data crlf %{ cmd = Replace };
    append = 'append' sp storage_header (sp noreply)? crlf storage_data crlf %{ cmd = Append };
    prepend = 'prepend' sp storage_header (sp noreply)? crlf storage_data crlf %{ cmd = Prepend };
    cas = 'cas' sp storage_header sp cas_unique (sp noreply)? crlf storage_data crlf %{ cmd = Cas };

    # Arithmetic commands

    delta = /[0-9]/+ >mark %{
        if parsed, err := strconv.ParseUint(m.text(), 10, 64); err == nil {
            delta = parsed
        }
    };

    incr = 'incr' sp key sp delta (sp noreply)? crlf %{ cmd = Incr };
    decr = 'decr' sp key sp delta (sp noreply)? crlf %{ cmd = Decr };
    delete = 'delete' sp key (sp noreply)? crlf %{ cmd = Delete };
    touch = 'touch' sp key sp exptime (sp noreply)? crlf %{ cmd = Touch };

    # Observability commands

    stats = 'stats' (sp subcommand)? crlf %{ cmd = Stats };
    watch = 'watch' (sp subcommand_list)? crlf %{ cmd = Watch };

    # Admin and metadata commands

    version = 'version' crlf %{ cmd = Version };
    shutdown = 'shutdown' (sp subcommand)? crlf %{ cmd = Shutdown };
    flush = 'flush_all' (sp exptime)? crlf %{ cmd = Flush };

    # Control commands

    quit = 'quit' crlf %{ cmd = Quit };

    # Main entry point

    command :=
        set |
        add |
        replace |
        append |
        prepend |
        cas |
        get |
        gets |
        gat |
        gats |
        incr |
        decr |
        delete |
        touch |
        stats |
        watch |
        version |
        shutdown |
        flush |
        quit;
}%%

%% write data noerror nofinal noentry noprefix;

var (
    // ErrInvalidParse is returned when the input data cannot be parsed according to the defined
    // ASCII protocol state machine.
    ErrInvalidParse = errors.New("protocol: command does not parse as a valid request")
)

type machine struct {
    data []byte
    cs int
    p int
    pb int
    pe int
    eof int
}

// NewASCIIParser creates a new state machine representing the memcache ASCII protocol.
func NewASCIIParser() Parser {
    m := &machine{}

    %% access m.;
    %% variable p m.p;
    %% variable pe m.pe;
    %% variable eof m.eof;
    %% variable data m.data;

    return m
}

func (m *machine) Parse(command []byte) (Request, error) {
    var cmd Command
    var keys []string
    var exptime time.Duration
    var flags uint16
    var size int
    var casUnique uint64
    var delta uint64
    var subcommands []string
    var noreply bool
    var data []byte

    m.data = command
    m.p = 0
    m.pb = 0
    m.pe = len(command)
    m.eof = len(command)

    %% write init;
    %% write exec;

    switch cmd {
    case Version:
        return &VersionRequest{}, nil

    case Shutdown:
        if len(subcommands) == 0 {
            return &ShutdownRequest{}, nil
        }

        return &ShutdownRequest{Type: subcommands[0]}, nil

    case Flush:
        return &FlushRequest{Delay: exptime}, nil

    case Quit:
        return &QuitRequest{}, nil

    case Stats:
        if len(subcommands) == 0 {
            return &StatsRequest{}, nil
        }

        return &StatsRequest{Type: subcommands[0]}, nil

    case Watch:
        return &WatchRequest{Loggers: subcommands}, nil

    case Touch:
        return &TouchRequest{
            Key: keys[0],
            Expiration: exptime,
            NoReply: noreply,
        }, nil

    case Delete:
        return &DeleteRequest{
            Key: keys[0],
            NoReply: noreply,
        }, nil

    case Incr:
        return &IncrRequest{
            Key: keys[0],
            Delta: delta,
            NoReply: noreply,
        }, nil

    case Decr:
        return &DecrRequest{
            Key: keys[0],
            Delta: delta,
            NoReply: noreply,
        }, nil

    case Get:
        return &GetRequest{Keys: keys}, nil

    case Gets:
        return &GetsRequest{Keys: keys}, nil

    case Gat:
        return &GatRequest{
            Expiration: exptime,
            Keys: keys,
        }, nil

    case Gats:
        return &GatsRequest{
            Expiration: exptime,
            Keys: keys,
        }, nil

    case Set, Add, Replace, Append, Prepend, Cas:
        payload := &Storage{
            Key: keys[0],
            Flags: flags,
            Expiration: exptime,
            Size: size,
            CasUnique: casUnique,
            Data: data,
            NoReply: noreply,
        }

        switch cmd {
            case Set:
                return &SetRequest{Payload: payload}, nil
            case Add:
                return &AddRequest{Payload: payload}, nil
            case Replace:
                return &ReplaceRequest{Payload: payload}, nil
            case Append:
                return &AppendRequest{Payload: payload}, nil
            case Prepend:
                return &PrependRequest{Payload: payload}, nil
            case Cas:
                return &CasRequest{Payload: payload}, nil
        }
    }

    return nil, ErrInvalidParse
}

func (m *machine) mark() {
    m.pb = m.p
}

func (m *machine) bytes() []byte {
    return m.data[m.pb:m.p]
}

func (m *machine) text() string {
    return string(m.bytes())
}
