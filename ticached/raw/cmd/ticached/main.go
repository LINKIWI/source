package main

import (
	"context"
	"crypto/tls"
	"flag"
	"log"
	"strings"

	"github.com/tikv/client-go/v2/config"
	"github.com/tikv/client-go/v2/rawkv"
	"lib.kevinlin.info/mcrpc"

	"ticached/internal/meta"
	"ticached/internal/server"
)

var (
	flagListenAddress       = flag.String("listen-address", "localhost:11211", "server listen address")
	flagPDAddress           = flag.String("pd-address", "http://localhost:2379", "comma-delimited PD URL addresses")
	flagServerTLSKey        = flag.String("server-tls-key", "", "path to a server TLS key for TLS termination")
	flagServerTLSCert       = flag.String("server-tls-cert", "", "path to a server TLS certificate for TLS termination")
	flagServerTLSCACert     = flag.String("server-tls-ca-cert", "", "path to a server TLS CA certificate for TLS termination")
	flagServerTLSClientAuth = flag.Bool("server-tls-client-auth", false, "enable and require valid client TLS authentication for TLS termination")
	flagClientTLSKey        = flag.String("client-tls-key", "", "path to a client TLS key for TiKV cluster TLS")
	flagClientTLSCert       = flag.String("client-tls-cert", "", "path to a client TLS certificate for TiKV cluster TLS")
	flagClientTLSCACert     = flag.String("client-tls-ca-cert", "", "path to a TLS CA certificate for TiKV cluster TLS")
)

func init() {
	flag.Parse()
}

func main() {
	log.Printf("main: starting ticached: version=%s", meta.Version)

	listener := &server.Listener{
		Network:       "tcp",
		Address:       *flagListenAddress,
		TLSKey:        *flagServerTLSKey,
		TLSCert:       *flagServerTLSCert,
		TLSCACert:     *flagServerTLSCACert,
		TLSClientAuth: tls.NoClientCert,
	}

	if *flagServerTLSClientAuth {
		listener.TLSClientAuth = tls.RequireAndVerifyClientCert
	}

	ln, err := listener.Listen()
	if err != nil {
		panic(err)
	}

	security := config.DefaultConfig().Security
	if *flagClientTLSKey != "" && *flagClientTLSCert != "" {
		log.Printf(
			"main: enabling TiKV cluster TLS: key=%s cert=%s ca=%s",
			*flagClientTLSKey,
			*flagClientTLSCert,
			*flagClientTLSCACert,
		)

		security = config.NewSecurity(
			*flagClientTLSCACert,
			*flagClientTLSCert,
			*flagClientTLSKey,
			[]string{},
		)
	}

	log.Printf("main: connecting to TiKV cluster: pd=%s", *flagPDAddress)
	kv, err := rawkv.NewClient(context.Background(), strings.Split(*flagPDAddress, ","), security)
	if err != nil {
		panic(err)
	}

	handler := server.NewTiKVHandler(kv)
	srv, err := mcrpc.NewServer(handler, mcrpc.WithErrorLog(log.Default()))
	if err != nil {
		panic(err)
	}

	log.Printf("main: starting memcache protocol listener: listen=%s", *flagListenAddress)
	if err := srv.Serve(ln); err != nil {
		panic(err)
	}
}
