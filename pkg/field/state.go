package field

const (
	StateUserCustom       = StateEnable // 用户自定义
	StateEnable     State = 1 << 1      // 停用/启用
	StateDel        State = 1 << 0      // 删除/存在
	StateInit       State = 0           // 初始
	StateInvisible  State = -1 << 0     // 可见/屏蔽(对外)
	StateBlack      State = -1 << 1     // 白名单/黑名单(登录等权限)
	StateSysCustom        = StateBlack  // 系统自定义
)

// State 状态
type State int64

func (s State) value() int64 {
	return int64(s)
}

func (s State) Add(states State) {
	s |= states
}

func (s State) Remove(states State) {
	s &= ^states
}

func (s State) Clear(ignores State) {
	s &= ignores
}

func (s State) HasAny(states State) bool {
	return s&states != 0
}

func (s State) HasAll(states State) bool {
	return (s & states) == states
}

func (s State) Equal(states State) bool {
	return s == states
}
