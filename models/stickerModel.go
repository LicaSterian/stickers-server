package models

import (
	"gopkg.in/mgo.v2/bson"
	"io"
	"mime/multipart"
	"./../helpers"
	"io/ioutil"
	"fmt"
)

type StickerModel struct {

}

type StickerImage struct {
	Data []byte
}

type StickerFiles struct {
	Filename string `bson:"filename"`
	Metadata StickerImageMeta `bson:"metadata"`
}

type StickerImageMeta struct {
	Title string `bson:"title"`
	Author string `bson:"author"`
}

func (stickerModel *StickerModel) GetStickersList(skip, limit int) ([]StickerFiles, error) {
	var stickerFiles []StickerFiles
	err := helpers.Session.DB("stickers").GridFS("stickerImages").Find(nil).Sort("-uploadDate").Skip(skip).Limit(limit).All(&stickerFiles)
	if err != nil {
		fmt.Println("GetSticker all error", err)
		return nil, err
	}
	return stickerFiles, err
}

//TODO remove fileName, use _id
func (stickerModel *StickerModel) AddSticker(title string, image multipart.File) (err error) {
	stickersDB := helpers.Session.DB("stickers")
	stickersImagesGridFS := stickersDB.GridFS("stickerImages")
	imageBsonId := bson.NewObjectId()
	imageGridFSFile, err := stickersImagesGridFS.Create(imageBsonId.Hex())
	if err != nil {
		fmt.Println("PostAddSticker create error")
		return
	}
	imageGridFSFile.SetId(imageBsonId)
	defer image.Close()
	_, err = io.Copy(imageGridFSFile, image)
	if err != nil {
		fmt.Println("PostAddSticker copy error")
		return
	}

	imageGridFSFile.SetContentType("image/png")
	imageGridFSFile.SetMeta(&StickerImageMeta{
		Title: title,
		Author: "lica",
	})

	err = imageGridFSFile.Close()
	if err != nil {
		fmt.Println("PostAddSticker close error", err)
	}
	return
}

func (stickerModel *StickerModel) GetStickerByFilename(filename string) (stickerBytes []byte, stickerMeta StickerImageMeta) {
	gridFile, err := helpers.Session.DB("stickers").GridFS("stickerImages").Open(filename)
	if err != nil {
		fmt.Println("GetStickerByFilename open error", err)
		return
	}

	stickerBytes, err = ioutil.ReadAll(gridFile)
	if err != nil {
		fmt.Println("GetStickerByFilename readAll error", err)
		return
	}

	err = gridFile.GetMeta(&stickerMeta)
	defer gridFile.Close()
	if err != nil {
		fmt.Println("GetStickerByFilename getMeta error", err)
		return
	}
	return
}