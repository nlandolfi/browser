package ui

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/spinsrv/browser"
	"github.com/spinsrv/browser/dom"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type theme interface {
	Button(label string) *browser.Node
	Card(children ...*browser.Node) *browser.Node
	H1(children ...*browser.Node) *browser.Node
	H2(children ...*browser.Node) *browser.Node
	H3(children ...*browser.Node) *browser.Node
	Link(href, text string) *browser.Node
	PassInput(value *string) *browser.Node
	Text(s string) *browser.Node
	TextArea(value *string) *browser.Node
	TextInput(value *string) *browser.Node
	Textf(format string, v ...interface{}) *browser.Node
}

var _ theme = (*Theme)(nil)

type Theme struct {
	FontFamily string

	BackgroundColor      string
	HoverBackgroundColor string
	TextColor            string
	LinkColor            string
}

var buttonBaseStyle = browser.Style{
	Padding:   browser.Size{Value: 5, Unit: browser.UnitPX},
	Cursor:    browser.CursorPointer,
	TextAlign: browser.TextAlignCenter,
	Color:     "blue",
}

func button(label string) *browser.Node {
	return &browser.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Div,
		Style:    buttonBaseStyle,
		Children: []*browser.Node{
			text(label),
		},
	}
}

func (t *Theme) Button(label string) *browser.Node {
	if t == nil {
		panic("called Button on nil theme")
	}
	return button(label).Color(t.LinkColor)
}

var cardBaseShadow = browser.BoxShadow{
	HOffset: browser.Size{},
	VOffset: browser.Size{Value: 2, Unit: browser.UnitPX},
	Blur:    browser.Size{Value: 4, Unit: browser.UnitPX},
	Spread:  browser.Size{},
	Color:   "rgba(0, 0, 0, 0.1)",
}

var cardLiftedShadow = browser.BoxShadow{
	HOffset: browser.Size{},
	VOffset: browser.Size{Value: 4, Unit: browser.UnitPX},
	Blur:    browser.Size{Value: 8, Unit: browser.UnitPX},
	Spread:  browser.Size{Value: 0, Unit: browser.UnitPX},
	Color:   "rgba(0, 0, 0, 0.1)",
}

var cardBaseStyle = browser.Style{
	Margin: browser.Size{Value: 10, Unit: browser.UnitPX},
	Border: browser.Border{
		Width: browser.Size{Value: 1, Unit: browser.UnitPX},
		Type:  browser.BorderSolid,
		Color: "lightgray",
	},
}

func card(children ...*browser.Node) *browser.Node {

	n := &browser.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Div,
		// TOOD: this doesn't copy the shadows, so if anyone mutates them we are in trouble
		Style:    cardBaseStyle,
		Children: children,
	}

	return n.Class("card")
}

func (t *Theme) Card(children ...*browser.Node) *browser.Node {
	if t == nil {
		panic("called Card on nil theme")
	}
	return card(children...).Background(t.BackgroundColor).Color(t.TextColor).FontFamily(t.FontFamily)
}

func h1(children ...*browser.Node) *browser.Node {
	return (&browser.Node{
		Type:     html.ElementNode,
		DataAtom: atom.H1,
		Children: children,
	})
}

func h2(children ...*browser.Node) *browser.Node {
	return &browser.Node{
		Type:     html.ElementNode,
		DataAtom: atom.H2,
		Children: children,
	}
}
func h3(children ...*browser.Node) *browser.Node {
	return &browser.Node{
		Type:     html.ElementNode,
		DataAtom: atom.H3,
		Children: children,
	}
}

func (t *Theme) H1(children ...*browser.Node) *browser.Node {
	if t == nil {
		panic("called H1 on nil theme")
	}
	return h1(children...).Color(t.TextColor).FontFamily(t.FontFamily)
}

func (t *Theme) H1f(format string, vs ...interface{}) *browser.Node {
	if t == nil {
		panic("called H1f on nil theme")
	}
	return t.H1(t.Textf(format, vs...))
}

func (t *Theme) H2(children ...*browser.Node) *browser.Node {
	if t == nil {
		panic("called H2 on nil theme")
	}
	return h2(children...).Color(t.TextColor).FontFamily(t.FontFamily)
}

func (t *Theme) H2f(format string, vs ...interface{}) *browser.Node {
	if t == nil {
		panic("called H2f on nil theme")
	}
	return t.H2(t.Textf(format, vs...))
}

func (t *Theme) H3(children ...*browser.Node) *browser.Node {
	if t == nil {
		panic("called H3 on nil theme")
	}
	return h3(children...).Color(t.TextColor).FontFamily(t.FontFamily)
}

func (t *Theme) H3f(format string, vs ...interface{}) *browser.Node {
	if t == nil {
		panic("called H3f on nil theme")
	}
	return t.H3(t.Textf(format, vs...))
}

func link(href, text string) *browser.Node {
	return &browser.Node{
		Type:     html.ElementNode,
		DataAtom: atom.A,
		Attr: []*html.Attribute{
			&html.Attribute{
				Key: atom.Href.String(),
				Val: href,
			},
		},
		Children: []*browser.Node{
			textNode(text),
		},
	}
}

func (t *Theme) Link(href, text string) *browser.Node {
	if t == nil {
		panic("called Link on nil theme")
	}
	return link(href, text).Color(t.LinkColor).FontFamily(t.FontFamily)
}

var inputBaseStyle = browser.Style{
	Padding: browser.Size{Value: 5, Unit: browser.UnitPX},
	Border:  browser.Border{Type: browser.BorderNone},
	BorderBottom: browser.Border{
		Width: browser.Size{Value: 1, Unit: browser.UnitPX},
		Type:  browser.BorderSolid,
		Color: "gray",
	},
	Outline: browser.Outline{Type: browser.OutlineNone},
}

func textInput(value *string) *browser.Node {
	return (&browser.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Input,
		Attr: []*html.Attribute{
			&html.Attribute{
				Key: atom.Type.String(),
				Val: "text",
			},
			&html.Attribute{
				Key: atom.Value.String(),
				Val: *value,
			},
		},
		Style: inputBaseStyle,
	}).OnInput(func(e dom.Event) {
		*value = e.Target().Value()
		go browser.Dispatch(nil)
	})
}

func (t *Theme) PassInput(value *string) *browser.Node {
	if t == nil {
		panic("called PassInput on nil theme")
	}
	return t.TextInput(value).AttrType("password")
}

func TextNode(s string) *browser.Node {
	return textNode(s)
}

func textNode(s string) *browser.Node {
	return &browser.Node{
		Type: html.TextNode,
		Data: s,
	}
}

func text(s string) *browser.Node {
	return Span(textNode(s))
}

func (t *Theme) Text(s string) *browser.Node {
	if t == nil {
		panic("called Text on nil theme")
	}
	return text(s).Color(t.TextColor).FontFamily(t.FontFamily)
}

func textArea(value *string) *browser.Node {
	return (&browser.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Textarea,
		Children: []*browser.Node{
			textNode(*value),
		},
		Attr: []*html.Attribute{
			&html.Attribute{
				Key: atom.Value.String(),
				Val: *value,
			},
		},
	}).OnInput(func(e dom.Event) {
		*value = e.Target().Value()
		go browser.Dispatch(nil)
	})
}

func (t *Theme) TextArea(value *string) *browser.Node {
	if t == nil {
		panic("called TextArea on nil theme")
	}
	return textArea(value).Color(t.TextColor).FontFamily(t.FontFamily).Background(t.BackgroundColor)
}

func (t *Theme) TextInput(value *string) *browser.Node {
	if t == nil {
		panic("called TextInput on nil theme")
	}
	return textInput(value).Color(t.TextColor).Background(t.BackgroundColor).FontFamily(t.FontFamily)
}

func (t *Theme) TimeInput(value *time.Time) *browser.Node {
	s := inputBaseStyle

	return (&browser.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Input,
		Attr: []*html.Attribute{
			&html.Attribute{
				Key: atom.Type.String(),
				Val: "time",
			},
			&html.Attribute{
				Key: atom.Value.String(),
				Val: value.Format("15:04"),
			},
		},
		Style: s,
	}).OnInput(func(e dom.Event) {
		s := e.Target().Value()
		t, err := time.Parse("15:04", s)
		if err != nil {
			log.Fatalf("error parsing time %q: %v", s, err)
		}
		cur := *value
		*value = time.Date(cur.Year(), cur.Month(), cur.Day(),
			t.Hour(), t.Minute(),
			cur.Second(), cur.Nanosecond(), cur.Location())
	}).Color(t.TextColor).Background(t.BackgroundColor).FontFamily(t.FontFamily)
}

func (t *Theme) DateInput(value *time.Time) *browser.Node {
	s := inputBaseStyle

	n := &browser.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Input,
		Attr: []*html.Attribute{
			&html.Attribute{
				Key: atom.Type.String(),
				Val: "date",
			},
		},
		Style: s,
	}

	if !(*value).IsZero() {
		n = n.AttrValue(value.Format("2006-01-02"))
	}

	n = n.OnInput(func(e dom.Event) {
		s := e.Target().Value()
		t, err := time.Parse("2006-01-02", s)
		if err != nil {
			log.Fatalf("error parsing date %q: %v", s, err)
		}
		cur := *value
		*value = time.Date(t.Year(), t.Month(), t.Day(),
			cur.Hour(), cur.Minute(),
			cur.Second(), cur.Nanosecond(), cur.Location())
	})

	return n.Color(t.TextColor).Background(t.BackgroundColor).FontFamily(t.FontFamily)
}

func (t *Theme) NumberInput(value *int64) *browser.Node {
	s := inputBaseStyle

	n := &browser.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Input,
		Attr: []*html.Attribute{
			&html.Attribute{
				Key: atom.Type.String(),
				Val: "number",
			},
		},
		Style: s,
	}

	n = n.AttrValue(fmt.Sprintf("%d", *value))

	n = n.OnInput(func(e dom.Event) {
		s := e.Target().Value()
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			log.Fatalf("error parsing number %q: %v", s, err)
		}
		*value = i
	})

	return n.Color(t.TextColor).Background(t.BackgroundColor).FontFamily(t.FontFamily)
}

func (t *Theme) Textf(format string, vs ...interface{}) *browser.Node {
	if t == nil {
		panic("called Textf on nil theme")
	}
	return t.Text(fmt.Sprintf(format, vs...))
}

func (t *Theme) Toggle(value *bool) *browser.Node {
	return If(*value,
		func() *browser.Node { return t.Text("☑") },
		func() *browser.Node { return t.Text("☐") },
	).Pointer().OnClick(func(e dom.Event) {
		*value = !*value
		go browser.Dispatch(nil) // force a re-render
	})
}

/*
type SelectOption struct {
	Key, Display string
}

func (t *Theme) TextSelect(selected *string, options []SelectOption) *browser.Node {
}
*/

var (
	DarkTheme = Theme{
		FontFamily:           "Avenir Next",
		BackgroundColor:      "black",
		HoverBackgroundColor: "rgb(30,30,30)",
		TextColor:            "antiquewhite",
		//		LinkColor:       "rgb(33,95,180)",
		LinkColor: "rgb(66,196,208)",
	}
	LightTheme = Theme{
		FontFamily:           "Avenir Next",
		BackgroundColor:      "white",
		HoverBackgroundColor: "rgb(238,238,238)",
		TextColor:            "black",
		LinkColor:            "blue",
	}
)

// DefaultTheme and Helpers (e.g. Link, H1) {{{

var DefaultTheme = LightTheme

func Link(href, text string) *browser.Node {
	return DefaultTheme.Link(href, text)
}

func H1(children ...*browser.Node) *browser.Node {
	return DefaultTheme.H1(children...)
}

func H2(children ...*browser.Node) *browser.Node {
	return DefaultTheme.H2(children...)
}

func H3(children ...*browser.Node) *browser.Node {
	return DefaultTheme.H3(children...)
}

func Card(children ...*browser.Node) *browser.Node {
	return DefaultTheme.Card(children...)
}

func Button(text string) *browser.Node {
	return DefaultTheme.Button(text)
}

func TextArea(text *string) *browser.Node {
	return DefaultTheme.TextArea(text)
}

// }}}
