package storage

import (
	"katydid-mp-account/internal/pkg/entity"
	"katydid-mp-account/internal/pkg/storage"
)

type (
	Organization struct {
		db *storage.Pgsql

		entity.Entity

		OwnAccId  int64   `json:"ownAccId" validate:"required,own-check" gorm:"comment:所属账号"`
		ParentIds []int64 `json:"parentIds" validate:"parent-check" gorm:"comment:父级组织"`

		Enable   bool     `json:"enable" gorm:"comment:是否可用"`
		IsPublic bool     `json:"isPublic" gorm:"comment:是否公开"`
		Kind     uint8    `json:"kind" validate:"kind-check" gorm:"comment:组织类型"`
		Become   uint8    `json:"become" validate:"become-check" gorm:"comment:加入方式"`
		Name     string   `json:"name" validate:"required,name-format" gorm:"comment:组织名称"`
		Display  string   `json:"display" validate:"display-format" gorm:"comment:组织显示名称"`
		Tags     []string `json:"tags" validate:"tags-format" gorm:"comment:组织标签们"`
	}

	OrgExtraKey string // 组织扩展键
)

var (
	IDFactory = entity.NewSnowflake()
)

const (
	OrgKindPhysical uint8 = 0 // 实体组织 (同时存在数受orgExtKeyMultiJob影响)
	OrgKindVirtual  uint8 = 1 // 虚拟组织 (能同时存在多个)

	OrgBecomePublic uint8 = 0 // 公开
	OrgBecomeApply  uint8 = 1 // 申请 (只有public有效?)
	OrgBecomeInvite uint8 = 2 // 邀请
)

func (o *Organization) IsTopParent() bool {
	return len(o.ParentIds) == 0
}

const (
	RootPwd OrgExtraKey = "rootPwd"

	// TODO:GG 有成员的时候，获取需要各种auth?登录不需要
	orgExtKeyRootPwd  = "rootPwd"  // 根密码
	orgExtKeyMultiJob = "multiJob" // 是否允许单用户多任职

	orgExtKeyWebsiteUrl = "websiteUrl" // 官网
	orgExtKeyFaviconUrl = "faviconUrl" // 图标
	orgExtKeyDesc       = "desc"       // 简介
	orgExtKeyAddresses  = "addresses"  // 地址
	orgExtKeyContacts   = "contacts"   // 联系方式

	// TODO:GG 支持的Account的认证方式? 支持的Permission的方式?
	// TODO:GG PasswordType, PasswordSalt
)

// TODO:GG 数据库增删改查
