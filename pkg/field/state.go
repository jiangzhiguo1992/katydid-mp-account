package field

// State 状态
type State int64

// StateNone 无状态
const StateNone State = 0

func (s State) value() int64 {
	return int64(s)
}

func (s State) Set(states int64) {
	s = State(states)
}

func (s State) Add(states int64) {
	s |= State(states)
}

func (s State) Remove(states int64) {
	s &= ^State(states)
}

func (s State) Clear(ignores int64) {
	s &= State(ignores)
}

func (s State) HasAny(states int64) bool {
	return s&State(states) != 0
}

func (s State) HasAll(states int64) bool {
	return s&State(states) == State(states)
}

func (s State) Equal(states int64) bool {
	return s == State(states)
}
