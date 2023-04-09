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

type loadBalance struct {
	db *gorm.DB
}

func NewLoadBalance(db *gorm.DB) LoadBalance {
	return &loadBalance{db}
}

func (repo *loadBalance) Get(c *gin.Context, search *model.LoadBalance) (*model.LoadBalance, error) {
	return Get(c, repo.db, search)
}

func (repo *loadBalance) Update(c *gin.Context, data *model.LoadBalance) error {
	return Update(c, repo.db, data)
}
