package repository

import (
	"gateway/internal/entity"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TcpRule interface {
	Getter[entity.TcpRule]
	Updater[entity.TcpRule]
}

type tcpRuleRepo struct {
	db *gorm.DB
}

func NewTcpRule(db *gorm.DB) TcpRule {
	return &tcpRuleRepo{
		db: db,
	}
}

func (repo *tcpRuleRepo) Get(c *gin.Context, search *entity.TcpRule) (*entity.TcpRule, error) {
	return Get(c, repo.db, search)
}

func (repo *tcpRuleRepo) Update(c *gin.Context, data *entity.TcpRule) error {
	return Update(c, repo.db, data)
}
