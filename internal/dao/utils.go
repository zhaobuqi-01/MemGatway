package dao

import (
	"fmt"
	"gateway/internal/dto"
	"gateway/pkg/log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Get查询单条数据
//
// Get retrieves a single record of type T based on the given search criteria.
// The input is a gin.Context, a gorm.DB instance, and a search object of type T.
// The output is a pointer to the found object of type T, and an error if the operation failed.
//
// Usage example:
//
//	search := &User{ID: 1}
//	user, err := Get(c, db, search)
//	if err != nil {
//	    // handle error
//	}
//	fmt.Println(user) // Output: User{ID: 1, Name: "John Doe", ...}
func Get[T Model](c *gin.Context, db *gorm.DB, search *T) (*T, error) {
	// log记录查询信息
	log.Info(fmt.Sprintf("searching for %v", search), zap.String("trace_id", c.GetString("TraceID")))

	var out T
	result := db.Where(search).First(&out)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			log.Error(fmt.Sprintf(" %v not found ", search), zap.Error(result.Error), zap.String("trace_id", c.GetString("TraceID")))
			return nil, result.Error
		}

		log.Error(fmt.Sprintf("error retrieving :%v ", search), zap.Error(result.Error), zap.String("trace_id", c.GetString("TraceID")))
		return nil, result.Error
	}

	log.Info(fmt.Sprintf(" %v was found", search), zap.String("trace_id", c.GetString("TraceID")))
	return &out, nil
}

// update更新对象
//
// Update updates an existing record of type T in the database.
// The input is a gin.Context, a gorm.DB instance, and a data object of type T to be updated.
// The output is an error if the operation failed.
//
// Usage example:
//
//	data := &User{ID: 1, Name: "John Doe", Age: 30}
//	err := Update(c, db, data)
//	if err != nil {
//	    // handle error
//	}
//	fmt.Println("User updated") // Output: User updated
func Update[T Model](c *gin.Context, db *gorm.DB, data *T) error {
	// log记录更新信息
	log.Info(fmt.Sprintf("updating %v", data), zap.String("trace_id", c.GetString("TraceID")))

	if err := db.Model(data).Updates(data).Error; err != nil {
		log.Error(fmt.Sprintf("error updating : %v ", data), zap.Error(err), zap.String("trace_id", c.GetString("TraceID")))
		return err
	}

	log.Info(fmt.Sprintf("%v updated", data), zap.String("trace_id", c.GetString("TraceID")))
	return nil
}

// Save保存对象
//
// Save saves a new record of type T in the database.
// The input is a gin.Context, a gorm.DB instance, and a data object of type T to be saved.
// The output is an error if the operation failed.
//
// Usage example:
//
//	data := &User{Name: "John Doe", Age: 25}
//	err := Save(c, db, data)
//	if err != nil {
//	    // handle error
//	}
//	fmt.Println("User saved") // Output: User saved
func Save[T Model](c *gin.Context, db *gorm.DB, data *T) error {
	// log记录保存信息
	log.Info(fmt.Sprintf("saving %v", data), zap.String("trace_id", c.GetString("TraceID")))

	if err := db.Save(data).Error; err != nil {
		log.Error(fmt.Sprintf("error saving : %v ", data), zap.Error(err), zap.String("trace_id", c.GetString("TraceID")))
		return err
	}

	log.Info(fmt.Sprintf("%v Saved", data), zap.String("trace_id", c.GetString("TraceID")))
	return nil
}

// delete删除对象
//
// Delete deletes an existing record of type T from the database.
// The input is a gin.Context, a gorm.DB instance, and a data object of type T to be deleted.
// The output is an error if the operation failed.
//
// Usage example:
//
//	data := &User{ID: 1}
//	err := Delete(c, db, data)
//	if err != nil {
//	    // handle error
//	}
//	fmt.Println("User deleted") // Output: User deleted
func Delete[T Model](c *gin.Context, db *gorm.DB, data *T) error {
	// log记录删除信息
	log.Info(fmt.Sprintf("deleting %v", data), zap.String("trace_id", c.GetString("TraceID")))

	if err := db.Delete(data).Error; err != nil {
		log.Error(fmt.Sprintf("error deleting : %v ", data), zap.Error(err), zap.String("trace_id", c.GetString("TraceID")))
		return err
	}

	log.Info(fmt.Sprintf("%v deleted", data), zap.String("trace_id", c.GetString("TraceID")))
	return nil
}

// ListByServiceID 根据服务ID查询列表
//
// ListByServiceID retrieves a list of records of type T based on the given serviceID.
// The input is a gin.Context, a gorm.DB instance, and an int64 serviceID.
// The output is a slice of records of type T, the total number of records, and an error if the operation failed.
//
// Usage example:
//
//	serviceID := int64(1)
//	services, count, err := ListByServiceID(c, db, serviceID)
//	if err != nil {
//	    // handle error
//	}
//	fmt.Println(services) // Output: [Service{ID: 1, Name: "Service1", ...} ...]
//	fmt.Println(count)    // Output: 5
func ListByServiceID[T Model](c *gin.Context, db *gorm.DB, serviceID int64) ([]T, int64, error) {
	// log记录查询信息
	log.Info(fmt.Sprintf("searching for %v", serviceID), zap.String("trace_id", c.GetString("TraceID")))

	var list []T
	var count int64
	query := db.Select("*")
	query = query.Where("service_id=?", serviceID)
	err := query.Order("id desc").Find(&list).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Error(fmt.Sprintf("error retrieving :%v ", serviceID), zap.Error(err), zap.String("trace_id", c.GetString("TraceID")))
		return nil, 0, err
	}
	errCount := query.Count(&count).Error
	if errCount != nil {
		log.Error(fmt.Sprintf("error retrieving :%v ", serviceID), zap.Error(errCount), zap.String("trace_id", c.GetString("TraceID")))
		return nil, 0, err
	}

	log.Info(fmt.Sprintf("%v was found", serviceID), zap.String("trace_id", c.GetString("TraceID")))
	return list, count, nil
}

// PageList 分页查询
//
// PageList retrieves a list of records of type T using pagination and query conditions.
// The input is a gin.Context, a gorm.DB instance, a slice of query conditions, and PageNo and PageSize for pagination.
// The output is a slice of records of type T, the total number of records, and an error if the operation failed.
//
// Usage example:
//
//	queryConditions := []func(db *gorm.DB) *gorm.DB{
//	    func(db *gorm.DB) *gorm.DB {
//	        return db.Where("age > ?", 18)
//	    },
//	}
//	users, total, err := PageList(c, db, queryConditions, 1, 10)
//	if err != nil {
//	    // handle error
//	}
//	fmt.Println(users) // Output: [User{ID: 2, Name: "Jane Doe", Age: 28} ...]
//	fmt.Println(total) // Output: 42
func PageList[T Model](c *gin.Context, db *gorm.DB, queryConditions []func(db *gorm.DB) *gorm.DB, PageNo, PageSize int) ([]T, int64, error) {
	// log记录查询信息
	log.Info(fmt.Sprintf("searching for %v", queryConditions), zap.String("trace_id", c.GetString("TraceID")))

	total := int64(0)
	list := []T{}
	offset := (PageNo - 1) * PageSize

	query := db.Where("is_delete=0")
	for _, condition := range queryConditions {
		query = condition(query)
	}
	if err := query.Limit(PageSize).Offset(offset).Order("id desc").Find(&list).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Error(fmt.Sprintf("error retrieving :%v ", query), zap.Error(err), zap.String("trace_id", c.GetString("TraceID")))
		return nil, 0, err
	}
	query.Limit(PageSize).Offset(offset).Count(&total)

	log.Info(fmt.Sprintf("%v was found", query), zap.String("trace_id", c.GetString("TraceID")))
	return list, total, nil
}

// GetServiceDetail 获取服务详情
//
// GetServiceDetail retrieves a detailed record of a service.
// The input is a gin.Context, a gorm.DB instance, and a search object of type ServiceInfo.
// The output is a pointer to the found ServiceDetail object, and an error if the operation failed.
//
// Usage example:
//
//	search := &ServiceInfo{ID: 1}
//	detail, err := GetServiceDetail(c, db, search)
//	if err != nil {
//	    // handle error
//	}
//	fmt.Println(detail) // Output: ServiceDetail{Info: {ID: 1, ...}, HTTPRule: {...}, ...}
func GetServiceDetail(c *gin.Context, db *gorm.DB, search *ServiceInfo) (*ServiceDetail, error) {
	// log记录查询信息
	log.Info(fmt.Sprintf("searching for %v", search), zap.String("trace_id", c.GetString("TraceID")))

	if search.ServiceName == "" {
		info, err := Get(c, db, search)
		if err != nil {
			return nil, err
		}
		search = info
	}

	httpRule, err := Get(c, db, &HttpRule{ServiceID: search.ID})
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	tcpRule, err := Get(c, db, &TcpRule{ServiceID: search.ID})
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	grpcRule, err := Get(c, db, &GrpcRule{ServiceID: search.ID})
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	accessControl, err := Get(c, db, &AccessControl{ServiceID: search.ID})
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	loadBalance, err := Get(c, db, &LoadBalance{ServiceID: search.ID})
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	detail := &ServiceDetail{
		Info:          search,
		HTTPRule:      httpRule,
		TCPRule:       tcpRule,
		GRPCRule:      grpcRule,
		LoadBalance:   loadBalance,
		AccessControl: accessControl,
	}

	// log记录成功取到信息
	log.Info(fmt.Sprintf("%v was found", search), zap.String("trace_id", c.GetString("TraceID")))
	return detail, nil
}

// GroupByLoadType 按照负载类型分组
//
// GroupByLoadType retrieves a list of records grouped by load type.
// The input is a gin.Context and a gorm.DB instance.
// The output is a slice of dto.DashServiceStatItemOutput objects, and an error if the operation failed.
//
// Usage example:
//
//	stats, err := GroupByLoadType(c, db)
//	if err != nil {
//	    // handle error
//	}
//	fmt.Println(stats) // Output: [DashServiceStatItemOutput{LoadType: "http", Value: 5} ...}]
func GroupByLoadType(c *gin.Context, tx *gorm.DB) ([]dto.DashServiceStatItemOutput, error) {
	// log记录开始查询
	log.Info("searching for group by load type", zap.String("trace_id", c.GetString("TraceID")))

	list := []dto.DashServiceStatItemOutput{}
	if err := tx.Table(ServiceInfo{}.TableName()).Where("is_delete=0").Select("load_type, count(*) as value").Group("load_type").Scan(&list).Error; err != nil {
		log.Error("error retrieving", zap.Error(err), zap.String("trace_id", c.GetString("TraceID")))
		return nil, err
	}

	// log记录成功取到信息
	log.Info("group by load type was found", zap.String("trace_id", c.GetString("TraceID")))
	return list, nil
}
