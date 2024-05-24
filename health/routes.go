package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Check(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Server is up and running")
}
