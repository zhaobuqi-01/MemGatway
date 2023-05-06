package dao

import "github.com/gin-gonic/gin"

type Creator[T Model] interface {
	Create(c *gin.Context, item *T) error
}

type Updater[T Model] interface {
	Update(c *gin.Context, item *T) error
}

type Deleter[T Model] interface {
	Delete(c *gin.Context, item *T) error
}

type Getter[T Model] interface {
	Get(c *gin.Context, item *T) (*T, error)
}

type Lister[T Model] interface {
	List(c *gin.Context, item *T) ([]T, error)
}

type Tabler interface {
	TablerName() string
}
