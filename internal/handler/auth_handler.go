package handler

import (
	"hackaton-service-api/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	UserUC *usecase.UserUseCase
}

func NewAuthHandler(userUC *usecase.UserUseCase) *AuthHandler {
	return &AuthHandler{UserUC: userUC}
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Register godoc
// @Summary Registra um novo usuário
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Dados do usuário"
// @Success 201 {object} map[string]string
// @Router /api/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos: " + err.Error()})
		return
	}

	if err := h.UserUC.Register(req.Username, req.Email, req.Password); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "usuário já existe" || err.Error() == "email já cadastrado" {
			status = http.StatusConflict
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Usuário criado com sucesso"})
}

// Login godoc
// @Summary Realiza login do usuário
// @Description Autentica o usuário e retorna um token JWT
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Credenciais de Login"
// @Success 200 {object} map[string]string
// @Router /api/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	token, username, err := h.UserUC.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":    token,
		"username": username,
	})
}