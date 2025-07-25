package field

const (
	StateWhite  State = 1 << 1 // 白名单
	StateEnable State = 1 << 0 // 启用
	StateInit   State = 0
	StateBlack  State = -1 << 0 // 黑名单
	StateDel    State = -1 << 1 // 删除
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
