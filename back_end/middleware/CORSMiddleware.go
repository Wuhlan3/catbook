package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
)

func CORSMidddleware() gin.HandlerFunc{
	return func(ctx *gin.Context){
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", viper.GetString("front_end.client"))
		ctx.Writer.Header().Set("Access-Control-Max-Age", "86400") //缓存时间
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "*") //GET,POST等等
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if ctx.Request.Method == http.MethodOptions {
			ctx.AbortWithStatus(200)
		}else {
			ctx.Next()
		}

	}
}