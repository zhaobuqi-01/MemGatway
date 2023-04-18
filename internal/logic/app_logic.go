package logic

import (
	"fmt"
	"gateway/internal/dao"
	"gateway/internal/dto"
	"gateway/internal/pkg"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type AppLogic interface {
	AppList(c *gin.Context, params *dto.APPListInput) ([]dto.APPListItemOutput, int64, error)
	AppDetail(c *gin.Context, params *dto.APPDetailInput) (*dao.App, error)
	AppDelete(c *gin.Context, params *dto.APPDetailInput) error
	AppAdd(c *gin.Context, params *dto.APPAddHttpInput) error
	AppUpdate(c *gin.Context, params *dto.APPUpdateHttpInput) error
	AppStat(c *gin.Context, params *dto.APPDetailInput) (*dto.StatisticsOutput, error)
}

type appLogic struct {
	db *gorm.DB
}

func NewAppLogic(tx *gorm.DB) *appLogic {
	return &appLogic{db: tx}
}

func (al *appLogic) AppList(c *gin.Context, params *dto.APPListInput) ([]dto.APPListItemOutput, int64, error) {
	queryConditions := []func(db *gorm.DB) *gorm.DB{
		func(db *gorm.DB) *gorm.DB {
			return db.Where("(name like ? or app_id like ?)", "%"+params.Info+"%", "%"+params.Info+"%")
		},
	}
	list, total, err := dao.PageList[dao.App](c, al.db, queryConditions, params.PageNo, params.PageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("获取分页数据失败")
	}

	outputList := []dto.APPListItemOutput{}
	for _, item := range list {
		// appCounter, err := public.FlowCounterHandler.GetCounter(public.FlowAppPrefix + item.AppID)
		// if err != nil {
		// 	return nil, 0, err
		// }
		outputList = append(outputList, dto.APPListItemOutput{
			ID:       item.ID,
			AppID:    item.AppID,
			Name:     item.Name,
			Secret:   item.Secret,
			WhiteIPS: item.WhiteIPS,
			Qpd:      item.Qpd,
			Qps:      item.Qps,
			RealQpd:  0,
			RealQps:  0,
		})
	}
	return outputList, total, nil
}

func (al *appLogic) AppDetail(c *gin.Context, params *dto.APPDetailInput) (*dao.App, error) {
	search := &dao.App{
		ID: params.ID,
	}
	detail, err := dao.Get(c, al.db, search)
	if err != nil {
		return nil, fmt.Errorf("获取app详情失败")
	}
	return detail, nil
}

func (al *appLogic) AppDelete(c *gin.Context, params *dto.APPDetailInput) error {
	search := &dao.App{
		ID: params.ID,
	}
	info, err := dao.Get(c, al.db, search)
	if err != nil {
		return fmt.Errorf("app不存在")
	}
	info.IsDelete = 1
	if err := dao.Save(c, al.db, info); err != nil {
		return errors.Wrap(err, "删除失败")
	}
	return nil
}

func (al *appLogic) AppAdd(c *gin.Context, params *dto.APPAddHttpInput) error {
	//验证app_id是否被占用
	search := &dao.App{
		AppID: params.AppID,
	}
	if _, err := dao.Get(c, al.db, search); err == nil {
		return fmt.Errorf("app_id已经被占用")
	}

	if params.Secret == "" {
		params.Secret, _ = pkg.HashPassword(params.AppID)
	}
	app := &dao.App{
		Name:     params.Name,
		Secret:   params.Secret,
		WhiteIPS: params.WhiteIPS,
		Qpd:      params.Qpd,
		Qps:      params.Qps,
	}
	if err := dao.Save(c, al.db, app); err != nil {
		return fmt.Errorf("添加app失败")
	}
	return nil
}

func (al *appLogic) AppUpdate(c *gin.Context, params *dto.APPUpdateHttpInput) error {
	search := &dao.App{
		ID: params.ID,
	}
	info, err := dao.Get(c, al.db, search)
	if err != nil {
		return fmt.Errorf("app不存在")
	}
	info.Name = params.Name
	info.WhiteIPS = params.WhiteIPS
	info.Qpd = params.Qpd
	info.Qps = params.Qps
	if err := dao.Save(c, al.db, info); err != nil {
		return fmt.Errorf("app信息更新失败")
	}
	return nil
}

func (al *appLogic) AppStat(c *gin.Context, params *dto.APPDetailInput) (*dto.StatisticsOutput, error) {
	return nil, nil
}
