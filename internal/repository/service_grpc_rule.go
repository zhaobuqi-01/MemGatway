package repository

import (
	"gateway/internal/entity"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GrpcRule interface {
	Getter[entity.GrpcRule]
	Updater[entity.GrpcRule]
}

type grpcRuleRepo struct {
	db *gorm.DB
}

func NewGrpcRulesitory(db *gorm.DB) GrpcRule {
	return &grpcRuleRepo{
		db: db,
	}
}

// Get
func (repo *grpcRuleRepo) Get(c *gin.Context, search *entity.GrpcRule) (*entity.GrpcRule, error) {
	return Get(c, repo.db, search)
}

// Updte
func (repo *grpcRuleRepo) Update(c *gin.Context, data *entity.GrpcRule) error {
	return Update(c, repo.db, data)
}
