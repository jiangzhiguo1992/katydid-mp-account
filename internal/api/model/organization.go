package model

import (
	"katydid-mp-account/internal/pkg/model"
	"katydid-mp-account/pkg/field"
	"katydid-mp-account/pkg/valid"
	"reflect"
	"unicode"
	"unicode/utf8"
)

type (
	Organization struct {
		model.Base

		OwnAccId  field.ID   `json:"ownAccId" validate:"required" gorm:"comment:所属账号"`
		ParentIds []field.ID `json:"parentIds" gorm:"comment:父级组织"`

		IsPrivate bool     `json:"isPrivate" gorm:"comment:是否私有"`
		Kind      uint8    `json:"kind" validate:"kind-check" gorm:"comment:组织类型"`
		Become    uint8    `json:"become" validate:"become-check" gorm:"comment:加入方式"`
		Name      string   `json:"name" validate:"name-format" gorm:"comment:组织名称"`
		Display   string   `json:"display" validate:"display-format" gorm:"comment:组织显示名称"`
		Tags      []string `json:"tags" validate:"tags-format,tag-format" gorm:"comment:组织标签们"`
	}
)

// 类型
const (
	OrgKindGroup   uint8 = 0 // 集团
	OrgKindCompany uint8 = 1 // 公司
	OrgKindStudio  uint8 = 2 // 工作室
	OrgKindTeam    uint8 = 3 // 团队 (虚拟组织)
)

// 加入方式
const (
	OrgBecomeDirect uint8 = 0 // 直接进
	OrgBecomeApply  uint8 = 1 // 申请进
	OrgBecomeInvite uint8 = 2 // 邀请制
)

func NewOrganizationEmpty() *Organization {
	return &Organization{
		Base:      model.NewBase(0),
		ParentIds: make([]field.ID, 0),
		Tags:      make([]string, 0),
	}
}

func NewOrganization(
	ownAccId field.ID, parentIds []field.ID,
	isPrivate bool, kind, become uint8, name, display string, tags []string,
) *Organization {
	return &Organization{
		Base:     model.NewBase(0),
		OwnAccId: ownAccId, ParentIds: parentIds,
		IsPrivate: isPrivate, Kind: kind, Become: become, Name: name, Display: display, Tags: tags,
	}
}

// 验证场景
const (
	OrgSceneUpdateName   valid.Scene = valid.SceneCustom + 1 // 更新名称
	OrgSceneUpdateBecome valid.Scene = valid.SceneCustom + 2 // 更新加入方式
)

func (o *Organization) ValidFieldRules() valid.FieldValidRules {
	return valid.FieldValidRules{
		valid.SceneInsert: valid.FieldValidRule{
			// 组织类型
			"kind-check": func(value reflect.Value, param string) bool {
				v := uint8(value.Uint())
				return v == OrgKindGroup || v == OrgKindCompany || v == OrgKindStudio || v == OrgKindTeam
			},
			// 加入方式
			"become-check": func(value reflect.Value, param string) bool {
				v := uint8(value.Uint())
				return v == OrgBecomeDirect || v == OrgBecomeApply || v == OrgBecomeInvite
			},
			// 名称(全) (1-50)
			"name-format": func(value reflect.Value, param string) bool {
				v := value.String()
				if utf8.RuneCountInString(v) < 1 || utf8.RuneCountInString(v) > 50 {
					return false
				}
				for _, r := range v {
					if !unicode.IsLetter(r) && !unicode.IsNumber(r) && r != '_' && r != '-' {
						return false
					}
				}
				return true
			},
			// 名称(简) (0-25)
			"display-format": func(value reflect.Value, param string) bool {
				v := value.String()
				if utf8.RuneCountInString(v) > 25 {
					return false
				}
				for _, r := range v {
					if !unicode.IsLetter(r) && !unicode.IsNumber(r) && r != '_' && r != '-' {
						return false
					}
				}
				return true
			},

			// TODO:GG 这里了
			// 组织标签们 (0-10)*(1-20)
			"tags-format": func(value reflect.Value, param string) bool {
				data := value.Interface().([]string)
				if len(data) > 10 {
					return false
				}
				for _, v := range data {
					if len(v) < 1 || len(v) > 20 {
						return false
					}
				}
				return true
			},
		},
	}
}

func (o *Organization) ValidExtraRules() (field.KMap, valid.ExtraValidRules) {
	return o.Extra, valid.ExtraValidRules{
		valid.SceneAll: valid.ExtraValidRule{
			// 官网 (<1000)
			orgExtKeyWebsiteUrl: valid.ExtraValidRuleInfo{
				Field: orgExtKeyWebsiteUrl,
				ValidFn: func(value any) bool {
					data, ok := value.(string)
					if !ok {
						return false
					}
					return len(data) <= 1000
				},
			},
			// 简介 (<1000)
			orgExtKeyDesc: valid.ExtraValidRuleInfo{
				Field: orgExtKeyDesc,
				ValidFn: func(value any) bool {
					data, ok := value.(string)
					if !ok {
						return false
					}
					return len(data) <= 1000
				},
			},
			// 地址 (<100)*(<1000)
			orgExtKeyAddresses: valid.ExtraValidRuleInfo{
				Field: orgExtKeyAddresses,
				ValidFn: func(value any) bool {
					data, ok := value.([]string)
					if !ok {
						return false
					}
					if len(data) > 100 {
						return false
					}
					for _, v := range data {
						if len(v) > 1000 {
							return false
						}
					}
					return true
				},
			},
			// 联系方式 (<100)*(<1000)
			orgExtKeyContacts: valid.ExtraValidRuleInfo{
				Field: orgExtKeyContacts,
				ValidFn: func(value any) bool {
					data, ok := value.([]string)
					if !ok {
						return false
					}
					if len(data) > 100 {
						return false
					}
					for _, v := range data {
						if len(v) > 1000 {
							return false
						}
					}
					return true
				},
			},
		},
	}
}

func (o *Organization) ValidLocalizeRules() valid.LocalizeValidRules {
	return valid.LocalizeValidRules{
		valid.SceneAll: valid.LocalizeValidRule{
			Rule1: map[valid.Tag]map[valid.FieldName]valid.LocalizeValidRuleParam{
				valid.TagRequired: {
					"OwnAccIds": {"format_s_input_required", false, []any{"own_accounts"}},
					"Name":      {"format_s_input_required", false, []any{"org_name"}},
				},
			}, Rule2: map[valid.Tag]valid.LocalizeValidRuleParam{
				"own-check":         {"format_org_own_accs_err", false, nil},
				"parent-check":      {"format_org_parents_err", false, nil},
				"name-format":       {"format_org_name_err", false, nil},
				"display-format":    {"format_org_display_err", false, nil},
				"kind-check":        {"format_org_kind_err", false, nil},
				"become-check":      {"format_org_become_err", false, nil},
				"tags-format":       {"format_org_tags_err", false, nil},
				orgExtKeyWebsiteUrl: {"format_website_err", false, nil},
				orgExtKeyDesc:       {"format_desc_err", false, nil},
				orgExtKeyAddresses:  {"format_addresses_err", false, nil},
				orgExtKeyContacts:   {"format_contacts_err", false, nil},
			},
		},
	}
}

// extra
const (
	// TODO:GG 有成员的时候，获取需要各种auth?登录不需要
	orgExtKeyRootPwd  = "rootPwd"  // 根密码
	orgExtKeyMultiJob = "multiJob" // 是否允许单用户多任职

	orgExtKeyWebsiteUrl = "websiteUrl" // 官网
	orgExtKeyFaviconUrl = "faviconUrl" // 图标
	orgExtKeyDesc       = "desc"       // 简介
	orgExtKeyAddresses  = "addresses"  // 地址
	orgExtKeyContacts   = "contacts"   // 联系方式
	orgExtKeyCertImgs   = "contacts"   // 联系方式

	// TODO:GG 支持的Account的认证方式? 支持的Permission的方式?
	// TODO:GG PasswordType, PasswordSalt
)

func (o *Organization) SetRootPwd(pwd *string) {
	o.Extra.SetString(orgExtKeyRootPwd, pwd)
}

func (o *Organization) GetRootPwd() string {
	data, _ := o.Extra.GetString(orgExtKeyRootPwd)
	return data
}

func (o *Organization) SetMultiJob(multiJob *bool) {
	o.Extra.SetBool(orgExtKeyMultiJob, multiJob)
}

func (o *Organization) GetMultiJob() bool {
	data, _ := o.Extra.GetBool(orgExtKeyMultiJob)
	return data
}

func (o *Organization) SetWebsiteUrl(website *string) {
	o.Extra.SetString(orgExtKeyWebsiteUrl, website)
}

func (o *Organization) GetWebsiteUrl() string {
	data, _ := o.Extra.GetString(orgExtKeyWebsiteUrl)
	return data
}

func (o *Organization) SetFaviconUrl(website *string) {
	o.Extra.SetString(orgExtKeyFaviconUrl, website)
}

func (o *Organization) GetFaviconUrl() string {
	data, _ := o.Extra.GetString(orgExtKeyFaviconUrl)
	return data
}

func (o *Organization) SetDesc(desc *string) {
	o.Extra.SetString(orgExtKeyDesc, desc)
}

func (o *Organization) GetDesc() string {
	data, _ := o.Extra.GetString(orgExtKeyDesc)
	return data
}

func (o *Organization) SetAddresses(addresses *[]string) {
	o.Extra.SetStringSlice(orgExtKeyAddresses, addresses)
}

func (o *Organization) GetAddresses() []string {
	data, _ := o.Extra.GetStringSlice(orgExtKeyAddresses)
	return data
}

func (o *Organization) SetContacts(contacts *[]string) {
	o.Extra.SetStringSlice(orgExtKeyContacts, contacts)
}

func (o *Organization) GetContacts() []string {
	data, _ := o.Extra.GetStringSlice(orgExtKeyContacts)
	return data
}
