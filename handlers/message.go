package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gragorther/epigo/apperrors"
	"github.com/gragorther/epigo/db"
	"github.com/gragorther/epigo/models"
	"github.com/gragorther/epigo/util"
)

type MessageHandler struct {
	M db.Messages
	A db.Auth
}

type messageInput struct {
	Title    string `json:"title" binding:"required"`
	Content  string `json:"content" binding:"required"`
	GroupIDs []uint `json:"groupIDs" binding:"required"`
}

func (h *MessageHandler) AddLastMessage(c *gin.Context) {
	currentUser, _ := c.Get("currentUser")
	user := currentUser.(models.User)
	var input messageInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": apperrors.ErrParsingFailed.Error()})
		return
	}
	authorized, err := h.A.CheckUserAuthorizationForGroup(input.GroupIDs, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !authorized {
		c.JSON(http.StatusUnauthorized, gin.H{"error": apperrors.ErrUnauthorized.Error()})
		return
	}

	groups, parseErr := util.ParseGroups(input.GroupIDs)
	if parseErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": parseErr.Error()})
		return
	}

	newLastMessage := models.LastMessage{
		Title:   input.Title,
		Content: input.Content,
		Groups:  groups,
		UserID:  user.ID,
	}

	err = h.M.CreateLastMessage(&newLastMessage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": apperrors.ErrServerError.Error()})
		return
	}
}

func (h *MessageHandler) ListLastMessages(c *gin.Context) {
	currentUser, _ := c.Get("currentUser")

	user := currentUser.(models.User)

	lastMessages, err := h.M.FindLastMessagesByUserID(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, lastMessages)

}

func (h *MessageHandler) EditLastMessage(c *gin.Context) {
	currentUser, _ := c.Get("currentUser")
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": apperrors.ErrWrongParam.Error()})
		return
	}

	user := currentUser.(models.User)
	var input messageInput
	err = c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": apperrors.ErrParsingFailed.Error()})
		return
	}
	authorized, autherr := h.A.CheckUserAuthorizationForLastMessage(uint(id), user.ID)
	if autherr != nil {
		c.JSON(http.StatusInternalServerError, apperrors.ErrServerError.Error())
		return
	}
	if !authorized {
		c.JSON(http.StatusUnauthorized, apperrors.ErrUnauthorized.Error())
	}
	authorizedToAddGroups, groupsAuthErr := h.A.CheckUserAuthorizationForGroup(input.GroupIDs, user.ID)
	if groupsAuthErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": apperrors.ErrServerError.Error()})
		return
	}
	if !authorizedToAddGroups {
		c.JSON(http.StatusUnauthorized, gin.H{"error": apperrors.ErrUnauthorizedToEdit.Error()})
		return
	}

	groups, parseErr := util.ParseGroups(input.GroupIDs)
	if parseErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": apperrors.ErrParsingFailed.Error()})
	}
	editedMessage := models.LastMessage{
		ID:      uint(id),
		Title:   input.Title,
		Content: input.Content,
		Groups:  groups,
	}
	err = h.M.UpdateLastMessage(&editedMessage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

}

func (h *MessageHandler) DeleteLastMessage(c *gin.Context) {
	lastMessageID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	currentUser, _ := c.Get("currentUser")
	user := currentUser.(models.User)
	authorized, authErr := h.A.CheckUserAuthorizationForLastMessage(uint(lastMessageID), user.ID)
	if authErr != nil {
		c.JSON(http.StatusInternalServerError, apperrors.ErrAuthCheckFailed.Error())
		return
	}
	if !authorized {
		c.JSON(http.StatusUnauthorized, gin.H{"error": apperrors.ErrUnauthorized.Error()})
		return
	}
	err = h.M.DeleteLastMessageByID(uint(lastMessageID))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

}
