package rest

import (
	"github.com/Arkosh744/simpleREST_blog/internal/domain"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
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
func (h *Handler) NewPost(c *gin.Context) {

	var post domain.Post
	if err := c.BindJSON(&post); err != nil {
		c.String(http.StatusBadRequest, "Bad Request: %s", err)
		return
	}

	if err := h.postServices.Create(c, post); err != nil {
		c.String(http.StatusBadRequest, "Bad Request: %s", err)
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
func (h *Handler) GetAllPosts(c *gin.Context) {
	posts, err := h.postServices.GetAll(c)
	if err != nil {
		log.Println("getAllPosts() error:", err)
		c.String(http.StatusInternalServerError, "InternalServerError: %s", err)
		return
	}

	//response, err := json.Marshal(posts)
	//if err != nil {
	//	log.Println("getAllPosts() error:", err)
	//	c.String(http.StatusInternalServerError, "InternalServerError: %s", err)
	//	return
	//}
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
func (h *Handler) GetPostById(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Println("GetPostById() error:", err)
		c.String(http.StatusBadRequest, "Invalid id - ensure it is a number")
		return
	}

	posts, err := h.postServices.GetById(c, id)
	if err != nil {
		log.Println("GetPostById() error:", err)
		c.String(http.StatusBadRequest, "getPostbyId() error: %s", err)
		return
	}
	c.JSON(http.StatusOK, posts)

	//response, err := json.Marshal(posts)
	//if err != nil {
	//	log.Println("getPostbyId() error:", err)
	//	w.WriteHeader(http.StatusInternalServerError)
	//	errorresponse := domain.PostError{"Error getting post: " + err.Error()}
	//	response, _ := json.Marshal(errorresponse)
	//	w.Write(response)
	//	return
	//}
	//
	//w.WriteHeader(http.StatusOK)
	//w.Write(response)
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
func (h *Handler) UpdatePostById(c *gin.Context) {
	var post *domain.UpdatePost
	if err := c.BindJSON(&post); err != nil {
		c.String(http.StatusBadRequest, "Bad Request: %s", err)
		return
	}

	if err := h.postServices.Update(c, post.Id, post); err != nil {
		c.String(http.StatusBadRequest, "update() error: %s", err)
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id":    post.Id,
		"title": post.Title,
		"body":  post.Body,
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
func (h *Handler) DeletePostById(c *gin.Context) {
	var post *domain.UpdatePost
	if err := c.BindJSON(&post); err != nil {
		c.String(http.StatusBadRequest, "Bad Request: %s", err)
		return
	}

	if err := h.postServices.Delete(c, post.Id); err != nil {
		c.String(http.StatusBadRequest, "delete() error: %s", err)
		return
	}
	c.JSON(http.StatusOK, "Deleted")
}
