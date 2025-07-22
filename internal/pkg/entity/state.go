package entity

// State 状态
type State uint64

func (s State) value() uint64 {
	return uint64(s)
}

func (s State) Add(states uint64) {
	s |= State(states)
}

func (s State) Remove(states uint64) {
	s &= ^State(states)
}

func (s State) HasAny(states uint64) bool {
	return s&State(states) != 0
}

func (s State) HasAll(states uint64) bool {
	return s&State(states) == State(states)
}

func (s State) Equal(states uint64) bool {
	return s == State(states)
}

func (s State) Clear(ignores uint64) State {
	return s & State(ignores)
}
