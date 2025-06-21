package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gragorther/epigo/models"
	"gorm.io/gorm"
)

type MessageHandler struct {
	DB *gorm.DB
}

func NewMessageHandler(db *gorm.DB) *MessageHandler {
	return &MessageHandler{DB: db}
}

type addMessageInput struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	GroupIDs []uint `json:"groupIDs"`
}

func (h *MessageHandler) AddLastMessage(c *gin.Context) {
	currentUser, _ := c.Get("currentUser")
	user := currentUser.(models.User)
	var input addMessageInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Couldn't bind JSON, wrong input?"})
		return
	}
	var authorizedGroups []models.Group
	if err := h.DB.
		Where("id IN ?", input.GroupIDs).
		Where("user_id = ?", user.ID).
		Find(&authorizedGroups).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't check if the user is authorized to add last messages to this group"})
		return
	}
	if len(authorizedGroups) != len(input.GroupIDs) {
		c.JSON(http.StatusForbidden, gin.H{"error": "invalid group selection"})
		return
	}
	var groups []models.Group
	for _, id := range input.GroupIDs {
		groups = append(groups, models.Group{ID: id})
	}

	newLastMessage := models.LastMessage{
		Title:   input.Title,
		Content: input.Content,
		Groups:  groups,
		UserID:  user.ID,
	}

	h.DB.Create(&newLastMessage)
}

type group struct {
	ID uint `gorm:"primarykey"`
}

type lastMessage struct {
	ID      uint    `gorm:"primarykey"`
	Title   string  `json:"title"`
	Groups  []group `json:"groups" gorm:"many2many:group_last_messages;"`
	Content string  `json:"content"`
}

type listLastMessagesOutput struct {
	ID       uint   `json:"ID"`
	Title    string `json:"title"`
	GroupIDs []uint `json:"groupIDs"`
	Content  string `json:"content"`
}

func (h *MessageHandler) ListLastMessages(c *gin.Context) {
	currentUser, _ := c.Get("currentUser")

	user := currentUser.(models.User)

	var lastMessages []lastMessage
	err := h.DB.Model(&models.LastMessage{}).Preload("Groups").Where("user_id = ?", user.ID).Find(&lastMessages)
	if err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list last messages"})
		log.Print(err.Error)
		return
	}

	c.JSON(http.StatusOK, lastMessages)

}
