package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gragorther/epigo/apperrors"
	"github.com/gragorther/epigo/db"
	"github.com/gragorther/epigo/models"
)

type MessageHandler struct {
	m db.Messages
	a db.Auth
}

func NewMessageHandler(m db.Messages, a db.Auth) *MessageHandler {
	return &MessageHandler{m: m, a: a}
}

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

func (h *MessageHandler) AddLastMessage(c *gin.Context) {
	currentUser, _ := c.Get("currentUser")
	user := currentUser.(*models.User)
	var input messageInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.AbortWithError(http.StatusUnprocessableEntity, apperrors.ErrParsingFailed)
		return
	}
	authorized, err := h.a.CheckUserAuthorizationForGroup(input.GroupIDs, user.ID)
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

	err = h.m.CreateLastMessage(&newLastMessage)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, apperrors.ErrCreationOfObjectFailed)
		return
	}
}

func (h *MessageHandler) ListLastMessages(c *gin.Context) {
	currentUser, _ := c.Get("currentUser")

	user := currentUser.(*models.User)

	lastMessages, err := h.m.FindLastMessagesByUserID(user.ID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, apperrors.ErrDatabaseFetchFailed)
		return
	}

	c.JSON(http.StatusOK, lastMessages)

}

func (h *MessageHandler) EditLastMessage(c *gin.Context) {
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
	authorized, autherr := h.a.CheckUserAuthorizationForLastMessage(uint(id), user.ID)
	if autherr != nil {
		c.AbortWithError(http.StatusInternalServerError, apperrors.ErrAuthCheckFailed)
		return
	}
	if !authorized {
		c.AbortWithError(http.StatusUnauthorized, apperrors.ErrUnauthorizedToEdit)
		return
	}
	authorizedToAddGroups, groupsAuthErr := h.a.CheckUserAuthorizationForGroup(input.GroupIDs, user.ID)
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
	err = h.m.UpdateLastMessage(&editedMessage)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, apperrors.ErrDatabaseFetchFailed)
		return
	}

}

func (h *MessageHandler) DeleteLastMessage(c *gin.Context) {
	lastMessageID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, apperrors.ErrTypeConversionFailed)
		return
	}
	currentUser, _ := c.Get("currentUser")
	user := currentUser.(*models.User)
	authorized, authErr := h.a.CheckUserAuthorizationForLastMessage(uint(lastMessageID), user.ID)
	if authErr != nil {
		c.AbortWithError(http.StatusUnauthorized, apperrors.ErrAuthCheckFailed)
		return
	}
	if !authorized {
		c.AbortWithError(http.StatusUnauthorized, apperrors.ErrUnauthorized)
		return
	}
	err = h.m.DeleteLastMessageByID(uint(lastMessageID))

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, apperrors.ErrDeleteFailed)
		return
	}

}
