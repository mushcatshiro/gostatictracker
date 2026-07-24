package images

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/mushcatshiro/gostatictracker/models"
)

type ImagesManager struct {
	Root      string
	SizeLimit int
	DB        *sql.DB
}

func NewImagesManager(r string, sizeLimit int, db *sql.DB) ImagesManager {
	return ImagesManager{
		Root:      r,
		SizeLimit: sizeLimit,
		DB:        db,
	}
}

func (i *ImagesManager) UploadImage(im []byte) (uuid.UUID, error) {
	var rv uuid.UUID
	// check size
	if len(im) > i.SizeLimit {
		return rv, errors.New("exceeded upload limit")
	}
	// save to db -> get UUID
	// first chart -> write to persistent
	return rv, nil
}

func (i *ImagesManager) ListImages(start, offset int) ([]models.ImageRecord, error) {
	var ims []models.ImageRecord
	// return list of images
	return ims, nil
}

func (i *ImagesManager) SearchImage() {
	// partial string, time, later on meta?
}
