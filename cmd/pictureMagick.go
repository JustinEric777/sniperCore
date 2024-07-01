package main

import (
	"github.com/sniperCore/app/uploader"
	"github.com/sniperCore/bootstrap"
	"github.com/sniperCore/core/server"
)

func main() {
	//init
	err := bootstrap.Bootstrap()
	if err != nil {
		panic(err)
	}
	// init server
	picMagick := server.Start()
	uploader.Register(picMagick.Engine)
	picMagick.Run()
	picMagick.ShutDown()
}
