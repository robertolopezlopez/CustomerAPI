package tracing

import (
	"api/logging"
	"api/tools"
	"time"

	"github.com/gin-gonic/gin"
)

const XRequestID = "X-RequestID"

func XRequestIDMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Request.Header.Get(XRequestID)
		if id == "" {
			var err error
			id, err = tools.GenerateUUID4()
			if err != nil {
				logging.ErrorLogger.Printf("cannot generate UUID4: %s", err.Error())
				return
			}

			ctx.Request.Header.Add(XRequestID, id)
		}

		ts := time.Now().Format(time.UnixDate)
		logging.InfoLogger.Printf("%s - %s\n", ts, id)
	}
}
