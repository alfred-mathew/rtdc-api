package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Check(context *gin.Context) {
	context.String(http.StatusOK, "Server is up and running")
}
