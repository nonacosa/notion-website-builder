package gui

import (
	"fyne.io/fyne/v2"
)

// Tutorial defines the data structure for a tutorial
type Info struct {
	Title, Intro string
	View         func(w fyne.Window, main *fyne.Container) fyne.CanvasObject
}

var (
	// Tutorials defines the metadata for each tutorial
	Tutorials = map[string]Info{
		"welcome":   {"welcome", "", welcomeScreen},
		"webManage": {"Web Manage", "", websiteScreen},
	}

	// TutorialIndex  defines how our tutorials should be laid out in the index tree
	TutorialIndex = map[string][]string{
		"": {"welcome", "webManage"},
	}
)
