//go:build js && wasm

package js

import (
	"sync"
	"syscall/js"

	"github.com/spinsrv/browser"
)

var (
	// this type checks the implementations
	_ browser.LocalStorage = (*localstorage)(nil)
)

var DefaultLocalStorage browser.LocalStorage = &localstorage{
	underlying: js.Global().Get("localStorage"),
}

type localstorage struct {
	underlying js.Value
	m          sync.Mutex
}

func (l localstorage) Put(k, v string) {
	l.m.Lock()
	defer l.m.Unlock()

	l.underlying.Call("setItem", k, v)
}

func (l localstorage) Get(k string) string {
	l.m.Lock()
	defer l.m.Unlock()

	v := l.underlying.Call("getItem", k)
	//	log.Printf("localStorage.getItem(%s) -> %v", k, v)

	if v.Equal(js.Null()) {
		return ""
	}

	if v.Equal(js.Undefined()) {
		return ""
	}

	return v.String()
}

func (l localstorage) Del(k string) {
	l.m.Lock()
	defer l.m.Unlock()

	l.underlying.Call("removeItem", k)
}
