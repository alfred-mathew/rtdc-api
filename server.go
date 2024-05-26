package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"signzy.com/rtdc-api/auth"
	"signzy.com/rtdc-api/health"
)

func createRouter(client *mongo.Client, signingKey []byte) *gin.Engine {
	router := gin.Default()

	router.SetTrustedProxies(nil)

	health.Register(router)
	auth.Register(router, client, signingKey)
	router.GET("/test", func(ctx *gin.Context) {
		ctx.JSON(http.StatusAccepted, gin.H{"message": "ok"})
	})

	return router
}

func startServer(addr string, client *mongo.Client, signingKey []byte) *http.Server {
	server := &http.Server{
		Addr:    addr,
		Handler: createRouter(client, signingKey).Handler(),
	}

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			fmt.Printf("HTTP server stopped with error : %s\n", err)
		}
	}()

	return server
}

func shutdownServer(ctx context.Context, server *http.Server, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return err
	}

	fmt.Println("HTTP server has shutdown")

	return nil
}
