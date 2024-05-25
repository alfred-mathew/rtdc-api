package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func check(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Server is up and running")
}
