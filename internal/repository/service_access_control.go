package repository

import (
	"gateway/internal/entity"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AccessControl interface {
	Getter[entity.AccessControl]
	Updater[entity.AccessControl]
}

type accessControlRepo struct {
	db *gorm.DB
}

func NewAccessControl(db *gorm.DB) AccessControl {
	return &accessControlRepo{
		db: db,
	}
}

// Get
func (repo *accessControlRepo) Get(c *gin.Context, search *entity.AccessControl) (*entity.AccessControl, error) {
	return Get(c, repo.db, search)
}

// Updte
func (repo *accessControlRepo) Update(c *gin.Context, data *entity.AccessControl) error {
	return Update(c, repo.db, data)
}
