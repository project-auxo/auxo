package router

import (
	"encoding/gob"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"github.com/project-auxo/auxo/hestia/controller/callback"
	"github.com/project-auxo/auxo/hestia/controller/login"
	"github.com/project-auxo/auxo/hestia/controller/logout"
	olympusCtrl "github.com/project-auxo/auxo/hestia/controller/olympus"
	"github.com/project-auxo/auxo/hestia/controller/user"
	hestiaCfg "github.com/project-auxo/auxo/hestia/internal/config"
	authenticator "github.com/project-auxo/auxo/hestia/platform/auth"
	"github.com/project-auxo/auxo/hestia/platform/middleware"
)

const (
	staticFilePath   = "hestia/web/static"
	templateFilePath = "hestia/web/template/*"
)

func New(auth *authenticator.Authenticator, cfg *hestiaCfg.Config) *gin.Engine {
	r := gin.Default()
	gob.Register(map[string]interface{}{})
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("auth-session", store))

	r.Static("public/", staticFilePath)
	r.LoadHTMLGlob(templateFilePath)

	r.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "home.html", nil)
	})
	r.GET("/login", login.Handler(auth))
	r.GET("/callback", callback.Handler(auth))
	r.GET("/user", middleware.IsAuthenticated, user.Handler)
	r.GET("/logout", logout.Handler)

	// APIs
	v1 := r.Group("/v1")
	{
		v1.GET("/", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})
	}
	addOlympus(v1, cfg)

	return r
}

func addOlympus(rg *gin.RouterGroup, cfg *hestiaCfg.Config) {
	client := olympusCtrl.GetClient(cfg)
	olympus := rg.Group("/olympus")
	{
		olympus.GET("/agents/num", olympusCtrl.GetNumberOfAgents(client))
	}
}
