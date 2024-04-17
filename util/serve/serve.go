package serve

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"gopkg.in/tylerb/graceful.v1"
)

var listenPort string
var cfgtestFlag bool

func init() {
	flag.StringVar(&listenPort, "p", "", "listener port")
	flag.BoolVar(&cfgtestFlag, "t", false, "config test")
}

func Serve(hport string, handler http.Handler) error {

	checkConfigTest()

	l, err := Listen(hport)
	if err != nil {
		log.Fatalln(err)
	}

	srv := &graceful.Server{
		Timeout: 10 * time.Second,
		Server: &http.Server{
			Handler:      handler,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}

	log.Println("starting serve on ", hport)
	return srv.Serve(l)
}

// This method can be used for any TCP Listener, e.g. non HTTP
func Listen(hport string) (net.Listener, error) {
	var l net.Listener

	fd := os.Getenv("EINHORN_FDS")
	if fd != "" {
		sock, err := strconv.Atoi(fd)
		if err == nil {
			hport = "socketmaster:" + fd
			log.Println("detected socketmaster, listening on", fd)
			file := os.NewFile(uintptr(sock), "listener")
			fl, err := net.FileListener(file)
			if err == nil {
				l = fl
			}
		}
	}

	if listenPort != "" {
		hport = ":" + listenPort
	}

	checkConfigTest()

	if l == nil {
		var err error
		l, err = net.Listen("tcp4", hport)
		if err != nil {
			return nil, err
		}
	}

	return l, nil
}

func checkConfigTest() {
	if cfgtestFlag == true {
		log.Println("config test mode, exiting")
		os.Exit(0)
	}
}
