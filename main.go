package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/olahol/go-imageupload"
)

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.File("index.html")
	})

	r.POST("/upload", func(c *gin.Context) {
		img, err := imageupload.Process(c.Request, "file")

		if err != nil {
			panic(err)
		}

		err = img.Save("upload/" + img.Filename)

		if err != nil {
			panic(err)
		}

		c.JSON(200, gin.H{
			"message": "uploadSuccess",
		})
	})
	r.Run(":" + os.Getenv("PORT"))
}
