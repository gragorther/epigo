package groups

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	dbHandler "github.com/gragorther/epigo/database/db"
	ginctx "github.com/gragorther/epigo/handlers/context"
	"github.com/guregu/null/v6"
)

type Recipient struct {
	Email string `json:"email"`
}

type AddGroupInput struct {
	Name           string      `json:"name" binding:"required"`
	Description    null.String `json:"description"`
	Recipients     []Recipient
	LastMessageIDs []uint `json:"lastMessageIDs"`
}

func Add(db interface {
	CreateGroup(context.Context, dbHandler.CreateGroup) error
},
) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := ginctx.GetUserID(c)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		var input AddGroupInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.AbortWithError(http.StatusUnprocessableEntity, fmt.Errorf("failed to bind add group JSON: %w", err))
			return
		}

		if input.Name == "" {
			c.AbortWithStatus(http.StatusUnprocessableEntity)
			return
		}

		err = db.CreateGroup(c, dbHandler.CreateGroup{UserID: userID, Name: input.Name, Description: input.Description, LastMessageIDs: input.LastMessageIDs})
		if err != nil {

			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}
}

func Delete(db interface {
	DeleteGroupByID(ctx context.Context, id uint) error
	UserAuthorizationForGroups(ctx context.Context, groupIDs []uint, userID uint) (bool, error)
},
) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := ginctx.GetUserID(c)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}

		authorized, err := db.UserAuthorizationForGroups(c, []uint{uint(id)}, userID)
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

func List(db interface {
	GroupsByUserID(ctx context.Context, userID uint) (groups []dbHandler.Group, err error)
},
) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := ginctx.GetUserID(c)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		groups, err := db.GroupsByUserID(c, userID)
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

type EditGroupInput struct {
	Name           null.String `json:"name"`
	Description    null.String `json:"description"`
	LastMessageIDs []uint      `json:"lastMessageIDs"`
}

func Edit(db interface {
	CanUserEditGroup(ctx context.Context, userID uint, groupID uint, lastMessageIDs []uint) (authorized bool, err error)
	UpdateGroup(ctx context.Context, id uint, group dbHandler.UpdateGroup) error
},
) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := ginctx.GetUserID(c)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		var input EditGroupInput
		err = c.ShouldBindJSON(&input)
		if err != nil {
			c.AbortWithError(http.StatusUnprocessableEntity, fmt.Errorf("failed to bind edit group json: %w", err))
			return
		}
		id, err := ginctx.GetID(c)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to get ID from context while editing group: %w", err))
			return
		}
		authorized, err := db.CanUserEditGroup(c, userID, id, input.LastMessageIDs)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to check user authorization for group: %w", err))
			return
		}
		if !authorized {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		err = db.UpdateGroup(c, id, dbHandler.UpdateGroup{Name: input.Name, Description: input.Description})
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}
}
