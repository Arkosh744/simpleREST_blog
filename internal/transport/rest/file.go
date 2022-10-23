package rest

import (
	"fmt"
	"github.com/Arkosh744/simpleREST_blog/internal/domain"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func (h *Handler) UploadFile(c *gin.Context) {
	var uploadedFile domain.UploadForm
	err := c.ShouldBind(&uploadedFile)

	if err != nil {
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
	uploadedFile.Json.Name = file.Filename
	// TODO: check folder, create sub folders by author id
	// TODO: check file size
	// TODO: check name for existing files
	// TODO: if filename exists => error
	// TODO: if ok => add info to db and save file

	fmt.Println(file.Header.Values("Content-Type")[0])
	fmt.Println(uploadedFile)
	dst := fmt.Sprintf("./file_folder/%s", file.Filename)
	err = c.SaveUploadedFile(file, dst)
	if err != nil {
		return
	}
	c.String(http.StatusOK, fmt.Sprintf("file '%s' uploaded!", file.Filename))
}
