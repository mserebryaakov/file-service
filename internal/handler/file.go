package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const MAX_UPLOAD_SIZE = 3000 << 20 //3GB

var FILE_TYPES = map[string]interface{}{
	"video/webm": nil,
}

// Отправка файла и создание на хосте
// @Summary      Uploading file
// @Description  Upload file and save in host
// @Accept       multipart/form-data
// @Produce      application/json
// @Param        File body string true "Upload file"
// @Success      200  {body}  handler.UploadResponse "Success upload"
// @Failure		 404  {body}  handler.UploadResponse "Bad request data"
// @Failure		 500  {body}  handler.UploadResponse "Server error"
// @Router       /v1/files [post]
func uploadFileHandler(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MAX_UPLOAD_SIZE)

	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Bad multipart")
		return
	}

	defer file.Close()

	buffer := make([]byte, fileHeader.Size)
	file.Read(buffer)
	fileType := http.DetectContentType(buffer)

	if _, ex := FILE_TYPES[fileType]; !ex {
		newErrorResponse(c, http.StatusBadRequest, "File type is not supported")
		return
	}

	newPath := "./files/" + uuid.New().String() + ".webm"

	err = c.SaveUploadedFile(fileHeader, newPath)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Error save file")
		return
	}

	// Успешный ответ
	c.JSON(http.StatusOK, &UploadResponse{
		Msg: "Files created",
	})
}
