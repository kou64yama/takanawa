package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/kou64yama/takanawa"
	"github.com/kou64yama/takanawa/internal/env"
)

const (
	synopsis = "takanawa [options] <upstream>"
)

var (
	version = "0.0.0"
)

func main() {
	err := run(os.Args[1:]...)
	if err == flag.ErrHelp {
		os.Exit(2)
	}
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func run(args ...string) error {
	f := flag.NewFlagSet("takanawa", flag.ContinueOnError)
	f.Usage = func() {
		o := f.Output()
		fmt.Fprintf(o, "Usage: %s\n", synopsis)
		fmt.Fprintln(o, "")
		fmt.Fprintln(o, "Options:")
		f.PrintDefaults()
	}

	port := f.Uint(
		"port",
		env.UintEnv("PORT", 5000),
		"specify the port",
	)
	host := f.String(
		"host",
		env.StringEnv("HOST", "127.0.0.1"),
		"specify the host address",
	)
	overwriteHost := f.Bool(
		"overwrite-host",
		env.BoolEnv("OVERWRITE_HOST", true),
		"overwrite Host header",
	)
	forwarded := f.Bool(
		"forwarded",
		env.BoolEnv("FORWARDED", true),
		"Add the Forwarded request header",
	)
	corsEnabled := f.Bool(
		"cors",
		env.BoolEnv("CORS_ENABLED", false),
		"enable CORS",
	)
	corsAllowOrigin := f.String(
		"cors-allow-origin",
		env.StringEnv("CORS_ALLOW_ORIGIN", "*"),
		"specify Access-Control-Allow-Origin",
	)
	corsAllowMethods := f.String(
		"cors-allow-methods",
		env.StringEnv("CORS_ALLOW_METHODS", "GET,POST,PUT,DELETE,PATCH,HEAD,OPTIONS,TRACE"),
		"specify Access-Control-Allow-Methods",
	)
	corsAllowHeaders := f.String(
		"cors-allow-headers",
		env.StringEnv("CORS_ALLOW_HEADERS", "Content-Type,Content-Length,Accept,Accept-Encoding"),
		"specify Access-Control-Allow-Headers",
	)
	corsExposeHeaders := f.String(
		"cors-expose-headers",
		env.StringEnv("CORS_EXPOSE_HEADERS", "Content-Type,Content-Length,Content-Encoding"),
		"specify Access-Control-Expose-Headers",
	)

	if err := f.Parse(args); err != nil {
		return err
	}
	if f.NArg() != 1 {
		f.Usage()
		return flag.ErrHelp
	}

	target, err := url.Parse(f.Arg(0))
	if err != nil {
		return err
	}

	middlewares := []takanawa.Middleware{
		takanawa.RequestID(),
	}
	if *forwarded {
		middlewares = append(middlewares, takanawa.ForwardedMiddleware())
	}
	if *corsEnabled {
		cors := &takanawa.Cors{
			AllowOrigin:   *corsAllowOrigin,
			AllowMethods:  strings.Split(*corsAllowMethods, ","),
			AllowHeaders:  strings.Split(*corsAllowHeaders, ","),
			ExposeHeaders: strings.Split(*corsExposeHeaders, ","),
		}
		middlewares = append(middlewares, takanawa.CorsMiddleware(cors))
	}

	middlewares = append(middlewares, takanawa.ProxyMiddleware(target, *overwriteHost))
	handler := takanawa.ComposeMiddlewares(middlewares...)

	addr := fmt.Sprintf("%s:%d", *host, *port)

	log.Printf("The server is running at http://localhost:%d/", *port)
	return http.ListenAndServe(addr, handler)
}
