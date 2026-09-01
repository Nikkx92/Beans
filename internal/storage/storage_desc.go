package storage

import (
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"

	_ "embed"

	"gioui.org/app"
)

//go:embed state.json
var DefaultState []byte

type State struct {
	StateCoffee StateCoffee `json:"state_coffee"`
}

func Path() (string, error) {
	dir, err := app.DataDir()
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	p := filepath.Join(dir, "Beans")
	os.MkdirAll(p+"/flags", 0755)

	return p, nil
}

func SaveState(s State) {
	dir, err := Path()
	if err != nil {
		slog.Info(err.Error())
		return
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		slog.Info(err.Error())
		return
	}

	data, err := json.Marshal(s)
	if err != nil {
		slog.Info(err.Error())
		return
	}

	filePath := filepath.Join(dir, "state.json")
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		slog.Info(err.Error())
		return
	}
}

func LoadState() State {
	dir, err := Path()
	if err != nil {
		slog.Info(err.Error())
		return State{}
	}
	filePath := filepath.Join(dir, "state.json")

	var st State

	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.WriteFile(filePath, DefaultState, 0644)
			if err != nil {
				slog.Info(err.Error())
				return State{}
			}

			data, err = os.ReadFile(filePath)
			if err != nil {
				slog.Info(err.Error())
				return State{}
			}

		} else {
			slog.Info(err.Error())
			return State{}
		}
	}

	if err = json.Unmarshal(data, &st); err != nil {
		slog.Info(err.Error())
		return State{}
	}

	return st
}
