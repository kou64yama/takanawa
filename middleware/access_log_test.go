package middleware_test

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/kou64yama/takanawa/internal/mock"
	"github.com/kou64yama/takanawa/internal/util"
	"github.com/kou64yama/takanawa/middleware"
)

func TestCommonLogFormat(t *testing.T) {
	tests := []struct {
		time       string
		request    string
		remoteAddr string
		statusCode int
		length     uint64
		out        string
	}{
		{
			time:       "02/Jan/2006:15:04:05 -0700",
			request:    "GET / HTTP/1.1",
			remoteAddr: "127.0.0.1:1234",
			statusCode: 200,
			length:     1024,
			out:        "127.0.0.1 - - [02/Jan/2006:15:04:05 -0700] \"GET / HTTP/1.1\" 200 1024",
		},
		{
			time:       "02/Jan/2006:15:04:05 -0700",
			request:    "POST /greeting HTTP/1.1",
			remoteAddr: "127.0.0.1:1234",
			statusCode: 400,
			out:        "127.0.0.1 - - [02/Jan/2006:15:04:05 -0700] \"POST /greeting HTTP/1.1\" 400 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.request, func(t *testing.T) {
			t.Helper()

			s := strings.SplitN(tt.request, " ", 3)
			time, _ := time.Parse(tt.time, "02/Jan/2006:15:04:05 -0700")
			w := &util.ResponseSniffer{
				StatusCode: tt.statusCode,
				Length:     tt.length,
			}
			r := &http.Request{
				Method:     s[0],
				RequestURI: s[1],
				Proto:      s[2],
				RemoteAddr: tt.remoteAddr,
			}
			out := middleware.CommonLogFormat(&time, w, r)

			if out != tt.out {
				t.Errorf("got %q, want %q", out, tt.out)
			}
		})
	}
}

func TestCombinedLogFormat(t *testing.T) {
	tests := []struct {
		time       string
		request    string
		remoteAddr string
		referer    string
		userAgent  string
		statusCode int
		length     uint64
		out        string
	}{
		{
			time:       "02/Jan/2006:15:04:05 -0700",
			request:    "GET / HTTP/1.1",
			remoteAddr: "127.0.0.1:1234",
			referer:    "/index.html",
			userAgent:  "go/test",
			statusCode: 200,
			length:     1024,
			out:        "127.0.0.1 - - [02/Jan/2006:15:04:05 -0700] \"GET / HTTP/1.1\" 200 1024 \"/index.html\" \"go/test\"",
		},
		{
			time:       "02/Jan/2006:15:04:05 -0700",
			request:    "GET / HTTP/1.1",
			remoteAddr: "127.0.0.1:1234",
			statusCode: 200,
			length:     1024,
			out:        "127.0.0.1 - - [02/Jan/2006:15:04:05 -0700] \"GET / HTTP/1.1\" 200 1024 \"-\" \"-\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.request+",Referer="+tt.referer+",User-Agent="+tt.userAgent, func(t *testing.T) {
			t.Helper()

			header := http.Header{}
			if len(tt.referer) > 0 {
				header.Set("Referer", tt.referer)
			}
			if len(tt.userAgent) > 0 {
				header.Set("User-Agent", tt.userAgent)
			}

			s := strings.SplitN(tt.request, " ", 3)
			time, _ := time.Parse(tt.time, "02/Jan/2006:15:04:05 -0700")
			w := &util.ResponseSniffer{
				StatusCode: tt.statusCode,
				Length:     tt.length,
			}
			r := &http.Request{
				Method:     s[0],
				RequestURI: s[1],
				Proto:      s[2],
				RemoteAddr: tt.remoteAddr,
				Header:     header,
			}
			out := middleware.CombinedLogFormat(&time, w, r)

			if out != tt.out {
				t.Errorf("got %q, want %q", out, tt.out)
			}
		})
	}
}

func TestAccessLog(t *testing.T) {
	defaultOpt := middleware.DefaultAccessLogOptions
	defer func() { middleware.DefaultAccessLogOptions = defaultOpt }()

	buf := bytes.NewBuffer(nil)
	logger := log.New(buf, "", 0)
	middleware.DefaultAccessLogOptions.Logger = logger
	middleware.DefaultAccessLogOptions.Format = func(t *time.Time, s *util.ResponseSniffer, r *http.Request) string {
		return "!!OUTPUT!!"
	}

	tests := []*middleware.AccessLogOptions{
		nil,
		{
			Logger: middleware.DefaultAccessLogOptions.Logger,
		},
		{
			Format: middleware.DefaultAccessLogOptions.Format,
		},
	}
	for _, opt := range tests {
		t.Run(fmt.Sprint(opt), func(t *testing.T) {
			t.Helper()

			buf.Reset()

			h := &mock.Handler{}
			w := &mock.ResponseWriter{}
			r := &http.Request{}
			middleware.AccessLog(opt).Apply(h).ServeHTTP(w, r)

			if string(buf.Bytes()) != "!!OUTPUT!!\n" {
				t.Errorf("got %q, want %q", string(buf.Bytes()), "!!OUTPUT!!\n")
			}
		})
	}
}
