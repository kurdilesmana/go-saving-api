package mid

import (
	"fmt"
	"net/http"
	"time"

	"github.com/kurdilesmana/go-saving-api/pkg/logging"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type reqResTrap struct {
	logger *logging.Logger
}

func NewReqResMiddleware(logger *logging.Logger) *reqResTrap {
	return &reqResTrap{
		logger: logger,
	}
}

func (rr *reqResTrap) Middle(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		req := c.Request()
		res := c.Response()
		start := time.Now()

		// traceID
		traceID := GetID(c)

		// referer
		referer := getFirstValueFromHeader(req, "Referer")
		clientID := getFirstValueFromHeader(req, "X-Client-Id")
		// source ip
		ipAddress := getFirstValueFromHeader(req, "X-Forwarded-For")

		// LOG for incoming request
		rr.logger.Info(
			logrus.Fields{
				"referer":   referer,
				"client_id": clientID,
				"source_ip": ipAddress,
				"path":      req.URL.Path,
				"trace_id":  traceID,
			}, nil, "request started",
		)

		// LOG after request end
		defer func() {
			rr.logger.Info(
				logrus.Fields{
					"referer":   referer,
					"client_id": clientID,
					"source_ip": ipAddress,
					"path":      req.URL.Path,
					"trace_id":  traceID,
					"latency":   fmt.Sprintf("%d ms", time.Since(start).Milliseconds()),
					"status":    res.Status,
				}, nil, "request completed",
			)
		}()

		if err = next(c); err != nil {
			c.Error(err)
		}
		return
	}
}

func getFirstValueFromHeader(req *http.Request, key string) string {
	vs, ok := req.Header[key]
	if ok {
		if len(vs) != 0 {
			return vs[0]
		}
	}
	return ""
}
