package routing

import "github.com/gin-gonic/gin"

type GET interface {
	GET(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes
}
