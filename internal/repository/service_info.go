package repository

import (
	"gateway/internal/dto"
	"gateway/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ServiceInfoRepository interface {
}

type serviceInfoRepo struct {
	DB *gorm.DB
}

func NewServiceInfoRepository(db *gorm.DB) ServiceInfoRepository {
	return &serviceInfoRepo{
		DB: db,
	}
}

func (repo *serviceInfoRepo) PageList(c *gin.Context, tx *gorm.DB, param *dto.ServiceListInput) ([]model.ServiceInfo, int64, error) {
	total := int64(0)
	list := []model.ServiceInfo{}
	offset := (param.PageNo - 1) * param.PageSize

	query := repo.DB.Where("is_delete=0")
	if param.Info != "" {
		query = query.Where("(service_name like ? or service_desc like ?)", "%"+param.Info+"%", "%"+param.Info+"%")
	}
	if err := query.Limit(param.PageSize).Offset(offset).Order("id desc").Find(&list).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}
	query.Limit(param.PageSize).Offset(offset).Count(&total)
	return list, total, nil
}
