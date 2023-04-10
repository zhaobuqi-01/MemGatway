package repository

import (
	"gateway/internal/entity"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LoadBalance interface {
	Getter[entity.LoadBalance]
	Updater[entity.LoadBalance]
}

type loadBalanceRepo struct {
	db *gorm.DB
}

func NewloadBalanceRepo(db *gorm.DB) LoadBalance {
	return &loadBalanceRepo{
		db: db,
	}
}

func (repo *loadBalanceRepo) Get(c *gin.Context, search *entity.LoadBalance) (*entity.LoadBalance, error) {
	return Get(c, repo.db, search)
}

func (repo *loadBalanceRepo) Update(c *gin.Context, data *entity.LoadBalance) error {
	return Update(c, repo.db, data)
}
