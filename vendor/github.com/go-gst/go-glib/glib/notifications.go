package glib

// #include <gio/gio.h>
// #include <glib.h>
// #include <glib-object.h>
// #include "glib.go.h"
import "C"
import "unsafe"

// Notification is a representation of GNotification.
type Notification struct {
	*Object
}

// native() returns a pointer to the underlying GNotification.
func (v *Notification) native() *C.GNotification {
	if v == nil || v.GObject == nil {
		return nil
	}
	return C.toGNotification(unsafe.Pointer(v.GObject))
}

func (v *Notification) Native() unsafe.Pointer {
	return unsafe.Pointer(v.native())
}

func marshalNotification(p unsafe.Pointer) (interface{}, error) {
	c := C.g_value_get_object((*C.GValue)(p))
	return wrapNotification(wrapObject(unsafe.Pointer(c))), nil
}

func wrapNotification(obj *Object) *Notification {
	return &Notification{obj}
}

// NotificationNew is a wrapper around g_notification_new().
func NotificationNew(title string) *Notification {
	cstr1 := (*C.gchar)(C.CString(title))
	defer C.free(unsafe.Pointer(cstr1))

	c := C.g_notification_new(cstr1)
	if c == nil {
		return nil
	}
	return wrapNotification(wrapObject(unsafe.Pointer(c)))
}

// SetTitle is a wrapper around g_notification_set_title().
func (v *Notification) SetTitle(title string) {
	cstr1 := (*C.gchar)(C.CString(title))
	defer C.free(unsafe.Pointer(cstr1))

	C.g_notification_set_title(v.native(), cstr1)
}

// SetBody is a wrapper around g_notification_set_body().
func (v *Notification) SetBody(body string) {
	cstr1 := (*C.gchar)(C.CString(body))
	defer C.free(unsafe.Pointer(cstr1))

	C.g_notification_set_body(v.native(), cstr1)
}

// SetDefaultAction is a wrapper around g_notification_set_default_action().
func (v *Notification) SetDefaultAction(detailedAction string) {
	cstr1 := (*C.gchar)(C.CString(detailedAction))
	defer C.free(unsafe.Pointer(cstr1))

	C.g_notification_set_default_action(v.native(), cstr1)
}

// AddButton is a wrapper around g_notification_add_button().
func (v *Notification) AddButton(label, detailedAction string) {
	cstr1 := (*C.gchar)(C.CString(label))
	defer C.free(unsafe.Pointer(cstr1))

	cstr2 := (*C.gchar)(C.CString(detailedAction))
	defer C.free(unsafe.Pointer(cstr2))

	C.g_notification_add_button(v.native(), cstr1, cstr2)
}

// SetIcon is a wrapper around g_notification_set_icon().
func (v *Notification) SetIcon(iconPath string) {
	fileIcon := FileIconNew(iconPath)

	C.g_notification_set_icon(v.native(), (*C.GIcon)(fileIcon.native()))
}

// void 	g_notification_set_default_action_and_target () // requires varargs
// void 	g_notification_set_default_action_and_target_value () // requires variant
// void 	g_notification_add_button_with_target () // requires varargs
// void 	g_notification_add_button_with_target_value () //requires variant
// void 	g_notification_set_urgent () // Deprecated, so not implemented
