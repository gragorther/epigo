package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gragorther/epigo/models"
)

func (h *UserHandler) AddSendToGroup(c *gin.Context) {
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
	var input models.SendToGroupInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sendToGroup := models.SendToGroup{
		UserID:      user.ID,
		Name:        input.Name,
		Description: input.Description,
	}
	h.DB.Create(&sendToGroup)

	for _, element := range input.RecipientEmails {
		// index is the index where we are
		// element is the element from someSlice for where we are
		object := models.RecipientEmail{
			Email:         element,
			SendToGroupID: sendToGroup.ID,
		}
		h.DB.Create(object)
	}

}
