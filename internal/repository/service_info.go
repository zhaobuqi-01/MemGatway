package repository

import (
	"gateway/internal/dto"
	"gateway/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ServiceInfo interface {
	// Getter[model.ServiceInfo]
	Updater[model.ServiceInfo]
	PageList(c *gin.Context, param *dto.ServiceListInput) ([]model.ServiceInfo, int64, error)
	ServiceDetail(c *gin.Context, search *model.ServiceInfo) (*model.ServiceDetail, error)
}

type ServiceInfoRepo struct {
	db *gorm.DB
}

func NewServiceInfoRepo(db *gorm.DB) ServiceInfo {
	return &ServiceInfoRepo{
		db: db,
	}
}

func (repo *ServiceInfoRepo) Get(c *gin.Context, search *model.ServiceInfo) (*model.ServiceInfo, error) {
	return Get(c, repo.db, search)
}

func (repo *ServiceInfoRepo) Update(c *gin.Context, data *model.ServiceInfo) error {
	return Update(c, repo.db, data)
}

func (repo *ServiceInfoRepo) PageList(c *gin.Context, param *dto.ServiceListInput) ([]model.ServiceInfo, int64, error) {
	total := int64(0)
	list := []model.ServiceInfo{}
	offset := (param.PageNo - 1) * param.PageSize

	query := repo.db.Where("is_delete=0")
	if param.Info != "" {
		query = query.Where("(service_name like ? or service_desc like ?)", "%"+param.Info+"%", "%"+param.Info+"%")
	}
	if err := query.Limit(param.PageSize).Offset(offset).Order("id desc").Find(&list).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}
	query.Limit(param.PageSize).Offset(offset).Count(&total)
	return list, total, nil
}

func (repo *ServiceInfoRepo) ServiceDetail(c *gin.Context, search *model.ServiceInfo) (*model.ServiceDetail, error) {
	if search.ServiceName == "" {
		info, err := repo.Get(c, search)
		if err != nil {
			return nil, err
		}
		search = info
	}

	httpRule, err := Get(c, repo.db, &model.HttpRule{ServiceID: search.ID})
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	tcpRule, err := Get(c, repo.db, &model.TcpRule{ServiceID: search.ID})
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	grpcRule, err := Get(c, repo.db, &model.GrpcRule{ServiceID: search.ID})
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	accessControl, err := Get(c, repo.db, &model.AccessControl{ServiceID: search.ID})
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	loadBalance, err := Get(c, repo.db, &model.LoadBalance{ServiceID: search.ID})
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	detail := &model.ServiceDetail{
		Info:          search,
		HTTPRule:      httpRule,
		TCPRule:       tcpRule,
		GRPCRule:      grpcRule,
		LoadBalance:   loadBalance,
		AccessControl: accessControl,
	}
	return detail, nil
}
