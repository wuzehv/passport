package base

import (
	"github.com/gin-gonic/gin"
	"github.com/wuzehv/passport/service/db"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type Model struct {
	Id        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"-"`
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
)

type Base interface {
	Base()
}

type Param struct {
	Page     int
	PageSize int
	Table    interface{}
	Where    string // 仅支持命名参数, e.g. name = @name
	Bind     map[string]interface{}
	Order    string
}

type PaginateResponse struct {
	Page     int                      `json:"page"`
	PageSize int                      `json:"page_size"`
	Total    int64                    `json:"total"`
	Items    []map[string]interface{} `json:"items"`
}

func Paginate(params *Param) (*PaginateResponse, error) {
	if params.Page <= 0 {
		params.Page = 1
	}

	switch {
	case params.PageSize > MaxPageSize:
		params.PageSize = MaxPageSize
	case params.PageSize <= 0:
		params.PageSize = DefaultPageSize
	}

	res := PaginateResponse{
		Page:     params.Page,
		PageSize: params.PageSize,
		Items:    []map[string]interface{}{},
	}

	dc := db.Db.Model(params.Table)
	dl := db.Db.Model(params.Table).Scopes(PaginateScopes(params.Page, params.PageSize))

	if params.Where != "" {
		if err := dc.Where(params.Where, params.Bind).Count(&res.Total).Error; err != nil {
			return nil, err
		}

		if err := dl.Where(params.Where, params.Bind).Find(&res.Items).Error; err != nil {
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

func Paginate2(c *gin.Context, params *Param) (*PaginateResponse, error) {
	page, _ := strconv.Atoi(c.PostForm("page"))
	pageSize, _ := strconv.Atoi(c.PostForm("page_size"))
	params.Page = page
	params.PageSize = pageSize
	return Paginate(params)
}

func PaginateScopes(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
