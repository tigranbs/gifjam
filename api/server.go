package api

import (
	"gifjam/gifGrabber"
	"gopkg.in/kataras/iris.v6"
	"gopkg.in/kataras/iris.v6/adaptors/cors"
	"gopkg.in/kataras/iris.v6/adaptors/httprouter"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"gifjam/config"
)

type obj map[string]interface{}

func StartApiServer() {
	app := iris.New()
	app.Adapt(iris.DevLogger(), httprouter.New(), cors.New(cors.Options{AllowedOrigins: []string{"*"}}))

	app.Get("/gif/:id", serveGif)
	app.Post("/gifs/visibility/:id/:visible", gifVisibility)
	app.Post("/gifs", GetGifs)

	host := os.Getenv("GIFJAM_SERVER_HOST")
	if len(config.GlobalConfig.Host) > 0 {
		host = config.GlobalConfig.Host
	}

	// Setting address from ENV
	app.Listen(host)
}

func GetGifs(ctx *iris.Context) {
	offset, err := ctx.URLParamInt("offset")
	if err != nil {
		offset = 0
	}

	limit, err := ctx.URLParamInt("offset")
	if err != nil {
		limit = 10
	}

	images, err := gifGrabber.GetItems(offset, limit)
	if err != nil {
		log.Println("Unable to get images from database -> ", err.Error())
		ctx.JSON(http.StatusOK, obj{"error": "Unable to get images!"})
		return
	}

	ctx.JSON(http.StatusOK, obj{"images": images})
}

func serveGif(ctx *iris.Context) {
	file_id := ctx.Param("id")
	if len(file_id) != 24 {
		ctx.JSON(http.StatusNotFound, obj{"error": "Image Not Found!"})
		return
	}

	length, r, err := gifGrabber.GetFileIO(file_id)
	if err != nil {
		log.Println("Unable to get image from database -> ", err.Error())
		ctx.JSON(http.StatusInternalServerError, obj{"error": "Unable to read Image file"})
		return
	}

	defer r.Close()

	ctx.SetHeader("Content-Length", strconv.FormatInt(length, 10))
	ctx.SetHeader("Content-Type", "image/gif")
	io.Copy(ctx, r)
}

func gifVisibility(ctx *iris.Context) {
	file_id := ctx.Param("id")
	if len(file_id) != 24 {
		ctx.JSON(http.StatusNotFound, obj{"error": "Image Not Found!"})
		return
	}

	visible := false

	visible_param := ctx.Param("visible")
	if visible_param == "1" {
		visible = true
	}

	err := gifGrabber.SetVisibility(file_id, visible)
	if err != nil {
		log.Println("Unable to set image to visible -> ", err.Error())
		ctx.JSON(http.StatusInternalServerError, obj{"error": "Unable to set image to visible"})
		return
	}

	ctx.JSON(http.StatusOK, obj{})
}
