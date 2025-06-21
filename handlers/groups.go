package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gragorther/epigo/db"
	"github.com/gragorther/epigo/models"
	"gorm.io/gorm"
)

type groupInput struct {
	RecipientEmails []string `json:"recipientEmails"`
	Name            string   `json:"name"`
	Description     string   `json:"description"`
}

type GroupHandler struct {
	db *db.DBHandler // functions for database operations
}

func NewGroupHandler(db *gorm.DB, dbHandler *db.DBHandler) *GroupHandler {
	return &GroupHandler{db: dbHandler}
}

func (h *GroupHandler) AddGroup(c *gin.Context) {
	currentUser, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	user, ok := currentUser.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assert user type"})
		return
	}

	var input groupInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sendToGroup := models.Group{
		UserID:      user.ID,
		Name:        input.Name,
		Description: input.Description,
	}
	var newRecipientEmails []models.RecipientEmail
	for _, e := range input.RecipientEmails {
		newRecipientEmails = append(newRecipientEmails, models.RecipientEmail{
			GroupID: sendToGroup.ID,
			Email:   e,
		})
	}
	h.db.CreateGroupAndRecipientEmails(&sendToGroup, &newRecipientEmails)
}

func (h *GroupHandler) DeleteGroup(c *gin.Context) {
	currentUser, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group ID"})
		return
	}
	user := currentUser.(models.User)
	authorized, err := h.db.CheckUserAuthorizationForGroup(uint(id), user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't check if the user is authorized to perform this action"})
		return
	}
	if !authorized {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authorized to delete this group"})
		return
	}

	err = h.db.DeleteGroupByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "group deleted"})
}

type listGroupsOutput struct {
	ID          uint
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

func (h *GroupHandler) ListGroups(c *gin.Context) {
	currentUser, _ := c.Get("currentUser")

	user, ok := currentUser.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assert user type"})
		return
	}
	groups, err := h.db.FindGroupsAndRecipientEmailsByUserID(user.ID) // gets the list of groups a user has via the association "Groups" on the User model
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	if len(groups) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No groups"})
		return
	}

	c.JSON(http.StatusOK, groups)
}

type editGroupInput struct {
	RecipientEmails []string `json:"recipientEmails"`
	Name            string   `json:"name"`
	Description     string   `json:"description"`
}

func (h *GroupHandler) EditGroup(c *gin.Context) {
	currentUser, _ := c.Get("currentUser")

	user := currentUser.(models.User)

	var input editGroupInput
	c.ShouldBindJSON(&input)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group ID"})
		return
	}

	authorized, err := h.db.CheckUserAuthorizationForGroup(uint(id), user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't check if the user is authorized to perform this action"})
		return
	}
	if !authorized {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authorized to edit this group"})
	}

	var group models.Group

	recipientEmails := make([]models.RecipientEmail, len(input.RecipientEmails))
	for i, email := range input.RecipientEmails {
		recipientEmails[i] = models.RecipientEmail{Email: email}
	}

	group.Name = input.Name
	group.Description = input.Description
	group.ID = uint(id)
	err = h.db.UpdateGroup(&group, &recipientEmails)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
	}

}
