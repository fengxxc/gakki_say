package api

import (
	"embed"
	"log"
	"net/http"
	"os"

	"github.com/fengxxc/gakki_say/policy"
	"github.com/gin-gonic/gin"
)

func Router(r *gin.Engine, imgDir embed.FS, fontDir embed.FS) {

	r.GET("/gakki_say/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"ok":      true,
			"message": "pang~",
		})
	})

	r.GET("/gakki_say/thumb/:image", func(c *gin.Context) {
		imageName := c.Param("image")
		baseDir, _ := os.Getwd()
		thumbImg := baseDir + "/.cache/thumb/" + imageName
		log.Printf("thumb image: %s", thumbImg)
		f, err := os.ReadFile(thumbImg)
		if err != nil {
			log.Printf("Error loading thumb image: %v", err)
			return
		}
		c.Header("Content-Type", "image/jpeg")
		c.Writer.Write(f)
	})

	r.GET("/gakki_say/image/:image", func(c *gin.Context) {
		imageName := c.Param("image")
		text := c.Query("text")
		c.Header("Content-Type", "image/jpeg")
		img, err := policy.ImgWriteText("img/"+imageName, text, policy.DrawStringConfig{
			Ax:          0.5,
			Ay:          0.5,
			FontFamily:  "SIMYOU.TTF",
			TextBgColor: &policy.RGBA{R: 89, G: 89, B: 89, A: 64},
		}, imgDir, fontDir)
		if err != nil {
			log.Println(err)
			return
		}
		c.Writer.Write(policy.ImgToBytes(img, policy.Jpeg))
	})
}
