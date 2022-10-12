package rest

import (
	"github.com/Arkosh744/simpleREST_blog/internal/domain"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func (h *Handler) signUp(c *gin.Context) {
	var inp domain.SignUpInput
	if err := c.BindJSON(&inp); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
		return
	}
	if err := inp.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"message": domain.ErrInvalidInput.Error(),
		})
		return
	}
	err := h.usersService.SignUp(c, inp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, "OK")
}

func (h *Handler) signIn(c *gin.Context) {
	var inp domain.SignInInput
	if err := c.BindJSON(&inp); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
		return
	}
	if err := inp.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"message": domain.ErrInvalidInput.Error(),
		})
		return
	}
	accessToken, refreshToken, err := h.usersService.SignIn(c, inp)
	if err != nil {
		log.Println("signIn", err)
		if err == domain.ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, map[string]string{
				"message": err.Error(),
			})
			return
		} else {
			c.JSON(http.StatusInternalServerError, map[string]string{
				"message": err.Error(),
			})
			return
		}
	}
	c.SetCookie("refresh-token", refreshToken, 2592000, "/", "", false, true)
	c.Writer.Header().Set("Content-Type", "application/json")
	c.JSON(http.StatusOK, map[string]string{
		"token": accessToken,
	})
}

func (h *Handler) refresh(c *gin.Context) {
	cookie, err := c.Cookie("refresh-token")
	if err != nil {
		log.Println("refresh", err)
		c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
		return
	}
	accessToken, refreshToken, err := h.usersService.RefreshTokens(c, cookie)
	if err != nil {
		log.Println("signIn", err)
		c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
		return
	}

	c.SetCookie("refresh-token", refreshToken, 2592000, "/", "", false, true)
	c.Writer.Header().Set("Content-Type", "application/json")
	c.JSON(http.StatusOK, map[string]string{
		"token": accessToken,
	})
}
