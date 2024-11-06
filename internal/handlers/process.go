package handlers

import (
	"errors"
	"finite-automatas/internal/automata"
	"fmt"
	"github.com/Mist3rBru/go-clack/prompts"
	"github.com/Mist3rBru/go-clack/third_party/picocolors"
	"github.com/Mist3rBru/go-clack/third_party/sisteransi"
	"os"
	"strings"
	"time"
)

func Process() {
	os.Stdout.Write([]byte(sisteransi.EraseDown()))
	prompts.Intro(picocolors.BgCyan(picocolors.White(" Finite Automata ")))

	path, err := prompts.Path(prompts.PathParams{
		Message:      "Where is your automata located?",
		InitialValue: ".",
		Validate: func(value string) error {
			if value == "" {
				return errors.New("please enter a path")
			}
			if !strings.HasPrefix(value, ".") {
				return errors.New("please enter a relative path")
			}
			return nil
		},
	})
	handleCancel(err)

	aut, err := automata.NewFromFile(path)
	if err != nil {
		prompts.Error(fmt.Sprintf("Error parsing automata: %v", err))
		return
	}

	start := time.Now()

	reachable := aut.FindReachableStates()
	reachableSt := strings.Builder{}
	for state := 0; state < aut.States; state++ {
		if reachable[state] {
			reachableSt.WriteString(fmt.Sprintf("[%d]: reachable", state))
		} else {
			reachableSt.WriteString(fmt.Sprintf("[%d]: not reachable", state))
		}

		if state != aut.States-1 {
			reachableSt.WriteString("\n")
		}
	}

	prompts.Note(reachableSt.String(), prompts.NoteOptions{Title: "Reachable States"})

	deadEnds := aut.FindDeadEndStates()
	deadEndsSt := strings.Builder{}
	if len(deadEnds) > 0 {
		for _, state := range deadEnds {
			deadEndsSt.WriteString(fmt.Sprintf("[%d]: dead end", state))

			if state != aut.States-1 {
				deadEndsSt.WriteString("\n")
			}
		}
	}

	if len(deadEnds) > 0 {
		prompts.Note(deadEndsSt.String(), prompts.NoteOptions{Title: "Dead Ends"})
	} else {
		prompts.Warn("No dead ends found")
	}

	dur := time.Since(start)

	prompts.Outro(
		fmt.Sprintf("Done in %s ✨",
			picocolors.Cyan(dur.String()),
		),
	)
}
