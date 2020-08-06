package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	. "search/common"
	"search/sonic"
	"time"
)

//搜索内容
func search(c *gin.Context) {
	//search?key=xxxx
	query := c.Query("key")

	if query == "" {
		c.JSON(200, gin.H{
			//0 请输入关键字
			//1 找到数据
			"status":  "0",
			"message": "请输入查询关键字",
		})
		return
	}

	posts := sonic.Search(query)

	c.JSON(200, gin.H{
		"status":  "1",
		"message": posts,
	})
}

func main() {

	// Disable Console Color, you don't need console color when writing the logs to file.
	gin.DisableConsoleColor()

	// Logging to a file.
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f)

	// Use the following code if you need to write the logs to file and console at the same time.
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	router := gin.New()

	// LoggerWithFormatter middleware will write the logs to gin.DefaultWriter
	// By default gin.DefaultWriter = os.Stdout
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {

		// your custom format
		return fmt.Sprintf("%s - [%s] \"%s * %s * %s * %d * %s * \"%s\" * %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	router.Use(gin.Recovery())

	router.GET("/search", search)

	router.Run(":" + Conf.SearchPort)

}
