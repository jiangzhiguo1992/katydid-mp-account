package model

import (
	"katydid-mp-account/pkg/field"
	"katydid-mp-account/pkg/valid"
	"time"
)

type (
	// Base 实体基类
	Base struct {
		//gorm.Model
		ID int64 `json:"id" gorm:"primarykey;comment:主键"`

		State    field.State `json:"state" gorm:"default:0;comment:状态"`
		CreateAt time.Time   `json:"createAt" gorm:"autoCreateTime:milli;comment:创建时间"`
		UpdateAt time.Time   `json:"updateAt" gorm:"autoUpdateTime:milli;comment:更新时间"`
		DeleteAt *time.Time  `json:"deleteAt;comment:删除时间"` // 删除人可以在 Extra 中设置

		// id
		// index
		// required

		Extra field.KMap `json:"extra" gorm:"serializer:json;comment:额外信息"` // (!索引+!必需)
	}
)

func NewBase(id int64) *Base {
	return &Base{
		ID:    id,
		Extra: make(field.KMap),
	}
}

const (
	extraKeyAdminNote = "adminNote" // 管理员备注
)

func (b *Base) GetAdminNote() (string, bool) {
	return b.Extra.GetString(extraKeyAdminNote)
}

func (b *Base) SetAdminNote(adminNote *string) {
	b.Extra.SetString(extraKeyAdminNote, adminNote)
}

// ValidFieldRules 字段验证规则
func (b *Base) ValidFieldRules() valid.FieldValidRules {
	return valid.FieldValidRules{
		valid.SceneAll:    valid.FieldValidRule{},
		valid.SceneBind:   valid.FieldValidRule{},
		valid.SceneSave:   valid.FieldValidRule{},
		valid.SceneInsert: valid.FieldValidRule{},
		valid.SceneUpdate: valid.FieldValidRule{},
		valid.SceneQuery:  valid.FieldValidRule{},
		valid.SceneReturn: valid.FieldValidRule{},
		valid.SceneCustom: valid.FieldValidRule{},
	}
}

// ValidExtraRules KMap/Extra验证规则
func (b *Base) ValidExtraRules() (field.KMap, valid.ExtraValidRules) {
	return b.Extra, valid.ExtraValidRules{
		valid.SceneAll: map[valid.Tag]valid.ExtraValidRuleInfo{
			// 管理员备注 (0-10000)
			extraKeyAdminNote: {
				Field: extraKeyAdminNote,
				ValidFn: func(value any) bool {
					if _, ok := value.(string); !ok {
						return false
					}
					return len(value.(string)) <= 10_000
				},
			},
		},
	}
}

// ValidStructRules 结构体验证规则
func (b *Base) ValidStructRules(scene valid.Scene, fn valid.FuncReportError) {
	switch scene {
	case valid.SceneAll:
		if b.CreateAt.Before(b.UpdateAt) {
			fn(b.CreateAt, "CreateAt", valid.TagCheck, "")
		}
	}
}

// ValidLocalizeRules 本地化验证规则
func (b *Base) ValidLocalizeRules() valid.LocalizeValidRules {
	return valid.LocalizeValidRules{
		valid.SceneAll: valid.LocalizeValidRule{
			Rule1: map[valid.Tag]map[valid.FieldName]valid.LocalizeValidRuleParam{
				valid.TagRequired: {},
				valid.TagFormat:   {},
				valid.TagRange:    {},
				valid.TagCheck: {
					"CreateAt": {"check_create_at_err", false, nil},
					"DeleteAt": {"check_delete_at_err", false, nil},
					"DeleteBy": {"check_delete_by_err", false, nil},
				},
			}, Rule2: map[valid.Tag]valid.LocalizeValidRuleParam{
				extraKeyAdminNote: {"format_admin_note_err", false, nil},
			},
		},
	}
}
