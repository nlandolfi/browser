package app

import (
	"log"
	"time"

	"github.com/spinsrv/browser"
	"github.com/spinsrv/browser/ui"
)

type State struct {
	ui.Theme
	ClientVersion string
	LastWrittenAt time.Time
}

func View(s *State) *browser.Node {
	return ui.VStack(
		ui.HStack(
			s.Theme.Text("Hello, spinsrv/browser!"),
		).JustifyContentCenter(),
	).HeightVH(100).WidthVW(100).JustifyContentCenter().Background(s.Theme.BackgroundColor)
}

func (s *State) Handle(e browser.Event) {
	switch e.(type) {
	case EventInitialize:
		log.Print("EventInitialize")
	}
}

type EventInitialize struct{}

func (s *State) Rewire() {
	// no op
	// if there were subcomponents, then here we could set their theme.
}
