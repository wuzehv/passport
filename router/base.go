package router

import (
	"github.com/gin-gonic/gin"
	"github.com/wuzehv/passport/middleware"
	"github.com/wuzehv/passport/util/config"
	"log"
	"path/filepath"
)

func InitRouter() *gin.Engine {
	gin.SetMode(config.App.RunMode)

	router := gin.New()

	router.Use(middleware.Log())

	router.Use(gin.Recovery())

	router.Delims("{[{", "}]}")

	router.LoadHTMLFiles(loadTemplates("template")...)

	router.Static("/static", "./static")

	construct(router)

	return router
}

func loadTemplates(templatesDir string) []string {
	other, err := filepath.Glob(templatesDir + "/**/*.html")
	if err != nil {
		log.Fatalf("load template error: %v\n", err)
	}

	admin, err := filepath.Glob(templatesDir + "/**/**/*.html")
	if err != nil {
		log.Fatalf("load template error: %v\n", err)
	}

	return append(admin, other...)
}
