package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/net/proxy"
	"lib.kevinlin.info/proton"

	"wg/internal/cli"
	"wg/internal/meta"
	"wg/pkg/webgrep"
)

const (
	envWebgrepURL     = "WG_WEBGREP_URL"
	envProxyAddr      = "WG_PROXY_ADDR"
	envTLSCertificate = "WG_TLS_CERTIFICATE"
	envTLSKey         = "WG_TLS_KEY"
)

var (
	flagWebgrepURL     = flag.String("webgrep-url", os.Getenv(envWebgrepURL), "base URL of the webgrep instance")
	flagRegex          = flag.Bool("regex", false, "interpret search query as a regular expression")
	flagCaseSensitive  = flag.Bool("case-sensitive", false, "respect search query case sensitivity")
	flagFile           = flag.String("file", "", "filter matches by file path pattern")
	flagMaxMatches     = flag.Int("max-matches", 50, "maximum number of matches in search results")
	flagContext        = flag.Int("context", 4, "number of surrounding context lines to include around matches")
	flagProxy          = flag.String("proxy", os.Getenv(envProxyAddr), "optional address of a SOCKS5 proxy server for establishing a connection")
	flagTLSCertificate = flag.String("tls-certificate", os.Getenv(envTLSCertificate), "optional path to a client TLS certificate for mutual authentication")
	flagTLSKey         = flag.String("tls-key", os.Getenv(envTLSKey), "optional path to a client TLS key for mutual authentication")
	flagAbout          = cli.NewChoicesFlag([]string{"index", "repos"}, "")
	flagRepos          = cli.NewArrayFlag()
	flagSearchType     = cli.NewChoicesFlag([]string{"files", "code"}, "code")
	flagRenderer       = cli.NewChoicesFlag([]string{"table", "stacked"}, "table")
)

func init() {
	flag.Var(flagAbout, "about", fmt.Sprintf("print current server-side index information; one of %v", flagAbout.Candidates()))
	flag.Var(flagRepos, "repo", "filter matches by repository name")
	flag.Var(flagSearchType, "search-type", fmt.Sprintf("search results type to print; one of %v", flagSearchType.Candidates()))
	flag.Var(flagRenderer, "renderer", fmt.Sprintf("renderer style for results; one of %v", flagRenderer.Candidates()))
	flag.Parse()
}

func main() {
	// Rudimentary input validation
	if *flagWebgrepURL == "" {
		panic(fmt.Errorf("main: no value specified for webgrep instance URL (see --help for docs)"))
	}

	// Client HTTP backend
	dialer := &net.Dialer{Timeout: 30 * time.Second}
	transport := &http.Transport{DialContext: dialer.DialContext}

	// Optional proxy server configuration
	if *flagProxy != "" {
		proxyDialer, err := proxy.SOCKS5("tcp", *flagProxy, nil, dialer)
		if err != nil {
			panic(err)
		}

		transport = &http.Transport{
			DialContext: func(ctx context.Context, network string, addr string) (net.Conn, error) {
				return proxyDialer.Dial(network, addr)
			},
		}
	}

	// Optional mutual TLS authentication configuration
	if *flagTLSCertificate != "" && *flagTLSKey != "" {
		cert, err := tls.LoadX509KeyPair(*flagTLSCertificate, *flagTLSKey)
		if err != nil {
			panic(err)
		}

		transport.TLSClientConfig = &tls.Config{Certificates: []tls.Certificate{cert}}
	}

	// Instantiate a webgrep client
	webgrepAddr, err := url.Parse(*flagWebgrepURL)
	if err != nil {
		panic(err)
	}
	client, err := webgrep.NewClient(&proton.Config{
		BaseURL:       webgrepAddr,
		ClientID:      meta.Identifier,
		ClientVersion: meta.Version,
		Backend:       &http.Client{Transport: transport},
	})
	if err != nil {
		panic(err)
	}

	// Application and index metadata
	if flagAbout.Choice() != "" {
		if err := about(client); err != nil {
			panic(err)
		}
		return
	}

	// Code search query
	if err := search(client); err != nil {
		panic(err)
	}
}

func about(client *webgrep.Client) error {
	table := cli.NewTable()

	metadata, err := client.Metadata()
	if err != nil {
		return err
	}

	switch flagAbout.Choice() {
	case "index":
		table.Add([]string{"wg client version:", meta.Version})
		table.Add([]string{"webgrep server version:", metadata.Version})
		table.Add([]string{"index name:", metadata.Name})
		table.Add([]string{"index timestamp:", time.Unix(int64(metadata.Timestamp), 0).String()})
		table.Add([]string{"index repositories:", strconv.Itoa(len(metadata.Repositories))})

	case "repos":
		sort.SliceStable(metadata.Repositories, func(i, j int) bool {
			// Sort in descending order of repository name
			return metadata.Repositories[i].Name < metadata.Repositories[j].Name
		})

		for _, repo := range metadata.Repositories {
			table.Add([]string{
				repo.Name,
				repo.Version,
				strings.Join(repo.Labels, ","),
				repo.Remote,
			})
		}

	default:
	}

	fmt.Println(table)

	return nil
}

func search(client *webgrep.Client) error {
	// Read search query as input from stdin
	input, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return err
	}

	// Execute the search, respecting CLI flags as parameters
	resp, searchErr := client.Search(&webgrep.SearchQueryRequest{
		Query:         strings.TrimSpace(input),
		File:          *flagFile,
		Repos:         flagRepos.Values(),
		Regex:         *flagRegex,
		CaseSensitive: *flagCaseSensitive,
		MaxMatches:    *flagMaxMatches,
		Context:       *flagContext,
	})
	if searchErr != nil {
		return searchErr
	}

	// Select a renderer consistent with the requested output style
	renderer, ok := renderers[flagRenderer.Choice()]
	if !ok {
		return fmt.Errorf("main: no renderer available for the requested output style")
	}

	// Format results as requested by parameters
	switch flagSearchType.Choice() {
	case "code":
		code, err := renderer.RenderCodeSearchResults(resp.Code)
		if err != nil {
			return err
		}

		if code != "" {
			fmt.Println(code)
		}

	case "files":
		files, err := renderer.RenderFileSearchResults(resp.Files)
		if err != nil {
			return err
		}

		if files != "" {
			fmt.Println(files)
		}

	default:
	}

	return nil
}
