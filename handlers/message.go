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
	db *db.DBHandler
}

func NewMessageHandler(db *db.DBHandler) *MessageHandler {
	return &MessageHandler{db: db}
}

type messageInput struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	GroupIDs []uint `json:"groupIDs"`
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
	authorized, err := h.db.CheckUserAuthorizationForGroup(input.GroupIDs, user.ID)
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

	h.db.CreateLastMessage(&newLastMessage)
}

func (h *MessageHandler) ListLastMessages(c *gin.Context) {
	currentUser, _ := c.Get("currentUser")

	user := currentUser.(models.User)

	lastMessages, err := h.db.FindLastMessagesByUserID(user.ID)
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
	c.ShouldBindJSON(&input)
	authorized, autherr := h.db.CheckUserAuthorizationForLastMessage(uint(id), user.ID)
	if autherr != nil {
		c.JSON(http.StatusInternalServerError, apperrors.ErrServerError.Error())
		return
	}
	if !authorized {
		c.JSON(http.StatusUnauthorized, apperrors.ErrUnauthorized.Error())
	}
	authorizedToAddGroups, groupsAuthErr := h.db.CheckUserAuthorizationForGroup(input.GroupIDs, user.ID)
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
	err = h.db.UpdateLastMessage(&editedMessage)
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
	authorized, authErr := h.db.CheckUserAuthorizationForLastMessage(uint(lastMessageID), user.ID)
	if authErr != nil {
		c.JSON(http.StatusInternalServerError, apperrors.ErrAuthCheckFailed.Error())
		return
	}
	if !authorized {
		c.JSON(http.StatusUnauthorized, gin.H{"error": apperrors.ErrUnauthorized.Error()})
		return
	}
	err = h.db.DeleteLastMessageByID(uint(lastMessageID))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

}
