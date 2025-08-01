package valid

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
	"sync"
)

var (
	valid *Validator
	vOnce sync.Once
)

// Validator 验证器
type Validator struct {
	validate *validator.Validate
	regTypes *sync.Map // 验证注册类型缓存
	regLocs  *sync.Map // 本地化文本缓存
}

// Scene 验证场景
type Scene uint64

const (
	SceneAll Scene = 0 // 所有的场景

	SceneBind Scene = 1 << 0 // 请求数据绑定
	SceneAdd  Scene = 1 << 1 // 添加/新增
	SceneDel  Scene = 1 << 2 // 删除/移除
	SceneUpd  Scene = 1 << 3 // 更新/修改
	SceneGet  Scene = 1 << 4 // 获取/查询
	SceneRes  Scene = 1 << 5 // 返回/响应

	SceneCustom Scene = SceneRes // 自定义 custom << ?
)

// Tag 字段标签
type Tag string

const (
	TagRequired Tag = "required"
	TagFormat   Tag = "format"
	TagRange    Tag = "range"
	TagCheck    Tag = "check"
)

// FieldName 字段名称
type FieldName string

// 字段验证
type (
	IFieldValidator interface {
		ValidFieldRules() FieldValidRules
	}

	FieldValidRules  = map[Scene]FieldValidRule
	FieldValidRule   = map[Tag]FieldValidRuleFn
	FieldValidRuleFn = func(value reflect.Value, param string) bool
)

// 额外(Extra)字段验证
type (
	IExtraValidator interface {
		ValidExtraRules() (map[string]any, ExtraValidRules)
	}

	ExtraValidRules    = map[Scene]ExtraValidRule
	ExtraValidRule     = map[Tag]ExtraValidRuleInfo
	ExtraValidRuleInfo struct {
		Field   string
		Param   string
		ValidFn func(value any) bool
	}
)

// 结构体(多字段关联)验证
type (
	IStructValidator interface {
		ValidStructRules(scene Scene, fn FuncReportError)
	}

	// FuncReportError validator.StructLevel.FuncReportError
	FuncReportError = func(field any, fieldName FieldName, tag Tag, param string)
)

// 本地化
type (
	ILocalizeValidator interface {
		ValidLocalizeRules() LocalizeValidRules
	}

	LocalizeValidRules = map[Scene]LocalizeValidRule
	LocalizeValidRule  struct {
		Rule1 map[Tag]map[FieldName]LocalizeValidRuleParam
		Rule2 map[Tag]LocalizeValidRuleParam
	}
	LocalizeValidRuleParam = [3]any // {msg, param, template([]any)}
)

// MsgErr 定义错误信息结构体
type MsgErr struct {
	Err    error
	Msg    string
	Params []any
}

func Get() *Validator {
	vOnce.Do(func() {
		opts := []validator.Option{
			validator.WithRequiredStructEnabled(),
		}

		valid = &Validator{
			validate: validator.New(opts...),
			regTypes: &sync.Map{},
			regLocs:  &sync.Map{},
		}

		// 设置Tag <- 默认json标签
		//validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		//	name := fld.Tag.Get("json")
		//	if name == "-" {
		//		return fld.Name
		//	}
		//	return name
		//})
	})
	return valid
}

// RegisterFieldRule 注册字段验证规则
func RegisterFieldRule(fieldRules FieldValidRule) {
	v := Get()
	// 注册字段验证规则
	for tag, rule := range fieldRules {
		_ = v.validate.RegisterValidation(string(tag), func(fl validator.FieldLevel) bool {
			return rule(fl.Field(), fl.Param())
		})
	}
}

// Check 根据场景执行验证，并返回本地化错误信息
func Check(obj any, scene Scene) []*MsgErr {
	if obj == nil {
		return []*MsgErr{{
			Err: errors.New("validation object cannot be nil"),
			Msg: "invalid_object_validation",
		}}
	}

	v := Get()
	typ := reflect.TypeOf(obj)

	// -- 注册验证(设置缓存) --
	if _, ok := v.regTypes.Load(typ); !ok {
		if e := v.registerValidations(obj, scene); e != nil {
			return []*MsgErr{{Err: e}}
		}
		v.regTypes.Store(typ, true)
	}

	// -- 执行验证(有缓存) --
	if e := v.validate.Struct(obj); e != nil {
		return v.handleValidationError(obj, scene, e)
	}
	return nil
}

// registerValidations 注册验证规则
func (v *Validator) registerValidations(obj any, scene Scene) error {
	// -- 字段验证注册 --
	if e := v.validFields(obj, scene); e != nil {
		return e
	}

	v.validate.RegisterStructValidation(func(sl validator.StructLevel) {
		cObj := sl.Current().Addr().Interface()

		// -- 额外验证注册 --
		v.validExtra(cObj, sl, scene)

		// -- 结构验证注册 --
		v.validStruct(cObj, sl, scene)
	}, obj)
	return nil
}

// validFields 注册字段验证规则
func (v *Validator) validFields(obj any, scene Scene) error {
	// 处理嵌入字段的验证规则
	if e := v.processEmbeddedValidations(obj, scene, 1, nil); e != nil {
		return e
	}

	fv, ok := obj.(IFieldValidator)
	if !ok {
		return nil
	}

	// 获取验证规则
	sceneRules := fv.ValidFieldRules()
	if sceneRules == nil {
		return nil
	}

	// 筛选出当前场景的验证规则
	scenes := make([]Scene, 0)
	for key := range sceneRules {
		if key == SceneAll {
			scenes = append(scenes, SceneAll) // 添加全局场景
		} else if (key & scene) == scene {
			scenes = append(scenes, key) // 添加当前场景(必须是全部包含)
		}
	}

	// 遍历所有场景的验证规则
	tagRules := make(map[Tag]FieldValidRuleFn)
	for _, s := range scenes {
		if tRules := sceneRules[s]; tRules != nil {
			for tag, rule := range tRules {
				tagRules[tag] = rule // 合并验证规则
			}
		}
	}

	// 注册验证规则
	for tag, rule := range tagRules {
		if e := v.validate.RegisterValidation(string(tag), func(fl validator.FieldLevel) bool {
			return rule(fl.Field(), fl.Param())
		}); e != nil {
			return e
		}
	}
	return nil
}

// validExtra 注册额外验证规则
func (v *Validator) validExtra(obj any, sl validator.StructLevel, scene Scene) {
	// 处理嵌入字段的验证规则
	_ = v.processEmbeddedValidations(obj, scene, 2, sl)

	ev, ok := obj.(IExtraValidator)
	if !ok {
		return
	}

	// 获取验证规则
	extra, sceneRules := ev.ValidExtraRules()
	if (extra == nil) || (sceneRules == nil) {
		return
	}

	// 筛选出当前场景的验证规则
	scenes := make([]Scene, 0)
	for key := range sceneRules {
		if key == SceneAll {
			scenes = append(scenes, SceneAll) // 添加全局场景
		} else if (key & scene) == scene {
			scenes = append(scenes, key) // 添加当前场景(必须是全部包含)
		}
	}

	// 遍历所有场景的验证规则
	tagRules := make(map[Tag]ExtraValidRuleInfo)
	for _, s := range scenes {
		if tRules := sceneRules[s]; tRules != nil {
			for tag, rule := range tRules {
				tagRules[tag] = rule // 合并验证规则
			}
		}
	}

	// 注册验证规则
	for tag, rule := range tagRules {
		value, exists := extra[string(tag)]
		if (tag == TagRequired) && !exists {
			sl.ReportError(value, rule.Field, rule.Field, string(tag), rule.Param)
			continue
		}
		if exists && !rule.ValidFn(value) {
			sl.ReportError(value, rule.Field, rule.Field, string(tag), rule.Param)
		}
	}
}

// validStruct 注册结构体验证规则
func (v *Validator) validStruct(obj any, sl validator.StructLevel, scene Scene) {
	// 筛选出当前场景的验证规则
	scenes := make([]Scene, 0)
	scenes = append(scenes, SceneAll) // 添加全局场景
	scenes = append(scenes, scene)    // 添加当前场景(实现类判断)

	// 处理嵌入字段的验证规则(全局+当前)
	for _, s := range scenes {
		_ = v.processEmbeddedValidations(obj, s, 3, sl)
	}

	sv, ok := obj.(IStructValidator)
	if !ok {
		return
	}

	// 获取验证规则(全局+当前)
	for _, s := range scenes {
		sv.ValidStructRules(s, func(field any, fieldName FieldName, tag Tag, param string) {
			sl.ReportError(field, string(fieldName), string(fieldName), string(tag), param)
		})
	}
}

// processEmbeddedValidations 递归注册组合类型的验证规则
// ttt: 处理类型 1=字段, 2=额外, 3=结构体
func (v *Validator) processEmbeddedValidations(
	obj any, scene Scene,
	ttt int, sl validator.StructLevel,
) error {
	val := reflect.ValueOf(obj)
	typ := reflect.TypeOf(obj)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil // 避免空指针引用
		}
		val = val.Elem()
		typ = typ.Elem()
	}

	// 遍历所有字段
	for i := 0; i < typ.NumField(); i++ {
		// 检查是否是组合类型的字段
		field := typ.Field(i)
		if !field.Anonymous {
			continue
		}

		fieldVal := val.Field(i)
		fieldType := field.Type
		var embedObj any

		if fieldType.Kind() == reflect.Ptr {
			// 处理指针类型的组合字段
			if fieldVal.IsNil() {
				continue
			}
			embedObj = fieldVal.Interface()
			fieldType = fieldType.Elem()
		} else {
			// 处理非指针类型的组合字段
			if !fieldVal.CanAddr() {
				continue // 避免不可寻址的字段
			}
			// 处理非指针类型的组合字段
			embedObj = fieldVal.Addr().Interface()
		}

		// 只处理结构体类型的组合字段
		if fieldType.Kind() != reflect.Struct || embedObj == nil {
			continue
		}

		// 递归处理嵌入字段
		if err := v.processEmbeddedValidations(embedObj, scene, ttt, sl); err != nil {
			return err
		}

		// 根据处理类型执行对应验证
		switch ttt {
		case 1: // 字段验证
			if fv, okk := embedObj.(IFieldValidator); okk {
				if err := v.validFields(fv, scene); err != nil {
					return err
				}
			}
		case 2: // 额外验证
			if ev, okk := embedObj.(IExtraValidator); okk {
				v.validExtra(ev, sl, scene)
			}
		case 3: // 结构体验证
			if sv, okk := embedObj.(IStructValidator); okk {
				v.validStruct(sv, sl, scene)
			}
		}
	}
	return nil
}

// handleValidationError 处理验证错误
func (v *Validator) handleValidationError(
	obj any, scene Scene, e error,
) []*MsgErr {
	var invalidErr *validator.InvalidValidationError
	if errors.As(e, &invalidErr) {
		// -- 验证失败 --
		return []*MsgErr{{Err: e, Msg: "invalid_object_validation"}}
	}

	var validateErrs validator.ValidationErrors
	if errors.As(e, &validateErrs) {
		// -- 本地化错误注册 --
		if rl, ok := obj.(ILocalizeValidator); ok {
			return v.validLocalize(scene, obj, rl, validateErrs, true)
		}
	}
	return []*MsgErr{{Err: e, Msg: "unknown_validator_err"}}
}

// validLocalize 验证本地化错误
func (v *Validator) validLocalize(
	scene Scene, obj any,
	rl ILocalizeValidator,
	validateErrs validator.ValidationErrors,
	first bool,
) []*MsgErr {
	var msgErrs []*MsgErr

	// 处理组合类型的验证规则
	if msgEs := v.processEmbeddedLocalizes(scene, obj, validateErrs); msgEs != nil {
		msgErrs = append(msgErrs, msgEs...)
	}

	var localRule LocalizeValidRule
	typ := reflect.TypeOf(obj)
	cacheRules, ok := v.regLocs.Load(typ)
	if !ok {
		// 没有就缓存，注册本地化规则
		sceneRules := rl.ValidLocalizeRules()
		if sceneRules == nil {
			return msgErrs
		}

		// 筛选出当前场景的验证规则
		scenes := make([]Scene, 0)
		for key := range sceneRules {
			if key == SceneAll {
				scenes = append(scenes, SceneAll) // 添加全局场景
			} else if (key & scene) == scene {
				scenes = append(scenes, key) // 添加当前场景(必须是全部包含)
			}
		}

		// 遍历所有场景的验证规则1
		tagFieldRules := make(map[Tag]map[FieldName]LocalizeValidRuleParam)
		for _, s := range scenes {
			if tRules := sceneRules[s]; tRules.Rule1 != nil {
				for tag, rule := range tRules.Rule1 {
					tagFieldRules[tag] = rule // 合并验证规则
				}
			}
		}

		// 遍历所有场景的验证规则2
		tagRules := make(map[Tag]LocalizeValidRuleParam)
		for _, s := range scenes {
			if tRules := sceneRules[s]; tRules.Rule2 != nil {
				for tag, rule := range tRules.Rule2 {
					tagRules[tag] = rule // 合并验证规则
				}
			}
		}

		localRule = LocalizeValidRule{Rule1: tagFieldRules, Rule2: tagRules}
		v.regLocs.Store(typ, localRule)
	} else {
		// 有就直接使用
		localRule = cacheRules.(LocalizeValidRule)
	}

	// 处理每个验证错误
	for _, ee := range validateErrs {
		// -- 本地化错误注册(Tag+Field) --
		for tag, fieldRules := range localRule.Rule1 {
			if ee.Tag() == string(tag) {
				for field, rules := range fieldRules {
					if ee.Field() == string(field) {
						var params []any
						if rules[2] != nil {
							params = append(params, rules[2].([]any)...)
						}
						if rules[1].(bool) {
							params = append(params, ee.Param())
						}
						msgErrs = append(msgErrs, &MsgErr{Msg: rules[0].(string), Params: params})
					}
				}
			}
		}
		// -- 本地化错误注册(Tag) --
		for tag, rules := range localRule.Rule2 {
			if ee.Tag() == string(tag) {
				var params []any
				if rules[2] != nil {
					params = append(params, rules[2].([]any)...)
				}
				if rules[1].(bool) {
					params = append(params, ee.Param())
				}
				msgErrs = append(msgErrs, &MsgErr{Msg: rules[0].(string), Params: params})
			}
		}
	}

	// 找不到就返回默认
	if (len(msgErrs) <= 0) && first {
		// 提供更具体的错误信息，包括字段和规则
		fieldErrors := make([]string, 0, len(validateErrs))
		for _, err := range validateErrs {
			fieldErrors = append(fieldErrors,
				fmt.Sprintf("field:%s, tag:%s, param:%s",
					err.Field(), err.Tag(), err.Param()))
		}

		msgErrs = append(msgErrs, &MsgErr{
			Msg:    "validation_failed",
			Params: []any{strings.Join(fieldErrors, "; ")},
		})
	}
	return msgErrs
}

// processEmbeddedLocalizes 递归注册组合类型的本地化规则
func (v *Validator) processEmbeddedLocalizes(
	scene Scene, obj any,
	validateErrs validator.ValidationErrors,
) []*MsgErr {
	var allMsgErrs []*MsgErr

	val := reflect.ValueOf(obj)
	typ := reflect.TypeOf(obj)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil // 避免空指针引用
		}
		val = val.Elem()
		typ = typ.Elem()
	}

	// 遍历所有字段
	for i := 0; i < typ.NumField(); i++ {
		// 检查是否是组合类型的字段
		field := typ.Field(i)
		if !field.Anonymous {
			continue
		}

		fieldVal := val.Field(i)
		fieldType := field.Type
		var embedObj any

		if fieldType.Kind() == reflect.Ptr {
			// 处理指针类型的组合字段
			if fieldVal.IsNil() {
				continue
			}
			embedObj = fieldVal.Interface()
			fieldType = fieldType.Elem()
		} else {
			// 处理非指针类型的组合字段
			if !fieldVal.CanAddr() {
				continue // 避免不可寻址的字段
			}
			// 处理非指针类型的组合字段
			embedObj = fieldVal.Addr().Interface()
		}

		// 只处理结构体类型的组合字段
		if fieldType.Kind() != reflect.Struct || embedObj == nil {
			continue
		}

		// 递归处理嵌入字段的本地化规则（无论是否实现接口）
		if embedMsgErrs := v.processEmbeddedLocalizes(scene, embedObj, validateErrs); embedMsgErrs != nil {
			allMsgErrs = append(allMsgErrs, embedMsgErrs...)
		}

		// 如果嵌入字段实现了ILocalizeValidator接口
		if embedLocValidator, ok := embedObj.(ILocalizeValidator); ok {
			if msgErrs := v.validLocalize(
				scene,
				embedObj,
				embedLocValidator,
				validateErrs,
				false,
			); msgErrs != nil {
				allMsgErrs = append(allMsgErrs, msgErrs...)
			}
		}
	}

	return allMsgErrs
}
