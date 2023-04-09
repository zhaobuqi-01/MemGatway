package repository

import (
	"gateway/internal/model"

	"github.com/gin-gonic/gin"
)

type Creator[T model.Model] interface {
	Create(c *gin.Context, item *T) error
}

type Updater[T model.Model] interface {
	Update(c *gin.Context, item *T) error
}

type Deleter[T model.Model] interface {
	Delete(c *gin.Context, item *T) error
}

type Getter[T model.Model] interface {
	Get(c *gin.Context, item *T) (*T, error)
}

type Lister[T model.Model] interface {
	List(c *gin.Context, item *T) ([]T, error)
}

type Tabler interface {
	TablerName() string
}
