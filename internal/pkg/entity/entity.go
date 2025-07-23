package entity

import (
	"katydid-mp-account/pkg/data"
	"time"
)

type (
	// Entity 实体基类
	Entity struct {
		//gorm.Model
		ID int64 `json:"id" gorm:"primarykey;comment:主键"`

		State    State      `json:"state" gorm:"default:0;comment:状态"`
		CreateAt time.Time  `json:"createAt" gorm:"autoCreateTime:milli;comment:创建时间"`
		UpdateAt time.Time  `json:"updateAt" gorm:"autoUpdateTime:milli;comment:更新时间"`
		DeleteAt *time.Time `json:"deleteAt;comment:删除时间"` // 删除人可以在 Extra 中设置

		// id
		// index
		// required

		Extra data.KSMap `json:"extra" gorm:"serializer:json;comment:额外信息"` // (!索引+!必需)
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
