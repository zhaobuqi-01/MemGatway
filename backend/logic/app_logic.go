package logic

import (
	"fmt"
	"gateway/backend/dto"
	"gateway/dao"
	"gateway/enity"
	"gateway/globals"
	"gateway/pkg/database/mysql"
	"gateway/pkg/log"
	"gateway/utils"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// AppLogic是应用程序逻辑的接口
type AppLogic interface {
	AppList(c *gin.Context, params *dto.APPListInput) ([]dto.APPListItemOutput, int64, error)
	AppDetail(c *gin.Context, params *dto.APPDetailInput) (*enity.App, error)
	AppDelete(c *gin.Context, params *dto.APPDetailInput) error
	AppAdd(c *gin.Context, params *dto.APPAddHttpInput) error
	AppUpdate(c *gin.Context, params *dto.APPUpdateHttpInput) error
	AppStat(c *gin.Context, params *dto.APPDetailInput) (*dto.StatisticsOutput, error)
}

// appLogic是实现AppLogic接口的结构体
type appLogic struct {
	dao.APP
	db *gorm.DB
}

// NewAppLogic创建一个新的appLogic实例
func NewAppLogic() *appLogic {
	return &appLogic{
		dao.NewApp(),
		mysql.GetDB(),
	}
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
	list, total, err := al.PageList(c, al.db, queryConditions, params.PageNo, params.PageSize)
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
func (al *appLogic) AppDetail(c *gin.Context, params *dto.APPDetailInput) (*enity.App, error) {
	// 使用dao中的Get方法获取指定ID的应用程序信息
	search := &enity.App{
		ID: params.ID,
	}
	detail, err := al.Get(c, al.db, search)
	if err != nil {
		return nil, fmt.Errorf("failed to get app details")
	}
	return detail, nil
}

// AppDelete删除指定的应用程序
func (al *appLogic) AppDelete(c *gin.Context, params *dto.APPDetailInput) error {
	// 使用dao中的Get方法获取指定ID的应用程序信息
	search := &enity.App{
		ID: params.ID,
	}
	info, err := al.Get(c, al.db, search)
	if err != nil {
		return fmt.Errorf("app not found")
	}
	// 将应用程序标记为已删除
	info.IsDelete = 1
	if err := al.Save(c, al.db, info); err != nil {
		return fmt.Errorf("failed to delete app")
	}

	// Publish data change message
	message := &globals.DataChangeMessage{
		Type:      "service",
		Payload:   info.AppID,
		Operation: globals.DataDelete,
	}
	if err := globals.MessageQueue.Publish(globals.DataChange, message); err != nil {
		log.Error("error publishing message", zap.Error(err), zap.String("trace_id", c.GetString("TraceID")))
		return fmt.Errorf("failed to publish save message")
	}
	log.Info("published save message successfully", zap.Any("data", params), zap.String("trace_id", c.GetString("TraceID")))

	return nil
}

// AppAdd添加一个新的应用程序
func (al *appLogic) AppAdd(c *gin.Context, params *dto.APPAddHttpInput) error {
	// 检查应用程序ID是否已经存在
	search := &enity.App{
		AppID: params.AppID,
	}
	if _, err := al.Get(c, al.db, search); err == nil {
		return fmt.Errorf("app ID is already taken")
	}
	// 如果没有指定密钥，则使用默认的appID哈希值作为密钥
	if params.Secret == "" {
		params.Secret, _ = utils.HashPassword(params.AppID)
	}
	// 创建新的应用程序对象
	app := &enity.App{
		AppID:    params.AppID,
		Name:     params.Name,
		Secret:   params.Secret,
		WhiteIPS: params.WhiteIPS,
		Qpd:      params.Qpd,
		Qps:      params.Qps,
	}
	// 保存应用程序对象
	if err := al.Save(c, al.db, app); err != nil {
		return fmt.Errorf("failed to add app")
	}
	// Publish data change message
	message := &globals.DataChangeMessage{
		Type:      "service",
		Payload:   params.AppID,
		Operation: globals.DataInsert,
	}
	if err := globals.MessageQueue.Publish(globals.DataChange, message); err != nil {
		log.Error("error publishing message", zap.Error(err), zap.String("trace_id", c.GetString("TraceID")))
		return fmt.Errorf("failed to publish save message")
	}
	log.Info("published save message successfully", zap.Any("data", params), zap.String("trace_id", c.GetString("TraceID")))

	return nil
}

// AppUpdate更新指定应用程序的信息
func (al *appLogic) AppUpdate(c *gin.Context, params *dto.APPUpdateHttpInput) error {
	// 使用dao中的Get方法获取指定ID的应用程序信息
	search := &enity.App{
		ID: params.ID,
	}
	info, err := al.Get(c, al.db, search)
	if err != nil {
		return fmt.Errorf("app not found")
	}
	// 更新应用程序信息
	info.Name = params.Name
	info.WhiteIPS = params.WhiteIPS
	info.Qpd = params.Qpd
	info.Qps = params.Qps
	if err := al.Save(c, al.db, info); err != nil {
		return fmt.Errorf("failed to Save app information")
	}

	// Publish data change message
	message := &globals.DataChangeMessage{
		Type:      "service",
		Payload:   params.AppID,
		Operation: globals.DataUpdate,
	}
	if err := globals.MessageQueue.Publish(globals.DataChange, message); err != nil {
		log.Error("error publishing message", zap.Error(err), zap.String("trace_id", c.GetString("TraceID")))
		return fmt.Errorf("failed to publish save message")
	}
	log.Info("published save message successfully", zap.Any("data", params), zap.String("trace_id", c.GetString("TraceID")))

	return nil
}

// AppStat返回指定应用程序的统计信息
func (al *appLogic) AppStat(c *gin.Context, params *dto.APPDetailInput) (*dto.StatisticsOutput, error) {
	// 使用dao中的Get方法获取指定ID的应用程序信息
	search := &enity.App{
		ID: params.ID,
	}
	detail, err := al.Get(c, al.db, search)
	if err != nil {
		return nil, fmt.Errorf("app not found")
	}
	counter, err := globals.FlowCounter.GetCounter(detail.AppID)
	if err != nil {
		return nil, fmt.Errorf("failed to get app flow counter")
	}
	todayList := []int64{}
	currentTime := time.Now()
	for i := 0; i <= currentTime.Hour(); i++ {
		dateTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), i, 0, 0, 0, time.UTC)
		hourData, _ := counter.GetHourData(dateTime)
		todayList = append(todayList, hourData)
	}

	yesterdayList := []int64{}
	yesterTime := currentTime.Add(-1 * time.Duration(time.Hour*24))
	for i := 0; i <= 23; i++ {
		dateTime := time.Date(yesterTime.Year(), yesterTime.Month(), yesterTime.Day(), i, 0, 0, 0, time.UTC)
		hourData, _ := counter.GetHourData(dateTime)
		yesterdayList = append(yesterdayList, hourData)
	}

	out := &dto.StatisticsOutput{
		Today:     todayList,
		Yesterday: yesterdayList,
	}
	return out, nil
}
