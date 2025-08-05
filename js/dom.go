//go:build js && wasm

package js

import (
	"fmt"
	"html/template"
	"log"

	"syscall/js"

	"github.com/nlandolfi/browser"
	"github.com/nlandolfi/browser/dom"
)

var (
	// this type checks all the implementations
	_ browser.Browser   = (*bbrowser)(nil)
	_ dom.Document      = (*document)(nil)
	_ dom.Element       = (*element)(nil)
	_ dom.EventListener = (*eventListener)(nil)
)

var DefaultBrowser browser.Browser = new(bbrowser)

type bbrowser struct {
	*document
}

func (b *bbrowser) doc() *document {
	if b.document == nil {
		b.document = &document{underlying: js.Global().Get("document")}
	}
	return b.document
}

func (b *bbrowser) Document() dom.Document {
	return b.doc()
}

func (b *bbrowser) NewElement(s template.HTML) dom.Element {
	t := b.doc().underlying.Call("createElement", "template")
	t.Set("innerHTML", string(s))
	return &element{
		underlying: t.Get("content").Get("firstChild"),
	}
}

type document struct {
	underlying js.Value
}

func (d *document) ReadyState() string {
	return d.underlying.Get("readyState").String()
}

func (d *document) Body() dom.Element {
	return &element{
		underlying: d.underlying.Get("body"),
	}
}

func (d *document) GetElementByID(id string) (dom.Element, error) {
	u := d.underlying.Call("getElementById", id)
	if u.Equal(js.Null()) {
		return nil, fmt.Errorf("element not found")
	}
	return &element{
		underlying: u,
	}, nil
}

func (d *document) CreateElement(tagName string) dom.Element {
	u := d.underlying.Call("createElement", tagName)
	return &element{
		underlying: u,
	}
}

func (e *document) CreateTextNode(s string) dom.Text {
	return &element{
		underlying: e.underlying.Call("createTextNode", s),
	}
}

func (e *document) Selection() dom.Selection {
	return &selection{
		underlying: e.underlying.Call("getSelection"),
	}
}

func (d *document) AddEventListener(on dom.EventType, h dom.EventHandler) dom.EventListener {
	c := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 1 {
			log.Printf("about to panic, args: %+v", args)
			panic("event listener called with more than one argument!")
		}
		h(&event{args[0]})
		return nil
	})

	el := &eventListener{Func: c}
	d.underlying.Call("addEventListener", string(on), c)

	return el
}

func (d *document) RemoveEventListener(on dom.EventType, l dom.EventListener) {
	if l == nil {
		panic("can not remove nil event listener")
	}
	el, ok := l.(*eventListener)
	if !ok {
		panic("bad event listener type")
	}

	d.underlying.Call("removeEventListener", string(on), el.Func)
}

type selection struct {
	underlying js.Value
}

func (s *selection) AnchorNode() dom.Node {
	return &element{underlying: s.underlying.Get("anchorNode")}
}

func (s *selection) AnchorOffset() int {
	return s.underlying.Get("anchorOffset").Int()
}

func (s *selection) FocusNode() dom.Node {
	return &element{underlying: s.underlying.Get("focusNode")}
}

func (s *selection) FocusOffset() int {
	return s.underlying.Get("focusOffset").Int()
}

func (s *selection) IsCollapsed() bool {
	return s.underlying.Get("isCollapsed").Bool()
}

func (s *selection) RangeCount() int {
	return s.underlying.Get("rangeCount").Int()
}

func (s *selection) Type() string {
	return s.underlying.Get("type").String()
}

type element struct {
	underlying js.Value

	// TODO: this may be broken? - NCL 2/19/2022
	s *js.Value // cached style object
}

func (e *element) SetInnerHTML(s template.HTML) {
	e.underlying.Set("innerHTML", js.ValueOf(string(s)))
}

func (e *element) ReplaceWith(o dom.Node) {
	j, ok := o.(*element)
	if !ok {
		panic("must be *element")
	}
	e.underlying.Call("replaceWith", j.underlying)
}

func (e *element) AddEventListener(on dom.EventType, h dom.EventHandler) dom.EventListener {
	c := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 1 {
			log.Printf("about to panic, args: %+v", args)
			panic("event listener called with more than one argument!")
		}
		h(&event{args[0]})
		return nil
	})

	el := &eventListener{Func: c}
	e.underlying.Call("addEventListener", string(on), c)

	return el
}

func (e *element) RemoveEventListener(on dom.EventType, l dom.EventListener) {
	if l == nil {
		panic("can not remove nil event listener")
	}
	el, ok := l.(*eventListener)
	if !ok {
		panic("bad event listener type")
	}

	e.underlying.Call("removeEventListener", string(on), el.Func)
}

func (e *element) AppendChild(n dom.Node) {
	j, ok := n.(*element)
	if !ok {
		panic("must be *element")
	}
	e.underlying.Call("appendChild", j.underlying)
}

func (e *element) RemoveChild(n dom.Node) dom.Node {
	j, ok := n.(*element)
	if !ok {
		panic("must be *element")
	}
	e.underlying.Call("removeChild", j.underlying)
	return n
}

func (e *element) SetAttribute(key, val string) {
	e.underlying.Call("setAttribute", key, val)
}

// https://developer.mozilla.org/en-US/docs/Web/API/Element/removeAttribute
func (e *element) RemoveAttribute(key string) {
	e.underlying.Call("removeAttribute", key)
}

func (e *element) style() js.Value {
	if e.s == nil {
		v := e.underlying.Get("style")
		e.s = &v
	}

	return *e.s
}

func (e *element) SetStyle(attr, val string) {
	e.style().Set(attr, val)
}

func (e *element) RemoveStyle(attr string) {
	e.style().Delete(attr)
}

func (e *element) LogSelf() {
	js.Global().Get("console").Call("log", e.underlying)
}

// https://developer.mozilla.org/en-US/docs/Web/API/Node/replaceChild
func (e *element) ReplaceChild(n, old dom.Node) dom.Node {
	un, ok := n.(*element)
	if !ok {
		panic("n not *element")
	}
	uold, ok := old.(*element)
	if !ok {
		panic("old not *element")
	}
	e.underlying.Call("replaceChild", un.underlying, uold.underlying)
	return old
}

func (e *element) SetValue(s string) {
	e.underlying.Set("value", s)
}

func (e *element) Value() string {
	return e.underlying.Get("value").String()
}

func (e *element) SetSelectionStart(i int) {
	e.underlying.Set("selectionStart", i)
}

func (e *element) SelectionStart() int {
	return e.underlying.Get("selectionStart").Int()
}

func (e *element) SelectionEnd() int {
	return e.underlying.Get("selectionEnd").Int()
}

func (e *element) SetSelectionEnd(i int) {
	e.underlying.Set("selectionEnd", i)
}

func (e *element) CanvasContext(width, height, dpm float64) dom.CanvasRenderingContext2D {
	return &canvascontext{underlying: e.underlying.Call("getContext", "2d")}
}

type canvascontext struct {
	underlying js.Value
}

// See: https://developer.mozilla.org/en-US/docs/Web/API/Node/parentElement
func (e *element) ParentElement() dom.Element {
	x := e.underlying.Get("parentElement")
	if x.IsNull() {
		return nil
	}
	return &element{underlying: x}
}

type eventListener struct {
	js.Func
}

func (el *eventListener) Release() {
	log.Print("not implemented")
}

type event struct {
	underlying js.Value
}

func (e *event) IsUndefined() bool {
	return e.underlying.IsUndefined()
}

func (e *event) DataTransfer() dom.DataTransfer {
	return &dataTransfer{underlying: e.underlying.Get("dataTransfer")}
}

func (e *event) Target() dom.Element {
	return &element{underlying: e.underlying.Get("target")}
}

func (e *event) OffsetX() int {
	return e.underlying.Get("clientX").Int()
}

func (e *event) OffsetY() int {
	return e.underlying.Get("clientY").Int()
}

func (e *event) PageX() int {
	return e.underlying.Get("pageX").Int()
}

func (e *event) PageY() int {
	return e.underlying.Get("pageY").Int()
}

func (e *event) ClientX() int {
	return e.underlying.Get("clientX").Int()
}

func (e *event) ClientY() int {
	return e.underlying.Get("clientY").Int()
}

func (e *event) MovementX() int {
	return e.underlying.Get("movementX").Int()
}

func (e *event) MovementY() int {
	return e.underlying.Get("movementY").Int()
}

func (e *event) Code() string {
	return e.underlying.Get("code").String()
}

func (e *event) KeyCode() int {
	u := e.underlying.Get("keyCode")
	// yikes
	if u.IsUndefined() {
		return -1 // TODO is this mellow?
	}
	return u.Int()
}

func (e *event) PreventDefault() {
	e.underlying.Call("preventDefault")
}

func (e *event) StopPropagation() {
	e.underlying.Call("stopPropagation")
}

type dataTransfer struct {
	underlying js.Value
}

func (d *dataTransfer) Items() []dom.DataTransferItem {
	log.Print("dataTransfer items not implemented!")
	return nil
}
