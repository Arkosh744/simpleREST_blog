package rest

import (
	"github.com/Arkosh744/simpleREST_blog/internal/domain"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func (h *Handler) signUp(c *gin.Context) {
	var inp domain.SignUpInput
	if err := c.ShouldBind(&inp); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	err := h.usersService.SignUp(c, inp)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, "OK")
}

func (h *Handler) signIn(c *gin.Context) {
	var inp domain.SignInInput
	if err := c.ShouldBind(&inp); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	log.Println(inp)
	accessToken, refreshToken, err := h.usersService.SignIn(c, inp)
	if err != nil {
		log.Println("signIn", err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.SetCookie("refresh-token", refreshToken, 2592000, "/", "", false, true)
	c.Writer.Header().Set("Content-Type", "application/json")
	c.JSON(http.StatusCreated, map[string]string{
		"token": accessToken,
	})
}

func (h *Handler) refresh(c *gin.Context) {
	cookie, err := c.Cookie("refresh-token")
	if err != nil {
		log.Println("refresh", err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	log.Println(cookie)

	accessToken, refreshToken, err := h.usersService.RefreshTokens(c, cookie)
	if err != nil {
		log.Println("signIn", err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.SetCookie("refresh-token", refreshToken, 2592000, "/", "", false, true)
	c.Writer.Header().Set("Content-Type", "application/json")
	c.JSON(http.StatusCreated, map[string]string{
		"token": accessToken,
	})
}