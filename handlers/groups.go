package handlers

import (
	"context"
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
	CreateGroup(group *models.Group) error
}) gin.HandlerFunc {
	return func(c *gin.Context) {

		userID, err := GetUserIDFromContext(c)
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
			UserID:      userID,
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
		err = db.CreateGroup(&sendToGroup)
		if err != nil {
			slog.Error("failed to create group and recipient emails",
				slog.Uint64("user_id", uint64(userID)),
			)
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}
}

func DeleteGroup(db interface {
	DeleteGroupByID(ctx context.Context, id uint) error
	CheckUserAuthorizationForGroups(ctx context.Context, groupIDs []uint, userID uint) (bool, error)
}) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := GetUserIDFromContext(c)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}

		authorized, err := db.CheckUserAuthorizationForGroups(c, []uint{uint(id)}, userID)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("failed to check user authorization for group: %w", err))
			return
		}
		if !authorized {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		err = db.DeleteGroupByID(c, uint(id))
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.Status(http.StatusOK)
	}
}

func ListGroups(db interface {
	FindGroupsAndRecipientsByUserID(ctx context.Context, userID uint) ([]models.Group, error)
}) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := GetUserIDFromContext(c)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		groups, err := db.FindGroupsAndRecipientsByUserID(c, userID) // gets the list of groups a user has via the association "Groups" on the User model
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
	CheckUserAuthorizationForGroups(ctx context.Context, groupIDs []uint, userID uint) (bool, error)
	UpdateGroup(ctx context.Context, group models.Group) error
}) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := GetUserIDFromContext(c)
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
		authorized, err := db.CheckUserAuthorizationForGroups(c, []uint{uint(id)}, userID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to check user authorization for group: %w", err))
			return
		}
		if !authorized {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		groupToUpdate := models.Group{
			ID:          uint(id),
			Name:        input.Name,
			Description: input.Description,
			Recipients:  parseAPIRecipients(input.Recipients),
		}

		err = db.UpdateGroup(c, groupToUpdate)
		if err != nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
	}
}
