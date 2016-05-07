package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"./../models"
	"./../helpers"
	"fmt"
	"strconv"
)

type StickersRouter struct {
	model *models.StickerModel
}

var stickersRouter *StickersRouter

func InitStickersRouter(r *gin.Engine) {
	router := &StickersRouter{
		model: &models.StickerModel{},
	}
	routerGroup := r.Group("/stickers")
	{
		routerGroup.GET("/getStickersList", router.GetStickersList)
		routerGroup.POST("/addSticker", stickersRouter.PostAddSticker)
		routerGroup.GET("/img/:filename", stickersRouter.GetSticker)
	}
}

func (router *StickersRouter) GetSticker(c *gin.Context) {
	filename := c.Param("filename")
	sticker, _ := router.model.GetStickerByFilename(filename)
	c.Data(http.StatusOK, "image/png", sticker)
}

func (router *StickersRouter) PostAddSticker(c *gin.Context) {
	stickerImage, _, formFileError := c.Request.FormFile("sticker")
	stickerTitle := c.Request.FormValue("title")
	if helpers.CheckError("PostAddSticker Error", formFileError) {
		return
	}
	router.model.AddSticker(stickerTitle, stickerImage)
	c.Redirect(http.StatusSeeOther, "/stickers/")
}


const (
	limitDefaultValue int = 12
	limitMaxValue int = 36
)


func (router *StickersRouter) GetStickersList(c *gin.Context) {
	skip := 0
	limit := limitDefaultValue
	var err error

	skipFormValue := c.Request.FormValue("skip")
	if len(skipFormValue) > 0 {
		skip, err = strconv.Atoi(skipFormValue)
		if err != nil {
			c.JSON(http.StatusBadRequest, "incorrect skip parameter")
			return
		}
	}
	limitFormValue := c.Request.FormValue("limit")
	if len(limitFormValue) > 0 {
		limit, err = strconv.Atoi(limitFormValue)
		if err != nil {
			c.JSON(http.StatusBadRequest, "incorrect limit parameter")
			return
		}
		if limit > limitMaxValue {
			c.JSON(http.StatusBadRequest, "limit parameter is too big")
		}
	}

	stickers, err := router.model.GetStickersList(skip, limit)
	if err != nil {
		fmt.Println("GetStickersList error \n\t", err.Error())
		return
	}

	c.JSON(http.StatusOK, stickers)
}