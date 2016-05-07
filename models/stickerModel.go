package models

import (
	"gopkg.in/mgo.v2/bson"
	"io"
	"mime/multipart"
	"./../helpers"
	"io/ioutil"
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

func (stickerModel *StickerModel) AddSticker(title string, image multipart.File) (err error) {
	stickersDB := helpers.Session.DB("stickers")
	stickersImagesGridFS := stickersDB.GridFS("stickerImages")
	imageBsonId := bson.NewObjectId()
	imageGridFSFile, err := stickersImagesGridFS.Create(imageBsonId.Hex())
	if helpers.CheckError("PostAddSticker create error", err) {
		return
	}
	imageGridFSFile.SetId(imageBsonId)
	defer image.Close()
	_, err = io.Copy(imageGridFSFile, image)
	if helpers.CheckError("PostAddSticker copy error", err) {
		return
	}

	imageGridFSFile.SetContentType("image/png")
	imageGridFSFile.SetMeta(&StickerImageMeta{
		Title: title,
		Author: "lica",
	})

	err = imageGridFSFile.Close()
	helpers.CheckError("PostAddSticker close error", err)
	return
}

func (stickerModel *StickerModel) GetStickerByFilename(filename string) (stickerBytes []byte, stickerMeta StickerImageMeta) {
	gridFile, err := helpers.Session.DB("stickers").GridFS("stickerImages").Open(filename)
	if helpers.CheckError("GetStickerByFilename open error", err) {
		return
	}
	defer gridFile.Close()

	stickerBytes, err = ioutil.ReadAll(gridFile)
	if helpers.CheckError("GetStickerByFilename readAll error", err) {
		return
	}

	err = gridFile.GetMeta(&stickerMeta)
	if helpers.CheckError("GetStickerByFilename getMeta error", err) {
		return
	}

	return
}

func (stickerModel *StickerModel) GetStickersList(skip, limit int) ([]StickerFiles, error) {
	var stickerFiles []StickerFiles
	err := helpers.Session.DB("stickers").GridFS("stickerImages").Find(nil).Sort("-uploadDate").Skip(skip).Limit(limit).All(&stickerFiles)
	if helpers.CheckError("GetSticker all error", err) {
		return nil, err
	}
	return stickerFiles, err
}