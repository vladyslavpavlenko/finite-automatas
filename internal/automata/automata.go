package automata

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type Transition struct {
	From  int
	To    int
	Input string
}

type Automata struct {
	States      int
	Initial     int
	FinalStates map[int]bool
	Transitions []Transition
}

func NewFromFile(filename string) (Automata, error) {
	file, err := os.Open(filename)
	if err != nil {
		return Automata{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// States count
	scanner.Scan()
	states, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return Automata{}, err
	}

	// Initial state
	scanner.Scan()
	initial, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return Automata{}, err
	}

	// Finite States
	scanner.Scan()
	finalStatesLine := scanner.Text()
	finalStates := make(map[int]bool)
	for _, s := range strings.Split(finalStatesLine, " ") {
		state, err := strconv.Atoi(s)
		if err != nil {
			return Automata{}, err
		}
		finalStates[state] = true
	}

	// Transitions
	var transitions []Transition
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " ")
		from, err := strconv.Atoi(parts[0])
		if err != nil {
			return Automata{}, err
		}
		input := parts[1]
		to, err := strconv.Atoi(parts[2])
		if err != nil {
			return Automata{}, err
		}
		transitions = append(transitions, Transition{from, to, input})
	}

	return Automata{
		States:      states,
		Initial:     initial,
		FinalStates: finalStates,
		Transitions: transitions,
	}, nil
}

func (a *Automata) FindReachableStates() map[int]bool {
	reachable := make(map[int]bool)
	reachable[a.Initial] = true

	queue := []int{a.Initial}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		for _, transition := range a.Transitions {
			if transition.From == current && !reachable[transition.To] {
				reachable[transition.To] = true
				queue = append(queue, transition.To)
			}
		}
	}
	return reachable
}

func (a *Automata) FindDeadEndStates() []int {
	nonDeadEnds := make(map[int]bool)
	for i := 0; i < a.States; i++ {
		nonDeadEnds[i] = false
	}

	for state := range a.FinalStates {
		nonDeadEnds[state] = true
	}

	for changed := true; changed; {
		changed = false
		for _, transition := range a.Transitions {
			if nonDeadEnds[transition.To] && !nonDeadEnds[transition.From] {
				nonDeadEnds[transition.From] = true
				changed = true
			}
		}
	}

	var deadEndStates []int
	for state := 0; state < a.States; state++ {
		if !nonDeadEnds[state] {
			deadEndStates = append(deadEndStates, state)
		}
	}
	return deadEndStates
}
