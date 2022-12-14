package img

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/minio/minio-go/v7"
	"github.com/noetarbouriech/storiesque/backend/internal/db"
	"github.com/noetarbouriech/storiesque/backend/internal/utils"
)

type Service struct {
	queries *db.Queries
	minio   *Minio
}

func NewService(queries *db.Queries, minio *Minio) *Service {
	return &Service{
		queries: queries,
		minio:   minio,
	}
}

func (s *Service) UserRoutes(r chi.Router) {
	r.Post("/image/upload", s.UploadImage)
	r.Delete("/image/{type}/{id}", s.DeleteImage)
}

func (s *Service) UploadImage(w http.ResponseWriter, r *http.Request) {

	// parse the multipart form from the request
	err := r.ParseMultipartForm(8 << 20) // max form size as 8MiB
	if err != nil {
		utils.Response(w, r, 400, "issue with form")
		return
	}

	// get the file in the multipart form
	file, _, err := r.FormFile("file")
	if err != nil {
		utils.Response(w, r, 400, "file not found in request")
		return
	}
	defer file.Close()

	// check the type of resource
	resType := r.FormValue("type")
	if resType != "user" && resType != "page" && resType != "story" {
		utils.Response(w, r, 400, "unknown resource type")
		return
	}

	// convert id to int64
	id, err := strconv.ParseInt(r.FormValue("id"), 10, 64)
	if err != nil {
		utils.Response(w, r, 400, "id bad format")
		return
	}

	// get owner of the resource
	owner, err := s.getResource(resType, id)
	if err != nil {
		utils.Response(w, r, 500, "internal error")
		return
	}

	// check if user is authorized
	if !utils.IsOwner(r, int(owner)) {
		utils.Response(w, r, 401, "user is not the owner of the given resource")
		return
	}

	// create a buffer to store the file content
	var buf bytes.Buffer
	_, err = io.Copy(&buf, file)
	if err != nil {
		utils.Response(w, r, 500, "internal error")
		return
	}

	// check the file size
	size := int64(buf.Len())
	if size > 2*1000*1000 { // 2 MB
		utils.Response(w, r, 400, "file too big")
		return
	}

	// check file type
	filetype := http.DetectContentType(buf.Bytes())
	if filetype != "image/png" && filetype != "image/jpeg" {
		utils.Response(w, r, 400, "file type is not accepted")
		return
	}

	// create the filename
	filename := fmt.Sprintf("%s/%s.png", r.FormValue("type"), r.FormValue("id"))

	// upload file as object in s3
	info, err := s.minio.S3.PutObject(context.Background(), "storiesque", filename, &buf, size, minio.PutObjectOptions{ContentType: "image/png"})
	if err != nil {
		utils.Response(w, r, 500, "internal error")
		return
	}

	err = s.setImgOnDB(resType, id, true)
	if err != nil {
		utils.Response(w, r, 500, "internal error")
		return
	}

	// response with image uri
	render.Status(r, 201)
	render.JSON(w, r, map[string]string{"uri": info.Key})
}

func (s *Service) DeleteImage(w http.ResponseWriter, r *http.Request) {

	// get type from uri
	resType := chi.URLParam(r, "type")
	if resType != "user" && resType != "page" && resType != "story" {
		utils.Response(w, r, 400, "unknown resource type")
		return
	}

	// get id in param
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.Response(w, r, 400, "id bad format")
		return
	}

	// get owner of the resource
	owner, err := s.getResource(resType, id)
	if err != nil {
		utils.Response(w, r, 500, "internal error")
		return
	}

	// check if user is authorized
	if !utils.IsOwner(r, int(owner)) {
		utils.Response(w, r, 401, "user is not the owner of the given resource")
		return
	}

	// create the filename
	filename := fmt.Sprintf("%s/%s.png", r.FormValue("type"), r.FormValue("id"))

	// upload file as object in s3
	err = s.minio.S3.RemoveObject(context.Background(), "storiesque", filename, minio.RemoveObjectOptions{})
	if err != nil {
		utils.Response(w, r, 500, "internal error")
		return
	}

	err = s.setImgOnDB(resType, id, false)
	if err != nil {
		utils.Response(w, r, 500, "internal error")
		return
	}

	utils.Response(w, r, 200, "image deleted")
}

// get the owner of DB resource of an object of type resType and of id id
func (s *Service) getResource(resType string, id int64) (int, error) {

	switch resType {
	case "user":
		res, err := s.queries.GetUserWithId(context.Background(), id)
		return int(res.ID), err
	case "story":
		res, err := s.queries.GetStory(context.Background(), id)
		return int(res.Author), err
	case "page":
		res, err := s.queries.GetPage(context.Background(), id)
		return int(res.Author), err
	default:
		return 0, errors.New("unknown resource type")
	}
}

// change img indicator on db
func (s *Service) setImgOnDB(resType string, id int64, has_img bool) error {
	var err error

	switch resType {
	case "user":
		err = s.queries.SetImgUser(context.Background(), db.SetImgUserParams{
			ID:     id,
			HasImg: has_img,
		})
	case "story":
		err = s.queries.SetImgStory(context.Background(), db.SetImgStoryParams{
			ID:     id,
			HasImg: has_img,
		})
	case "page":
		err = s.queries.SetImgPage(context.Background(), db.SetImgPageParams{
			ID:     id,
			HasImg: has_img,
		})
	default:
		err = errors.New("unknown resource type")
	}

	return err
}
