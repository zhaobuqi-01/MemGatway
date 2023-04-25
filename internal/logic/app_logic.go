package logic

import (
	"fmt"
	"gateway/internal/dao"
	"gateway/internal/dto"
	"gateway/internal/pkg"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AppLogic是应用程序逻辑的接口
type AppLogic interface {
	AppList(c *gin.Context, params *dto.APPListInput) ([]dto.APPListItemOutput, int64, error)
	AppDetail(c *gin.Context, params *dto.APPDetailInput) (*dao.App, error)
	AppDelete(c *gin.Context, params *dto.APPDetailInput) error
	AppAdd(c *gin.Context, params *dto.APPAddHttpInput) error
	AppUpdate(c *gin.Context, params *dto.APPUpdateHttpInput) error
}

// appLogic是实现AppLogic接口的结构体
type appLogic struct {
	db *gorm.DB
}

// NewAppLogic创建一个新的appLogic实例
func NewAppLogic(tx *gorm.DB) *appLogic {
	return &appLogic{db: tx}
}

// AppList返回应用程序列表
func (al *appLogic) AppList(c *gin.Context, params *dto.APPListInput) ([]dto.APPListItemOutput, int64, error) {
	// 构造查询条件
	queryConditions := []func(db *gorm.DB) *gorm.DB{
		func(db *gorm.DB) *gorm.DB {
			return db.Where("(name like ? or app_id like ?)", "%"+params.Info+"%", "%"+params.Info+"%")
		},
	}
	// 使用dao中的PageList方法获取分页的应用程序列表
	list, total, err := dao.PageList[dao.App](c, al.db, queryConditions, params.PageNo, params.PageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get all app data")
	}
	// 转换为输出DTO对象
	outputList := []dto.APPListItemOutput{}
	for _, item := range list {
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

// AppDetail返回应用程序的详细信息
func (al *appLogic) AppDetail(c *gin.Context, params *dto.APPDetailInput) (*dao.App, error) {
	// 使用dao中的Get方法获取指定ID的应用程序信息
	search := &dao.App{
		ID: params.ID,
	}
	detail, err := dao.Get(c, al.db, search)
	if err != nil {
		return nil, fmt.Errorf("failed to get app details")
	}
	return detail, nil
}

// AppDelete删除指定的应用程序
func (al *appLogic) AppDelete(c *gin.Context, params *dto.APPDetailInput) error {
	// 使用dao中的Get方法获取指定ID的应用程序信息
	search := &dao.App{
		// ID: params.ID,
	}
	info, err := dao.Get(c, al.db, search)
	if err != nil {
		return fmt.Errorf("app not found")
	}
	// 将应用程序标记为已删除
	info.IsDelete = 1
	if err := dao.Update(c, al.db, info); err != nil {
		return fmt.Errorf("failed to delete app")
	}
	return nil
}

// AppAdd添加一个新的应用程序
func (al *appLogic) AppAdd(c *gin.Context, params *dto.APPAddHttpInput) error {
	// 检查应用程序ID是否已经存在
	search := &dao.App{
		AppID: params.AppID,
	}
	if _, err := dao.Get(c, al.db, search); err == nil {
		return fmt.Errorf("app ID is already taken")
	}
	// 如果没有指定密钥，则使用默认的appID哈希值作为密钥
	if params.Secret == "" {
		params.Secret, _ = pkg.HashPassword(params.AppID)
	}
	// 创建新的应用程序对象
	app := &dao.App{
		AppID:    params.AppID,
		Name:     params.Name,
		Secret:   params.Secret,
		WhiteIPS: params.WhiteIPS,
		Qpd:      params.Qpd,
		Qps:      params.Qps,
	}
	// 保存应用程序对象
	if err := dao.Save(c, al.db, app); err != nil {
		return fmt.Errorf("failed to add app")
	}
	return nil
}

// AppUpdate更新指定应用程序的信息
func (al *appLogic) AppUpdate(c *gin.Context, params *dto.APPUpdateHttpInput) error {
	// 使用dao中的Get方法获取指定ID的应用程序信息
	search := &dao.App{
		ID: params.ID,
	}
	info, err := dao.Get(c, al.db, search)
	if err != nil {
		return fmt.Errorf("app not found")
	}
	// 更新应用程序信息
	info.Name = params.Name
	info.WhiteIPS = params.WhiteIPS
	info.Qpd = params.Qpd
	info.Qps = params.Qps
	if err := dao.Save(c, al.db, info); err != nil {
		return fmt.Errorf("failed to update app information")
	}
	return nil
}
