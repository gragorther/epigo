package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gragorther/epigo/email"
	"github.com/gragorther/epigo/models"
)

type GroupInput struct {
	Name        string                `json:"name" binding:"required"`
	Description *string               `json:"description"`
	Recipients  []models.APIRecipient `json:"recipients"`
}

func parseAPIRecipients(recipients []models.APIRecipient) []models.Recipient {
	var newRecipients []models.Recipient
	for _, recipient := range recipients {
		newRecipients = append(newRecipients, models.Recipient{APIRecipient: recipient})
	}
	return newRecipients
}

func AddGroup(db interface {
	CreateGroupAndRecipientEmails(group *models.Group) error
}) gin.HandlerFunc {
	return func(c *gin.Context) {

		user, err := GetUserFromContext(c)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		var input GroupInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.AbortWithError(http.StatusUnprocessableEntity, fmt.Errorf("failed to bind add group JSON: %w", err))
			return
		}
		sendToGroup := models.Group{
			UserID:      user.ID,
			Name:        input.Name,
			Description: input.Description,
		}
		if sendToGroup.Name == "" {
			c.AbortWithStatus(http.StatusUnprocessableEntity)
			return
		}

		for _, e := range input.Recipients {
			if !email.Validate(e.Email) {
				c.AbortWithStatus(http.StatusUnprocessableEntity)
				return
			}
		}

		sendToGroup.Recipients = parseAPIRecipients(input.Recipients)
		err = db.CreateGroupAndRecipientEmails(&sendToGroup)
		if err != nil {
			slog.Error("failed to create group and recipient emails",
				slog.Uint64("user_id", uint64(user.ID)),
			)
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}
}

func DeleteGroup(db interface {
	DeleteGroupByID(id uint) error
	CheckUserAuthorizationForGroup(groupIDs []uint, userID uint) (bool, error)
}) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := GetUserFromContext(c)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}

		authorized, err := db.CheckUserAuthorizationForGroup([]uint{uint(id)}, user.ID)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("failed to check user authorization for group: %w", err))
			return
		}
		if !authorized {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		err = db.DeleteGroupByID(uint(id))
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "group deleted"})
	}
}

func ListGroups(db interface {
	FindGroupsAndRecipientsByUserID(userID uint) ([]models.Group, error)
}) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := GetUserFromContext(c)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		groups, err := db.FindGroupsAndRecipientsByUserID(user.ID) // gets the list of groups a user has via the association "Groups" on the User model
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to find groups and recipients by user ID during ListGroups: %w", err))
			return
		}

		if len(groups) == 0 {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.JSON(http.StatusOK, groups)
	}
}

func EditGroup(db interface {
	CheckUserAuthorizationForGroup(groupIDs []uint, userID uint) (bool, error)
	UpdateGroup(group *models.Group) error
}) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := GetUserFromContext(c)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		var input GroupInput
		err = c.ShouldBindJSON(&input)
		if err != nil {
			c.AbortWithError(http.StatusUnprocessableEntity, fmt.Errorf("failed to bind edit group json: %w", err))
			return
		}
		id, err := GetIDFromContext(c)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to get ID from context while editing group: %w", err))
			return
		}
		authorized, err := db.CheckUserAuthorizationForGroup([]uint{uint(id)}, user.ID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to check user authorization for group: %w", err))
			return
		}
		if !authorized {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		groupToUpdate := &models.Group{
			ID:          uint(id),
			Name:        input.Name,
			Description: input.Description,
			Recipients:  parseAPIRecipients(input.Recipients),
		}

		err = db.UpdateGroup(groupToUpdate)
		if err != nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
	}
}
