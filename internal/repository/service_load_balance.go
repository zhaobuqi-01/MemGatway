package repository

import (
	"gateway/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LoadBalance interface {
	Getter[model.LoadBalance]
	Updater[model.LoadBalance]
}

type LoadBalanceRepo struct {
	db *gorm.DB
}

func NewLoadBalanceRepo(db *gorm.DB) LoadBalance {
	return &LoadBalanceRepo{
		db: db,
	}
}

func (repo *LoadBalanceRepo) Get(c *gin.Context, search *model.LoadBalance) (*model.LoadBalance, error) {
	return Get(c, repo.db, search)
}

func (repo *LoadBalanceRepo) Update(c *gin.Context, data *model.LoadBalance) error {
	return Update(c, repo.db, data)
}
