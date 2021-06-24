package model

type Action struct {
	Model
	Url    string `gorm:"unique" json:"url"`
	Remark string `gorm:"not null" json:"remark"`
}

func (u *Action) Base() {}
