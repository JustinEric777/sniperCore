package uploader

import (
	"github.com/gin-gonic/gin"
	"github.com/sniperCore/app/uploader/http/middleware"
	"github.com/sniperCore/app/uploader/route"
)

func Register(engine *gin.Engine) {
	// load middleware
	engine.Use(middleware.LogToFile())

	// init route
	route.InitRoute(engine)
}
