package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"./../models"
	"fmt"
	"strconv"
)

type StickersRouter struct {
	model *models.StickerModel
}

var stickersRouter *StickersRouter

const (
	limitDefaultValue int = 12
	limitMaxValue int = 36
)

func InitStickersRouter(r *gin.Engine) {
	router := &StickersRouter{
		model: &models.StickerModel{},
	}
	routerGroup := r.Group("/stickers")
	{
		routerGroup.GET("/list", router.GetStickersList)
		routerGroup.POST("/add", router.PostAddSticker)
	}
	r.GET("/sticker/:filename", router.GetStickerByFilename)
}

//
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
	fmt.Println("stickers", stickers)

	c.JSON(http.StatusOK, stickers)
}

//
func (router *StickersRouter) GetStickerByFilename(c *gin.Context) {
	filename := c.Param("filename")
	sticker, _ := router.model.GetStickerByFilename(filename)
	c.Data(http.StatusOK, "image/png", sticker)
}

//
func (router *StickersRouter) PostAddSticker(c *gin.Context) {
	stickerImage, _, formFileError := c.Request.FormFile("sticker")
	stickerTitle := c.Request.FormValue("title")
	if formFileError != nil {
		fmt.Println("PostAddSticker Error", formFileError)
		return
	}
	fmt.Println("router", router)
	router.model.AddSticker(stickerTitle, stickerImage)
	c.Redirect(http.StatusSeeOther, "/stickers/")
}