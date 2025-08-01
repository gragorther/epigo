package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gragorther/epigo/models"
	"github.com/gragorther/epigo/types"
)

type MessageInput struct {
	Title    string `json:"title" binding:"required"`
	Content  string `json:"content" binding:"required"`
	GroupIDs []uint `json:"groupIDs" binding:"required"`
}

func parseGroups(groupIDs []uint) ([]models.Group, error) {
	var groups []models.Group
	for _, id := range groupIDs {
		groups = append(groups, models.Group{ID: id})
	}
	if len(groups) != len(groupIDs) {
		return nil, fmt.Errorf("failed to parse groups")
	}
	return groups, nil
}

func AddLastMessage(db interface {
	CheckUserAuthorizationForGroup(groupIDs []uint, userID uint) (bool, error)
	CreateLastMessage(lastMessage *models.LastMessage) error
}) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := GetUserFromContext(c)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		var input MessageInput

		err = c.ShouldBindJSON(&input)
		if err != nil {
			c.AbortWithError(http.StatusUnprocessableEntity, fmt.Errorf("failed to bind json while creating last message: %w", err))
			return
		}
		authorized, err := db.CheckUserAuthorizationForGroup(input.GroupIDs, user.ID)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}
		if !authorized {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		groups, err := parseGroups(input.GroupIDs)
		if err != nil {
			c.AbortWithError(http.StatusUnprocessableEntity, err)
			return
		}

		newLastMessage := models.LastMessage{
			Title:   input.Title,
			Content: input.Content,
			Groups:  groups,
			UserID:  user.ID,
		}

		err = db.CreateLastMessage(&newLastMessage)
		c.Status(http.StatusOK)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to create last message: %w", err))
			return
		}
	}
}

func ListLastMessages(db interface {
	FindLastMessagesByUserID(userID uint) ([]types.LastMessageOut, error)
}) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := GetUserFromContext(c)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		lastMessages, err := db.FindLastMessagesByUserID(user.ID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to find last message: %w", err))
			return
		}

		c.JSON(http.StatusOK, lastMessages)
	}
}

func EditLastMessage(db interface {
	CheckUserAuthorizationForLastMessage(messageID uint, userID uint) (bool, error)
	CheckUserAuthorizationForGroup(groupIDs []uint, userID uint) (bool, error)
	UpdateLastMessage(newMessage *models.LastMessage) error
}) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := GetUserFromContext(c)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		id, err := GetIDFromContext(c)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		var input MessageInput
		err = c.ShouldBindJSON(&input)
		if err != nil {
			c.AbortWithError(http.StatusUnprocessableEntity, fmt.Errorf("failed to bind edit last message json: %w", err))
			return
		}
		authorized, autherr := db.CheckUserAuthorizationForLastMessage(uint(id), user.ID)
		if autherr != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to check user authorization for last message: %w", err))
			return
		}
		if !authorized {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		authorizedToAddGroups, groupsAuthErr := db.CheckUserAuthorizationForGroup(input.GroupIDs, user.ID)
		if groupsAuthErr != nil {
			c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("failed to check user auth for group during edit last message: %w", groupsAuthErr))
			return
		}
		if !authorizedToAddGroups {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		groups, parseErr := parseGroups(input.GroupIDs)
		if parseErr != nil {
			c.AbortWithError(http.StatusUnprocessableEntity, fmt.Errorf("failed to edit last message: %w", parseErr))
			return
		}
		editedMessage := models.LastMessage{
			ID:      uint(id),
			Title:   input.Title,
			Content: input.Content,
			Groups:  groups,
		}
		err = db.UpdateLastMessage(&editedMessage)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to update last message: %w", err))
			return
		}
	}

}

func DeleteLastMessage(db interface {
	CheckUserAuthorizationForLastMessage(messageID uint, userID uint) (bool, error)
	DeleteLastMessageByID(lastMessageID uint) error
}) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := GetUserFromContext(c)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to get user from context when deleting last message: %w", err))
			return
		}
		lastMessageID, err := GetIDFromContext(c)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("failed to delete last message: %w", err))
			return
		}

		authorized, authErr := db.CheckUserAuthorizationForLastMessage(uint(lastMessageID), user.ID)
		if authErr != nil {
			c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("failed to check auth for last message: %w", err))
			return
		}
		if !authorized {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		err = db.DeleteLastMessageByID(uint(lastMessageID))

		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to delete last message: %w", err))
			return
		}
	}

}
