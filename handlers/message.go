package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gragorther/epigo/apperrors"
	"github.com/gragorther/epigo/models"
	"github.com/gragorther/epigo/types"
)

type messageInput struct {
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
		return nil, apperrors.ErrParsingFailed
	}
	return groups, nil
}

func AddLastMessage(db interface {
	CheckUserAuthorizationForGroup(groupIDs []uint, userID uint) (bool, error)
	CreateLastMessage(lastMessage *models.LastMessage) error
}) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUser, _ := c.Get("currentUser")
		user := currentUser.(*models.User)
		var input messageInput

		err := c.ShouldBindJSON(&input)
		if err != nil {
			c.AbortWithError(http.StatusUnprocessableEntity, apperrors.ErrParsingFailed)
			return
		}
		authorized, err := db.CheckUserAuthorizationForGroup(input.GroupIDs, user.ID)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, apperrors.ErrAuthCheckFailed)
			return
		}
		if !authorized {
			c.JSON(http.StatusUnauthorized, gin.H{"error": apperrors.ErrUnauthorized.Error()})
			return
		}

		groups, parseErr := parseGroups(input.GroupIDs)
		if parseErr != nil {
			c.AbortWithError(http.StatusUnprocessableEntity, apperrors.ErrParsingFailed)
			return
		}

		newLastMessage := models.LastMessage{
			Title:   input.Title,
			Content: input.Content,
			Groups:  groups,
			UserID:  user.ID,
		}

		err = db.CreateLastMessage(&newLastMessage)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, apperrors.ErrCreationOfObjectFailed)
			return
		}
	}
}

func ListLastMessages(db interface {
	FindLastMessagesByUserID(userID uint) ([]types.LastMessageOut, error)
}) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUser, _ := c.Get("currentUser")

		user := currentUser.(*models.User)

		lastMessages, err := db.FindLastMessagesByUserID(user.ID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, apperrors.ErrDatabaseFetchFailed)
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
		currentUser, _ := c.Get("currentUser")
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, apperrors.ErrTypeConversionFailed)
			return
		}

		user := currentUser.(*models.User)
		var input messageInput
		err = c.ShouldBindJSON(&input)
		if err != nil {
			c.AbortWithError(http.StatusUnprocessableEntity, apperrors.ErrParsingFailed)
			return
		}
		authorized, autherr := db.CheckUserAuthorizationForLastMessage(uint(id), user.ID)
		if autherr != nil {
			c.AbortWithError(http.StatusInternalServerError, apperrors.ErrAuthCheckFailed)
			return
		}
		if !authorized {
			c.AbortWithError(http.StatusUnauthorized, apperrors.ErrUnauthorizedToEdit)
			return
		}
		authorizedToAddGroups, groupsAuthErr := db.CheckUserAuthorizationForGroup(input.GroupIDs, user.ID)
		if groupsAuthErr != nil {
			c.AbortWithError(http.StatusUnauthorized, apperrors.ErrAuthCheckFailed)
			return
		}
		if !authorizedToAddGroups {
			c.AbortWithError(http.StatusUnauthorized, apperrors.ErrUnauthorizedToEdit)
			return
		}

		groups, parseErr := parseGroups(input.GroupIDs)
		if parseErr != nil {
			c.AbortWithError(http.StatusUnprocessableEntity, apperrors.ErrParsingFailed)
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
			c.AbortWithError(http.StatusInternalServerError, apperrors.ErrDatabaseFetchFailed)
			return
		}
	}

}

func DeleteLastMessage(db interface {
	CheckUserAuthorizationForLastMessage(messageID uint, userID uint) (bool, error)
	DeleteLastMessageByID(lastMessageID uint) error
}) gin.HandlerFunc {
	return func(c *gin.Context) {
		lastMessageID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, apperrors.ErrTypeConversionFailed)
			return
		}
		currentUser, _ := c.Get("currentUser")
		user := currentUser.(*models.User)
		authorized, authErr := db.CheckUserAuthorizationForLastMessage(uint(lastMessageID), user.ID)
		if authErr != nil {
			c.AbortWithError(http.StatusUnauthorized, apperrors.ErrAuthCheckFailed)
			return
		}
		if !authorized {
			c.AbortWithError(http.StatusUnauthorized, apperrors.ErrUnauthorized)
			return
		}
		err = db.DeleteLastMessageByID(uint(lastMessageID))

		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, apperrors.ErrDeleteFailed)
			return
		}
	}

}
