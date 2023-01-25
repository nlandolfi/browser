package browser

import (
	"fmt"
	"log"

	"github.com/spinsrv/browser/dom"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// A Mounter attaches to a node in the DOM. Subsequent calls to
// Mount will render the given browser.Node into the DOM element Root.
type Mounter struct {
	// Root is the DOM element into which this Mounter renders.
	//
	// This is often the document body, but need not be.
	Root dom.Element

	// Document is the document object. We only use CreateElement and CreateTextNode.
	//
	// This is often js.DefaultBrowser. See the js package herein.
	Document dom.Document

	// last is the Node last mounted, it is used to diff a new mount
	// with the old, and decide on DOM changes. see `mount` below
	last *Node
}

// Use Mount to mount the Node to the DOM element.
func (m *Mounter) Mount(n *Node) error {
	return m.mount(n)
}

// mount is a hidden helper; currently the indirection is not used, as the public
// facing Mount simply calls this function; but we leave open the possibility for
// an API change.
func (m *Mounter) mount(n *Node) error {
	if m.Root == nil || m.Document == nil {
		return fmt.Errorf("browser.Mount: Mounter requires non-nil Root and Document")
	}

	changes := reconcileWalker(m.Root, m.last, n)
	// commented out 1/30/22
	//s := make([]string, len(changes))
	//for i := range changes {
	//	s[i] = changes[i].Type.String()
	//}
	m.last = n
	for _, c := range changes {
		m.apply(c)
	}

	return nil
}

func (m *Mounter) apply(c *change) {
	switch c.Type {
	case insert:
		log.Print("insert")
		m.create(c.Ref)           // creates the DOM element
		m.insert(c.Parent, c.Ref) // inserts it
	case replace:
		log.Print("replace")
		m.create(c.Ref) // creates the DOM element
		// I don't understand how this is better than just mutating the current node? - NCL 1/30/22
		m.replace(c.Parent, c.Old, c.Ref)
	case remove:
		m.remove(c.Parent, c.Ref)
	case attrSet:
		m.attrSet(c.Ref, c.Key, c.Val)
	case attrDelete:
		m.attrDelete(c.Ref, c.Key)
	case listenerAdd:
		m.listenerAdd(c.Ref, c.EventType, c.NewListener)
	case listenerDelete:
		m.listenerDelete(c.Ref, c.EventType, c.OldListener)
	case canvasDraw:
		m.canvasDraw(c.Ref, c.CanvasDraw)
	default:
		panic(fmt.Sprintf("unknown change type: %s", c.Type))
	}
}

// create calls into JS to make the DOM node
func (m *Mounter) create(ref *Node) {
	if ref == nil {
		panic("creating a nil ref!")
	}

	if ref.rendered != nil {
		panic("calling create on already rendered *Node")
	}

	switch ref.Type {
	case html.ElementNode:
		var tagName = ref.Data
		if tagName == "" { // allow only defining the atom
			tagName = ref.DataAtom.String()
		}
		if tagName == "" {
			panic("trying to mount ElementNode with empty tag name: must define DataAtom or Data (or both)")
		}
		ref.renderedElement = m.Document.CreateElement(tagName)
		ref.rendered = ref.renderedElement
	case html.TextNode:
		ref.rendered = m.Document.CreateTextNode(ref.Data)
	default:
		panic(fmt.Sprintf("unknown Node.Type: %#v", ref.Type))
	}
}

// insert calls into JS to add a node as a child to `into`
func (m *Mounter) insert(into, ref *Node) {
	if into.rendered == nil {
		panic("insterting into an unrendered node")
	}
	if into.Type != html.ElementNode {
		panic("inserting into something that isn't an element node!")
	}

	into.rendered.AppendChild(ref.rendered)
}

// replace calls into JS to swap a child
func (m *Mounter) replace(parent, old, new *Node) {
	if old.rendered == nil {
		panic("mount replace: old has nil rendered")
	}
	if new.rendered == nil {
		panic("mount replace: new has nil rendered")
	}
	if parent == nil {
		panic("mount replace: parents is nil")
	}
	if parent.rendered == nil {
		panic("mount replace: parent has nil rendered")
	}

	// yes, the old should be second - NCL 1/30/22
	parent.rendered.ReplaceChild(new.rendered, old.rendered)
}

// remove calls into JS to drop a node
func (m *Mounter) remove(parent, child *Node) {
	if parent.rendered == nil {
		panic("parent rendered nil")
	}
	if child.rendered == nil {
		panic("child rendered nil")
	}

	parent.rendered.RemoveChild(child.rendered)
}

func (m *Mounter) attrSet(ref *Node, key, val string) {
	if ref.rendered == nil {
		panic("attrSet on a node with nil rendered")
	}
	if ref.Type != html.ElementNode {
		//log.Printf("%+v", ref)
		//log.Printf("key: %q, val %q", key, val)
		panic("can't set attribute of a node that is not element node")
	}
	if ref.renderedElement == nil {
		panic("attrSet on a node with a nil renderedElement")
	}
	// unfortunately there are special cases
	switch key {
	case "value":
		// I read online you need to do this, and empirically verified in safari
		ref.renderedElement.SetValue(val)
		/*
			case "selectionStart":
				i, err := strconv.ParseInt(val, 10, 64)
				if err != nil {
					panic(err)
				}
				ref.renderedElement.SetSelectionStart(int(i))
			case "selectionEnd":
				i, err := strconv.ParseInt(val, 10, 64)
				if err != nil {
					panic(err)
				}
				ref.renderedElement.SetSelectionEnd(int(i))
		*/
	}
	ref.renderedElement.SetAttribute(key, val)
}

func (m *Mounter) attrDelete(ref *Node, key string) {
	if ref.rendered == nil {
		panic("attrDelete on a node with nil rendered")
	}
	if ref.Type != html.ElementNode {
		panic("can't delete attribute of a node that is not element node")
	}
	if ref.renderedElement == nil {
		panic("attrSet on a node with a nil renderedElement")
	}
	ref.renderedElement.RemoveAttribute(key)
}

func (m *Mounter) listenerAdd(r *Node, t dom.EventType, l dom.EventHandler) {
	if r.rendered == nil {
		panic("listenerAdd on a node with nil rendered")
	}

	el := r.rendered.AddEventListener(t, l)

	switch t {
	case dom.Click:
		r.Handlers.click = el
	case dom.DoubleClick:
		r.Handlers.doubleClick = el
	case dom.MouseOver:
		r.Handlers.mouseOver = el
	case dom.MouseOut:
		r.Handlers.mouseOut = el
	case dom.Input:
		r.Handlers.input = el
	case dom.MouseDown:
		r.Handlers.mouseDown = el
	case dom.MouseMove:
		r.Handlers.mouseMove = el
	case dom.MouseUp:
		r.Handlers.mouseUp = el
	case dom.MouseEnter:
		r.Handlers.mouseEnter = el
	case dom.MouseLeave:
		r.Handlers.mouseLeave = el
	case dom.KeyUp:
		r.Handlers.keyUp = el
	case dom.KeyDown:
		r.Handlers.keyDown = el
	case dom.Drop:
		r.Handlers.drop = el
	case dom.DragOver:
		r.Handlers.dragOver = el
	default:
		panic("unknown listener")
	}
}

func (m *Mounter) canvasDraw(r *Node, draw func(c dom.CanvasRenderingContext2D)) {
	if r.renderedElement == nil {
		panic("canvasDraw on a node with nil renderedElemenet")
	}

	// need to figure out proper height and width...
	c := r.renderedElement.CanvasContext(200, 100, 10.0)

	// should be new context?
	draw(c)
}

func (m *Mounter) listenerDelete(r *Node, t dom.EventType, l dom.EventListener) {
	if r.rendered == nil {
		panic("listenerDelete on a node with nil rendered")
	}
	if l == nil {
		panic("listenerDelete event listener shouldn't be nil")
	}

	r.rendered.RemoveEventListener(t, l)
}

type changeType int

const (
	none changeType = iota
	insert
	replace
	remove
	attrSet
	attrDelete
	styleSet
	styleDelete
	listenerAdd
	listenerDelete
	canvasDraw
)

func (t changeType) String() string {
	switch t {
	case none:
		return "NONE"
	case insert:
		return "INSERT"
	case replace:
		return "REPLACE"
	case remove:
		return "REMOVE"
	case attrSet:
		return "ATTR_SET"
	case attrDelete:
		return "ATTR_DELETE"
	case styleSet:
		return "STYLE_SET"
	case styleDelete:
		return "STYLE_DELETE"
	case listenerAdd:
		return "LISTENER_ADD"
	case listenerDelete:
		return "LISTENER_DELETE"
	case canvasDraw:
		return "CANVAS_DRAW"
	default:
		panic("unknown type")
	}
}

type change struct {
	Type changeType

	// both set for replace, New set for insert
	Parent, Old, Ref *Node

	Key string // set for attrSet, attrDelete, styleSet, styleDelete
	Val string // set for attrDelete, styleDelete

	dom.EventType // set for listenerAdd, listenerDelete

	OldListener dom.EventListener // set for listenerDelete
	NewListener dom.EventHandler  // set for listenerAdd

	CanvasDraw func(dom.CanvasRenderingContext2D)
}

type nodePair struct {
	parent, old, new *Node
	level            int
}

func reconcileWalker(base dom.Element, oldRoot, newRoot *Node) (changes []*change) {
	if base == nil {
		panic("reconcileWalker base element can not be nil")
	}

	fakedParent := &Node{Type: html.ElementNode, rendered: base}

	if oldRoot == nil {
		changes = inserts(fakedParent, newRoot, 0)
		return
	}

	var stack []*nodePair
	stack = append(stack, &nodePair{
		parent: fakedParent,
		old:    oldRoot,
		new:    newRoot,
		level:  0,
	})

	for len(stack) > 0 {
		// pop element from the stack
		top := stack[len(stack)-1]
		stack = stack[0 : len(stack)-1]
		parent, old, new := top.parent, top.old, top.new

		localChanges, replacedTree := reconcile(parent, old, new, top.level)
		changes = append(changes, localChanges...)
		//printl(top.level, "level %d (%s -> %s)", top.level, old.DataAtom.String(), new.DataAtom.String())

		// TODO REMOVE DEAD NODES FROM TREE!

		if replacedTree {
			//			log.Print("replaced tree")
			continue
		}

		// TODO: can this be optimized with 'insert before' / 'insert after'
		// may then handle common diff cases like adding an alerty box above some other UI

		for i := 0; i < len(old.Children) && i < len(new.Children); i++ {
			stack = append(stack, &nodePair{
				parent: new, old: old.Children[i], new: new.Children[i],
				level: top.level + 1,
			})
		}

		for i := len(old.Children); i < len(new.Children); i++ {
			changes = append(changes, inserts(new, new.Children[i], top.level+1)...) // with new.Children[i]
		}

		for i := len(new.Children); i < len(old.Children); i++ {
			changes = append(changes, &change{
				Type:   remove,
				Parent: old,
				Ref:    old.Children[i],
			})
		}
	}

	return
}

func reconcileAttr(ref *Node, old, new []*html.Attribute) (changes []*change) {
	oldM := make(map[string]string, len(old))
	newM := make(map[string]string, len(new))

	for _, a := range old {
		oldM[a.Key] = a.Val
	}

	for _, a := range new {
		newM[a.Key] = a.Val
	}

	for k, v := range newM {
		oldV, ok := oldM[k]
		if !ok || oldV != v { // lacks attribute or has diff value
			changes = append(changes, &change{
				Type: attrSet,
				Ref:  ref,
				Key:  k,
				Val:  v,
			})
		}

		// have handled skipping/resetting
		delete(oldM, k)
	}

	// any letfovers in the old map should be removed
	for k := range oldM {
		changes = append(changes, &change{
			Type: attrDelete,
			Ref:  ref,
			Key:  k,
		})
	}

	return
}

// the history here is complicated, but a summary
// - I think you can set some of the .style properties of an element and it will work, but
// potentiall some of the values do not. Also the code was super long to check each of
func diffStyles(ref *Node, from, to *Style, level int) (changes []*change) {
	if from == nil && to == nil {
		return
	}

	if to == nil && from != nil {
		changes = append(changes, &change{
			Type: attrSet,
			Ref:  ref,
			Key:  "style",
			Val:  "",
		})
	}

	if from == nil && to != nil || *from != *to {
		changes = append(changes, &change{
			Type: attrSet,
			Ref:  ref,
			Key:  "style",
			Val:  to.Val(),
		})
	}

	return
}

// reconcileHandlers {{{

func getClick(n *Node) dom.EventHandler {
	if n == nil {
		return nil
	}

	return n.Handlers.Click
}

func getDoubleClick(n *Node) dom.EventHandler {
	if n == nil {
		return nil
	}

	return n.Handlers.DoubleClick
}

func getMouseOver(n *Node) dom.EventHandler {
	if n == nil {
		return nil
	}

	return n.Handlers.MouseOver
}

func getMouseOut(n *Node) dom.EventHandler {
	if n == nil {
		return nil
	}

	return n.Handlers.MouseOut
}

func getMouseDown(n *Node) dom.EventHandler {
	if n == nil {
		return nil
	}

	return n.Handlers.MouseDown
}

func getMouseMove(n *Node) dom.EventHandler {
	if n == nil {
		return nil
	}

	return n.Handlers.MouseMove
}

func getMouseUp(n *Node) dom.EventHandler {
	if n == nil {
		return nil
	}

	return n.Handlers.MouseUp
}

func getMouseEnter(n *Node) dom.EventHandler {
	if n == nil {
		return nil
	}

	return n.Handlers.MouseEnter
}

func getMouseLeave(n *Node) dom.EventHandler {
	if n == nil {
		return nil
	}

	return n.Handlers.MouseLeave
}

func getInput(n *Node) dom.EventHandler {
	if n == nil {
		return nil
	}

	return n.Handlers.Input
}

func getKeyUp(n *Node) dom.EventHandler {
	if n == nil {
		return nil
	}

	return n.Handlers.KeyUp
}

func getKeyDown(n *Node) dom.EventHandler {
	if n == nil {
		return nil
	}

	return n.Handlers.KeyDown
}

func getDrop(n *Node) dom.EventHandler {
	if n == nil {
		return nil
	}

	return n.Handlers.Drop
}

func getDragOver(n *Node) dom.EventHandler {
	if n == nil {
		return nil
	}

	return n.Handlers.DragOver
}

const cacheListeners bool = true

func reconcileHandlers(old, new *Node) (changes []*change) {
	oldClick, newClick := getClick(old), getClick(new)
	if oldClick == nil && newClick == nil {
		// no op
	}
	if oldClick != nil && newClick == nil {
		if old.Handlers.click == nil {
			panic("listener has Click but no underlying event listener")
		}

		changes = append(changes, &change{
			Type:        listenerDelete,
			Ref:         old,
			EventType:   dom.Click,
			OldListener: old.Handlers.click,
		})
	}
	if oldClick == nil && newClick != nil {
		changes = append(changes, &change{
			Type:        listenerAdd,
			Ref:         new,
			EventType:   dom.Click,
			NewListener: new.Handlers.Click,
		})
	}
	if oldClick != nil && newClick != nil && &oldClick != &newClick {
		//	&& ((new.Handlers.ClickCacheKey != "" || old.Handlers.ClickCacheKey != "") && new.Handlers.ClickCacheKey != old.Handlers.ClickCacheKey) {

		// last attempt, check the cache
		if cacheListeners && new.Handlers.ClickCacheKey != "" && old.Handlers.ClickCacheKey != "" && new.Handlers.ClickCacheKey == old.Handlers.ClickCacheKey {
			// no need to syscall listener adding!
			new.Handlers.Click = old.Handlers.Click // not sure if we need to move this
			new.Handlers.click = old.Handlers.click // definitely need to move this, was a bug
		} else {
			if old.Handlers.click == nil {
				panic("listener has Click but no underlying event listener")
			}

			changes = append(changes, &change{
				Type:        listenerDelete,
				Ref:         old,
				EventType:   dom.Click,
				OldListener: old.Handlers.click,
			})
			changes = append(changes, &change{
				Type:        listenerAdd,
				Ref:         new,
				EventType:   dom.Click,
				NewListener: new.Handlers.Click,
			})
		}
	}

	oldDoubleClick, newDoubleClick := getDoubleClick(old), getDoubleClick(new)
	if oldDoubleClick == nil && newDoubleClick == nil {
		// no op
	}
	if oldDoubleClick != nil && newDoubleClick == nil {
		if old.Handlers.doubleClick == nil {
			panic("listener has DoubleClick but no underlying event listener")
		}

		changes = append(changes, &change{
			Type:        listenerDelete,
			Ref:         old,
			EventType:   dom.DoubleClick,
			OldListener: old.Handlers.doubleClick,
		})
	}
	if oldDoubleClick == nil && newDoubleClick != nil {
		changes = append(changes, &change{
			Type:        listenerAdd,
			Ref:         new,
			EventType:   dom.DoubleClick,
			NewListener: new.Handlers.DoubleClick,
		})
	}
	if oldDoubleClick != nil && newDoubleClick != nil && &oldDoubleClick != &newDoubleClick {
		// last attempt, check the cache
		if cacheListeners && new.Handlers.DoubleClickCacheKey != "" && old.Handlers.DoubleClickCacheKey != "" && new.Handlers.DoubleClickCacheKey == old.Handlers.DoubleClickCacheKey {
			// no need to syscall listener adding!
			new.Handlers.DoubleClick = old.Handlers.DoubleClick // not sure if we need to move this
			new.Handlers.doubleClick = old.Handlers.doubleClick // definitely need to move this, was a bug
		} else {

			if old.Handlers.doubleClick == nil {
				panic("listener has Click but no underlying event listener")
			}
			changes = append(changes, &change{
				Type:        listenerDelete,
				Ref:         old,
				EventType:   dom.DoubleClick,
				OldListener: old.Handlers.doubleClick,
			})
			changes = append(changes, &change{
				Type:        listenerAdd,
				Ref:         new,
				EventType:   dom.DoubleClick,
				NewListener: new.Handlers.DoubleClick,
			})
		}
	}

	oldMouseOver, newMouseOver := getMouseOver(old), getMouseOver(new)
	if oldMouseOver == nil && newMouseOver == nil {
		// no op
	}
	if oldMouseOver != nil && newMouseOver == nil {
		if old.Handlers.mouseOver == nil {
			panic("listener has MouseOver but no underlying event listener")
		}

		changes = append(changes, &change{
			Type:        listenerDelete,
			Ref:         old,
			EventType:   dom.MouseOver,
			OldListener: old.Handlers.mouseOver,
		})
	}
	if oldMouseOver == nil && newMouseOver != nil {
		changes = append(changes, &change{
			Type:        listenerAdd,
			Ref:         new,
			EventType:   dom.MouseOver,
			NewListener: new.Handlers.MouseOver,
		})
	}
	if oldMouseOver != nil && newMouseOver != nil && &oldMouseOver != &newMouseOver {
		// last attempt, check the cache
		if cacheListeners && new.Handlers.MouseOverCacheKey != "" && old.Handlers.MouseOverCacheKey != "" && new.Handlers.MouseOverCacheKey == old.Handlers.MouseOverCacheKey {
			// no need to syscall listener adding!
			new.Handlers.MouseOver = old.Handlers.MouseOver // not sure if we need to move this
			new.Handlers.mouseOver = old.Handlers.mouseOver // definitely need to move this, was a bug
		} else {
			if old.Handlers.mouseOver == nil {
				panic("listener has MouseOver but no underlying event listener")
			}

			changes = append(changes, &change{
				Type:        listenerDelete,
				Ref:         old,
				EventType:   dom.MouseOver,
				OldListener: old.Handlers.mouseOver,
			})

			changes = append(changes, &change{
				Type:        listenerAdd,
				Ref:         new,
				EventType:   dom.MouseOver,
				NewListener: new.Handlers.MouseOver,
			})
		}
	}

	oldMouseOut, newMouseOut := getMouseOut(old), getMouseOut(new)
	if oldMouseOut == nil && newMouseOut == nil {
		// no op
	}
	if oldMouseOut != nil && newMouseOut == nil {
		if old.Handlers.mouseOut == nil {
			panic("listener has MouseOut but no underlying event listener")
		}

		changes = append(changes, &change{
			Type:        listenerDelete,
			Ref:         old,
			EventType:   dom.MouseOut,
			OldListener: old.Handlers.mouseOut,
		})
	}
	if oldMouseOut == nil && newMouseOut != nil {
		changes = append(changes, &change{
			Type:        listenerAdd,
			Ref:         new,
			EventType:   dom.MouseOut,
			NewListener: new.Handlers.MouseOut,
		})
	}
	if oldMouseOut != nil && newMouseOut != nil && &oldMouseOut != &newMouseOut {
		// last attempt, check the cache
		if cacheListeners && new.Handlers.MouseOutCacheKey != "" && old.Handlers.MouseOutCacheKey != "" && new.Handlers.MouseOutCacheKey == old.Handlers.MouseOutCacheKey {
			// no need to syscall listener adding!
			new.Handlers.MouseOut = old.Handlers.MouseOut // not sure if we need to move this
			new.Handlers.mouseOut = old.Handlers.mouseOut // definitely need to move this, was a bug
		} else {
			if old.Handlers.mouseOut == nil {
				panic("listener has MouseOut but no underlying event listener")
			}

			changes = append(changes, &change{
				Type:        listenerDelete,
				Ref:         old,
				EventType:   dom.MouseOut,
				OldListener: old.Handlers.mouseOut,
			})

			changes = append(changes, &change{
				Type:        listenerAdd,
				Ref:         new,
				EventType:   dom.MouseOut,
				NewListener: new.Handlers.MouseOut,
			})
		}
	}

	oldMouseDown, newMouseDown := getMouseDown(old), getMouseDown(new)
	if oldMouseDown == nil && newMouseDown == nil {
		// no op
	}
	if oldMouseDown != nil && newMouseDown == nil {
		if old.Handlers.mouseDown == nil {
			panic("listener has MouseDown but no underlying event listener")
		}

		changes = append(changes, &change{
			Type:        listenerDelete,
			Ref:         old,
			EventType:   dom.MouseDown,
			OldListener: old.Handlers.mouseDown,
		})
	}
	if oldMouseDown == nil && newMouseDown != nil {
		changes = append(changes, &change{
			Type:        listenerAdd,
			Ref:         new,
			EventType:   dom.MouseDown,
			NewListener: new.Handlers.MouseDown,
		})
	}
	if oldMouseDown != nil && newMouseDown != nil && &oldMouseDown != &newMouseDown {
		// last attempt, check the cache
		if cacheListeners && new.Handlers.MouseDownCacheKey != "" && old.Handlers.MouseDownCacheKey != "" && new.Handlers.MouseDownCacheKey == old.Handlers.MouseDownCacheKey {
			// no need to syscall listener adding!
			new.Handlers.MouseDown = old.Handlers.MouseDown // not sure if we need to move this
			new.Handlers.mouseDown = old.Handlers.mouseDown // definitely need to move this, was a bug
		} else {
			if old.Handlers.mouseDown == nil {
				panic("listener has MouseDown but no underlying event listener")
			}

			changes = append(changes, &change{
				Type:        listenerDelete,
				Ref:         old,
				EventType:   dom.MouseDown,
				OldListener: old.Handlers.mouseDown,
			})

			changes = append(changes, &change{
				Type:        listenerAdd,
				Ref:         new,
				EventType:   dom.MouseDown,
				NewListener: new.Handlers.MouseDown,
			})
		}
	}

	oldMouseMove, newMouseMove := getMouseMove(old), getMouseMove(new)
	if oldMouseMove == nil && newMouseMove == nil {
		// no op
	}
	if oldMouseMove != nil && newMouseMove == nil {
		if old.Handlers.mouseMove == nil {
			panic("listener has MouseMove but no underlying event listener")
		}

		changes = append(changes, &change{
			Type:        listenerDelete,
			Ref:         old,
			EventType:   dom.MouseMove,
			OldListener: old.Handlers.mouseMove,
		})
	}
	if oldMouseMove == nil && newMouseMove != nil {
		changes = append(changes, &change{
			Type:        listenerAdd,
			Ref:         new,
			EventType:   dom.MouseMove,
			NewListener: new.Handlers.MouseMove,
		})
	}
	if oldMouseMove != nil && newMouseMove != nil && &oldMouseMove != &newMouseMove {
		// last attempt, check the cache
		if cacheListeners && new.Handlers.MouseMoveCacheKey != "" && old.Handlers.MouseMoveCacheKey != "" && new.Handlers.MouseMoveCacheKey == old.Handlers.MouseMoveCacheKey {
			// no need to syscall listener adding!
			new.Handlers.MouseMove = old.Handlers.MouseMove // not sure if we need to move this
			new.Handlers.mouseMove = old.Handlers.mouseMove // definitely need to move this, was a bug
		} else {
			if old.Handlers.mouseMove == nil {
				panic("listener has MouseMove but no underlying event listener")
			}

			changes = append(changes, &change{
				Type:        listenerDelete,
				Ref:         old,
				EventType:   dom.MouseMove,
				OldListener: old.Handlers.mouseMove,
			})

			changes = append(changes, &change{
				Type:        listenerAdd,
				Ref:         new,
				EventType:   dom.MouseMove,
				NewListener: new.Handlers.MouseMove,
			})
		}
	}

	oldMouseUp, newMouseUp := getMouseUp(old), getMouseUp(new)
	if oldMouseUp == nil && newMouseUp == nil {
		// no op
	}
	if oldMouseUp != nil && newMouseUp == nil {
		if old.Handlers.mouseUp == nil {
			panic("listener has MouseUp but no underlying event listener")
		}

		changes = append(changes, &change{
			Type:        listenerDelete,
			Ref:         old,
			EventType:   dom.MouseUp,
			OldListener: old.Handlers.mouseUp,
		})
	}
	if oldMouseUp == nil && newMouseUp != nil {
		changes = append(changes, &change{
			Type:        listenerAdd,
			Ref:         new,
			EventType:   dom.MouseUp,
			NewListener: new.Handlers.MouseUp,
		})
	}
	if oldMouseUp != nil && newMouseUp != nil && &oldMouseUp != &newMouseUp {
		// last attempt, check the cache
		if cacheListeners && new.Handlers.MouseUpCacheKey != "" && old.Handlers.MouseUpCacheKey != "" && new.Handlers.MouseUpCacheKey == old.Handlers.MouseUpCacheKey {
			// no need to syscall listener adding!
			new.Handlers.MouseUp = old.Handlers.MouseUp // not sure if we need to move this
			new.Handlers.mouseUp = old.Handlers.mouseUp // definitely need to move this, was a bug
		} else {
			if old.Handlers.mouseUp == nil {
				panic("listener has MouseUp but no underlying event listener")
			}

			changes = append(changes, &change{
				Type:        listenerDelete,
				Ref:         old,
				EventType:   dom.MouseUp,
				OldListener: old.Handlers.mouseUp,
			})

			changes = append(changes, &change{
				Type:        listenerAdd,
				Ref:         new,
				EventType:   dom.MouseUp,
				NewListener: new.Handlers.MouseUp,
			})
		}
	}

	oldMouseEnter, newMouseEnter := getMouseEnter(old), getMouseEnter(new)
	if oldMouseEnter == nil && newMouseEnter == nil {
		// no op
	}
	if oldMouseEnter != nil && newMouseEnter == nil {
		if old.Handlers.mouseEnter == nil {
			panic("listener has MouseEnter but no underlying event listener")
		}

		changes = append(changes, &change{
			Type:        listenerDelete,
			Ref:         old,
			EventType:   dom.MouseEnter,
			OldListener: old.Handlers.mouseEnter,
		})
	}
	if oldMouseEnter == nil && newMouseEnter != nil {
		changes = append(changes, &change{
			Type:        listenerAdd,
			Ref:         new,
			EventType:   dom.MouseEnter,
			NewListener: new.Handlers.MouseEnter,
		})
	}
	if oldMouseEnter != nil && newMouseEnter != nil && &oldMouseEnter != &newMouseEnter {
		// last attempt, check the cache
		if cacheListeners && new.Handlers.MouseEnterCacheKey != "" && old.Handlers.MouseEnterCacheKey != "" && new.Handlers.MouseEnterCacheKey == old.Handlers.MouseEnterCacheKey {
			// no need to syscall listener adding!
			new.Handlers.MouseEnter = old.Handlers.MouseEnter // not sure if we need to move this
			new.Handlers.mouseEnter = old.Handlers.mouseEnter // definitely need to move this, was a bug
		} else {
			if old.Handlers.mouseEnter == nil {
				panic("listener has MouseEnter but no underlying event listener")
			}

			changes = append(changes, &change{
				Type:        listenerDelete,
				Ref:         old,
				EventType:   dom.MouseEnter,
				OldListener: old.Handlers.mouseEnter,
			})

			changes = append(changes, &change{
				Type:        listenerAdd,
				Ref:         new,
				EventType:   dom.MouseEnter,
				NewListener: new.Handlers.MouseEnter,
			})
		}
	}

	oldMouseLeave, newMouseLeave := getMouseLeave(old), getMouseLeave(new)
	if oldMouseLeave == nil && newMouseLeave == nil {
		// no op
	}
	if oldMouseLeave != nil && newMouseLeave == nil {
		if old.Handlers.mouseLeave == nil {
			panic("listener has MouseLeave but no underlying event listener")
		}

		changes = append(changes, &change{
			Type:        listenerDelete,
			Ref:         old,
			EventType:   dom.MouseLeave,
			OldListener: old.Handlers.mouseLeave,
		})
	}
	if oldMouseLeave == nil && newMouseLeave != nil {
		changes = append(changes, &change{
			Type:        listenerAdd,
			Ref:         new,
			EventType:   dom.MouseLeave,
			NewListener: new.Handlers.MouseLeave,
		})
	}
	if oldMouseLeave != nil && newMouseLeave != nil && &oldMouseLeave != &newMouseLeave {
		// last attempt, check the cache
		if cacheListeners && new.Handlers.MouseLeaveCacheKey != "" && old.Handlers.MouseLeaveCacheKey != "" && new.Handlers.MouseLeaveCacheKey == old.Handlers.MouseLeaveCacheKey {
			// no need to syscall listener adding!
			new.Handlers.MouseLeave = old.Handlers.MouseLeave // not sure if we need to move this
			new.Handlers.mouseLeave = old.Handlers.mouseLeave // definitely need to move this, was a bug
		} else {
			if old.Handlers.mouseLeave == nil {
				panic("listener has MouseLeave but no underlying event listener")
			}

			changes = append(changes, &change{
				Type:        listenerDelete,
				Ref:         old,
				EventType:   dom.MouseLeave,
				OldListener: old.Handlers.mouseLeave,
			})

			changes = append(changes, &change{
				Type:        listenerAdd,
				Ref:         new,
				EventType:   dom.MouseLeave,
				NewListener: new.Handlers.MouseLeave,
			})
		}
	}

	oldInput, newInput := getInput(old), getInput(new)
	if oldInput == nil && newInput == nil {
		// no op
	}
	if oldInput != nil && newInput == nil {
		if old.Handlers.input == nil {
			panic("listener has Input but no underlying event listener")
		}

		changes = append(changes, &change{
			Type:        listenerDelete,
			Ref:         old,
			EventType:   dom.Input,
			OldListener: old.Handlers.input,
		})
	}
	if oldInput == nil && newInput != nil {
		changes = append(changes, &change{
			Type:        listenerAdd,
			Ref:         new,
			EventType:   dom.Input,
			NewListener: new.Handlers.Input,
		})
	}
	if oldInput != nil && newInput != nil && &oldInput != &newInput {
		// last attempt, check the cache
		if cacheListeners && new.Handlers.InputCacheKey != "" && old.Handlers.InputCacheKey != "" && new.Handlers.InputCacheKey == old.Handlers.InputCacheKey {
			// no need to syscall listener adding!
			new.Handlers.Input = old.Handlers.Input // not sure if we need to move this
			new.Handlers.input = old.Handlers.input // definitely need to move this, was a bug
		} else {
			if old.Handlers.input == nil {
				panic("listener has Input but no underlying event listener")
			}

			changes = append(changes, &change{
				Type:        listenerDelete,
				Ref:         old,
				EventType:   dom.Input,
				OldListener: old.Handlers.input,
			})

			changes = append(changes, &change{
				Type:        listenerAdd,
				Ref:         new,
				EventType:   dom.Input,
				NewListener: new.Handlers.Input,
			})
		}
	}

	oldKeyUp, newKeyUp := getKeyUp(old), getKeyUp(new)
	if oldKeyUp == nil && newKeyUp == nil {
		// no op
	}
	if oldKeyUp != nil && newKeyUp == nil {
		if old.Handlers.keyUp == nil {
			panic("listener has KeyUp but no underlying event listener")
		}

		changes = append(changes, &change{
			Type:        listenerDelete,
			Ref:         old,
			EventType:   dom.KeyUp,
			OldListener: old.Handlers.keyUp,
		})
	}
	if oldKeyUp == nil && newKeyUp != nil {
		changes = append(changes, &change{
			Type:        listenerAdd,
			Ref:         new,
			EventType:   dom.KeyUp,
			NewListener: new.Handlers.KeyUp,
		})
	}
	if oldKeyUp != nil && newKeyUp != nil && &oldKeyUp != &newKeyUp {
		// last attempt, check the cache
		if cacheListeners && new.Handlers.KeyUpCacheKey != "" && old.Handlers.KeyUpCacheKey != "" && new.Handlers.KeyUpCacheKey == old.Handlers.KeyUpCacheKey {
			// no need to syscall listener adding!
			new.Handlers.KeyUp = old.Handlers.KeyUp // not sure if we need to move this
			new.Handlers.keyUp = old.Handlers.keyUp // definitely need to move this, was a bug
		} else {
			if old.Handlers.keyUp == nil {
				panic("listener has KeyUp but no underlying event listener")
			}

			changes = append(changes, &change{
				Type:        listenerDelete,
				Ref:         old,
				EventType:   dom.KeyUp,
				OldListener: old.Handlers.keyUp,
			})

			changes = append(changes, &change{
				Type:        listenerAdd,
				Ref:         new,
				EventType:   dom.KeyUp,
				NewListener: new.Handlers.KeyUp,
			})
		}
	}

	oldKeyDown, newKeyDown := getKeyDown(old), getKeyDown(new)
	if oldKeyDown == nil && newKeyDown == nil {
		// no op
	}
	if oldKeyDown != nil && newKeyDown == nil {
		if old.Handlers.keyDown == nil {
			panic("listener has KeyDown but no underlying event listener")
		}

		changes = append(changes, &change{
			Type:        listenerDelete,
			Ref:         old,
			EventType:   dom.KeyDown,
			OldListener: old.Handlers.keyDown,
		})
	}
	if oldKeyDown == nil && newKeyDown != nil {
		changes = append(changes, &change{
			Type:        listenerAdd,
			Ref:         new,
			EventType:   dom.KeyDown,
			NewListener: new.Handlers.KeyDown,
		})
	}
	if oldKeyDown != nil && newKeyDown != nil && &oldKeyDown != &newKeyDown {
		// last attempt, check the cache
		if cacheListeners && new.Handlers.KeyDownCacheKey != "" && old.Handlers.KeyDownCacheKey != "" && new.Handlers.KeyDownCacheKey == old.Handlers.KeyDownCacheKey {
			// no need to syscall listener adding!
			new.Handlers.KeyDown = old.Handlers.KeyDown // not sure if we need to move this
			new.Handlers.keyDown = old.Handlers.keyDown // definitely need to move this, was a bug
		} else {
			if old.Handlers.keyDown == nil {
				panic("listener has KeyDown but no underlying event listener")
			}

			changes = append(changes, &change{
				Type:        listenerDelete,
				Ref:         old,
				EventType:   dom.KeyDown,
				OldListener: old.Handlers.keyDown,
			})

			changes = append(changes, &change{
				Type:        listenerAdd,
				Ref:         new,
				EventType:   dom.KeyDown,
				NewListener: new.Handlers.KeyDown,
			})
		}
	}

	oldDrop, newDrop := getDrop(old), getDrop(new)
	if oldDrop == nil && newDrop == nil {
		// no op
	}
	if oldDrop != nil && newDrop == nil {
		if old.Handlers.drop == nil {
			panic("listener has Drop but no underlying event listener")
		}

		changes = append(changes, &change{
			Type:        listenerDelete,
			Ref:         old,
			EventType:   dom.Drop,
			OldListener: old.Handlers.drop,
		})
	}
	if oldDrop == nil && newDrop != nil {
		changes = append(changes, &change{
			Type:        listenerAdd,
			Ref:         new,
			EventType:   dom.Drop,
			NewListener: new.Handlers.Drop,
		})
	}
	if oldDrop != nil && newDrop != nil && &oldDrop != &newDrop {
		// last attempt, check the cache
		if cacheListeners && new.Handlers.DropCacheKey != "" && old.Handlers.DropCacheKey != "" && new.Handlers.DropCacheKey == old.Handlers.DropCacheKey {
			// no need to syscall listener adding!
			new.Handlers.Drop = old.Handlers.Drop // not sure if we need to move this
			new.Handlers.drop = old.Handlers.drop // definitely need to move this, was a bug
		} else {
			if old.Handlers.drop == nil {
				panic("listener has Drop but no underlying event listener")
			}

			changes = append(changes, &change{
				Type:        listenerDelete,
				Ref:         old,
				EventType:   dom.Drop,
				OldListener: old.Handlers.drop,
			})

			changes = append(changes, &change{
				Type:        listenerAdd,
				Ref:         new,
				EventType:   dom.Drop,
				NewListener: new.Handlers.Drop,
			})
		}
	}

	oldDragOver, newDragOver := getDragOver(old), getDragOver(new)
	if oldDragOver == nil && newDragOver == nil {
		// no op
	}
	if oldDragOver != nil && newDragOver == nil {
		if old.Handlers.dragOver == nil {
			panic("listener has DragOver but no underlying event listener")
		}

		changes = append(changes, &change{
			Type:        listenerDelete,
			Ref:         old,
			EventType:   dom.DragOver,
			OldListener: old.Handlers.dragOver,
		})
	}
	if oldDragOver == nil && newDragOver != nil {
		changes = append(changes, &change{
			Type:        listenerAdd,
			Ref:         new,
			EventType:   dom.DragOver,
			NewListener: new.Handlers.DragOver,
		})
	}
	if oldDragOver != nil && newDragOver != nil && &oldDragOver != &newDragOver {
		// last attempt, check the cache
		if cacheListeners && new.Handlers.DragOverCacheKey != "" && old.Handlers.DragOverCacheKey != "" && new.Handlers.DragOverCacheKey == old.Handlers.DragOverCacheKey {
			// no need to syscall listener adding!
			new.Handlers.DragOver = old.Handlers.DragOver // not sure if we need to move this
			new.Handlers.dragOver = old.Handlers.dragOver // definitely need to move this, was a bug
		} else {
			if old.Handlers.dragOver == nil {
				panic("listener has DragOver but no underlying event listener")
			}

			changes = append(changes, &change{
				Type:        listenerDelete,
				Ref:         old,
				EventType:   dom.DragOver,
				OldListener: old.Handlers.dragOver,
			})

			changes = append(changes, &change{
				Type:        listenerAdd,
				Ref:         new,
				EventType:   dom.DragOver,
				NewListener: new.Handlers.DragOver,
			})
		}
	}

	return
}

// }}}

// reconcileCanvasDraw {{{

func reconcileCanvasDraw(old, new *Node) (changes []*change) {
	// only draw canvases
	if new.DataAtom != atom.Canvas {
		return
	}

	// for now, nothing fancy, but lots of redrawing... - NCL 2/19/22
	changes = append(changes, &change{
		Type:       canvasDraw,
		Ref:        new,
		CanvasDraw: new.CanvasDraw,
	})
	return
}

// }}}

func printl(l int, f string, vs ...interface{}) {
	sout := ""
	for i := 0; i < l; i++ {
		sout += " "
	}
	log.Printf(sout+f, vs...)
}

func reconcile(parent, old, new *Node, level int) (changes []*change, replaced bool) {
	if old.Type != new.Type {
		//printl(level, "different types")
		changes = replaces(parent, old, new, level)
		replaced = true
		return
	}

	switch old.Type {
	case html.ElementNode:
		// The second check here is for custom components (i.e., WebComponents) (e.g., sl-button)
		if (old.DataAtom != new.DataAtom) || (old.Data != new.Data) {
			changes = replaces(parent, old, new, level)
			replaced = true
			return
		}

		// at this point, we have decided to mutate the node.
		// i.e., the new node's underlying dom element is the same
		// as the old one
		new.rendered = old.rendered
		new.renderedElement = old.renderedElement

		//log.Printf("old node! %+v with style %s", old, old.Style.Val())
		//log.Printf("new node! %+v with style %s", new, new.Style.Val())
		changes = append(changes, diffStyles(new, &old.Style, &new.Style, level)...)
		changes = append(changes, reconcileHandlers(old, new)...)
		changes = append(changes, reconcileAttr(new, old.Attr, new.Attr)...)
		changes = append(changes, reconcileCanvasDraw(old, new)...)
	case html.TextNode:
		if old.Data != new.Data {
			changes = append(changes, &change{
				Type:   replace,
				Parent: parent,
				Old:    old,
				Ref:    new,
			})
			replaced = true
			return
		}

		new.rendered = old.rendered
	default:
		panic("unknown node type")
	}

	return
}

func replaces(parent, old, new *Node, level int) (changes []*change) {
	changes = append(changes, &change{
		Type:   replace,
		Parent: parent,
		Old:    old,
		Ref:    new,
	})

	// I believe both text and elements can have listeners - NCL 2/4/22
	changes = append(changes, reconcileHandlers(nil, new)...)

	// only elements can have style or attrs or are canvases
	if new.Type == html.ElementNode {
		changes = append(changes, diffStyles(new, nil, &new.Style, level)...)
		changes = append(changes, reconcileAttr(new, nil, new.Attr)...)
		changes = append(changes, reconcileCanvasDraw(nil, new)...)
	}

	// insert each of the new node's children
	for _, c := range new.Children {
		changes = append(changes, inserts(new, c, level+1)...)
	}

	return
}

func inserts(into, root *Node, level int) (changes []*change) {
	if into == nil {
		panic("into can't be nil")
	}
	if root == nil {
		panic("root can't be nil")
	}

	changes = append(changes, &change{
		Type:   insert,
		Parent: into,
		Ref:    root,
	})

	if root.Type == html.TextNode {
		// it is unclear to me if this will work, but it can be worked
		// around by wrapping text nodes in spans or divs - NCL 1/30/22
		changes = append(changes, reconcileHandlers(nil, root)...) // reconcile the root's listeners against empty listeners
		return
	}

	changes = append(changes, diffStyles(root, nil, &root.Style, level)...) // reconcile the root's styles against empty styles
	changes = append(changes, reconcileHandlers(nil, root)...)              // reconcile the root's listeners against empty listeners
	changes = append(changes, reconcileAttr(root, nil, root.Attr)...)       // reconcile the root's attr's against empty attrs
	changes = append(changes, reconcileCanvasDraw(nil, root)...)

	for _, c := range root.Children { // recurse, for each child
		changes = append(changes, inserts(root, c, level+1)...)
	}

	return
}
