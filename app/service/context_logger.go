package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"strings"
	"time"
)

// provides an interface for logging to the request context
type ContextLogger interface {
	Log(message string)
	Debug(message string)
	Error(err error)
}

var (
	green  = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	white  = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	yellow = string([]byte{27, 91, 57, 55, 59, 52, 51, 109})
	red    = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	cyan   = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	reset  = string([]byte{27, 91, 48, 109})
)

type ContextLoggerImplementation struct {
	logs   []string
	debugs []string
	errors []error
}

func NewContextLogger() *ContextLoggerImplementation {
	return &ContextLoggerImplementation{
		logs:   make([]string, 0),
		debugs: make([]string, 0),
		errors: make([]error, 0),
	}
}

func (cli *ContextLoggerImplementation) Log(message string) {
	cli.logs = append(cli.logs, message)
}

func (cli *ContextLoggerImplementation) Debug(message string) {
	cli.debugs = append(cli.logs, message)
}

func (cli *ContextLoggerImplementation) Error(err error) {
	cli.errors = append(cli.errors, err)
}

func (cli *ContextLoggerImplementation) LoggerWithWriter(out io.Writer) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()
		ctx := GetContext(c)

		logger := ctx.Logger.(*ContextLoggerImplementation)

		var logs, debugs, errors string

		if len(logger.logs) > 0 {
			logs = fmt.Sprintf("\n[LOGS]\n    %s", strings.Join(logger.logs, "\n    "))
		}

		if len(logger.debugs) > 0 {
			debugs = fmt.Sprintf("\n[DEBUG]\n    %s", strings.Join(logger.debugs, "\n    "))
		}

		if len(logger.errors) > 0 {
			errors = "\n[ERRORS]\n"
			for _, err := range logger.errors {
				errors = errors + "    " + err.Error() + "\n"
			}
		}

		end := time.Now()
		latency := end.Sub(start)
		path := c.Request.URL.Path
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		statusColor := cli.colorForStatus(statusCode)
		methodColor := cli.colorForMethod(method)

		fmt.Fprintf(out, "[GIN] %v |%s %3d %s| latency %5.6f | %15s |%s %-7s %s %s %s\n",
			end.Format("2006/01/02 - 15:04:05.999"),
			statusColor,
			statusCode,
			reset,
			latency.Seconds(),
			clientIP,
			methodColor,
			method,
			reset,
			path,
			logs+debugs+errors,
		)
	}
}

func (cli *ContextLoggerImplementation) colorForStatus(code int) string {
	switch {
	case code >= 200 && code < 300:
		return green
	case code >= 300 && code < 400:
		return white
	case code >= 400 && code < 500:
		return yellow
	default:
		return red
	}
}

func (cli *ContextLoggerImplementation) colorForMethod(method string) string {
	switch method {
	case "POST":
		return cyan
	default:
		return reset
	}
}
