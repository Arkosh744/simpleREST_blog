package rest

import (
	"fmt"
	"github.com/Arkosh744/simpleREST_blog/internal/domain"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func (h *Handler) UploadFile(c *gin.Context) {
	var uploadedFile domain.UploadForm

	if err := c.ShouldBind(&uploadedFile); err != nil {
		log.WithFields(log.Fields{"handler": "FileUpload"}).Error(err)
		c.JSON(http.StatusBadRequest, map[string]string{
			"message": "invalid input body",
		})
		return
	} else if uploadedFile.File == nil {
		log.WithFields(log.Fields{"handler": "FileUpload"}).Error(err)
		c.JSON(http.StatusBadRequest, map[string]string{
			"message": "no file uploaded",
		})
		return
	}

	cookie, err := c.Cookie("refresh-token")
	if err != nil {
		log.WithFields(log.Fields{"handler": "NewPost"}).Error(err)
		c.JSON(http.StatusUnauthorized, map[string]string{
			"message": err.Error(),
		})
		return
	}
	uploadedFile.Json.AuthorId, err = h.usersService.GetIdByToken(c, cookie)

	file, _ := c.FormFile("file")
	if file.Size>>20 > 10 {
		c.JSON(http.StatusBadRequest, map[string]string{
			"message": "file is too big. Limit is 10MB",
		})
		return
	}

	uploadedFile.Json.Name = file.Filename
	path := fmt.Sprintf("./saved_files/%v", uploadedFile.Json.AuthorId)
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		log.Println(err)
	}
	dst := fmt.Sprintf("%s/%s", path, file.Filename)
	if _, err = os.Stat(dst); err == nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"message": "file with this name already exists",
		})
		return
	}

	if err = h.filesService.Upload(c, uploadedFile.Json); err != nil {
		log.WithFields(log.Fields{"handler": "FileUpload"}).Error(err)
		c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
		return
	}
	if err = c.SaveUploadedFile(file, dst); err != nil {
		log.WithFields(log.Fields{"handler": "FileUpload"}).Error(err)
		c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusBadRequest, map[string]string{
		"message": file.Filename + " uploaded",
	})
}
