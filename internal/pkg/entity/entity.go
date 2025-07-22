package entity

import (
	"katydid-mp-account/pkg/data"
	"time"
)

type (
	// Entity 实体基类
	Entity struct {
		//gorm.Model
		ID int64 `json:"id" gorm:"primarykey"` // 主键

		State    State      `json:"state" gorm:"default:0"`               // 状态
		CreateAt time.Time  `json:"createAt" gorm:"autoCreateTime:milli"` // 创建时间
		UpdateAt time.Time  `json:"updateAt" gorm:"autoUpdateTime:milli"` // 更新时间
		DeleteAt *time.Time `json:"deleteAt"`                             // 删除时间

		// id
		// index
		// required

		Extra data.KSMap `json:"extra" gorm:"serializer:json"` // 额外信息 (!索引/!必需)
	}
)

func NewEntity(id int64) *Entity {
	return &Entity{
		ID:    id,
		Extra: make(data.KSMap),
	}
}

const (
	extraKeyAdminNote = "adminNote" // 管理员备注
)

func (e *Entity) GetAdminNote() (string, bool) {
	return e.Extra.GetString(extraKeyAdminNote)
}

func (e *Entity) SetAdminNote(adminNote *string) {
	e.Extra.SetString(extraKeyAdminNote, adminNote)
}
