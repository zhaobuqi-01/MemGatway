package repository

import (
	"gateway/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GrpcRule interface {
	Getter[model.GrpcRule]
	Updater[model.GrpcRule]
}

type GrpcRuleRepo struct {
	db *gorm.DB
}

func NewGrpcRulesitory(db *gorm.DB) GrpcRule {
	return &GrpcRuleRepo{
		db: db,
	}
}

// Get
func (repo *GrpcRuleRepo) Get(c *gin.Context, search *model.GrpcRule) (*model.GrpcRule, error) {
	return Get(c, repo.db, search)
}

// Updte
func (repo *GrpcRuleRepo) Update(c *gin.Context, data *model.GrpcRule) error {
	return Update(c, repo.db, data)
}
