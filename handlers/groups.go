package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gragorther/epigo/apperrors"
	"github.com/gragorther/epigo/email"
	"github.com/gragorther/epigo/models"
	"github.com/gragorther/epigo/types"
)

type groupInput struct {
	RecipientEmails []string `json:"recipientEmails"`
	Name            string   `json:"name"`
	Description     string   `json:"description"`
}

func AddGroup(db interface {
	CreateGroupAndRecipientEmails(group *models.Group, recipientEmails *[]models.RecipientEmail) error
}) gin.HandlerFunc {
	return func(c *gin.Context) {

		currentUser, _ := c.Get("currentUser")
		user, ok := currentUser.(*models.User)
		if !ok {
			c.AbortWithError(http.StatusInternalServerError, apperrors.ErrTypeConversionFailed)
			return
		}

		var input groupInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.AbortWithError(http.StatusInternalServerError, apperrors.ErrTypeConversionFailed)
			return
		}
		sendToGroup := models.Group{
			UserID:      user.ID,
			Name:        input.Name,
			Description: input.Description,
		}
		var newRecipientEmails []models.RecipientEmail
		for _, e := range input.RecipientEmails {
			if !email.Validate(e) {
				c.Error(apperrors.ErrInvalidEmail)
				return
			}
			newRecipientEmails = append(newRecipientEmails, models.RecipientEmail{
				GroupID: sendToGroup.ID,
				Email:   e,
			})
		}
		err := db.CreateGroupAndRecipientEmails(&sendToGroup, &newRecipientEmails)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, apperrors.ErrCreationOfObjectFailed)
			return
		}
	}
}

func DeleteGroup(db interface {
	DeleteGroupByID(id uint) error
	CheckUserAuthorizationForGroup(groupIDs []uint, userID uint) (bool, error)
}) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUser, _ := c.Get("currentUser")

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.AbortWithError(http.StatusNotFound, apperrors.ErrUserNotFound)
			return
		}
		user := currentUser.(*models.User)
		authorized, err := db.CheckUserAuthorizationForGroup([]uint{uint(id)}, user.ID)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, apperrors.ErrAuthCheckFailed)
			return
		}
		if !authorized {
			c.AbortWithError(http.StatusUnauthorized, apperrors.ErrUnauthorized)
			return
		}

		err = db.DeleteGroupByID(uint(id))
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, apperrors.ErrCreationOfObjectFailed)
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "group deleted"})
	}
}

func ListGroups(db interface {
	FindGroupsAndRecipientEmailsByUserID(userID uint) ([]types.GroupWithEmails, error)
}) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUser, _ := c.Get("currentUser")

		user, ok := currentUser.(*models.User)
		if !ok {
			c.AbortWithError(http.StatusInternalServerError, apperrors.ErrTypeConversionFailed)
			return
		}
		groups, err := db.FindGroupsAndRecipientEmailsByUserID(user.ID) // gets the list of groups a user has via the association "Groups" on the User model
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, apperrors.ErrDatabaseFetchFailed)
			return
		}

		if len(groups) == 0 {
			c.AbortWithError(http.StatusNotFound, apperrors.ErrNoGroups)
			return
		}

		c.JSON(http.StatusOK, groups)
	}
}

type editGroupInput struct {
	RecipientEmails []string `json:"recipientEmails"`
	Name            string   `json:"name"`
	Description     string   `json:"description"`
}

func EditGroup(db interface {
	CheckUserAuthorizationForGroup(groupIDs []uint, userID uint) (bool, error)
	UpdateGroup(group *models.Group, recipientEmails *[]models.RecipientEmail) error
}) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUser, _ := c.Get("currentUser")

		user := currentUser.(*models.User)

		var input editGroupInput
		err := c.ShouldBindJSON(&input)
		if err != nil {
			c.AbortWithError(http.StatusUnprocessableEntity, apperrors.ErrParsingFailed)
			return
		}
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, apperrors.ErrTypeConversionFailed)
			return
		}

		authorized, err := db.CheckUserAuthorizationForGroup([]uint{uint(id)}, user.ID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, apperrors.ErrAuthCheckFailed)
			return
		}
		if !authorized {
			c.AbortWithError(http.StatusUnauthorized, apperrors.ErrUnauthorizedToEdit)
			return
		}

		var group models.Group

		recipientEmails := make([]models.RecipientEmail, len(input.RecipientEmails))
		for i, address := range input.RecipientEmails {
			valid := email.Validate(address)
			if !valid {
				c.AbortWithError(http.StatusBadRequest, apperrors.ErrInvalidEmail)
				return
			}
			recipientEmails[i] = models.RecipientEmail{Email: address}
		}

		group.Name = input.Name
		group.Description = input.Description
		group.ID = uint(id)
		err = db.UpdateGroup(&group, &recipientEmails)
		if err != nil {
			c.AbortWithError(http.StatusNotFound, apperrors.ErrNoGroups)
			return
		}
	}

}
