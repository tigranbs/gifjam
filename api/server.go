package api

import (
	"gopkg.in/kataras/iris.v6"
	"gopkg.in/kataras/iris.v6/adaptors/httprouter"
	"os"
	"gifjam/gifGrabber"
	"net/http"
	"log"
	"gopkg.in/kataras/iris.v6/adaptors/cors"
	"strconv"
	"io"
)

type obj map[string]interface{}

func StartApiServer() {
	app := iris.New()
	app.Adapt(iris.DevLogger(), httprouter.New(), cors.New(cors.Options{AllowedOrigins: []string{"*"}}))

	app.Post("/gifs", GetGifs)
	app.Post("/gif/:id", serveGif)

	// Setting address from ENV
	app.Listen(os.Getenv("GIFJAM_SERVER_HOST"))
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
		ctx.JSON(http.StatusInternalServerError, obj{"error": "Unable to read Image file"})
		return
	}

	defer r.Close()

	ctx.SetHeader("Content-Length", strconv.FormatInt(length, 10))
	ctx.SetHeader("Content-Type", "image/gif")
	io.Copy(ctx, r)
}