// Package dom provides interfaces for interacting with a browser's
// Document Object Model (DOM). For more information on the DOM, see
// https://developer.mozilla.org/en-US/docs/Web/API/Document_Object_Model.
package dom

import (
	"html/template"
)

// Document is an interface for a browser's Document object.
// See: https://developer.mozilla.org/en-US/docs/Web/API/Document.
type Document interface {
	// Body is an interface for the Document's body, as in the javascript `document.body`.
	// See: https://developer.mozilla.org/en-US/docs/Web/API/Document/body.
	Body() Element

	// GetElementByID retrieves a DOM element by its id, as in the javascript `document.getElementById`.
	// See: https://developer.mozilla.org/en-US/docs/Web/API/Document/getElementById.
	GetElementByID(string) (Element, error)

	// CreateElement creates a DOM element, as in the javacsript `document.createElement`.
	// See: https://developer.mozilla.org/en-US/docs/Web/API/Document/createElement.
	CreateElement(string) Element

	// CreateTextNode creates a DOM texdt node, as in the javascript `document.createTextNode`.
	// See: https://developer.mozilla.org/en-US/docs/Web/API/Document/createTextNode
	CreateTextNode(s string) Text

	// GetSelection gets the Selection object representing the range of text selected by the user,
	// as in `document.getSelection()`.
	// See: https://developer.mozilla.org/en-US/docs/Web/API/Document/getSelection
	Selection() Selection
}

// See https://developer.mozilla.org/en-US/docs/Web/API/Selection
type Selection interface {
	// AnchorNode is where the user began the selection
	AnchorNode() Node
	AnchorOffset() int
	// FocusNode is where the user ends the selection
	FocusNode() Node
	FocusOffset() int
	IsCollapsed() bool
	RangeCount() int
	Type() string
}

type Event interface {
	Target() Element
	ClientX() int
	ClientY() int
	MovementX() int
	MovementY() int
	Code() string
	KeyCode() int
	PreventDefault()
	StopPropagation()
	IsUndefined() bool
	DataTransfer() DataTransfer
}

type DataTransfer interface {
	Items() []DataTransferItem
}

type DataTransferItem interface {
}

type EventListener interface {
	Release()
}

type EventHandler func(Event)

// EventTarget exposes a subset of functionality of an EventTarget object.
// See: https://developer.mozilla.org/en-US/docs/Web/API/EventTarget.
type EventTarget interface {
	// AddEventListener adds the function as a listener to events of type EventType.
	// See: https://developer.mozilla.org/en-US/docs/Web/API/EventTarget/addEventListener.
	AddEventListener(EventType, EventHandler) EventListener

	RemoveEventListener(EventType, EventListener)
}

// Node exposes a subset of functionality of a Node object.
// See https://developer.mozilla.org/en-US/docs/Web/API/Node
type Node interface {
	EventTarget

	AppendChild(Node)

	ReplaceChild(Node, Node) Node

	ReplaceWith(e Node)

	RemoveChild(Node) Node

	LogSelf() // useful for debugging

	ParentElement() Element
}

type Text interface {
	Node
}

// Element exposes a subset of functionality of an Element object.
// See: https://developer.mozilla.org/en-US/docs/Web/API/Element.
type Element interface {
	// SetInnerHTML sets the innerHTML property of the element.
	// See: https://developer.mozilla.org/en-US/docs/Web/API/Element/innerHTML.
	SetInnerHTML(template.HTML)

	SetAttribute(key, val string)

	RemoveAttribute(key string)

	SetStyle(attr, val string)
	RemoveStyle(attr string)

	// Useful for text inputs
	SetValue(s string)
	Value() string

	SetSelectionStart(int)
	SelectionStart() int
	SetSelectionEnd(int)
	SelectionEnd() int

	CanvasContext(width, height, dpm float64) CanvasRenderingContext2D

	Node
}

// ElemResolver finds a DOM element by it's ID
//type ElemResolver func(id string) (Element, error)

type ElemConstructor func(template.HTML) Element

type EventType string

const (
	Click       EventType = "click"
	DoubleClick           = "dblclick"
	Input                 = "input"
	MouseDown             = "mousedown"
	MouseMove             = "mousemove"
	MouseOut              = "mouseout"
	MouseOver             = "mouseover"
	MouseUp               = "mouseup"
	MouseEnter            = "mouseenter"
	MouseLeave            = "mouseleave"
	KeyUp                 = "keyup"
	KeyDown               = "keydown"
	Drop                  = "drop"
	DragOver              = "dragover"
)

type CanvasRenderingContext2D interface {
}
