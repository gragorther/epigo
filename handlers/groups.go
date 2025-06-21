package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gragorther/epigo/models"
	"gorm.io/gorm"
)

type groupInput struct {
	RecipientEmails []string `json:"recipientEmails"`
	Name            string   `json:"name"`
	Description     string   `json:"description"`
}

type GroupHandler struct {
	DB *gorm.DB
}

func NewGroupHandler(db *gorm.DB) *GroupHandler {
	return &GroupHandler{DB: db}
}

func (h *GroupHandler) AddGroup(c *gin.Context) {
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

	var input groupInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sendToGroup := models.Group{
		UserID:      user.ID,
		Name:        input.Name,
		Description: input.Description,
	}
	var newRecipientEmails []models.RecipientEmail
	for _, e := range input.RecipientEmails {
		newRecipientEmails = append(newRecipientEmails, models.RecipientEmail{
			GroupID: sendToGroup.ID,
			Email:   e,
		})
	}
	h.DB.Create(&sendToGroup).Association("RecipientEmails").Append(&newRecipientEmails)
}

func (h *GroupHandler) DeleteGroup(c *gin.Context) {
	_, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group ID"})
		return
	}

	res := h.DB.Delete(&models.Group{}, id)
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": res.Error})
	}
	if res.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "group not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "group deleted"})
}

type listGroupsOutput struct {
	ID          uint
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Name        string
	Description string
}

func (h *GroupHandler) ListGroups(c *gin.Context) {
	currentUser, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	user, ok := currentUser.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assert user type"})
		return
	}
	var groups []listGroupsOutput
	h.DB.Model(&user).Association("Groups").Find(&groups) // gets the list of groups a user has via the association "Groups" on the User model

	if len(groups) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No groups"})
		return
	}

	c.JSON(http.StatusOK, groups)
}

type editGroupInput struct {
	RecipientEmails []string `json:"recipientEmails"`
	Name            string   `json:"name"`
	Description     string   `json:"description"`
}

func (h *GroupHandler) EditGroup(c *gin.Context) {
	currentUser, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}
	user := currentUser.(models.User)

	var input editGroupInput
	c.ShouldBindJSON(&input)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group ID"})
		return
	}
	var authorizedGroup int64
	if err := h.DB.Model(&models.Group{}).
		Where("id = ?", id).
		Where("user_id = ?", user.ID).
		Count(&authorizedGroup).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't check if the user is authorized to add last messages to this group"})
	}
	if authorizedGroup != 1 {
		c.JSON(http.StatusForbidden, gin.H{"error": "invalid group selection"})
		return
	}

	var group models.Group

	recipientEmails := make([]models.RecipientEmail, len(input.RecipientEmails))
	for i, email := range input.RecipientEmails {
		recipientEmails[i] = models.RecipientEmail{Email: email}
	}

	group.Name = input.Name
	group.Description = input.Description
	group.RecipientEmails = recipientEmails
	group.ID = uint(id)
	output := h.DB.Updates(&group)
	if output.RowsAffected < 1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "group not found"})
		return
	}
}
