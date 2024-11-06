package automata_test

import (
	"finite-automatas/internal/automata"
	"testing"
)

func newAutomata(
	states int, initial int, finalStates []int, transitions []automata.Transition,
) automata.Automata {
	finalStateMap := make(map[int]bool)
	for _, s := range finalStates {
		finalStateMap[s] = true
	}
	return automata.Automata{
		States:      states,
		Initial:     initial,
		FinalStates: finalStateMap,
		Transitions: transitions,
	}
}

func TestFindReachableStates(t *testing.T) {
	tests := []struct {
		name                    string
		automaton               automata.Automata
		expectedReachableStates map[int]bool
	}{
		{
			name: "AllReachable",
			automaton: newAutomata(4, 0, []int{1, 2}, []automata.Transition{
				{0, 1, "a"},
				{1, 2, "b"},
				{2, 3, "c"},
				{3, 1, "d"},
			}),
			expectedReachableStates: map[int]bool{0: true, 1: true, 2: true, 3: true},
		},
		{
			name: "OneUnreachable",
			automaton: newAutomata(3, 0, []int{1}, []automata.Transition{
				{0, 1, "a"},
				{1, 0, "b"},
			}),
			expectedReachableStates: map[int]bool{0: true, 1: true, 2: false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reachables := tt.automaton.FindReachableStates()
			for state, expected := range tt.expectedReachableStates {
				if reachables[state] != expected {
					t.Errorf("State %d: expected reachable=%v, got=%v", state, expected, reachables[state])
				}
			}
		})
	}
}

func TestFindDeadEndStates(t *testing.T) {
	tests := []struct {
		name             string
		automata         automata.Automata
		expectedDeadEnds []int
	}{
		{
			name: "DeadEnd",
			automata: newAutomata(5, 0, []int{3}, []automata.Transition{
				{0, 1, "a"},
				{1, 2, "b"},
				{2, 3, "c"},
				{4, 4, "d"},
			}),
			expectedDeadEnds: []int{4},
		},
		{
			name: "NoDeadEnds",
			automata: newAutomata(3, 0, []int{1, 2}, []automata.Transition{
				{0, 1, "a"},
				{1, 2, "b"},
				{2, 0, "c"},
			}),
			expectedDeadEnds: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deadEnds := tt.automata.FindDeadEndStates()
			if len(deadEnds) != len(tt.expectedDeadEnds) {
				t.Errorf("Expected %d dead-end States, got %d", len(tt.expectedDeadEnds), len(deadEnds))
			}

			for _, expectedState := range tt.expectedDeadEnds {
				found := false
				for _, state := range deadEnds {
					if state == expectedState {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected dead-end state %d not found", expectedState)
				}
			}
		})
	}
}
