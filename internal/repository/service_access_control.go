package repository

import (
	"gateway/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AccessControl interface {
	Getter[model.AccessControl]
	Updater[model.AccessControl]
}

type AccessControlRepo struct {
	db *gorm.DB
}

func NewAccessControl(db *gorm.DB) AccessControl {
	return &AccessControlRepo{
		db: db,
	}
}

// Get
func (repo *AccessControlRepo) Get(c *gin.Context, search *model.AccessControl) (*model.AccessControl, error) {
	return Get(c, repo.db, search)
}

// Updte
func (repo *AccessControlRepo) Update(c *gin.Context, data *model.AccessControl) error {
	return Update(c, repo.db, data)
}
