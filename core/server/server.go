package server

import (
	"bytes"
	"context"
	"github.com/sniperCore/core/config"
	"github.com/sniperCore/core/consul"
	"github.com/sniperCore/core/helper"
	"github.com/sniperCore/core/log"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	serviceConfig *ServiceConfig
	Engine        *gin.Engine
	Server        *http.Server
}

func Start() *HttpServer {
	config, _ := InitConfig()

	// debug setting
	if !config.Http.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	// init gin
	engine := gin.New()

	// set default middleware
	buffer := new(bytes.Buffer)
	handleRecovery := func(ctx *gin.Context, err interface{}) {
		log.Error(err)
		ctx.AbortWithStatus(http.StatusBadRequest)
	}
	engine.Use(gin.CustomRecoveryWithWriter(buffer, handleRecovery))

	// set static file dir
	if !isEmptyStaticDir() {
		initStatic(engine)
	}

	// set resource view dir
	if !isEmptyResourceDir() {
		initResource(engine)
		engine.NoRoute(func(c *gin.Context) {
			c.HTML(404, "404.html", gin.H{})
		})
	} else {
		engine.NoRoute(func(c *gin.Context) {
			c.JSON(404, gin.H{})
		})
	}

	return &HttpServer{
		Engine:        engine,
		serviceConfig: config,
	}
}

func isEmptyStaticDir() bool {
	path := config.GetAppPath()
	dir, _ := ioutil.ReadDir(path + "/resources/static")
	if len(dir) == 0 {
		return true
	}

	return false
}

func initStatic(engine *gin.Engine) {
	dir := config.GetAppPath()
	engine.StaticFS("static", http.Dir(dir+"/resources/static"))
}

func isEmptyResourceDir() bool {
	path := config.GetAppPath()
	dir, _ := ioutil.ReadDir(path + "/resources/views/")
	if len(dir) == 0 {
		return true
	}

	return false
}

func initResource(engine *gin.Engine) {
	dir := config.GetAppPath()
	engine.LoadHTMLGlob(dir + "/resources/views/**/*")
}

func (hp *HttpServer) Run() {
	server := &http.Server{
		Addr:         hp.serviceConfig.Http.Addr + ":" + hp.serviceConfig.Http.Port,
		Handler:      hp.Engine,
		IdleTimeout:  hp.serviceConfig.Http.IdleTimeout * time.Second,
		ReadTimeout:  hp.serviceConfig.Http.ReadTimeOut * time.Second,
		WriteTimeout: hp.serviceConfig.Http.WriteTimeOut * time.Second,
	}
	hp.Server = server

	//server connect
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("server connect error, err =", err)
			panic(err)
		}
	}()

	if hp.serviceConfig.Consul.Enabled {
		consul.GetConsul("consul").Register()
	}

	helper.Println("server init success, server =" + hp.serviceConfig.Http.Addr + ":" + hp.serviceConfig.Http.Port)
}

func (hp *HttpServer) ShutDown() {
	//send stop signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// server quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := hp.Server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	if hp.serviceConfig.Consul.Enabled {
		consul.GetConsul("consul").DisRegister()
	}

	helper.Println("server quit success, stop server ...")
}
