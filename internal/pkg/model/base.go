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
		ID    field.ID    `json:"id" gorm:"primarykey;comment:主键"`
		State field.State `json:"state" gorm:"default:0;comment:状态"`

		CreateAt time.Time  `json:"createAt" gorm:"autoCreateTime:milli;comment:创建时间"`
		UpdateAt time.Time  `json:"updateAt" gorm:"autoUpdateTime:milli;comment:更新时间"`
		DeleteAt *time.Time `json:"deleteAt;comment:删除时间"` // 删除人可以在 Extra 中设置

		// id
		// index
		// required

		Extra field.KMap `json:"extra" gorm:"serializer:json;comment:额外信息"` // (!索引+!必需)
	}
)

func NewBase(id field.ID) Base {
	return Base{
		ID:    id,
		State: field.StateInit,
		// ...times
		Extra: make(field.KMap),
	}
}

const (
	extKeyAdminNote = "adminNote" // 管理员备注
)

func (b *Base) GetAdminNote() (string, bool) {
	return b.Extra.GetString(extKeyAdminNote)
}

func (b *Base) SetAdminNote(adminNote *string) {
	b.Extra.SetString(extKeyAdminNote, adminNote)
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

// ValidExtraRules KMap/Extra验证规则 TODO:GG 父类验证完，这里可以执行吗？
func (b *Base) ValidExtraRules() (field.KMap, valid.ExtraValidRules) {
	return b.Extra, valid.ExtraValidRules{
		valid.SceneAll: map[valid.Tag]valid.ExtraValidRuleInfo{
			// 管理员备注 (0-10000)
			extKeyAdminNote: {
				Field: extKeyAdminNote,
				ValidFn: func(value any) bool {
					if _, ok := value.(string); !ok {
						return false
					}
					return len(value.(string)) <= 10_000
				},
			},
		},
		valid.SceneBind:   map[valid.Tag]valid.ExtraValidRuleInfo{},
		valid.SceneSave:   map[valid.Tag]valid.ExtraValidRuleInfo{},
		valid.SceneInsert: map[valid.Tag]valid.ExtraValidRuleInfo{},
		valid.SceneUpdate: map[valid.Tag]valid.ExtraValidRuleInfo{},
		valid.SceneQuery:  map[valid.Tag]valid.ExtraValidRuleInfo{},
		valid.SceneReturn: map[valid.Tag]valid.ExtraValidRuleInfo{},
		valid.SceneCustom: map[valid.Tag]valid.ExtraValidRuleInfo{},
	}
}

// ValidStructRules 结构体验证规则
func (b *Base) ValidStructRules(scene valid.Scene, fn valid.FuncReportError) {
	switch scene {
	case valid.SceneAll:
	case valid.SceneBind:
	case valid.SceneSave:
		// TODO:GG 这里检查是不是多余了?
		if b.CreateAt.After(b.UpdateAt) {
			if b.UpdateAt.Unix() == 0 {
				fn(b.UpdateAt, "UpdateAt", valid.TagCheck, "")
			} else {
				fn(b.CreateAt, "CreateAt", valid.TagCheck, "")
			}
		}
	case valid.SceneInsert:
	case valid.SceneUpdate:
	case valid.SceneQuery:
	case valid.SceneReturn:
	case valid.SceneCustom:
	default:
		return
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
					"UpdateAt": {"check_update_at_err", false, nil},
					"DeleteAt": {"check_delete_at_err", false, nil},
				},
			},
			Rule2: map[valid.Tag]valid.LocalizeValidRuleParam{
				extKeyAdminNote: {"format_admin_note_err", false, nil},
			},
		},
		valid.SceneBind: valid.LocalizeValidRule{
			Rule1: map[valid.Tag]map[valid.FieldName]valid.LocalizeValidRuleParam{
				valid.TagRequired: {},
				valid.TagFormat:   {},
				valid.TagRange:    {},
				valid.TagCheck:    {},
			},
			Rule2: map[valid.Tag]valid.LocalizeValidRuleParam{},
		},
		valid.SceneSave: valid.LocalizeValidRule{
			Rule1: map[valid.Tag]map[valid.FieldName]valid.LocalizeValidRuleParam{
				valid.TagRequired: {},
				valid.TagFormat:   {},
				valid.TagRange:    {},
				valid.TagCheck:    {},
			},
			Rule2: map[valid.Tag]valid.LocalizeValidRuleParam{},
		},
		valid.SceneInsert: valid.LocalizeValidRule{
			Rule1: map[valid.Tag]map[valid.FieldName]valid.LocalizeValidRuleParam{
				valid.TagRequired: {},
				valid.TagFormat:   {},
				valid.TagRange:    {},
				valid.TagCheck:    {},
			},
			Rule2: map[valid.Tag]valid.LocalizeValidRuleParam{},
		},
		valid.SceneUpdate: valid.LocalizeValidRule{
			Rule1: map[valid.Tag]map[valid.FieldName]valid.LocalizeValidRuleParam{
				valid.TagRequired: {},
				valid.TagFormat:   {},
				valid.TagRange:    {},
				valid.TagCheck:    {},
			},
			Rule2: map[valid.Tag]valid.LocalizeValidRuleParam{},
		},
		valid.SceneQuery: valid.LocalizeValidRule{
			Rule1: map[valid.Tag]map[valid.FieldName]valid.LocalizeValidRuleParam{
				valid.TagRequired: {},
				valid.TagFormat:   {},
				valid.TagRange:    {},
				valid.TagCheck:    {},
			},
			Rule2: map[valid.Tag]valid.LocalizeValidRuleParam{},
		},
		valid.SceneReturn: valid.LocalizeValidRule{
			Rule1: map[valid.Tag]map[valid.FieldName]valid.LocalizeValidRuleParam{
				valid.TagRequired: {},
				valid.TagFormat:   {},
				valid.TagRange:    {},
				valid.TagCheck:    {},
			},
			Rule2: map[valid.Tag]valid.LocalizeValidRuleParam{},
		},
		valid.SceneCustom: valid.LocalizeValidRule{
			Rule1: map[valid.Tag]map[valid.FieldName]valid.LocalizeValidRuleParam{
				valid.TagRequired: {},
				valid.TagFormat:   {},
				valid.TagRange:    {},
				valid.TagCheck:    {},
			},
			Rule2: map[valid.Tag]valid.LocalizeValidRuleParam{},
		},
	}
}
