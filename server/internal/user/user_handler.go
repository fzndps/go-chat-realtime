package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		Service: service,
	}
}

func (h *Handler) CreateUser(ctx *gin.Context) {
	var u CreateUserReq

	if err := ctx.ShouldBindJSON(&u); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.Service.CreateUser(ctx.Request.Context(), &u)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (h *Handler) Login(ctx *gin.Context) {
	var u LoginUserReq

	if err := ctx.ShouldBindJSON(&u); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.Service.Login(ctx, &u)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.SetCookie("jwt", user.accessToken, 3600, "/", "localhost", false, true)
	res := &LoginUserRes{
		ID:       user.ID,
		Username: user.Username,
	}

	ctx.JSON(http.StatusOK, res)
}

func (h *Handler) Logout(ctx *gin.Context) {
	ctx.SetCookie("jwt", "", -1, "", "", false, true)
	ctx.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}
