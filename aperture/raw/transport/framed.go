package transport

// Framed is an abstraction over a Transport that simulates message framing appending a constant
// tailer sequence to every outgoing message.
type Framed struct {
	tailer []byte
	Transport
}

// NewFramed creates a new framed transport over an existing Transport backend.
func NewFramed(backendFactory Factory, tailer []byte) (Transport, error) {
	tport, err := backendFactory()
	if err != nil {
		return nil, err
	}

	return &Framed{
		Transport: tport,
		tailer:    tailer,
	}, nil
}

// Send appends the constant binary tailer sequence to the message and dispatches the write request
// to the underlying transport. Note that this assumes that every Send invocation corresponds to a
// fully formed "message" in the protocol leveraging this transport.
func (t *Framed) Send(message []byte) (int, error) {
	return t.Transport.Send(append(message, t.tailer...))
}
