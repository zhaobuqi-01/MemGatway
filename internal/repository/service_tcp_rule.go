package repository

import (
	"gateway/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TcpRule interface {
	Getter[model.TcpRule]
	Updater[model.TcpRule]
}

type tcpRule struct {
	db *gorm.DB
}

func NewTcpRule(db *gorm.DB) TcpRule {
	return &tcpRule{
		db: db,
	}
}

func (repo *tcpRule) Get(c *gin.Context, search *model.TcpRule) (*model.TcpRule, error) {
	return Get(c, repo.db, search)
}

func (repo *tcpRule) Update(c *gin.Context, data *model.TcpRule) error {
	return Update(c, repo.db, data)
}
