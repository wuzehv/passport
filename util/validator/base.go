// 定义控制器公共的结构体和校验方法

package validator

type Pager struct {
	Page     int `json:"page" minimum:"1"`      // 页码
	PageSize int `json:"page_size" minimum:"1"` // 每页数量
}
