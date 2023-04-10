package repository

import (
	"gateway/internal/dto"
	"gateway/internal/entity"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ServiceInfo interface {
	// Getter[entity.ServiceInfo]
	Updater[entity.ServiceInfo]
	PageList(c *gin.Context, param *dto.ServiceListInput) ([]entity.ServiceInfo, int64, error)
	ServiceDetail(c *gin.Context, search *entity.ServiceInfo) (*entity.ServiceDetail, error)
}

type serviceInfoRepo struct {
	db *gorm.DB
}

func NewServiceInfoRepo(db *gorm.DB) ServiceInfo {
	return &serviceInfoRepo{
		db: db,
	}
}

func (repo *serviceInfoRepo) Get(c *gin.Context, search *entity.ServiceInfo) (*entity.ServiceInfo, error) {
	return Get(c, repo.db, search)
}

func (repo *serviceInfoRepo) Update(c *gin.Context, data *entity.ServiceInfo) error {
	return Update(c, repo.db, data)
}

func (repo *serviceInfoRepo) PageList(c *gin.Context, param *dto.ServiceListInput) ([]entity.ServiceInfo, int64, error) {
	total := int64(0)
	list := []entity.ServiceInfo{}
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

func (repo *serviceInfoRepo) ServiceDetail(c *gin.Context, search *entity.ServiceInfo) (*entity.ServiceDetail, error) {
	if search.ServiceName == "" {
		info, err := repo.Get(c, search)
		if err != nil {
			return nil, err
		}
		search = info
	}

	httpRule, err := Get(c, repo.db, &entity.HttpRule{ServiceID: search.ID})
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	tcpRule, err := Get(c, repo.db, &entity.TcpRule{ServiceID: search.ID})
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	grpcRule, err := Get(c, repo.db, &entity.GrpcRule{ServiceID: search.ID})
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	accessControl, err := Get(c, repo.db, &entity.AccessControl{ServiceID: search.ID})
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	loadBalance, err := Get(c, repo.db, &entity.LoadBalance{ServiceID: search.ID})
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	detail := &entity.ServiceDetail{
		Info:          search,
		HTTPRule:      httpRule,
		TCPRule:       tcpRule,
		GRPCRule:      grpcRule,
		LoadBalance:   loadBalance,
		AccessControl: accessControl,
	}
	return detail, nil
}
