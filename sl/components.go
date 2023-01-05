package sl

import (
	"fmt"
	"time"

	"github.com/spinsrv/browser"
	"github.com/spinsrv/browser/dom"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// Alert {{{

type AlertAttr struct {
	Open     bool
	Closable bool
	Variant  string
	Duration string
}

func (a *AlertAttr) Attr() (out []*html.Attribute) {
	if a == nil {
		return
	}

	if a.Open {
		out = append(out, &html.Attribute{Key: "open", Val: "true"})
	}
	if a.Closable {
		out = append(out, &html.Attribute{Key: "closable", Val: "true"})
	}
	if a.Variant != "" {
		out = append(out, &html.Attribute{Key: "variant", Val: a.Variant})
	}
	if a.Duration != "" {
		out = append(out, &html.Attribute{Key: "duration", Val: a.Duration})
	}

	return
}

func Alert(a *AlertAttr, children ...*browser.Node) *browser.Node {
	return &browser.Node{
		Type:     html.ElementNode,
		Data:     "sl-alert",
		Children: children,
		Attr:     a.Attr(),
	}
}

// }}}

// Breadcrumb {{{

type BreadcrumbAttr struct {
	Label string
}

func (a *BreadcrumbAttr) Attr() (out []*html.Attribute) {
	if a == nil {
		return
	}

	if a.Label != "" {
		out = append(out, &html.Attribute{Key: "label", Val: a.Label})
	}
	return
}

func Breadcrumb(a *BreadcrumbAttr, children ...*browser.Node) *browser.Node {
	return &browser.Node{
		Type:     html.ElementNode,
		Data:     "sl-breadcrumb",
		Children: children,
		Attr:     a.Attr(),
	}
}

// }}}

// BreadcrumbItemAttr {{{

type BreadcrumbItemAttr struct {
}

func (a *BreadcrumbItemAttr) Attr() (out []*html.Attribute) {
	if a == nil {
		return
	}

	return
}

func BreadcrumbItem(a *BreadcrumbItemAttr, children ...*browser.Node) *browser.Node {
	return &browser.Node{
		Type:     html.ElementNode,
		Data:     "sl-breadcrumb-item",
		Children: children,
		Attr:     a.Attr(),
	}
}

// }}}

// BreadcrumpSeparator {{{

func BreadcrumbSeparator(children ...*browser.Node) *browser.Node {
	return &browser.Node{
		Type:     html.ElementNode,
		Data:     "span",
		Children: children,
		Attr: []*html.Attribute{
			&html.Attribute{
				Key: "slot",
				Val: "separator",
			},
		},
	}
}

// }}}

// Button {{{

type ButtonAttr struct {
	Variant  string
	Size     string
	Caret    bool
	Disabled bool
	Loading  bool
	Outline  bool
}

func (a *ButtonAttr) Attr() (out []*html.Attribute) {
	if a == nil {
		return
	}
	if a.Caret {
		out = append(out, &html.Attribute{Key: "caret", Val: "true"})
	}
	if a.Disabled {
		out = append(out, &html.Attribute{Key: "disabled", Val: "true"})
	}
	if a.Loading {
		out = append(out, &html.Attribute{Key: "loading", Val: "true"})
	}
	if a.Outline {
		out = append(out, &html.Attribute{Key: "outline", Val: "true"})
	}

	return
}

func Button(a *ButtonAttr, children ...*browser.Node) *browser.Node {
	return &browser.Node{
		Type:     html.ElementNode,
		Data:     "sl-button",
		Children: children,
		Attr:     a.Attr(),
	}
}

// }}}

// ButtonGroup {{{

type ButtonGroupAttr struct {
	Label string
}

func (a *ButtonGroupAttr) Attr() (out []*html.Attribute) {
	if a == nil {
		return
	}

	if a.Label != "" {
		out = append(out, &html.Attribute{Key: "label", Val: a.Label})
	}
	return
}

func ButtonGroup(a *ButtonGroupAttr, children ...*browser.Node) *browser.Node {
	return &browser.Node{
		Type:     html.ElementNode,
		Data:     "sl-button-group",
		Children: children,
		Attr:     a.Attr(),
	}
}

// }}}

// Card {{{

type CardAttr struct{}

func (a *CardAttr) Attr() (out []*html.Attribute) {
	if a == nil {
		return
	}

	return
}

func Card(a *CardAttr, children ...*browser.Node) *browser.Node {
	return &browser.Node{
		Type:     html.ElementNode,
		Data:     "sl-card",
		Children: children,
		Attr:     a.Attr(),
	}
}

// }}}

// Details {{{

type DetailsAttr struct {
	Open     bool
	Summary  string
	Disabled bool
}

func (a *DetailsAttr) Attr() (out []*html.Attribute) {
	if a == nil {
		return
	}
	if a.Open {
		out = append(out, &html.Attribute{Key: "open", Val: "true"})
	}
	if a.Summary != "" {
		out = append(out, &html.Attribute{Key: "summary", Val: a.Summary})
	}
	if a.Disabled {
		out = append(out, &html.Attribute{Key: "disabled", Val: "true"})
	}

	return
}

func Details(a *DetailsAttr, children ...*browser.Node) *browser.Node {
	return &browser.Node{
		Type:     html.ElementNode,
		Data:     "sl-details",
		Children: children,
		Attr:     a.Attr(),
	}
}

// }}}

// Divider {{{

type DividerAttr struct {
	Vertical bool
}

func (a *DividerAttr) Attr() (out []*html.Attribute) {
	if a == nil {
		return
	}

	if a.Vertical {
		out = append(out, &html.Attribute{Key: "vertical", Val: "true"})
	}
	return
}

func Divider(a *DividerAttr, children ...*browser.Node) *browser.Node {
	return &browser.Node{
		Type:     html.ElementNode,
		Data:     "sl-divider",
		Children: children,
		Attr:     a.Attr(),
	}
}

// }}}

// FormatDate {{{

type FormatDateAttr struct {
	Date         time.Time
	Weekday      string
	Era          string
	Year         string
	Month        string
	Day          string
	Hour         string
	Minute       string
	TimeZoneName string
	TimeZone     string
	HourFormat   string
}

func (a *FormatDateAttr) Attr() (out []*html.Attribute) {
	if a == nil {
		return
	}

	out = append(out, &html.Attribute{Key: "date", Val: a.Date.Format(time.RFC3339)})

	if a.Weekday != "" {
		out = append(out, &html.Attribute{Key: "weekday", Val: a.Weekday})
	}
	if a.Era != "" {
		out = append(out, &html.Attribute{Key: "era", Val: a.Era})
	}
	if a.Year != "" {
		out = append(out, &html.Attribute{Key: "year", Val: a.Year})
	}
	if a.Month != "" {
		out = append(out, &html.Attribute{Key: "month", Val: a.Month})
	}
	if a.Day != "" {
		out = append(out, &html.Attribute{Key: "day", Val: a.Day})
	}
	if a.Hour != "" {
		out = append(out, &html.Attribute{Key: "hour", Val: a.Hour})
	}
	if a.Minute != "" {
		out = append(out, &html.Attribute{Key: "minute", Val: a.Minute})
	}
	if a.TimeZoneName != "" {
		out = append(out, &html.Attribute{Key: "time-zone-name", Val: a.TimeZoneName})
	}
	if a.TimeZone != "" {
		out = append(out, &html.Attribute{Key: "time-zone", Val: a.TimeZone})
	}
	if a.HourFormat != "" {
		out = append(out, &html.Attribute{Key: "hour-format", Val: a.HourFormat})
	}
	return
}

func FormatDate(a *FormatDateAttr, children ...*browser.Node) *browser.Node {
	return &browser.Node{
		Type:     html.ElementNode,
		Data:     "sl-format-date",
		Children: children,
		Attr:     a.Attr(),
	}
}

// }}}

// FormatNumber {{{

type FormatNumberAttr struct {
	Value           float64
	Type            string
	NoGrouping      bool
	Currency        string
	CurrencyDisplay string
	// pointers here give "optional" semantics
	MinimumIntegerDigits     *uint8
	MinimumFractionDigits    *uint8
	MaximumFractionDigits    *uint8
	MinimumSignificantDigits *uint8
	MaximumSignificantDigits *uint8
}

func (a *FormatNumberAttr) Attr() (out []*html.Attribute) {
	if a == nil {
		return
	}

	out = append(out, &html.Attribute{Key: "value", Val: fmt.Sprintf("%f", a.Value)})

	if a.Type != "" {
		out = append(out, &html.Attribute{Key: "type", Val: a.Type})
	}
	if a.NoGrouping {
		out = append(out, &html.Attribute{Key: "no-grouping", Val: "true"})
	}
	if a.Currency != "" {
		out = append(out, &html.Attribute{Key: "currency", Val: a.Currency})
	}
	if a.CurrencyDisplay != "" {
		out = append(out, &html.Attribute{Key: "currency-display", Val: a.CurrencyDisplay})
	}
	if a.MinimumIntegerDigits != nil {
		out = append(out, &html.Attribute{Key: "minimum-integer-digits", Val: fmt.Sprintf("%d", *a.MinimumIntegerDigits)})
	}
	if a.MinimumFractionDigits != nil {
		out = append(out, &html.Attribute{Key: "minimum-fraction-digits", Val: fmt.Sprintf("%d", *a.MinimumFractionDigits)})
	}
	if a.MaximumFractionDigits != nil {
		out = append(out, &html.Attribute{Key: "maximum-fraction-digits", Val: fmt.Sprintf("%d", *a.MaximumFractionDigits)})
	}
	if a.MinimumSignificantDigits != nil {
		out = append(out, &html.Attribute{Key: "minimum-significant-digits", Val: fmt.Sprintf("%d", *a.MinimumSignificantDigits)})
	}
	if a.MaximumSignificantDigits != nil {
		out = append(out, &html.Attribute{Key: "maximum-significant-digits", Val: fmt.Sprintf("%d", *a.MaximumSignificantDigits)})
	}
	return
}

func FormatNumber(a *FormatNumberAttr, children ...*browser.Node) *browser.Node {
	return &browser.Node{
		Type:     html.ElementNode,
		Data:     "sl-format-number",
		Children: children,
		Attr:     a.Attr(),
	}
}

// }}}

// Icon {{{

type IconAttr struct {
	Slot    string
	Name    string
	Src     string
	Label   string
	Library string
}

func (a *IconAttr) Attr() (out []*html.Attribute) {
	if a == nil {
		return
	}

	if a.Slot != "" {
		out = append(out, &html.Attribute{Key: "slot", Val: a.Slot})
	}
	if a.Name != "" {
		out = append(out, &html.Attribute{Key: "name", Val: a.Name})
	}
	if a.Src != "" {
		out = append(out, &html.Attribute{Key: "src", Val: a.Src})
	}
	if a.Label != "" {
		out = append(out, &html.Attribute{Key: "label", Val: a.Label})
	}
	if a.Library != "" {
		out = append(out, &html.Attribute{Key: "library", Val: a.Library})
	}

	return
}

func Icon(a *IconAttr, children ...*browser.Node) *browser.Node {
	return &browser.Node{
		Type:     html.ElementNode,
		Data:     "sl-icon",
		Children: children,
		Attr:     a.Attr(),
	}
}

// }}}

// IconButton {{{

type IconButtonAttr struct {
	Name     string
	Library  string
	Src      string
	Href     string
	Target   string
	Download string
	Label    string
	Disabled bool
}

func (a *IconButtonAttr) Attr() (out []*html.Attribute) {
	if a == nil {
		return
	}

	if a.Name != "" {
		out = append(out, &html.Attribute{Key: "name", Val: a.Name})
	}
	if a.Library != "" {
		out = append(out, &html.Attribute{Key: "library", Val: a.Library})
	}
	if a.Src != "" {
		out = append(out, &html.Attribute{Key: "src", Val: a.Src})
	}
	if a.Href != "" {
		out = append(out, &html.Attribute{Key: "href", Val: a.Href})
	}
	if a.Target != "" {
		out = append(out, &html.Attribute{Key: "target", Val: a.Target})
	}
	if a.Download != "" {
		out = append(out, &html.Attribute{Key: "download", Val: a.Download})
	}
	if a.Label != "" {
		out = append(out, &html.Attribute{Key: "label", Val: a.Label})
	}
	if a.Disabled {
		out = append(out, &html.Attribute{Key: "disabled", Val: "true"})
	}

	return
}

func IconButton(a *IconButtonAttr, children ...*browser.Node) *browser.Node {
	return &browser.Node{
		Type:     html.ElementNode,
		Data:     "sl-icon-button",
		Children: children,
		Attr:     a.Attr(),
	}
}

// }}}

// Input {{{

type InputAttr struct {
	Type            string
	Size            string
	Name            string
	Value           *string
	Filled          bool
	Pill            bool
	Label           string
	HelpText        string
	Clearable       bool
	PasswordToggle  bool
	PasswordVisible bool
	Required        bool
	Placeholder     string
}

func (a *InputAttr) Attr() (out []*html.Attribute) {
	if a == nil {
		return
	}

	if a.Type != "" {
		out = append(out, &html.Attribute{Key: "type", Val: a.Type})
	}
	if a.Size != "" {
		out = append(out, &html.Attribute{Key: "size", Val: a.Size})
	}
	if a.Name != "" {
		out = append(out, &html.Attribute{Key: "name", Val: a.Name})
	}
	if *a.Value != "" {
		out = append(out, &html.Attribute{Key: "value", Val: *a.Value})
	}
	if a.Filled {
		out = append(out, &html.Attribute{Key: "filled", Val: "true"})
	}
	if a.Pill {
		out = append(out, &html.Attribute{Key: "pill", Val: "true"})
	}
	if a.Label != "" {
		out = append(out, &html.Attribute{Key: "label", Val: a.Label})
	}
	if a.HelpText != "" {
		out = append(out, &html.Attribute{Key: "help-text", Val: a.HelpText})
	}
	if a.Clearable {
		out = append(out, &html.Attribute{Key: "clearable", Val: "true"})
	}
	if a.PasswordToggle {
		out = append(out, &html.Attribute{Key: "password-toggle", Val: "true"})
	}
	if a.Clearable {
		out = append(out, &html.Attribute{Key: "password-visible", Val: "true"})
	}
	if a.Required {
		out = append(out, &html.Attribute{Key: "required", Val: "true"})
	}
	if a.Placeholder != "" {
		out = append(out, &html.Attribute{Key: "placeholder", Val: a.Placeholder})
	}

	return
}

func Input(a *InputAttr, children ...*browser.Node) *browser.Node {
	return (&browser.Node{
		Type:     html.ElementNode,
		Data:     "sl-input",
		Children: children,
		Attr:     a.Attr(),
	}).OnInput(func(e dom.Event) {
		*a.Value = e.Target().Value()
		go browser.Dispatch(nil)
	})
}

// }}}

// Menu {{{

type MenuAttr struct{}

func (a *MenuAttr) Attr() (out []*html.Attribute) {
	if a == nil {
		return
	}

	return
}

func Menu(a *MenuAttr, children ...*browser.Node) *browser.Node {
	return &browser.Node{
		Type:     html.ElementNode,
		Data:     "sl-menu",
		Children: children,
		Attr:     a.Attr(),
	}
}

// }}}

// MenuItem {{{

type MenuItemAttr struct{}

func (a *MenuItemAttr) Attr() (out []*html.Attribute) {
	if a == nil {
		return
	}

	return
}

func MenuItem(a *MenuItemAttr, children ...*browser.Node) *browser.Node {
	return &browser.Node{
		Type:     html.ElementNode,
		Data:     "sl-menu-item",
		Children: children,
		Attr:     a.Attr(),
	}
}

// }}}

// MenuLabelAttr {{{

type MenuLabelAttr struct{}

func (a *MenuLabelAttr) Attr() (out []*html.Attribute) {
	if a == nil {
		return
	}

	return
}

func MenuLabel(a *MenuLabelAttr, children ...*browser.Node) *browser.Node {
	return &browser.Node{
		Type:     html.ElementNode,
		Data:     "sl-menu-label",
		Children: children,
		Attr:     a.Attr(),
	}
}

// }}}

// SkeletonAttr {{{

type SkeletonAttr struct {
	Effect string
}

func (a *SkeletonAttr) Attr() (out []*html.Attribute) {
	if a == nil {
		return nil
	}

	if a.Effect != "" {
		out = append(out, &html.Attribute{Key: "effect", Val: a.Effect})
	}

	return
}

func Skeleton(a *SkeletonAttr, children ...*browser.Node) *browser.Node {
	return &browser.Node{
		Type:     html.ElementNode,
		Data:     "sl-skeleton",
		Children: children,
		Attr:     a.Attr(),
	}
}

// }}}

// Tab {{{

type TabAttr struct {
	Slot     string
	Panel    string
	Active   bool
	Closable bool
	Disabled bool
}

func (a *TabAttr) Attr() (out []*html.Attribute) {
	if a == nil {
		return
	}

	if a.Slot != "" {
		out = append(out, &html.Attribute{Key: "slot", Val: a.Slot})
	}
	if a.Panel != "" {
		out = append(out, &html.Attribute{Key: "panel", Val: a.Panel})
	}
	if a.Active {
		out = append(out, &html.Attribute{Key: "active", Val: "true"})
	}
	if a.Closable {
		out = append(out, &html.Attribute{Key: "closable", Val: "true"})
	}
	if a.Disabled {
		out = append(out, &html.Attribute{Key: "disabled", Val: "true"})
	}
	return
}

func Tab(a *TabAttr, children ...*browser.Node) *browser.Node {
	return &browser.Node{
		Type:     html.ElementNode,
		Data:     "sl-tab",
		Children: children,
		Attr:     a.Attr(),
	}
}

// }}}

// TabGroup {{{

type TabGroupAttr struct {
	Placement        string
	Activation       string
	NoScrollControls bool
}

func (a *TabGroupAttr) Attr() (out []*html.Attribute) {
	if a == nil {
		return
	}

	if a.Placement != "" {
		out = append(out, &html.Attribute{Key: "placement", Val: a.Placement})
	}
	if a.Activation != "" {
		out = append(out, &html.Attribute{Key: "activation", Val: a.Activation})
	}
	if a.NoScrollControls {
		out = append(out, &html.Attribute{Key: "no-scroll-controls", Val: "true"})
	}

	return
}

func TabGroup(a *TabGroupAttr, children ...*browser.Node) *browser.Node {
	return &browser.Node{
		Type:     html.ElementNode,
		Data:     "sl-tab-group",
		Children: children,
		Attr:     a.Attr(),
	}
}

// }}}

// TabPanel {{{

type TabPanelAttr struct {
	Name string
}

func (a *TabPanelAttr) Attr() (out []*html.Attribute) {
	if a.Name != "" {
		out = append(out, &html.Attribute{Key: "name", Val: a.Name})
	}
	return
}

func TabPanel(a *TabPanelAttr, children ...*browser.Node) *browser.Node {
	return &browser.Node{
		Type:     html.ElementNode,
		Data:     "sl-tab-panel",
		Children: children,
		Attr:     a.Attr(),
	}
}

// }}}

func Tabs(names []string, panels []*browser.Node) *browser.Node {
	if len(names) != len(panels) {
		panic(fmt.Sprintf("sl.Tabs: must have same number of names as panels: got %d and %d", len(names), len(panels)))
	}

	var children []*browser.Node = make([]*browser.Node, 2*len(names))

	for _, name := range names {
		children = append(children, Tab(&TabAttr{Slot: "nav", Panel: name}))
	}

	for i, panel := range panels {
		children = append(children, TabPanel(&TabPanelAttr{Name: names[i]}, panel))
	}

	return TabGroup(&TabGroupAttr{Placement: "start"}, children...)
}

// Tag {{{

type TagAttr struct {
	Variant   string
	Size      string
	Pill      bool
	Removable bool
}

func (a *TagAttr) Attr() (out []*html.Attribute) {
	if a == nil {
		return
	}

	if a.Variant != "" {
		out = append(out, &html.Attribute{Key: "variant", Val: a.Variant})
	}
	if a.Size != "" {
		out = append(out, &html.Attribute{Key: "size", Val: a.Size})
	}
	if a.Pill {
		out = append(out, &html.Attribute{Key: "pill", Val: "true"})
	}
	if a.Removable {
		out = append(out, &html.Attribute{Key: "removable", Val: "true"})
	}
	return
}

func Tag(a *TagAttr, children ...*browser.Node) *browser.Node {
	return &browser.Node{
		Type:     html.ElementNode,
		Data:     "sl-tag",
		Children: children,
		Attr:     a.Attr(),
	}
}

// }}}

// A TextArea is a <sl-textarea> component.
// See: https://shoelace.style/components/textarea
func TextArea(value *string) *browser.Node {
	return textArea(value)
}

func textArea(value *string) *browser.Node {
	return (&browser.Node{
		Type: html.ElementNode,
		Data: "sl-textarea",
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

func textNode(s string) *browser.Node {
	return &browser.Node{
		Type: html.TextNode,
		Data: s,
	}
}

// Tooltip {{{

type TooltipAttr struct {
	Content   string
	Placement string
	Disabled  bool
	Distance  *int // handles optional
	Open      bool
	Skidding  *int // handles optional
	Trigger   string
	Hoist     bool
}

func (a *TooltipAttr) Attr() (out []*html.Attribute) {
	if a == nil {
		return
	}

	if a.Content != "" {
		out = append(out, &html.Attribute{Key: "content", Val: a.Content})
	}
	if a.Placement != "" {
		out = append(out, &html.Attribute{Key: "placement", Val: a.Placement})
	}
	if a.Disabled {
		out = append(out, &html.Attribute{Key: "disabled", Val: "true"})
	}
	if a.Distance != nil {
		out = append(out, &html.Attribute{Key: "distance", Val: fmt.Sprintf("%d", a.Distance)})
	}
	if a.Open {
		out = append(out, &html.Attribute{Key: "open", Val: "true"})
	}
	if a.Skidding != nil {
		out = append(out, &html.Attribute{Key: "skidding", Val: fmt.Sprintf("%d", a.Skidding)})
	}
	if a.Trigger != "" {
		out = append(out, &html.Attribute{Key: "trigger", Val: a.Trigger})
	}
	if a.Hoist {
		out = append(out, &html.Attribute{Key: "hoist", Val: "true"})
	}
	return
}

func Tooltip(a *TooltipAttr, children ...*browser.Node) *browser.Node {
	return &browser.Node{
		Type:     html.ElementNode,
		Data:     "sl-tooltip",
		Children: children,
		Attr:     a.Attr(),
	}
}

// }}}

// Tree {{{

type TreeAttr struct {
	Selection string
}

func (a *TreeAttr) Attr() (out []*html.Attribute) {
	if a == nil {
		return
	}

	if a.Selection != "" {
		out = append(out, &html.Attribute{Key: "selection", Val: a.Selection})
	}
	return
}

// Tree constructs a <sl-tree> component.
func Tree(a *TreeAttr, children ...*browser.Node) *browser.Node {
	return &browser.Node{
		Type:     html.ElementNode,
		Data:     "sl-tree",
		Children: children,
		Attr:     a.Attr(),
	}
}

// }}}

// TreeItemAttr {{{

type TreeItemAttr struct {
	Expanded bool
	Selected bool
	Disabled bool
	Lazy     bool
}

func (a *TreeItemAttr) Attr() (out []*html.Attribute) {
	if a == nil {
		return
	}

	if a.Expanded {
		out = append(out, &html.Attribute{Key: "expanded", Val: "true"})
	}
	if a.Selected {
		out = append(out, &html.Attribute{Key: "selected", Val: "true"})
	}
	if a.Disabled {
		out = append(out, &html.Attribute{Key: "disabled", Val: "true"})
	}
	if a.Lazy {
		out = append(out, &html.Attribute{Key: "lazy", Val: "true"})
	}
	return
}

// TreeItem constructs a <sl-tree-item> component.
func TreeItem(a *TreeItemAttr, children ...*browser.Node) *browser.Node {
	return &browser.Node{
		Type:     html.ElementNode,
		Data:     "sl-tree-item",
		Children: children,
		Attr:     a.Attr(),
	}
}

// }}}
