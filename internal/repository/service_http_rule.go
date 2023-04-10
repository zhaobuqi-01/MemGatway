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

type HttpRuleRepo struct {
	db *gorm.DB
}

func NewHttpRulesitory(db *gorm.DB) HTTPRule {
	return &HttpRuleRepo{
		db: db,
	}
}

func (repo *HttpRuleRepo) Get(c *gin.Context, search *model.HttpRule) (*model.HttpRule, error) {
	return Get(c, repo.db, search)
}

func (repo *HttpRuleRepo) Update(c *gin.Context, data *model.HttpRule) error {
	return Update(c, repo.db, data)
}
