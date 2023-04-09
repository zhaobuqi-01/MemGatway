package repository

import (
	"gateway/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HTTPRule interface {
	Getter[model.HttpRule]
	Updater[model.HttpRule]
}

type httpRule struct {
	db *gorm.DB
}

func NewHttpRulesitory(db *gorm.DB) HTTPRule {
	return &httpRule{
		db: db,
	}
}

func (repo *httpRule) Get(c *gin.Context, search *model.HttpRule) (*model.HttpRule, error) {
	return Get(c, repo.db, search)
}

func (repo *httpRule) Update(c *gin.Context, data *model.HttpRule) error {
	return Update(c, repo.db, data)
}
