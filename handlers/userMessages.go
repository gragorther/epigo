package handlers

import (
	"log"
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

	var recipientEmails []models.RecipientEmail

	for _, e := range input.RecipientEmails {
		recipientEmails = append(recipientEmails, models.RecipientEmail{
			UserID:        user.ID,
			SendToGroupID: sendToGroup.ID,
			Email:         e,
		})
	}

	// batch insert
	if err := h.DB.Create(&recipientEmails).Error; err != nil {
		log.Println("failed to insert recipient emails:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create recipients"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "SendToGroup and recipients created successfully",
	})
}
