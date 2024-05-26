package auth

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type Controller interface {
	signUp(c *gin.Context)
	signIn(c *gin.Context)
	middleware(c *gin.Context)
}

type controller struct {
	service Service
}

func NewController(client *mongo.Client, signingKey []byte) Controller {
	return controller{
		service: newService(client, signingKey),
	}
}

func (c controller) signUp(ctx *gin.Context) {
	var creds Credentials
	if err := ctx.BindJSON(&creds); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validateCredentials(creds); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exists, err := c.service.userExists(creds.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to query database for existing users: %s", err.Error())})
		return
	}
	if exists {
		ctx.JSON(http.StatusConflict, gin.H{"error": "username is already taken"})
		return
	}

	err = c.service.createUser(creds)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func (c controller) signIn(ctx *gin.Context) {
	var creds Credentials
	if err := ctx.BindJSON(&creds); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := c.service.findUser(creds.Username)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	err = passwordMatches(creds.Password, user.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	token, err := c.service.createToken(creds.Username)
	if err == nil {
		ctx.JSON(http.StatusOK, gin.H{"token": token})
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func (c controller) middleware(ctx *gin.Context) {
	if ctx.FullPath() == "" {
		return
	}

	claims, err := c.service.parseClaimsFromAuthHeader(ctx.GetHeader("Authorization"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("could not validate user : %s", err.Error())})
		return
	}

	exists, err := c.service.userExists(claims.Username)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to query database for existing users: %s", err.Error())})
		return
	}

	if !exists {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "token is for an unregistered user",
		})
		return
	}

	ctx.Set(UserNameConetxtKey, claims.Username)
}
