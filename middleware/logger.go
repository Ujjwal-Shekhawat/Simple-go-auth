package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger - Userdefined logger to use as a middleware for logging events (developement mode)
func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(params gin.LogFormatterParams) string {
		return fmt.Sprintf("\033[97;41mIP-%s\033[0m-\033[97;45mTIME-[%s]\033[0m-%sMETHOD-%s%s-\033[90;47mPATH-%s\033[0m-%sSTATUSCODE-%d%s-\033[35mLATENCY-%s\033[0m\n",
			params.ClientIP,
			params.TimeStamp.Format(time.RFC822),
			params.MethodColor(),
			params.Method,
			params.ResetColor(),
			params.Path,
			params.StatusCodeColor(),
			params.StatusCode,
			params.ResetColor(),
			params.Latency,
		)
	})
}

// ProductionLogger - Userdefined logger to use as a middleware for logging events (production mode)
func ProductionLogger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(params gin.LogFormatterParams) string {
		return fmt.Sprintf("\033[97;41mIP-%s\033[0m-\033[97;45mTIME-[%s]\033[0m-%sMETHOD-%s%s-\033[90;47mPATH-%s\033[0m-%sSTATUSCODE-%d%s-\033[35mLATENCY-%s\033[0m\n",
			params.ClientIP,
			params.TimeStamp.Format(time.RFC822),
			params.MethodColor(),
			params.Method,
			params.ResetColor(),
			params.Path,
			params.StatusCodeColor(),
			params.StatusCode,
			params.ResetColor(),
			params.Latency,
		)
	})
}
