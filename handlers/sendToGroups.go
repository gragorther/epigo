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
	var recipientEmails []models.RecipientEmail
	for _, e := range input.RecipientEmails {
		recipientEmails = append(recipientEmails, models.RecipientEmail{
			SendToGroupID: sendToGroup.ID,
			Email:         e,
		})
	}
	h.DB.Create(&sendToGroup).Association("RecipientEmails").Append(&recipientEmails)

}

// func (h *UserHandler) removeUserHandler(c *gin.Context) {
// 	currentUser, exists := c.Get("currentUser")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
// 	}
// 	user, ok := currentUser.(models.User)
// 	if !ok {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assert user type"})
// 		return
// 	}
// }
