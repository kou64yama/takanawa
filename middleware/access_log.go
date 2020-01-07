package middleware

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/kou64yama/takanawa"
	"github.com/kou64yama/takanawa/internal/util"
)

var (
	// CommonLogFormat formats access logs according to the
	// Apache's Common Log Format.
	//
	// Common Log Format: https://httpd.apache.org/docs/2.4/en/logs.html#common
	CommonLogFormat = func(t *time.Time, w *util.ResponseSniffer, r *http.Request) string {
		remoteAddr := strings.SplitN(r.RemoteAddr, ":", 2)
		return fmt.Sprintf(
			"%s - - [%s] \"%s %s %s\" %d %d",
			remoteAddr[0],
			t.Format("02/Jan/2006:15:04:05 -0700"),
			r.Method,
			r.RequestURI,
			r.Proto,
			w.StatusCode,
			w.Length,
		)
	}

	// CombinedLogFormat formats access logs according to the
	// Apache's Combined Log Format.
	//
	// Combined Log Format: https://httpd.apache.org/docs/2.4/en/logs.html#combined
	CombinedLogFormat = func(t *time.Time, w *util.ResponseSniffer, r *http.Request) string {
		referer := r.Referer()
		if len(referer) == 0 {
			referer = "-"
		}
		userAgent := r.UserAgent()
		if len(userAgent) == 0 {
			userAgent = "-"
		}
		return CommonLogFormat(t, w, r) + fmt.Sprintf(" \"%s\" \"%s\"", referer, userAgent)
	}
)

var (
	// DefaultAccessLogOptions is used by AccessLog as default options.
	DefaultAccessLogOptions = AccessLogOptions{
		Logger: log.New(ioutil.Discard, "", 0),
		Format: CommonLogFormat,
	}
)

// AccessLogFormat formats access logs.
type AccessLogFormat func(t *time.Time, s *util.ResponseSniffer, r *http.Request) string

// AccessLogOptions is options of AccessLog.
type AccessLogOptions struct {
	Logger *log.Logger
	Format AccessLogFormat
}

// AccessLog returns the middleware.
func AccessLog(opt *AccessLogOptions) takanawa.Middleware {
	if opt == nil {
		opt = &DefaultAccessLogOptions
	}
	logger := opt.Logger
	if logger == nil {
		logger = DefaultAccessLogOptions.Logger
	}
	format := opt.Format
	if format == nil {
		format = DefaultAccessLogOptions.Format
	}

	return takanawa.MiddlewareFunc(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t := time.Now().Local()
			s := &util.ResponseSniffer{Writer: w, StatusCode: 200}
			next.ServeHTTP(s, r)
			logger.Println(format(&t, s, r))
		})
	})
}
