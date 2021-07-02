// 定义控制器公共的结构体和校验方法

package validator

type Pager struct {
	Page     int `json:"page" binding:"omitempty,gt=0"`      // 页码
	PageSize int `json:"page_size" binding:"omitempty,gt=0"` // 每页数量
}
