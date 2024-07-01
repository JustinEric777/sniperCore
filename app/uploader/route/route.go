package route

import (
	"github.com/gin-gonic/gin"
	"github.com/sniperCore/app/uploader/http/handler"
)

func InitRoute(engine *gin.Engine) {
	app := engine.Group("/api/picture")
	{
		app.POST("/process", handler.PictureProcess)
	}
}
