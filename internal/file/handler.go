package file

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	fileURL         = "/file"
	filesURL        = "/files"
	MAX_UPLOAD_SIZE = 3000 << 20 //3GB
)

type filesHandler struct {
	log         *logrus.Entry
	fileService *service
}

func NewFilesHandler(fileService *service, log *logrus.Entry) *filesHandler {
	return &filesHandler{
		log:         log,
		fileService: fileService,
	}
}

func (h *filesHandler) Register(router *gin.Engine) {
	v1_API := router.Group("/v1")
	{
		v1_API.GET(fileURL, h.GetFile)
		v1_API.DELETE(fileURL, h.DeleteFile)
		v1_API.POST(fileURL, h.CreateFile)
		v1_API.GET(filesURL, h.GetFilesByBucketName)
	}
}

// Получение файла из bucket по id
// @Summary      Get file
// @Tags file
// @Description  Get file from bucket by id
// @Produce      text/plain
// @Param        id query string true "file id"
// @Param        bucket query string true "bucket name"
// @Success      200  {object}  string "Success get file"
// @Failure		 400  {object}  messageResponse "Invalid parameters"
// @Failure		 500  {object}  messageResponse "Server error"
// @Router       /v1/file [get]
func (h *filesHandler) GetFile(c *gin.Context) {
	bucket := c.Query("bucket")
	if bucket == "" {
		NewErrorResponse(c, http.StatusBadRequest, "Invalid parameters: missing - 'bucket'")
		return
	}

	id := c.Query("id")
	if id == "" {
		NewErrorResponse(c, http.StatusBadRequest, "Invalid parameters: missing - 'id'")
		return
	}

	h.log.Debugf("Get file from %s bucket with %s (id)", bucket, id)

	f, err := h.fileService.GetFile(c.Request.Context(), bucket, id)
	if err != nil {
		h.log.Errorf("Error get file from %s bucket with %s (id) with err: %v", bucket, id, err)
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", f.Name))
	c.Header("Content-Type", c.Request.Header.Get("Content-Type"))
	c.Writer.Write(f.Bytes)
}

// Получение файлов из bucket
// @Summary      Get files
// @Tags file
// @Description  Get files from bucket
// @Produce      json
// @Param        bucket query string true "bucket name"
// @Success      200  {object}  file.File[] "Success get files"
// @Failure		 400  {object}  messageResponse "Invalid parameters"
// @Failure		 500  {object}  messageResponse "Server error"
// @Router       /v1/files [get]
func (h *filesHandler) GetFilesByBucketName(c *gin.Context) {
	h.log.Debug("Get file by bucket name")

	bucket := c.Query("bucket")
	if bucket == "" {
		NewErrorResponse(c, http.StatusBadRequest, "Invalid parameters: missing - 'bucket'")
		return
	}

	h.log.Debugf("Get files from %s bucket", bucket)

	files, err := h.fileService.GetFilesByBucketName(c.Request.Context(), bucket)
	if err != nil {
		h.log.Errorf("Error get files from %s bucket with err: %v", bucket, err)
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(200, files)
}

// Создание файла в bucket
// @Summary      Create file
// @Tags file
// @Description  Create files in bucket
// @Accept		 multipart/form-data
// @Produce      json
// @Param        bucket query string true "bucket name"
// @Success      201  {object}  newIdResponse "Success create file"
// @Failure		 400  {object}  messageResponse "Invalid parameters or multipart data"
// @Failure		 500  {object}  messageResponse "Server error"
// @Router       /v1/file [post]
func (h *filesHandler) CreateFile(c *gin.Context) {
	bucket := c.Query("bucket")
	if bucket == "" {
		NewErrorResponse(c, http.StatusBadRequest, "Invalid parameters: missing - 'bucket'")
	}

	h.log.Debugf("Create files in %s bucket", bucket)

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MAX_UPLOAD_SIZE)

	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		h.log.Errorf("Bad multipart parsing with err: %v", err)
		NewErrorResponse(c, http.StatusBadRequest, "Bad multipart")
		return
	}

	defer file.Close()

	dto := CreateFileDTO{
		Name:   fileHeader.Filename,
		Size:   fileHeader.Size,
		Reader: file,
	}

	id, err := h.fileService.Create(c.Request.Context(), bucket, dto)
	if err != nil {
		h.log.Errorf("Error create file in %s bucket with err: %v", bucket, err)
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(201, &newIdResponse{
		Id: id,
	})
}

// Удаление файла в bucket по id
// @Summary      Delete file
// @Tags file
// @Description  Delete files in bucket
// @Produce      json
// @Param        id query string true "file id"
// @Param        bucket query string true "bucket name"
// @Success      200  {object}  string "Success create file"
// @Failure		 400  {object}  messageResponse "Invalid parameters"
// @Failure		 500  {object}  messageResponse "Server error"
// @Router       /v1/file [delete]
func (h *filesHandler) DeleteFile(c *gin.Context) {
	h.log.Debug("Delete file")

	bucket := c.Query("bucket")
	if bucket == "" {
		NewErrorResponse(c, http.StatusBadRequest, "Invalid parameters: missing - 'bucket'")
		return
	}

	id := c.Query("id")
	if id == "" {
		NewErrorResponse(c, http.StatusBadRequest, "Invalid parameters: missing - 'id'")
		return
	}

	h.log.Debugf("Delete file from %s bucket with %s (id)", bucket, id)

	err := h.fileService.Delete(c.Request.Context(), bucket, id)
	if err != nil {
		h.log.Errorf("Error delete file in %s bucket %s with (id) with err: %v", bucket, id, err)
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(200, "Success")
}
