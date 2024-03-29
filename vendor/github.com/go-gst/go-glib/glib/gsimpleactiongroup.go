package glib

// #include <gio/gio.h>
// #include <glib.h>
// #include <glib-object.h>
// #include "glib.go.h"
import "C"
import (
	"unsafe"
)

// SimpleActionGroup is a representation of glib's GSimpleActionGroup
type SimpleActionGroup struct {
	*Object

	// Interfaces
	IActionMap
	IActionGroup
}

// deprecated since 2.38:
// g_simple_action_group_lookup()
// g_simple_action_group_insert()
// g_simple_action_group_remove()
// g_simple_action_group_add_entries()
// -> See implementations in ActionMap

// native() returns a pointer to the underlying GSimpleActionGroup.
func (v *SimpleActionGroup) native() *C.GSimpleActionGroup {
	if v == nil || v.GObject == nil {
		return nil
	}
	return C.toGSimpleActionGroup(unsafe.Pointer(v.GObject))
}

func (v *SimpleActionGroup) Native() unsafe.Pointer {
	return unsafe.Pointer(v.native())
}

func marshalSimpleActionGroup(p unsafe.Pointer) (interface{}, error) {
	c := C.g_value_get_object((*C.GValue)(p))
	return wrapSimpleActionGroup(wrapObject(unsafe.Pointer(c))), nil
}

func wrapSimpleActionGroup(obj *Object) *SimpleActionGroup {
	am := wrapActionMap(obj)
	ag := wrapActionGroup(obj)
	return &SimpleActionGroup{obj, am, ag}
}

// SimpleActionGroupNew is a wrapper around g_simple_action_group_new
func SimpleActionGroupNew() *SimpleActionGroup {
	c := C.g_simple_action_group_new()
	if c == nil {
		return nil
	}
	return wrapSimpleActionGroup(wrapObject(unsafe.Pointer(c)))
}
