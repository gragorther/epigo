package util

import (
	"github.com/gragorther/epigo/apperrors"
	"github.com/gragorther/epigo/models"
)

func ParseGroups(groupIDs []uint) ([]models.Group, error) {
	var groups []models.Group
	for _, id := range groupIDs {
		groups = append(groups, models.Group{ID: id})
	}
	if len(groups) != len(groupIDs) {
		return nil, apperrors.ErrParsingFailed
	}
	return groups, nil
}
