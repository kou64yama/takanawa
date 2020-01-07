package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	"github.com/kou64yama/takanawa"
	"github.com/kou64yama/takanawa/internal/util"
	"github.com/kou64yama/takanawa/middleware"
)

const (
	usage = `Usage: takanawa [options] [path:]upstream [[path:]upstream ...]

Takanawa is a reverse proxy for HTTP services for development.

For example, execute the following command to proxy / to UI server and
/api to the API server:

  $ takanawa /api:http://localhost:8080/v1 http://localhost:3000

Takanawa runs on port 5000 by default.

Options:`
)

var (
	version = "0.0.0"
	logger  = log.New(os.Stdout, "", 0)

	showVersion          bool
	host                 string
	port                 uint
	accessLog            string
	overwriteHost        bool
	corsAllowedOrigins   string
	corsAllowedMethods   string
	corsAllowedHeaders   string
	corsExposedHeaders   string
	corsAllowCredentials bool
)

func main() {
	args, err := flags(os.Args[1:]...)
	if err != nil {
		os.Exit(exit(err))
		return
	}

	err = run(args...)
	if err != nil {
		logger.Println(err)
	}
	code := exit(err)
	os.Exit(code)
}

func exit(err error) int {
	if err == flag.ErrHelp {
		return 2
	}
	if err != nil {
		return 1
	}
	return 0
}

func flags(args ...string) ([]string, error) {
	f := flag.NewFlagSet("takanawa", flag.ContinueOnError)
	f.Usage = func() {
		o := f.Output()
		fmt.Fprintln(o, usage)
		f.PrintDefaults()
	}
	f.BoolVar(
		&showVersion,
		"version",
		false,
		"show version number",
	)
	f.StringVar(
		&host,
		"host",
		util.DefaultHost(),
		"specify listening host address",
	)
	f.UintVar(
		&port,
		"port",
		util.DefaultPort(),
		"specify listening port number",
	)
	f.StringVar(
		&accessLog,
		"access-log",
		"",
		"show access log (common|combined)",
	)
	f.BoolVar(
		&overwriteHost,
		"overwrite-host",
		true,
		"overwrite the 'Host' request header",
	)
	f.StringVar(
		&corsAllowedOrigins,
		"cors-allowed-origins",
		"",
		"specify allowed origins for CORS",
	)
	f.StringVar(
		&corsAllowedMethods,
		"cors-allowed-methods",
		"",
		"specify allowed methods for CORS",
	)
	f.StringVar(
		&corsAllowedHeaders,
		"cors-allowed-headers",
		"",
		"specify allowed headers for CORS",
	)
	f.StringVar(
		&corsExposedHeaders,
		"cors-exposed-headers",
		"",
		"specify exposed headers for CORS",
	)
	f.BoolVar(
		&corsAllowCredentials,
		"cors-allow-credentials",
		false,
		"allow credentials for CORS",
	)
	if err := f.Parse(args); err != nil {
		return nil, err
	}
	if f.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "no upstream arguments")
		f.Usage()
		return nil, errors.New("no upstream arguments")
	}
	return f.Args(), nil
}

func run(args ...string) error {
	if showVersion {
		fmt.Println("takanawa", version)
		return nil
	}

	t := &takanawa.Takanawa{}
	t.Middleware(middleware.RequestID())

	if len(accessLog) > 0 {
		var format middleware.AccessLogFormat
		switch accessLog {
		case "common":
			format = middleware.CommonLogFormat
		case "combined":
			format = middleware.CombinedLogFormat
		default:
			return fmt.Errorf("invalid access log format: %s", accessLog)
		}
		opt := &middleware.AccessLogOption{
			Log:    logger,
			Format: format,
		}
		t.Middleware(middleware.AccessLog(opt))
	}

	allowedOrigins := util.SplitAndTrimSpace(corsAllowedOrigins, ",")
	if len(allowedOrigins) > 0 {
		opt := &middleware.CorsOption{}
		opt.AllowedOrigins = allowedOrigins
		opt.AllowedMethods = util.SplitAndTrimSpace(corsAllowedMethods, ",")
		opt.AllowedHeaders = util.SplitAndTrimSpace(corsAllowedHeaders, ",")
		opt.ExposedHeaders = util.SplitAndTrimSpace(corsExposedHeaders, ",")
		opt.AllowCredentials = corsAllowCredentials
		t.Middleware(middleware.Cors(opt))
	}

	for _, v := range args {
		u, opt, err := middleware.ParseReverseProxyOption(v)
		if err != nil {
			return err
		}

		opt.OverwriteHost = overwriteHost
		opt.ErrorLog = log.New(ioutil.Discard, "", 0)
		t.Middleware(middleware.ReverseProxy(u, opt))
		logger.Printf("Proxy: %q -> %s", opt.Path, u)
	}

	addr := host + ":" + strconv.Itoa(int(port))
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	srv := &http.Server{Handler: t.Handler()}
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		// We received an interrupt signal, shut down.
		if err := srv.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			logger.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	logger.Printf("The server is running at http://%s", ln.Addr().String())
	if err := srv.Serve(ln); err != nil {
		return err
	}

	<-idleConnsClosed
	return nil
}
