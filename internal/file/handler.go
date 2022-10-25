package file

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mserebryaakov/file-service/pkg/logger"
)

const (
	fileURL         = "/file"
	filesURL        = "/files"
	MAX_UPLOAD_SIZE = 3000 << 20 //3GB
)

type filesHandler struct {
	log         *logger.Logger
	fileService *service
}

func NewFilesHandler(log *logger.Logger, fileService *service) *filesHandler {
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

func (h *filesHandler) GetFile(c *gin.Context) {
	h.log.Debug("Get file")

	bucket := c.Query("bucket")
	if bucket == "" {
		NewErrorResponse(c, http.StatusBadRequest, "Invalid parameters: missing - 'bucket'")
	}

	id := c.Query("id")
	if id == "" {
		NewErrorResponse(c, http.StatusBadRequest, "Invalid parameters: missing - 'id'")
	}

	f, err := h.fileService.GetFile(c.Request.Context(), bucket, id)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", f.Name))
	c.Header("Content-Type", c.Request.Header.Get("Content-Type"))
	c.Writer.Write(f.Bytes)
}

func (h *filesHandler) GetFilesByBucketName(c *gin.Context) {
	h.log.Debug("Get file by bucket name")

	bucket := c.Query("bucket")
	if bucket == "" {
		NewErrorResponse(c, http.StatusBadRequest, "Invalid parameters: missing - 'bucket'")
	}

	file, err := h.fileService.GetFilesByBucketName(c.Request.Context(), bucket)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(200, file)
}

func (h *filesHandler) CreateFile(c *gin.Context) {
	h.log.Debug("Create file")

	bucket := c.Query("bucket")
	if bucket == "" {
		NewErrorResponse(c, http.StatusBadRequest, "Invalid parameters: missing - 'bucket'")
	}

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MAX_UPLOAD_SIZE)

	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "Bad multipart")
		return
	}

	defer file.Close()

	dto := CreateFileDTO{
		Name:   fileHeader.Filename,
		Size:   fileHeader.Size,
		Reader: file,
	}

	err = h.fileService.Create(c.Request.Context(), bucket, dto)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(201, "Success")
}

func (h *filesHandler) DeleteFile(c *gin.Context) {
	h.log.Debug("Delete file")

	bucket := c.Query("bucket")
	if bucket == "" {
		NewErrorResponse(c, http.StatusBadRequest, "Invalid parameters: missing - 'bucket'")
	}

	id := c.Query("id")
	if id == "" {
		NewErrorResponse(c, http.StatusBadRequest, "Invalid parameters: missing - 'id'")
	}

	err := h.fileService.Delete(c.Request.Context(), bucket, id)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(200, "Success")
}
