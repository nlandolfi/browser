package browser

import (
	"bytes"
	"fmt"
	"io"

	"github.com/nlandolfi/browser/dom"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Node struct {
	Type     html.NodeType
	DataAtom atom.Atom
	Data     string
	Attr     []*html.Attribute

	Style      Style
	Handlers   Handlers
	CanvasDraw func(ctx dom.CanvasRenderingContext2D)

	Children []*Node

	// these are used by the Mounter
	rendered        dom.Node
	renderedElement dom.Element
}

// Style {{{

type AlignItemsType int

const (
	AlignItemsUnset AlignItemsType = iota
	AlignItemsCenter
)

func (t AlignItemsType) String() string {
	switch t {
	case AlignItemsUnset:
		return ""
	case AlignItemsCenter:
		return "center"
	}

	panic(fmt.Sprintf("unknown AlignItemsType: %#v", t))
}

type BorderType int

const (
	BorderUnset BorderType = iota
	BorderNone
	BorderSolid
)

func (t BorderType) String() string {
	switch t {
	case BorderUnset:
		return ""
	case BorderNone:
		return "none"
	case BorderSolid:
		return "solid"
	}

	panic(fmt.Sprintf("unknown BorderType: %#v", t))
}

type Border struct {
	Width Size
	Type  BorderType
	Color string
}

func (b *Border) String() string {
	if b.Type == BorderNone {
		return b.Type.String()
	}
	return fmt.Sprintf("%s %s %s", &b.Width, b.Type, b.Color)
}

type BoxShadow struct {
	HOffset, VOffset, Blur, Spread Size
	Color                          string
}

func (b *BoxShadow) IsZero() bool {
	return b.HOffset == Size{} && b.VOffset == Size{} && b.Blur == Size{} && b.Spread == Size{} && b.Color == ""
}

func (b *BoxShadow) Encode(w io.Writer) {
	fmt.Fprintf(w, "%s %s %s %s %s", &b.HOffset, &b.VOffset, &b.Blur, &b.Spread, b.Color)
}

func (b *BoxShadow) String() string {
	var buf bytes.Buffer
	b.Encode(&buf)
	return buf.String()
}

type CursorType int

const (
	CursorUnset CursorType = iota
	CursorDefault
	CursorPointer
	CursorMove
	CursorNWSEResize
	CursorEWResize
)

func (t CursorType) String() string {
	switch t {
	case CursorUnset:
		return ""
	case CursorPointer:
		return "pointer"
	case CursorMove:
		return "move"
	case CursorDefault:
		return "default"
	case CursorNWSEResize:
		return "nwse-resize"
	case CursorEWResize:
		return "ew-resize"
	}

	panic(fmt.Sprintf("unknown CursorType: %#v", t))
}

type DisplayType int

const (
	DisplayUnset = iota
	DisplayNone
	DisplayFlex
	DisplayGrid
)

func (t DisplayType) String() string {
	switch t {
	case DisplayUnset:
		return ""
	case DisplayNone:
		return "none"
	case DisplayFlex:
		return "flex"
	case DisplayGrid:
		return "grid"
	}

	panic(fmt.Sprintf("unknown DisplayType: %#v", t))
}

type FlexDirectionType int

const (
	FlexDirectionUnset FlexDirectionType = iota
	FlexDirectionRow
	FlexDirectionColumn
	FlexDirectionRowReverse
	FlexDirectionColumnReverse
)

func (t FlexDirectionType) String() string {
	switch t {
	case FlexDirectionUnset:
		return ""
	case FlexDirectionRow:
		return "row"
	case FlexDirectionColumn:
		return "column"
	case FlexDirectionRowReverse:
		return "row-reverse"
	case FlexDirectionColumnReverse:
		return "column-reverse"
	}

	panic(fmt.Sprintf("unknown FlexDirectionType: %#v", t))
}

type FlexWrapType int

const (
	FlexWrapUnset FlexWrapType = iota
	FlexWrapWrap
)

func (t FlexWrapType) String() string {
	switch t {
	case FlexWrapUnset:
		return ""
	case FlexWrapWrap:
		return "wrap"
	}

	panic(fmt.Sprintf("unknown FlexWrapType: %#v", t))
}

type UnitType int

const (
	UnitUnset UnitType = iota
	UnitDefault
	UnitEM
	UnitPC
	UnitPG
	UnitPT
	UnitPX
	UnitVH
	UnitVW
)

func (t UnitType) String() string {
	switch t {
	case UnitUnset, UnitDefault:
		return ""
	case UnitEM:
		return "em"
	case UnitPC:
		return "pc"
	case UnitPG:
		return "%"
	case UnitPX:
		return "px"
	case UnitVH:
		return "vh"
	case UnitVW:
		return "vw"
	}

	panic(fmt.Sprintf("unknown UnitType: %#v", t))
}

type Size struct {
	Value          float64
	Unit           UnitType
	StringOverride string // Use for calcs and vars
}

func (s *Size) IsZero() bool {
	return s.Value == 0 && s.Unit == UnitUnset && s.StringOverride == ""
}

func (s *Size) String() string {
	if s.StringOverride != "" {
		return s.StringOverride
	} else {
		return fmt.Sprintf("%f%s", s.Value, s.Unit)
	}
}

type JustifyContentType int

const (
	JustifyContentUnset JustifyContentType = iota
	JustifyContentCenter
	JustifyContentSpaceBetween
	JustifyContentFlexEnd
)

func (t JustifyContentType) String() string {
	switch t {
	case JustifyContentUnset:
		// some wierd bug where this is getting called in (*Style.Val) -2/4/2022
		//panic("called String on JustifyContentUnset")
		//		log.Print("need to fix JustifyContentUnset weirdness")
		// TODO this does not make sense to me - NCL -2/3/2022
		//return "unset" // hopefully changing this to empty string doesn't break everything
		return ""
	case JustifyContentCenter:
		return "center"
	case JustifyContentSpaceBetween:
		return "space-between"
	case JustifyContentFlexEnd:
		return "flex-end"
	}

	panic(fmt.Sprintf("unknown JustifyContentType: %#v", t))
}

type JustifySelfType int

const (
	JustifySelfUnset JustifySelfType = iota
	JustifySelfRight
	JustifySelfLeft
	JustifySelfFlexEnd
	JustifySelfFlexStart
)

func (t JustifySelfType) String() string {
	switch t {
	case JustifySelfUnset:
		return ""
	case JustifySelfRight:
		return "right"
	case JustifySelfLeft:
		return "left"
	case JustifySelfFlexEnd:
		return "flex-end"
	case JustifySelfFlexStart:
		return "flex-start"
	}

	panic(fmt.Sprintf("unknown JustifySelfType: %#v", t))
}

type OverflowType int

const (
	OverflowUnset OverflowType = iota
	OverflowHidden
	OverflowScroll
	OverflowAuto
)

func (t OverflowType) String() string {
	switch t {
	case OverflowUnset:
		//panic("called String on OverflowUnset")
		return ""
	case OverflowHidden:
		return "hidden"
	case OverflowScroll:
		return "scroll"
	case OverflowAuto:
		return "auto"
	}

	panic(fmt.Sprintf("unknown OverflowType: %#v", t))
}

type PositionType int

const (
	PositionUnset PositionType = iota
	PositionRelative
	PositionAbsolute
)

func (t PositionType) String() string {
	switch t {
	case PositionUnset:
		return ""
	case PositionRelative:
		return "relative"
	case PositionAbsolute:
		return "absolute"
	}

	panic(fmt.Sprintf("unknown PositionType: %#v", t))
}

type TextAlignType int

const (
	TextAlignUnset TextAlignType = iota
	TextAlignCenter
	TextAlignRight
	TextAlignLeft
	TextAlignJustify
)

func (t TextAlignType) String() string {
	switch t {
	case TextAlignUnset:
		return ""
	case TextAlignCenter:
		return "center"
	case TextAlignRight:
		return "right"
	case TextAlignLeft:
		return "left"
	case TextAlignJustify:
		return "justify"
	}

	panic(fmt.Sprintf("unknown TextAlignType: %#v", t))
}

type TextDecorationType int

const (
	TextDecorationUnset TextDecorationType = iota
	TextDecorationLineThrough
	TextDecorationOverline
	TextDecorationUnderline
)

func (t TextDecorationType) String() string {
	switch t {
	case TextDecorationUnset:
		return ""
	case TextDecorationLineThrough:
		return "line-through"
	case TextDecorationOverline:
		return "overline"
	case TextDecorationUnderline:
		return "underline"
	}

	panic(fmt.Sprintf("unknown TextDecorationType: %#v", t))
}

type OutlineType int

const (
	OutlineUnset OutlineType = iota
	OutlineNone
	OutlineSolid
)

func (t OutlineType) String() string {
	switch t {
	case OutlineUnset:
		return ""
	case OutlineNone:
		return "none"
	case OutlineSolid:
		return "solid"
	}

	panic(fmt.Sprintf("unknown OutlineType: %#v", t))
}

type Outline struct {
	Width Size
	Type  OutlineType
	Color string
}

func (o *Outline) IsZero() bool {
	return o.Width.IsZero() && o.Type == OutlineUnset && o.Color == ""
}

func (o *Outline) String() string {
	if o.Type == OutlineNone {
		return o.Type.String()
	}

	return fmt.Sprintf("%s %s %s", &o.Width, o.Type, o.Color)
}

type Style struct {
	AlignItems      AlignItemsType
	Background      string
	BackgroundColor string
	Border          Border
	BorderBottom    Border
	BorderColor     string
	BorderLeft      Border
	BorderRadius    Size
	BorderRight     Border
	BorderTop       Border
	BoxShadow       BoxShadow
	Color           string
	Cursor          CursorType
	Display         DisplayType
	FlexBasis       string
	FlexDirection   FlexDirectionType
	FlexGrow        string
	FlexShrink      string
	FlexWrap        FlexWrapType
	FontFamily      string
	FontSize        Size
	FontWeight      string
	GridArea        string
	Height          Size
	JustifyContent  JustifyContentType
	JustifySelf     JustifySelfType
	Left            Size
	Margin          Size
	MarginBottom    Size
	MarginLeft      Size
	MarginRight     Size
	MarginTop       Size
	MaxHeight       Size
	MaxWidth        Size
	MinHeight       Size
	MinWidth        Size
	Outline         Outline
	Overflow        OverflowType
	Padding         Size
	PaddingBottom   Size
	PaddingLeft     Size
	PaddingRight    Size
	PaddingTop      Size
	Position        PositionType
	TextAlign       TextAlignType
	TextDecoration  TextDecorationType
	Top             Size
	Transform       string
	Transition      string
	UserSelect      string
	Width           Size
}

func (s *Style) Val() string {
	if s == nil || (*s == Style{}) {
		return ""
	}

	var buf bytes.Buffer
	w := &buf

	if s.AlignItems != AlignItemsUnset {
		fmt.Fprintf(w, "align-items:%s;", s.AlignItems)
	}

	if s.Background != "" {
		fmt.Fprintf(w, "background:%s;", s.Background)
	}

	if s.BackgroundColor != "" {
		fmt.Fprintf(w, "background-color:%s;", s.BackgroundColor)
	}

	if s.Border.Type != BorderUnset {
		fmt.Fprintf(w, "border:%s;", &s.Border) // maybe use Border.Encode
	}

	if s.BorderTop.Type != BorderUnset {
		fmt.Fprintf(w, "border-top:%s;", &s.BorderTop) // maybe use Border.Encode
	}

	if s.BorderBottom.Type != BorderUnset {
		fmt.Fprintf(w, "border-bottom:%s;", &s.BorderBottom) // maybe use Border.Encode
	}

	if s.BorderLeft.Type != BorderUnset {
		fmt.Fprintf(w, "border-left:%s;", &s.BorderLeft) // maybe use Border.Encode
	}

	if s.BorderRight.Type != BorderUnset {
		fmt.Fprintf(w, "border-right:%s;", &s.BorderRight) // maybe use Border.Encode
	}

	if s.BorderColor != "" {
		fmt.Fprintf(w, "border-color:%s;", s.BorderColor)
	}

	if !s.BorderRadius.IsZero() {
		fmt.Fprintf(w, "border-radius:%s;", &s.BorderRadius)
	}

	if !s.BoxShadow.IsZero() {
		fmt.Fprintf(w, "box-shadow:%s;", &s.BoxShadow)
	}

	if s.Color != "" {
		fmt.Fprintf(w, "color:%s;", s.Color)
	}

	if s.Cursor != CursorUnset {
		fmt.Fprintf(w, "cursor:%s;", s.Cursor)
	}

	if s.Display != DisplayUnset {
		fmt.Fprintf(w, "display:%s;", s.Display)
	}

	if s.FlexDirection != FlexDirectionUnset {
		fmt.Fprintf(w, "flex-direction:%s;", s.FlexDirection)
	}

	if s.FlexGrow != "" {
		fmt.Fprintf(w, "flex-grow:%s;", s.FlexGrow)
	}

	if s.FlexBasis != "" {
		fmt.Fprintf(w, "flex-basis:%s;", s.FlexBasis)
	}

	if s.FlexShrink != "" {
		fmt.Fprintf(w, "flex-shrink:%s;", s.FlexShrink)
	}

	if s.FlexWrap != FlexWrapUnset {
		fmt.Fprintf(w, "flex-wrap:%s;", s.FlexWrap)
	}

	if s.FontFamily != "" {
		fmt.Fprintf(w, "font-family:%s;", s.FontFamily)
	}

	if !s.FontSize.IsZero() {
		fmt.Fprintf(w, "font-size:%s;", &s.FontSize)
	}

	if s.FontWeight != "" {
		fmt.Fprintf(w, "font-weight:%s;", s.FontWeight)
	}

	if s.GridArea != "" {
		fmt.Fprintf(w, "grid-area:%s;", s.GridArea)
	}

	if s.JustifyContent != JustifyContentUnset {
		fmt.Fprintf(w, "justify-content:%s;", s.JustifyContent)
	}

	if s.JustifySelf != JustifySelfUnset {
		fmt.Fprintf(w, "justify-self:%s;", s.JustifySelf)
	}

	if !s.Height.IsZero() {
		fmt.Fprintf(w, "height:%s;", &s.Height)
	}

	if !s.Left.IsZero() {
		fmt.Fprintf(w, "left:%s;", &s.Left)
	}

	if !s.Margin.IsZero() {
		fmt.Fprintf(w, "margin:%s;", &s.Margin)
	}

	if !s.MarginBottom.IsZero() {
		fmt.Fprintf(w, "margin-bottom:%s;", &s.MarginBottom)
	}

	if !s.MarginLeft.IsZero() {
		fmt.Fprintf(w, "margin-left:%s;", &s.MarginLeft)
	}

	if !s.MarginRight.IsZero() {
		fmt.Fprintf(w, "margin-right:%s;", &s.MarginRight)
	}

	if !s.MarginTop.IsZero() {
		fmt.Fprintf(w, "margin-top:%s;", &s.MarginTop)
	}

	if !s.MaxHeight.IsZero() {
		fmt.Fprintf(w, "max-height:%s;", &s.MaxHeight)
	}

	if !s.MinHeight.IsZero() {
		fmt.Fprintf(w, "min-height:%s;", &s.MinHeight)
	}

	if !s.MaxWidth.IsZero() {
		fmt.Fprintf(w, "max-width:%s;", &s.MaxWidth)
	}

	if !s.MinWidth.IsZero() {
		fmt.Fprintf(w, "min-width:%s;", &s.MinWidth)
	}

	if !s.Outline.IsZero() {
		fmt.Fprintf(w, "outline:%s;", &s.Outline)
	}

	if s.Overflow != OverflowUnset {
		fmt.Fprintf(w, "overflow:%s;", s.Overflow)
	}

	if !s.Padding.IsZero() {
		fmt.Fprintf(w, "padding:%s;", &s.Padding)
	}

	if !s.PaddingBottom.IsZero() {
		fmt.Fprintf(w, "padding-bottom:%s;", &s.PaddingBottom)
	}

	if !s.PaddingLeft.IsZero() {
		fmt.Fprintf(w, "padding-left:%s;", &s.PaddingLeft)
	}

	if !s.PaddingRight.IsZero() {
		fmt.Fprintf(w, "padding-right:%s;", &s.PaddingRight)
	}

	if !s.PaddingTop.IsZero() {
		fmt.Fprintf(w, "padding-top:%s;", &s.PaddingTop)
	}

	if s.Position != PositionUnset {
		fmt.Fprintf(w, "position:%s;", s.Position)
	}

	if s.TextAlign != TextAlignUnset {
		fmt.Fprintf(w, "text-align:%s;", s.TextAlign)
	}

	if s.TextDecoration != TextDecorationUnset {
		fmt.Fprintf(w, "text-decoration:%s;", s.TextDecoration)
	}

	if !s.Top.IsZero() {
		fmt.Fprintf(w, "top:%s;", &s.Top)
	}

	if s.Transform != "" {
		fmt.Fprintf(w, "transform:%s;", s.Transform)
	}

	if s.Transition != "" {
		fmt.Fprintf(w, "transition:%s;", s.Transition)
	}

	if s.UserSelect != "" {
		fmt.Fprintf(w, "user-select:%s;", s.UserSelect)
		fmt.Fprintf(w, "-webkit-user-select:%s;", s.UserSelect)
		fmt.Fprintf(w, "-moz-user-select:%s;", s.UserSelect)
		fmt.Fprintf(w, "-ms-user-select:%s;", s.UserSelect)
	}

	if !s.Width.IsZero() {
		fmt.Fprintf(w, "width:%s;", &s.Width)
	}

	return buf.String()
}

// }}}

// Style Helpers (e.g., WidthPX, BorderRadiusPX) {{{

func (n *Node) AlignItems(s AlignItemsType) *Node { n.Style.AlignItems = s; return n }
func (n *Node) AlignItemsCenter() *Node           { return n.AlignItems(AlignItemsCenter) }
func (n *Node) Background(s string) *Node         { n.Style.Background = s; return n }
func (n *Node) BackgroundColor(s string) *Node    { n.Style.BackgroundColor = s; return n }
func (n *Node) Border(b Border) *Node             { n.Style.Border = b; return n }
func (n *Node) BorderBottom(b Border) *Node       { n.Style.BorderBottom = b; return n }
func (n *Node) BorderLeft(b Border) *Node         { n.Style.BorderLeft = b; return n }
func (n *Node) BorderRadius(s Size) *Node         { n.Style.BorderRadius = s; return n }
func (n *Node) BorderRadiusPX(s float64) *Node    { return n.BorderRadius(Size{Value: s, Unit: UnitPX}) }
func (n *Node) BorderRight(b Border) *Node        { n.Style.BorderRight = b; return n }
func (n *Node) BorderTop(b Border) *Node          { n.Style.BorderTop = b; return n }
func (n *Node) BoxShadow(b BoxShadow) *Node       { n.Style.BoxShadow = b; return n }
func (n *Node) Color(s string) *Node              { n.Style.Color = s; return n }
func (n *Node) Display(t DisplayType) *Node       { n.Style.Display = t; return n }
func (n *Node) FlexBasis(s string) *Node          { n.Style.FlexBasis = s; return n }
func (n *Node) FlexCenter() *Node {
	return n.Display(DisplayFlex).JustifyContentCenter().AlignItemsCenter()
}
func (n *Node) FlexDirection(s FlexDirectionType) *Node   { n.Style.FlexDirection = s; return n }
func (n *Node) FlexGrow(s string) *Node                   { n.Style.FlexGrow = s; return n }
func (n *Node) FlexShrink(s string) *Node                 { n.Style.FlexShrink = s; return n }
func (n *Node) FlexWrap(s FlexWrapType) *Node             { n.Style.FlexWrap = s; return n }
func (n *Node) FontFamily(s string) *Node                 { n.Style.FontFamily = s; return n }
func (n *Node) FontSize(s Size) *Node                     { n.Style.FontSize = s; return n }
func (n *Node) FontSizeEM(v float64) *Node                { return n.FontSize(Size{Value: v, Unit: UnitEM}) }
func (n *Node) FontSizePX(v float64) *Node                { return n.FontSize(Size{Value: v, Unit: UnitPX}) }
func (n *Node) FontWeight(s string) *Node                 { n.Style.FontWeight = s; return n }
func (n *Node) GridArea(s string) *Node                   { n.Style.GridArea = s; return n }
func (n *Node) JustifyContent(s JustifyContentType) *Node { n.Style.JustifyContent = s; return n }
func (n *Node) JustifyContentCenter() *Node               { return n.JustifyContent(JustifyContentCenter) }
func (n *Node) JustifySelf(s JustifySelfType) *Node       { n.Style.JustifySelf = s; return n }
func (n *Node) MarginBottom(s Size) *Node                 { n.Style.MarginBottom = s; return n }
func (n *Node) MarginBottomPX(f float64) *Node            { return n.MarginBottom(Size{Value: f, Unit: UnitPX}) }
func (n *Node) MarginLeft(s Size) *Node                   { n.Style.MarginLeft = s; return n }
func (n *Node) MarginLeftPX(f float64) *Node              { return n.MarginLeft(Size{Value: f, Unit: UnitPX}) }
func (n *Node) MarginPX(pxs float64) *Node                { return n.Margin(Size{Value: pxs, Unit: UnitPX}) }
func (n *Node) MarginRight(s Size) *Node                  { n.Style.MarginRight = s; return n }
func (n *Node) MarginRightPX(f float64) *Node             { return n.MarginRight(Size{Value: f, Unit: UnitPX}) }
func (n *Node) MarginTop(s Size) *Node                    { n.Style.MarginTop = s; return n }
func (n *Node) MarginTopPX(f float64) *Node               { return n.MarginTop(Size{Value: f, Unit: UnitPX}) }
func (n *Node) JustifyContentSpaceBetween() *Node {
	return n.JustifyContent(JustifyContentSpaceBetween)
}
func (n *Node) Height(s Size) *Node           { n.Style.Height = s; return n }
func (n *Node) HeightPX(v float64) *Node      { return n.Height(Size{Value: v, Unit: UnitPX}) }
func (n *Node) HeightPG(v float64) *Node      { return n.Height(Size{Value: v, Unit: UnitPG}) }
func (n *Node) HeightVH(v float64) *Node      { return n.Height(Size{Value: v, Unit: UnitVH}) }
func (n *Node) MaxHeight(s Size) *Node        { n.Style.MaxHeight = s; return n }
func (n *Node) Width(s Size) *Node            { n.Style.Width = s; return n }
func (n *Node) WidthPX(f float64) *Node       { n.Style.Width = Size{Value: f, Unit: UnitPX}; return n }
func (n *Node) WidthPG(f float64) *Node       { n.Style.Width = Size{Value: f, Unit: UnitPG}; return n }
func (n *Node) WidthVW(f float64) *Node       { n.Style.Width = Size{Value: f, Unit: UnitVW}; return n }
func (n *Node) Left(s Size) *Node             { n.Style.Left = s; return n }
func (n *Node) MaxWidth(s Size) *Node         { n.Style.MaxWidth = s; return n }
func (n *Node) MaxWidthPX(v float64) *Node    { return n.MaxWidth(Size{Value: v, Unit: UnitPX}) }
func (n *Node) MinHeight(s Size) *Node        { n.Style.MinHeight = s; return n }
func (n *Node) MinWidth(s Size) *Node         { n.Style.MinWidth = s; return n }
func (n *Node) MinWidthPX(v float64) *Node    { return n.MinWidth(Size{Value: v, Unit: UnitPX}) }
func (n *Node) Cursor(t CursorType) *Node     { n.Style.Cursor = t; return n }
func (n *Node) CursorPointer() *Node          { return n.Cursor(CursorPointer) }
func (n *Node) CursorMove() *Node             { return n.Cursor(CursorMove) }
func (n *Node) Outline(o Outline) *Node       { n.Style.Outline = o; return n }
func (n *Node) Overflow(t OverflowType) *Node { n.Style.Overflow = t; return n }
func (n *Node) OverflowHidden() *Node         { return n.Overflow(OverflowHidden) }
func (n *Node) OverflowScroll() *Node         { return n.Overflow(OverflowScroll) }
func (n *Node) Margin(s Size) *Node           { n.Style.Margin = s; return n }
func (n *Node) PaddingBottom(s Size) *Node    { n.Style.PaddingBottom = s; return n }
func (n *Node) Pointer() *Node                { return n.Cursor(CursorPointer) }
func (n *Node) Position(t PositionType) *Node { n.Style.Position = t; return n }
func (n *Node) PositionAbsolute() *Node       { n.Style.Position = PositionAbsolute; return n }
func (n *Node) PositionRelative() *Node       { n.Style.Position = PositionRelative; return n }
func (n *Node) Padding(s Size) *Node          { n.Style.Padding = s; return n }
func (n *Node) PaddingPX(pxs float64) *Node   { return n.Padding(Size{Value: pxs, Unit: UnitPX}) }
func (n *Node) PaddingRight(s Size) *Node     { n.Style.PaddingRight = s; return n }
func (n *Node) PaddingRightPX(pxs float64) *Node {
	return n.PaddingRight(Size{Value: pxs, Unit: UnitPX})
}
func (n *Node) PaddingLeft(s Size) *Node { n.Style.PaddingLeft = s; return n }
func (n *Node) PaddingLeftPX(pxs float64) *Node {
	return n.PaddingLeft(Size{Value: pxs, Unit: UnitPX})
}
func (n *Node) PaddingTop(s Size) *Node        { n.Style.PaddingTop = s; return n }
func (n *Node) PaddingTopPX(pxs float64) *Node { return n.PaddingTop(Size{Value: pxs, Unit: UnitPX}) }
func (n *Node) PaddingBottomPX(pxs float64) *Node {
	return n.PaddingBottom(Size{Value: pxs, Unit: UnitPX})
}
func (n *Node) TextAlign(t TextAlignType) *Node           { n.Style.TextAlign = t; return n }
func (n *Node) TextAlignCenter() *Node                    { return n.TextAlign(TextAlignCenter) }
func (n *Node) TextAlignRight() *Node                     { return n.TextAlign(TextAlignRight) }
func (n *Node) TextDecoration(t TextDecorationType) *Node { n.Style.TextDecoration = t; return n }
func (n *Node) TextDecorationUnderline() *Node {
	n.Style.TextDecoration = TextDecorationUnderline
	return n
}
func (n *Node) Top(s Size) *Node { n.Style.Top = s; return n }

func (n *Node) TopPX(s float64) *Node     { return n.Top(Size{Value: s, Unit: UnitPX}) }
func (n *Node) Transform(s string) *Node  { n.Style.Transform = s; return n }
func (n *Node) Transition(s string) *Node { n.Style.Transition = s; return n }
func (n *Node) UserSelect(s string) *Node { n.Style.UserSelect = s; return n }
func (n *Node) LeftPX(s float64) *Node    { return n.Left(Size{Value: s, Unit: UnitPX}) }

func (n *Node) OnlyIf(b bool, f func(n *Node) *Node) *Node {
	if b {
		return f(n)
	}

	return n
}

// }}}

// Attr Helpers (e.g. ID) {{{

func (n *Node) ID(id string) *Node {
	for _, a := range n.Attr {
		if a.Key == atom.Id.String() {
			a.Val = id
			return n
		}
	}

	n.Attr = append(n.Attr,
		&html.Attribute{
			Key: atom.Id.String(),
			Val: id,
		},
	)

	return n
}

func (n *Node) Draggable() *Node {
	for _, a := range n.Attr {
		if a.Key == atom.Draggable.String() {
			a.Val = "true"
			return n
		}
	}

	n.Attr = append(n.Attr,
		&html.Attribute{
			Key: atom.Draggable.String(),
			Val: "true",
		},
	)

	return n
}

func (n *Node) NotDraggable() *Node {
	for _, a := range n.Attr {
		if a.Key == atom.Draggable.String() {
			a.Val = "false"
			return n
		}
	}

	n.Attr = append(n.Attr,
		&html.Attribute{
			Key: atom.Draggable.String(),
			Val: "false",
		},
	)

	return n
}

func (n *Node) Min(min int) *Node {
	s := fmt.Sprintf("%d", min)

	for _, a := range n.Attr {
		if a.Key == atom.Min.String() {
			a.Val = s
			return n
		}
	}

	n.Attr = append(n.Attr,
		&html.Attribute{
			Key: atom.Min.String(),
			Val: s,
		},
	)
	return n
}

func (n *Node) Max(max int) *Node {
	s := fmt.Sprintf("%d", max)

	for _, a := range n.Attr {
		if a.Key == atom.Max.String() {
			a.Val = s
			return n
		}
	}

	n.Attr = append(n.Attr,
		&html.Attribute{
			Key: atom.Max.String(),
			Val: s,
		},
	)
	return n
}

func (n *Node) Class(id string) *Node {
	for _, a := range n.Attr {
		if a.Key == atom.Class.String() {
			a.Val = id
			return n
		}
	}

	n.Attr = append(n.Attr,
		&html.Attribute{
			Key: atom.Class.String(),
			Val: id,
		},
	)

	return n
}

func (n *Node) Placeholder(p string) *Node {
	for _, a := range n.Attr {
		if a.Key == atom.Placeholder.String() {
			a.Val = p
			return n
		}
	}

	n.Attr = append(n.Attr,
		&html.Attribute{
			Key: atom.Placeholder.String(),
			Val: p,
		},
	)

	return n
}

func (n *Node) AttrType(p string) *Node {
	for _, a := range n.Attr {
		if a.Key == atom.Type.String() {
			a.Val = p
			return n
		}
	}

	n.Attr = append(n.Attr,
		&html.Attribute{
			Key: atom.Type.String(),
			Val: p,
		},
	)

	return n
}

func (n *Node) AttrValue(v string) *Node {
	for _, a := range n.Attr {
		if a.Key == atom.Value.String() {
			a.Val = v
			return n
		}
	}

	n.Attr = append(n.Attr,
		&html.Attribute{
			Key: atom.Value.String(),
			Val: v,
		},
	)

	return n
}

func (n *Node) AddAttr(a *html.Attribute) *Node { n.Attr = append(n.Attr, a); return n }

// }}}

// Handlers {{{

type Handlers struct {
	Click       dom.EventHandler `json:"-"`
	DoubleClick dom.EventHandler `json:"-"`
	Drag        dom.EventHandler `json:"-"`
	Input       dom.EventHandler `json:"-"`
	MouseOut    dom.EventHandler `json:"-"`
	MouseOver   dom.EventHandler `json:"-"`
	MouseDown   dom.EventHandler `json:"-"`
	MouseEnter  dom.EventHandler `json:"-"`
	MouseLeave  dom.EventHandler `json:"-"`
	MouseUp     dom.EventHandler `json:"-"`
	MouseMove   dom.EventHandler `json:"-"`
	KeyUp       dom.EventHandler `json:"-"`
	KeyDown     dom.EventHandler `json:"-"`
	Drop        dom.EventHandler `json:"-"`
	DragOver    dom.EventHandler `json:"-"`

	ClickCacheKey       string
	DragCacheKey        string
	DoubleClickCacheKey string
	InputCacheKey       string
	MouseOverCacheKey   string
	MouseOutCacheKey    string
	MouseDownCacheKey   string
	MouseEnterCacheKey  string
	MouseLeaveCacheKey  string
	MouseUpCacheKey     string
	MouseMoveCacheKey   string
	KeyUpCacheKey       string
	KeyDownCacheKey     string
	DropCacheKey        string
	DragOverCacheKey    string

	click, drag, doubleClick, input, mouseOut, mouseOver, mouseDown, mouseUp, mouseMove, mouseEnter, mouseLeave, keyUp, keyDown, drop, dragOver dom.EventListener
}

// }}}

// Handler Helpers (e.g., OnInput) {{{

func (n *Node) OnInput(f dom.EventHandler) *Node {
	n.Handlers.Input = f
	return n
}

// See OnClickKeyCache
func (n *Node) OnInputKeyCache(s string) *Node { n.Handlers.InputCacheKey = s; return n }

func (n *Node) OnClick(f dom.EventHandler) *Node {
	n.Handlers.Click = f
	return n
}

// This is a key which will be checked when diffing the nodes.
// If it is the same between two nodes, their Handlers.OnClick
// function values will be considered equivalent.
func (n *Node) OnClickKeyCache(s string) *Node {
	n.Handlers.ClickCacheKey = s
	return n
}

func (n *Node) OnDrag(f dom.EventHandler) *Node {
	n.Handlers.Drag = f
	return n
}

func (n *Node) OnDragKeyCache(s string) *Node {
	n.Handlers.DragCacheKey = s
	return n
}

func (n *Node) OnDragDispatch(e Event) *Node {
	return n.OnDrag(func(_ dom.Event) { go Dispatch(e) })
}

func (n *Node) OnMouseOver(f dom.EventHandler) *Node {
	n.Handlers.MouseOver = f
	return n
}

// See OnClickKeyCache
func (n *Node) OnMouseOverKeyCache(s string) *Node {
	n.Handlers.MouseOverCacheKey = s
	return n
}

func (n *Node) OnMouseOut(f dom.EventHandler) *Node {
	n.Handlers.MouseOut = f
	return n
}

// See OnClickKeyCache
func (n *Node) OnMouseOutKeyCache(s string) *Node {
	n.Handlers.MouseOutCacheKey = s
	return n
}

// See OnClickKeyCached
func (n *Node) OnInputCached(s string, f dom.EventHandler) *Node {
	n.Handlers.InputCacheKey = s
	n.Handlers.Input = f
	return n
}

// This is a key which will be checked when diffing the nodes.
// If it is the same between two nodes, their Handlers.OnClick
// function values will be considered equivalent.
func (n *Node) OnClickCached(s string, f dom.EventHandler) *Node {
	n.Handlers.ClickCacheKey = s
	n.Handlers.Click = f
	return n
}

// See OnClickCached
func (n *Node) OnMouseOverCached(s string, f dom.EventHandler) *Node {
	n.Handlers.MouseOverCacheKey = s
	n.Handlers.MouseOver = f
	return n
}

// See OnClickCached
func (n *Node) OnMouseOutCached(s string, f dom.EventHandler) *Node {
	n.Handlers.MouseOutCacheKey = s
	n.Handlers.MouseOut = f
	return n
}

func (n *Node) OnMouseDown(f dom.EventHandler) *Node {
	n.Handlers.MouseDown = f
	return n
}

func (n *Node) OnMouseDownCached(k string, f dom.EventHandler) *Node {
	n.Handlers.MouseDownCacheKey = k
	return n.OnMouseDown(f)
}

func (n *Node) OnMouseDownDispatch(e Event) *Node {
	return n.OnMouseDown(func(_ dom.Event) { go Dispatch(e) })
}

func (n *Node) OnMouseUp(f dom.EventHandler) *Node {
	n.Handlers.MouseUp = f
	return n
}

func (n *Node) OnMouseUpCached(k string, f dom.EventHandler) *Node {
	n.Handlers.MouseUpCacheKey = k
	return n.OnMouseUp(f)
}

func (n *Node) OnMouseUpDispatch(e Event) *Node {
	return n.OnMouseUp(func(_ dom.Event) { go Dispatch(e) })
}

func (n *Node) OnMouseMove(f dom.EventHandler) *Node {
	n.Handlers.MouseMove = f
	return n
}

func (n *Node) OnMouseMoveCached(k string, f dom.EventHandler) *Node {
	n.Handlers.MouseMoveCacheKey = k
	return n.OnMouseMove(f)
}

func (n *Node) OnMouseMoveDispatch(e Event) *Node {
	return n.OnMouseMove(func(_ dom.Event) { go Dispatch(e) })
}

func (n *Node) OnClickDispatch(e Event) *Node {
	return n.OnClick(func(_ dom.Event) { go Dispatch(e) })
}

func (n *Node) OnMouseOverDispatch(e Event) *Node {
	return n.OnMouseOver(func(_ dom.Event) { go Dispatch(e) })
}

func (n *Node) OnMouseEnter(f dom.EventHandler) *Node {
	n.Handlers.MouseEnter = f
	return n
}

func (n *Node) OnMouseLeave(f dom.EventHandler) *Node {
	n.Handlers.MouseLeave = f
	return n
}

func (n *Node) OnMouseOutDispatch(e Event) *Node {
	return n.OnMouseOut(func(_ dom.Event) { go Dispatch(e) })
}

func (n *Node) OnMouseLeaveDispatch(e Event) *Node {
	return n.OnMouseLeave(func(_ dom.Event) { go Dispatch(e) })
}

func (n *Node) OnMouseEnterDispatch(e Event) *Node {
	return n.OnMouseEnter(func(_ dom.Event) { go Dispatch(e) })
}

func (n *Node) OnClickDispatchCached(key string, e Event) *Node {
	return n.OnClickCached(key, func(_ dom.Event) { go Dispatch(e) })
}

func Dispatcher(e Event) func(dom.Event) {
	return func(_ dom.Event) {
		go Dispatch(e)
	}
}

func (n *Node) OnKeyUp(f dom.EventHandler) *Node {
	n.Handlers.KeyUp = f
	return n
}

func (n *Node) OnKeyDown(f dom.EventHandler) *Node {
	n.Handlers.KeyDown = f
	return n
}

func (n *Node) OnDrop(f dom.EventHandler) *Node {
	n.Handlers.Drop = f
	return n
}

func (n *Node) OnDragOver(f dom.EventHandler) *Node {
	n.Handlers.DragOver = f
	return n
}

/*
func (n *Node) OnEnter(f dom.EventHandler) *Node {
	n.Handlers.KeyUp = func(e dom.Event) {
		if e.KeyCode() == 13 { //enter?
			f(e)
		}
	}
	return n
}

func (n *Node) OnEnterDispatch(e Event) *Node {
	return n.OnEnter(func(_ dom.Event) {
		go Dispatch(e)
	})
}
*/

// }}}

// Canvas (e.g. WithDraw) {{{

func (n *Node) WithDraw(d func(ctx dom.CanvasRenderingContext2D)) {
	n.CanvasDraw = d
}

// }}}
