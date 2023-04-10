package repository

import (
	"gateway/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Admin interface {
	Getter[model.Admin]
	Updater[model.Admin]
}

type AdminRepo struct {
	db *gorm.DB
}

func NewAdmin(db *gorm.DB) Admin {
	return &AdminRepo{
		db: db,
	}
}

// FindAdminByID finds an admin by their ID using GORM
func (repo *AdminRepo) Get(c *gin.Context, search *model.Admin) (*model.Admin, error) {
	return Get(c, repo.db, search)
}

func (repo *AdminRepo) Update(c *gin.Context, data *model.Admin) error {
	return Update(c, repo.db, data)
}
