package repository

import (
	"gateway/internal/entity"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HTTPRule interface {
	Getter[entity.HttpRule]
	Updater[entity.HttpRule]
}

type httpRuleRepo struct {
	db *gorm.DB
}

func NewHttpRulesitory(db *gorm.DB) HTTPRule {
	return &httpRuleRepo{
		db: db,
	}
}

func (repo *httpRuleRepo) Get(c *gin.Context, search *entity.HttpRule) (*entity.HttpRule, error) {
	return Get(c, repo.db, search)
}

func (repo *httpRuleRepo) Update(c *gin.Context, data *entity.HttpRule) error {
	return Update(c, repo.db, data)
}
