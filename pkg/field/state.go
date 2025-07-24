package field

// State 状态
type State int64

// StateNone 无状态
const StateNone State = 0

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
