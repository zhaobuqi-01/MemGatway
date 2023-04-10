package repository

import (
	"gateway/internal/entity"

	"github.com/gin-gonic/gin"
)

type Creator[T entity.Model] interface {
	Create(c *gin.Context, item *T) error
}

type Updater[T entity.Model] interface {
	Update(c *gin.Context, item *T) error
}

type Deleter[T entity.Model] interface {
	Delete(c *gin.Context, item *T) error
}

type Getter[T entity.Model] interface {
	Get(c *gin.Context, item *T) (*T, error)
}

type Lister[T entity.Model] interface {
	List(c *gin.Context, item *T) ([]T, error)
}

type Tabler interface {
	TablerName() string
}
