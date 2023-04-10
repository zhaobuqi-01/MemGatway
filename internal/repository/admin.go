package repository

import (
	"gateway/internal/entity"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Admin interface {
	Getter[entity.Admin]
	Updater[entity.Admin]
}

type adminRepo struct {
	db *gorm.DB
}

func NewAdmin(db *gorm.DB) Admin {
	return &adminRepo{
		db: db,
	}
}

// FindAdminByID finds an admin by their ID using GORM
func (repo *adminRepo) Get(c *gin.Context, search *entity.Admin) (*entity.Admin, error) {
	return Get(c, repo.db, search)
}

func (repo *adminRepo) Update(c *gin.Context, data *entity.Admin) error {
	return Update(c, repo.db, data)
}
