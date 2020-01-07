package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/kou64yama/takanawa"
	"github.com/kou64yama/takanawa/internal/util"
	"github.com/kou64yama/takanawa/middleware"
)

const (
	usageText = `
	Usage: %s [OPTIONS]

	Takanawa is a reverse proxy for HTTP services for development.

	For example, execute the following command to proxy / to UI server and
	/api to the API server:

	  $ takanawa -access-log=common \
	      http://localhost:3000 -change-origin \
	      http://localhost:8080/v1 -change-origin -path=/api
	`
)

var (
	version = "0.0.0"
)

var (
	errExitUsage = errors.New("exit usage")
)

func main() {
	err := run(os.Args[1:]...)
	if err == errExitUsage {
		os.Exit(2)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func usage(f *flag.FlagSet, u string) func() {
	return func() {
		u = strings.ReplaceAll(u, "\n\t", "\n")
		u = strings.TrimSpace(u)
		o := f.Output()
		fmt.Printf(u, f.Name())
		fmt.Fprintln(o, "\n\nOptions:")
		f.PrintDefaults()
	}
}

func run(args ...string) error {
	logger := log.New(os.Stdout, "", 0)

	var (
		showVersion bool
		listenAddr  string
		port        uint
		accessLog   string
	)
	globalFlags := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	globalFlags.Usage = usage(globalFlags, usageText)
	globalFlags.BoolVar(
		&showVersion,
		"version",
		false,
		"show version number",
	)
	globalFlags.StringVar(
		&listenAddr,
		"listen-address",
		util.DefaultHost(),
		"specify listening address",
	)
	globalFlags.UintVar(
		&port,
		"port",
		util.DefaultPort(),
		"specify listening port number",
	)
	globalFlags.StringVar(
		&accessLog,
		"access-log",
		"",
		"show access log (common|combined)",
	)
	if err := globalFlags.Parse(args); err != nil {
		return errExitUsage
	}
	if showVersion {
		fmt.Printf("Takanawa %s\n", version)
		return nil
	}

	globalMids := []takanawa.Middleware{middleware.RequestID()}
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
		opt := &middleware.AccessLogOptions{
			Logger: logger,
			Format: format,
		}
		globalMids = append(globalMids, middleware.AccessLog(opt))
	}

	mu := http.NewServeMux()
	srv := &http.Server{Handler: takanawa.ComposeMiddleware(globalMids...).Apply(mu)}

	args = globalFlags.Args()
	for len(args) > 0 {
		var (
			upstream             = args[0]
			path                 string
			changeOrigin         bool
			corsAllowedOrigins   string
			corsAllowedMethods   string
			corsAllowedHeaders   string
			corsExposedHeaders   string
			corsAllowCredentials bool
		)
		proxyFlags := flag.NewFlagSet(os.Args[0]+" UPSTREAM", flag.ContinueOnError)
		proxyFlags.Usage = usage(proxyFlags, usageText)
		proxyFlags.StringVar(
			&path,
			"path",
			"/",
			"specify path to proxy",
		)
		proxyFlags.BoolVar(
			&changeOrigin,
			"change-origin",
			false,
			"overwrite the 'Host' request header",
		)
		proxyFlags.StringVar(
			&corsAllowedOrigins,
			"cors-allowed-origins",
			"",
			"specify allowed origins for CORS",
		)
		proxyFlags.StringVar(
			&corsAllowedMethods,
			"cors-allowed-methods",
			"",
			"specify allowed methods for CORS",
		)
		proxyFlags.StringVar(
			&corsAllowedHeaders,
			"cors-allowed-headers",
			"",
			"specify allowed headers for CORS",
		)
		proxyFlags.StringVar(
			&corsExposedHeaders,
			"cors-exposed-headers",
			"",
			"specify exposed headers for CORS",
		)
		proxyFlags.BoolVar(
			&corsAllowCredentials,
			"cors-allow-credentials",
			false,
			"allow credentials for CORS",
		)
		if err := proxyFlags.Parse(args[1:]); err != nil {
			return errExitUsage
		}
		u, err := url.Parse(upstream)
		if err != nil {
			return err
		}

		mids := []takanawa.Middleware{}
		if changeOrigin {
			mids = append(mids, middleware.ChangeOrigin(u.Host))
		}
		if allowedOrigins := util.SplitAndTrimSpace(corsAllowedOrigins, ","); len(allowedOrigins) > 0 {
			opt := &middleware.CorsOption{}
			opt.AllowedOrigins = allowedOrigins
			opt.AllowedMethods = util.SplitAndTrimSpace(corsAllowedMethods, ",")
			opt.AllowedHeaders = util.SplitAndTrimSpace(corsAllowedHeaders, ",")
			opt.ExposedHeaders = util.SplitAndTrimSpace(corsExposedHeaders, ",")
			opt.AllowCredentials = corsAllowCredentials
			mids = append(mids, middleware.Cors(opt))
		}
		mids = append(mids, middleware.StripPrefix(path), middleware.ReverseProxy(u))
		mu.Handle(path, takanawa.ComposeMiddleware(mids...).Apply(http.NotFoundHandler()))

		args = proxyFlags.Args()
	}

	logger.Printf("Takanawa %s", version)
	logger.Printf("PID %d", os.Getpid())

	idleConnsClosed := make(chan struct{})
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
		logger.Printf("Shutdown: %s", <-sig)

		// We received an SIGHUP, SIGINT, SIGQUIT or SIGTERM signal, shut down.
		if err := srv.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			logger.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()
	ln, err := net.Listen("tcp", listenAddr+":"+strconv.Itoa(int(port)))
	if err != nil {
		return err
	}

	logger.Printf("Listen on %s", ln.Addr().String())
	if err := srv.Serve(ln); err != http.ErrServerClosed {
		return err
	}

	<-idleConnsClosed
	return nil
}
