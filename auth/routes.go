package auth

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func Register(router *gin.Engine, client *mongo.Client, signingKey []byte) {
	controller := NewController(client, signingKey)
	api := router.Group("/auth")
	{
		api.POST("/signup", controller.SignUp)
		api.POST("/signin", controller.SignIn)
	}
}
