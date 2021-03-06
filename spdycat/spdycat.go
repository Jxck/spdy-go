package main

import (
	"flag"
	"fmt"
	"github.com/shykes/spdy-go"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type Server struct{}

func headersString(headers http.Header) string {
	var s string
	for key := range headers {
		if len(s) != 0 {
			s = s + " "
		}
		s += fmt.Sprintf("%s=%s", key, headers.Get(key))
	}
	return s
}

func (server *Server) ServeSPDY(stream *spdy.Stream) {
	stream.Output.Headers().Add(":status", "200")
	stream.Output.SendHeaders(false)
	processStream(stream)
}

func processStream(stream *spdy.Stream) {
	if stream.Id == 1 {
		go func() {
			_, err := io.Copy(stream.Output, os.Stdin)
			if err != nil {
				fmt.Printf("Error while sending to stream: %v\n", err)
				stream.Input.Error(err)
				stream.Output.Error(err)
			} else {
				os.Exit(0)
			}
		}()
	}
	_, err := io.Copy(os.Stdout, stream.Input)
	if err != nil {
		fmt.Printf("Error while printing stream: %v\n", err)
		stream.Input.Error(err)
		stream.Output.Error(err)
	}
}

func main() {
	listen := flag.Bool("l", false, "Listen to <addr>")
	tls := flag.Bool("t", false, "Enable TLS")
	cert := flag.String("cert", "cert.pem", "Filename to a TLS certificate (use in combination with -t and -l)")
	key  := flag.String("key", "key.pem", "Filename to a TLS private key (use in combination with -t and -l)")

	flag.Parse()
	addr := flag.Args()[0]
	headers := extractHeaders(flag.Args()[1:])
	server := &Server{} // FIXME: find another name for Server since it is used by both sides
	if *listen {
		var err error
		if *tls {
                    err = spdy.ListenAndServeTLS(addr, *cert, *key, server)
		} else {
		    err = spdy.ListenAndServeTCP(addr, server)
		}
		if err != nil {
			log.Fatal("Listen: %s", err)
		}
	} else {
		var err error
		var session *spdy.Session
		if *tls {
		    session, err = spdy.DialTLS(addr, server)
		} else {
		    session, err = spdy.DialTCP(addr, server)
		}
		if err != nil {
			log.Fatal("Error connecting: %s", err)
		}
		stream, err := session.OpenStream(headers)
		if err != nil {
			log.Fatal("Error opening stream: %s", err)
		}
		processStream(stream)
	}
}

func extractHeaders(args []string) *http.Header {
	headers := http.Header{}
	for _, keyvalue := range args {
		pair := strings.SplitN(keyvalue, "=", 2)
		headers.Set(pair[0], pair[1])
	}
	return &headers
}
