package ui

import (
	"github.com/spinsrv/browser"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// Layout (e.g. Div, HStack, VStack) {{{

func Div(children ...*browser.Node) *browser.Node {
	return &browser.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Div,
		Children: children,
	}
}

func Span(children ...*browser.Node) *browser.Node {
	return &browser.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Span,
		Children: children,
	}
}

var hstackBaseStyle = browser.Style{
	Display:       browser.DisplayFlex,
	FlexDirection: browser.FlexDirectionRow,
}

func HStack(children ...*browser.Node) *browser.Node {
	return &browser.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Div,
		Style:    hstackBaseStyle,
		Children: children,
	}
}

var vstackBaseStyle = browser.Style{
	Display:       browser.DisplayFlex,
	FlexDirection: browser.FlexDirectionColumn,
}

func VStack(children ...*browser.Node) *browser.Node {
	return &browser.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Div,
		Style:    vstackBaseStyle,
		Children: children,
	}
}

func HSpace(width browser.Size) *browser.Node {
	return &browser.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Div,
		Style: browser.Style{
			Width: width,
		},
	}
}

func VSpace(height browser.Size) *browser.Node {
	return &browser.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Div,
		Style: browser.Style{
			Height: height,
		},
	}
}

func Spacer() *browser.Node {
	return Div().FlexGrow("1")
}

var zstackBaseStyle = browser.Style{
	Display: browser.DisplayGrid,
}

func ZStack(children ...*browser.Node) *browser.Node {
	for i, c := range children {
		//https://codepen.io/everdimension/pen/BaNpeWe
		children[i] = c.GridArea("1/1/1/1")
	}
	return &browser.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Div,
		Style:    zstackBaseStyle,
		Children: children,
	}
}

// }}}

// Conditionals (e.g. OnlyIf, If) {{{

// originally these had argument types that were not functions, but the
// actual node. But this defeats the concept. It pays the computation
// and causes bugs when something is nil. So we require the bit more verbose
// function syntax since then we can guarantee that the node is only computed
// in that case. - NCL 2/4/2022

func OnlyIf(b bool, show func() *browser.Node) *browser.Node {
	// TODO need a zero node type...
	return If(b, show, func() *browser.Node { return textNode("") })
}

func If(b bool, true, false func() *browser.Node) *browser.Node {
	if b {
		return true()
	} else {
		return false()
	}
}

// }}}
