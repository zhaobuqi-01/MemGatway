package dao

import (
	"github.com/gin-gonic/gin"
)

type Saver[T Model] interface {
	Save(c *gin.Context, item *T) error
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
	List(c *gin.Context, serviceID int64) ([]T, error)
}

type PageLister[T Model] interface {
	PageList(c *gin.Context, search string, PageNo, PageSize int) ([]T, int64, error)
}

type Tabler interface {
	TablerName() string
}
