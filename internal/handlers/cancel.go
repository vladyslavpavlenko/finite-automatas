package handlers

import (
	"github.com/Mist3rBru/go-clack/prompts"
	"os"
)

func handleCancel(err error) {
	if err != nil {
		prompts.Cancel("Operation cancelled.")
		os.Exit(0)
	}
}
