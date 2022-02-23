package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
)

func GetFile(ctx *gin.Context){
	file, err := ctx.FormFile("file")
	if  err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	fmt.Println(file.Filename)
	dst := fmt.Sprintf("%s%s", viper.GetString("img.path") , file.Filename)
	// 上传文件到指定的目录
	ctx.SaveUploadedFile(file, dst)
	ctx.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("'%s' uploaded!", file.Filename),
	})
}