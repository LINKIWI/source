package main

import (
	"context"
	"flag"
	"log"
	"net"
	"strings"

	"github.com/tikv/client-go/v2/config"
	"github.com/tikv/client-go/v2/rawkv"
	"lib.kevinlin.info/mcrpc"

	"ticached/internal/meta"
	"ticached/internal/server"
)

var (
	flagListenAddress   = flag.String("listen-address", "localhost:11211", "server listen address")
	flagPDAddress       = flag.String("pd-address", "localhost:2379", "comma-delimited PD addresses")
	flagClientTLSKey    = flag.String("client-tls-key", "", "path to a client TLS key for TiKV cluster TLS")
	flagClientTLSCert   = flag.String("client-tls-cert", "", "path to a client TLS certificate for TiKV cluster TLS")
	flagClientTLSCACert = flag.String("client-tls-ca-cert", "", "path to a CA certificate for TiKV cluster TLS")
)

func init() {
	flag.Parse()
}

func main() {
	log.Printf("main: starting ticached: version=%s", meta.Version)

	ln, err := net.Listen("tcp", *flagListenAddress)
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
