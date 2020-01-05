package middleware

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/kou64yama/takanawa"
	"github.com/kou64yama/takanawa/internal/util"
)

var (
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

type AccessLogFormat func(t *time.Time, s *util.ResponseSniffer, r *http.Request) string

type AccessLogOption struct {
	Log    *log.Logger
	Format AccessLogFormat
}

func AccessLog(opt *AccessLogOption) takanawa.Middleware {
	if opt == nil {
		opt = &AccessLogOption{}
	}
	logger := opt.Log
	if logger == nil {
		logger = log.New(os.Stdout, "", 0)
	}
	format := opt.Format
	if format == nil {
		format = CommonLogFormat
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t := time.Now().Local()
			s := &util.ResponseSniffer{Writer: w, StatusCode: 200}
			next.ServeHTTP(s, r)
			logger.Println(format(&t, s, r))
		})
	}
}
