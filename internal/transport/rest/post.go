package rest

import (
	"github.com/Arkosh744/simpleREST_blog/internal/domain"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

// New Post godoc
// @Summary Create new post
// @Description Create new post with title and content
// @Tags posts
// @Accept  json
// @Produce  json
// @Param new post body domain.PostQuery true "new post"
// @Success 200 {object} domain.Post
// @Router /post/new [post]
func (h *Handler) Create(c *gin.Context) {
	var post domain.Post
	err := c.BindJSON(&post)
	if post.Body == "" || post.Title == "" || err != nil {
		log.WithFields(log.Fields{"handler": "NewPost"}).Error(err)
		c.JSON(http.StatusBadRequest, map[string]string{
			"message": "invalid input post body",
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
	post.AuthorId, err = h.usersService.GetIdByToken(c, cookie)

	if err := h.postsService.Create(c, post); err != nil {
		log.WithFields(log.Fields{"handler": "NewPost"}).Error(err)
		c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, map[string]interface{}{
		"title": post.Title,
		"body":  post.Body,
	})
}

// Get posts by ID godoc
// @Summary Get details of a post
// @Description Get details of a post by ID
// @Tags posts
// @Accept  json
// @Produce  json
// @Success 200 {array} []domain.Post
// @Router /post/all [get]
func (h *Handler) List(c *gin.Context) {
	cookie, err := c.Cookie("refresh-token")
	if err != nil {
		log.WithFields(log.Fields{"handler": "List"}).Error(err)
		c.JSON(http.StatusUnauthorized, map[string]string{
			"message": err.Error(),
		})
		return
	}
	userId, _ := h.usersService.GetIdByToken(c, cookie)

	posts, err := h.postsService.List(c, userId)
	if err != nil {
		log.WithFields(log.Fields{"handler": "List"}).Error(err)
		c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, posts)
}

// Get post by ID godoc
// @Summary Get details of a post
// @Description Get details of a post by ID
// @Tags posts
// @Accept  json
// @Produce  json
// @Param id path int true "Post ID"
// @Success 200 {object} domain.Post
// @Router /post/get/{id} [get]
func (h *Handler) GetById(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "GetPostById",
		}).Error(err)
		log.WithFields(log.Fields{"handler": "GetPostById"}).Error(err)
		c.JSON(http.StatusBadRequest, map[string]string{
			"message": "invalid input post id",
		})
		return
	}
	cookie, err := c.Cookie("refresh-token")
	if err != nil {
		log.WithFields(log.Fields{"handler": "GetPostById"}).Error(err)
		c.JSON(http.StatusUnauthorized, map[string]string{
			"message": err.Error(),
		})
		return
	}
	userId, _ := h.usersService.GetIdByToken(c, cookie)

	posts, err := h.postsService.GetById(c, id, userId)
	if err != nil {
		log.WithFields(log.Fields{"handler": "GetPostById"}).Error(err)
		c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, posts)
}

// Update post by ID godoc
// @Summary Get details of a post
// @Description Get details of a post by ID
// @Tags posts
// @Accept  json
// @Produce  json
// @Param updatePost body domain.UpdatePost true "update post"
// @Success 200 {object} domain.Post
// @Router /post/update [post]
func (h *Handler) UpdateById(c *gin.Context) {
	var post *domain.UpdatePost
	err := c.BindJSON(&post)
	if post.Body == "" || post.Title == "" || err != nil {
		log.WithFields(log.Fields{"handler": "UpdatePostById"}).Error(err)
		c.JSON(http.StatusBadRequest, map[string]string{
			"message": "invalid input post body",
		})
		return
	}

	cookie, err := c.Cookie("refresh-token")
	if err != nil {
		log.WithFields(log.Fields{"handler": "UpdatePostById"}).Error(err)
		c.JSON(http.StatusUnauthorized, map[string]string{
			"message": err.Error(),
		})
		return
	}
	userId, _ := h.usersService.GetIdByToken(c, cookie)

	if err := h.postsService.Update(c, post.Id, post, userId); err != nil {
		log.WithFields(log.Fields{"handler": "UpdatePostById"}).Error(err)
		c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "updated",
	})
}

// Delete post by ID godoc
// @Summary Delete a post
// @Description Delete a post by ID
// @Tags posts
// @Accept  json
// @Produce  json
// @Param id body domain.Post true "id"
// @Success 200 {string} string "Post deleted"
// @Router /post/delete [post]
func (h *Handler) DeleteById(c *gin.Context) {
	var post *domain.UpdatePost
	err := c.BindJSON(&post)
	if post.Id == 0 || err != nil {
		log.WithFields(log.Fields{"handler": "DeletePostById"}).Error(err)
		c.JSON(http.StatusBadRequest, map[string]string{
			"message": "invalid input post body",
		})
		return
	}

	cookie, err := c.Cookie("refresh-token")
	if err != nil {
		log.WithFields(log.Fields{"handler": "DeletePostById"}).Error(err)
		c.JSON(http.StatusUnauthorized, map[string]string{
			"message": err.Error(),
		})
		return
	}
	userId, _ := h.usersService.GetIdByToken(c, cookie)

	if err := h.postsService.Delete(c, post.Id, userId); err != nil {
		log.WithFields(log.Fields{"handler": "DeletePostById"}).Error(err)
		c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "deleted",
	})
}
