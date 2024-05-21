package route

import "github.com/gin-gonic/gin"

type Router struct {
	*gin.Engine
}

func NewRouter(engine *gin.Engine) *Router {
	return &Router{engine}
}
