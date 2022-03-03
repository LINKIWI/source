package proxy

import (
	"crypto/tls"
	"encoding/json"
	"io"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"lib.kevinlin.info/aperture"
	"lib.kevinlin.info/aperture/lib"
)

// Proxy is a server that acts as a TCP proxy between a listener and a target.
type Proxy struct {
	// Name is a human-readable identifier for this proxy instance.
	Name string
	// DialTarget is a function used to dial a connection to the target server. It returns a
	// net.Conn to the destination endpoint, as well as a Pipe implementation used for copying
	// bytes between the source and destination.
	DialTarget func(src net.Conn) (dst net.Conn, pipe Pipe, err error)
	// MaxActiveConnections is the total number of permitted active, concurrent client
	// connections. Accepted connections in excess of this limit are rejected and closed.
	MaxActiveConnections int
	// EnableProxyHeader enables use of the HAProxy PROXY protocol header.
	EnableProxyHeader bool
	// Metrics is an aperture.Statsd for reporting of internal metrics.
	Metrics aperture.Statsd
	// DebugLog is a standard library logger for debug messages.
	DebugLog *log.Logger
	// ErrorLog is a standard library logger for error messages.
	ErrorLog *log.Logger
	// ConnectionLog is an optional io.Writer for writing JSON-formatted connection access logs.
	ConnectionLog io.Writer

	closing bool
}

// Serve the TCP proxy on a listener.
func (p *Proxy) Serve(ln net.Listener) error {
	var connections sync.WaitGroup

	srcSentTotal := new(int64)
	dstRecvTotal := new(int64)
	serveTotal := new(int64)
	activeConns := new(int64)
	tags := map[string]interface{}{
		"proxy":    p.Name,
		"listener": ln.Addr().String(),
	}

	for {
		src, err := ln.Accept()
		if err != nil {
			if p.closing {
				p.Metrics.Incr("proxy.serve.accept.shutdown", tags)
				p.DebugLog.Printf(
					"proxy: closing listener: name=%s addr=%v",
					p.Name,
					ln.Addr(),
				)
				break
			}

			p.Metrics.Incr("proxy.serve.accept.error", tags)
			p.ErrorLog.Printf(
				"proxy: error accepting connection from listener: name=%s err=%v",
				p.Name,
				err,
			)
			continue
		}

		if p.MaxActiveConnections > 0 && atomic.LoadInt64(activeConns) >= int64(p.MaxActiveConnections) {
			src.Close()
			p.Metrics.Incr("proxy.serve.accept.reject", tags)
			p.ErrorLog.Printf(
				"proxy: rejecting new client connection due to circuit breaker: "+
					"name=%s remote=%v max_active_conns=%d",
				p.Name,
				src.RemoteAddr(),
				p.MaxActiveConnections,
			)
			continue
		}

		p.Metrics.Incr("proxy.serve.accept.success", tags)
		p.DebugLog.Printf(
			"proxy: accepted new connection: name=%s local=%v remote=%v",
			p.Name,
			src.LocalAddr(),
			src.RemoteAddr(),
		)

		go func() {
			connections.Add(1)
			defer connections.Done()

			timestamp := time.Now()
			stopwatch := lib.NewStopwatch()

			atomic.AddInt64(activeConns, 1)
			p.Metrics.Gauge(
				"proxy.serve.active_connections",
				float64(atomic.LoadInt64(activeConns)),
				tags,
			)

			defer func() {
				atomic.AddInt64(activeConns, -1)
				p.Metrics.Gauge(
					"proxy.serve.active_connections",
					float64(atomic.LoadInt64(activeConns)),
					tags,
				)
			}()

			dst, pipe, err := p.DialTarget(src)
			if err != nil {
				src.Close()
				p.Metrics.Incr("proxy.serve.dial.error", tags)
				p.ErrorLog.Printf(
					"proxy: error dialing target: name=%s err=%v",
					p.Name,
					err,
				)
				return
			}

			defer src.Close()
			defer dst.Close()

			p.Metrics.Incr("proxy.serve.dial.success", tags)
			p.Metrics.Timing("proxy.serve.dial.latency", stopwatch.Elapsed(), tags)
			stopwatch.Reset()

			if p.EnableProxyHeader {
				header := &protocolHeader{src: src, dst: dst}

				if _, err := dst.Write([]byte(header.String())); err != nil {
					p.Metrics.Incr("proxy.protocol.error", tags)
					p.ErrorLog.Printf(
						"proxy: error writing proxy protocol header: "+
							"name=%s err=%v",
						p.Name,
						err,
					)
					return
				}
			}

			if pipe == nil {
				pipe = DefaultDirectPipe
			}

			srcSent, dstRecv, err := pipe.Do(dst, src)
			if err != nil {
				p.Metrics.Incr("proxy.pipe.error", tags)
				p.ErrorLog.Printf(
					"proxy: error piping data to target: name=%s err=%v",
					p.Name,
					err,
				)
				return
			}

			atomic.AddInt64(srcSentTotal, srcSent)
			atomic.AddInt64(dstRecvTotal, dstRecv)
			atomic.AddInt64(serveTotal, 1)

			p.Metrics.Incr("proxy.pipe.success", tags)
			p.Metrics.Size("proxy.pipe.src_sent.size", srcSent, tags)
			p.Metrics.Size("proxy.pipe.dst_recv.size", dstRecv, tags)
			p.Metrics.Gauge("proxy.pipe.src_sent.total", float64(atomic.LoadInt64(srcSentTotal)), tags)
			p.Metrics.Gauge("proxy.pipe.dst_recv.total", float64(atomic.LoadInt64(dstRecvTotal)), tags)
			p.Metrics.Timing("proxy.pipe.duration", stopwatch.Elapsed(), tags)
			p.Metrics.Gauge("proxy.serve.total", float64(atomic.LoadInt64(serveTotal)), tags)

			if p.ConnectionLog != nil {
				srcNet, srcIP, srcPort, _ := parseAddress(src.RemoteAddr())
				dstNet, dstIP, dstPort, _ := parseAddress(dst.RemoteAddr())

				srcTransport := "plaintext"
				if _, ok := src.(*tls.Conn); ok {
					srcTransport = "tls"
				}

				dstTransport := "plaintext"
				if _, ok := dst.(*tls.Conn); ok {
					dstTransport = "tls"
				}

				entry := &connLog{
					Name:                    p.Name,
					Timestamp:               timestamp,
					Duration:                float64(stopwatch.Elapsed().Microseconds()) / 1000.0,
					SourceIP:                srcIP,
					SourcePort:              srcPort,
					SourceProtocol:          srcNet,
					SourceTransport:         srcTransport,
					DestinationIP:           dstIP,
					DestinationPort:         dstPort,
					DestinationProtocol:     dstNet,
					DestinationTransport:    dstTransport,
					SourceSentBytes:         srcSent,
					DestinationReceiveBytes: dstRecv,
				}

				serialized, err := entry.MarshalBinary()
				if err == nil {
					p.ConnectionLog.Write(append(serialized, []byte("\n")...))
				}
			}
		}()
	}

	p.DebugLog.Printf(
		"proxy: draining active connections before shutdown: name=%s connections=%d",
		p.Name,
		atomic.LoadInt64(activeConns),
	)

	connections.Wait()

	return nil
}

// Close requests asynchronous graceful shutdown of the proxy by atomically setting a shutdown flag.
// It is the responsibility of the caller to close the associated listener after closing the proxy.
func (p *Proxy) Close() {
	p.closing = true
}

// connLog describes a connection access log entry.
type connLog struct {
	Name                    string    `json:"name"`
	Timestamp               time.Time `json:"timestamp"`
	Duration                float64   `json:"duration"`
	SourceIP                string    `json:"source_ip"`
	SourcePort              int       `json:"source_port"`
	SourceProtocol          string    `json:"source_protocol"`
	SourceTransport         string    `json:"source_transport"`
	DestinationIP           string    `json:"destination_ip"`
	DestinationPort         int       `json:"destination_port"`
	DestinationProtocol     string    `json:"destination_protocol"`
	DestinationTransport    string    `json:"destination_transport"`
	SourceSentBytes         int64     `json:"source_sent_bytes"`
	DestinationReceiveBytes int64     `json:"destination_receive_bytes"`
}

// MarshalBinary serializes the connection log entry as a JSON payload.
func (c *connLog) MarshalBinary() ([]byte, error) {
	return json.Marshal(c)
}
