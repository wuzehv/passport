package model

import (
	"github.com/gin-gonic/gin"
	"github.com/wuzehv/passport/service/db"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type Model struct {
	Id        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"-"`
}

const (
	// 正常
	StatusNormal = iota + 1
	// 已禁用
	StatusDisabled
)

const (
	DefaultPageSize = 10
	MaxPageSize     = 100

	GenderMale   = 1
	GenderFemale = 2
)

type Param struct {
	Page       int
	PageSize   int
	Table      interface{}
	Where      string // 仅支持命名参数, e.g. name = @name
	Bind       map[string]interface{}
	OrderField string
	OrderType  string
	JoinObj    *gorm.DB
}

type PaginateResponse struct {
	Page     int                      `json:"page"`
	PageSize int                      `json:"page_size"`
	Total    int64                    `json:"total"`
	Items    []map[string]interface{} `json:"items"`
}

func PaginateContext(c *gin.Context, param *Param) (*PaginateResponse, error) {
	page, _ := strconv.Atoi(c.Query("page"))
	pageSize, _ := strconv.Atoi(c.Query("page_size"))
	param.Page = page
	param.PageSize = pageSize
	return Paginate(param)
}

func PaginateScopes(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func Paginate(param *Param) (*PaginateResponse, error) {
	if param.Page <= 0 {
		param.Page = 1
	}

	switch {
	case param.PageSize > MaxPageSize:
		param.PageSize = MaxPageSize
	case param.PageSize <= 0:
		param.PageSize = DefaultPageSize
	}

	res := PaginateResponse{
		Page:     param.Page,
		PageSize: param.PageSize,
		Items:    []map[string]interface{}{},
	}

	if param.OrderField == "" {
		param.OrderField = "id"
	}

	if param.OrderType == "" {
		param.OrderType = "desc"
	}

	order := param.OrderField + " " + param.OrderType

	var dc, dl *gorm.DB
	if param.JoinObj == nil {
		dc = db.Db.Model(param.Table)
		// 注意这里不能用dc进行调用
		dl = db.Db.Model(param.Table).Order(order).Scopes(PaginateScopes(param.Page, param.PageSize))
	} else {
		dc = param.JoinObj
		dl = param.JoinObj.Order(order).Scopes(PaginateScopes(param.Page, param.PageSize))
	}

	if param.Where != "" && len(param.Bind) > 0 {
		if err := dc.Where(param.Where, param.Bind).Count(&res.Total).Error; err != nil {
			return nil, err
		}

		if err := dl.Where(param.Where, param.Bind).Find(&res.Items).Error; err != nil {
			return nil, err
		}
	} else {
		if err := dc.Count(&res.Total).Error; err != nil {
			return nil, err
		}

		if err := dl.Find(&res.Items).Error; err != nil {
			return nil, err
		}
	}

	return &res, nil
}
