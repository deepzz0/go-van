// Package main provides ...
package main

import (
	"time"

	"github.com/deepzz0/go-van"
	"github.com/deepzz0/go-van/pkg/logx"
	"github.com/deepzz0/go-van/pkg/registry"
	"github.com/deepzz0/go-van/pkg/registry/etcd"
	"github.com/deepzz0/go-van/pkg/server"
	"github.com/deepzz0/go-van/pkg/server/http"

	"github.com/gin-gonic/gin"
)

func main() {
	reg := etcd.NewRegistry(
		registry.WithTTL(time.Second*10),
		registry.WithAddress("localhost:2379"),
	)
	// gin server
	e := gin.New()
	e.GET("/hello", func(c *gin.Context) {
		c.String(200, "hello world")
	})

	srv := http.NewServer(
		server.WithRecover(true),
		server.WithEndpoint(":9000"),
		http.WithHandler(e),
	)
	service := van.NewService(
		van.WithName("gin"),
		van.WithServer(srv),
		van.WithRegistry(reg),
	)
	if err := service.Run(); err != nil {
		logx.Fatal(err)
	}
}
