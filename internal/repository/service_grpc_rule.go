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

type grpcRule struct {
	db *gorm.DB
}

func NewGrpcRulesitory(db *gorm.DB) GrpcRule {
	return &grpcRule{
		db: db,
	}
}

// Get
func (repo *grpcRule) Get(c *gin.Context, search *model.GrpcRule) (*model.GrpcRule, error) {
	return Get(c, repo.db, search)
}

// Updte
func (repo *grpcRule) Update(c *gin.Context, data *model.GrpcRule) error {
	return Update(c, repo.db, data)
}
