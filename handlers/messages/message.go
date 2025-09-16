package messages

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	dbHandler "github.com/gragorther/epigo/database/db"
	ginctx "github.com/gragorther/epigo/handlers/context"
	"github.com/guregu/null/v6"
)

type AddMessageInput struct {
	Title    string `json:"title" binding:"required"`
	Content  string `json:"content"`
	GroupIDs []uint `json:"groupIDs"`
}

func Add(db interface {
	UserAuthorizationForGroups(ctx context.Context, groupIDs []uint, userID uint) (bool, error)
	CreateLastMessage(ctx context.Context, message dbHandler.CreateLastMessage) error
},
) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := ginctx.GetUserID(c)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		var input AddMessageInput

		err = c.ShouldBindJSON(&input)
		if err != nil {
			c.AbortWithError(http.StatusUnprocessableEntity, fmt.Errorf("failed to bind json while creating last message: %w", err))
			return
		}
		authorized, err := db.UserAuthorizationForGroups(c, input.GroupIDs, userID)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}
		if !authorized {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		err = db.CreateLastMessage(c, dbHandler.CreateLastMessage{})
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to create last message: %w", err))
			return
		}
		c.Status(http.StatusOK)
	}
}

func List(db interface {
	LastMessagesByUserID(ctx context.Context, userID uint) (lastMessages []dbHandler.LastMessage, err error)
},
) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := ginctx.GetUserID(c)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		lastMessages, err := db.LastMessagesByUserID(c, userID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to find last message: %w", err))
			return
		}

		c.JSON(http.StatusOK, lastMessages)
	}
}

type EditMessageInput struct {
	Title    null.String `json:"title"`
	Content  null.String `json:"content"`
	GroupIDs []uint      `json:"groupIDs"`
}

func Edit(db interface {
	CanUserEditLastmessage(ctx context.Context, userID uint, messageID uint, groupIDs []uint) (authorized bool, err error)
	UpdateLastMessage(ctx context.Context, id uint, m dbHandler.UpdateLastMessage) error
},
) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := ginctx.GetUserID(c)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		messageID, err := ginctx.GetID(c)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		var input EditMessageInput
		err = c.ShouldBindJSON(&input)
		if err != nil {
			c.AbortWithError(http.StatusUnprocessableEntity, fmt.Errorf("failed to bind edit last message json: %w", err))
			return
		}
		authorizedToEdit, err := db.CanUserEditLastmessage(c, userID, messageID, input.GroupIDs)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		if !authorizedToEdit {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		err = db.UpdateLastMessage(c, messageID, dbHandler.UpdateLastMessage{
			Title:    input.Title,
			GroupIDs: input.GroupIDs,
			Content:  input.Content,
		})
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to update last message: %w", err))
			return
		}

		c.Status(http.StatusNoContent)
	}
}

func Delete(db interface {
	UserAuthorizationForLastMessage(ctx context.Context, messageID uint, userID uint) (bool, error)
	DeleteLastMessageByID(ctx context.Context, id uint) error
},
) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := ginctx.GetUserID(c)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to get user from context when deleting last message: %w", err))
			return
		}
		lastMessageID, err := ginctx.GetID(c)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("failed to delete last message: %w", err))
			return
		}

		authorized, authErr := db.UserAuthorizationForLastMessage(c, uint(lastMessageID), userID)
		if authErr != nil {
			c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("failed to check auth for last message: %w", authErr))
			return
		}
		if !authorized {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		err = db.DeleteLastMessageByID(c, lastMessageID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to delete last message: %w", err))
			return
		}
	}
}
