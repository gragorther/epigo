package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gragorther/epigo/models"
	"gorm.io/gorm"
)

type MessageHandler struct {
	DB *gorm.DB
}

func NewMessageHandler(db *gorm.DB) *GroupHandler {
	return &GroupHandler{DB: db}
}

type addMessageInput struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	GroupIDs []uint `json:"groupIDs"`
}

func (h *MessageHandler) AddLastMessage(c *gin.Context) {
	currentUser, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	user, ok := currentUser.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assert user type"})
		return
	}
	var input addMessageInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Couldn't bind JSON, wrong input?"})
	}
}
