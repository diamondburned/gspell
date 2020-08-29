package gspell

import (
	"github.com/diamondburned/gspell/internal/callback"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"unsafe"
)

// #cgo pkg-config: gspell-1 gtk+-3.0 glib-2.0 gio-2.0 glib-2.0 gobject-2.0
// #include <gspell/gspell.h>
// #include <gtk/gtk.h>
// #include <gio/gio.h>
// #include <glib.h>
// #include <glib-object.h>
// extern void callbackDelete(gpointer ptr);
//
import "C"

//export callbackDelete
func callbackDelete(ptr C.gpointer) {
	callback.Delete(uintptr(ptr))
}

// objector is used internally for other interfaces.
type objector interface {
	glib.IObject
	Connect(string, interface{}, ...interface{}) (glib.SignalHandle, error)
	ConnectAfter(string, interface{}, ...interface{}) (glib.SignalHandle, error)
	GetProperty(name string) (interface{}, error)
	SetProperty(name string, value interface{}) error
	Native() uintptr
}

// asserting objector interface
var _ objector = (*glib.Object)(nil)

// Caster is the interface that allows casting objects to widgets.
type Caster interface {
	objector
	Cast() (gtk.IWidget, error)
}

func init() {
	glib.RegisterGValueMarshalers([]glib.TypeMarshaler{
		// Enums
		{glib.Type(C.gspell_checker_error_get_type()), marshalCheckerError},

		// Objects/Classes
		{glib.Type(C.gspell_checker_get_type()), marshalChecker},
		{glib.Type(C.gspell_checker_dialog_get_type()), marshalCheckerDialog},
		{glib.Type(C.gspell_entry_get_type()), marshalEntry},
		{glib.Type(C.gspell_entry_buffer_get_type()), marshalEntryBuffer},
		{glib.Type(C.gspell_language_chooser_button_get_type()), marshalLanguageChooserButton},
		{glib.Type(C.gspell_language_chooser_dialog_get_type()), marshalLanguageChooserDialog},
		{glib.Type(C.gspell_navigator_text_view_get_type()), marshalNavigatorTextView},
		{glib.Type(C.gspell_text_buffer_get_type()), marshalTextBuffer},
		{glib.Type(C.gspell_text_view_get_type()), marshalTextView},

		// Boxed
		{glib.Type(C.gspell_language_get_type()), marshalLanguage},
	})
}

// CheckerError an error code used with GSPELL_CHECKER_ERROR in a #GError
// returned from a spell-checker-related function.
type CheckerError int

func marshalCheckerError(p uintptr) (interface{}, error) {
	return CheckerError(C.g_value_get_enum((*C.GValue)(unsafe.Pointer(p)))), nil
}

const (
	// CheckerErrorDictionary dictionary error.
	CheckerErrorDictionary CheckerError = 0
	// CheckerErrorNoLanguageSet no language set.
	CheckerErrorNoLanguageSet CheckerError = 1
)

type LanguageChooserer interface {
	objector
	GetLanguage() *Language
	GetLanguageCode() string
	// SetLanguage sets the selected language.
	SetLanguage(language *Language)
	SetLanguageCode(languageCode string)
}

type LanguageChooser struct {
	*glib.Object
}

// native turns the current *LanguageChooser into the native C pointer type.
func (l *LanguageChooser) native() *C.GspellLanguageChooser {
	return (*C.GspellLanguageChooser)(unsafe.Pointer(l.Native()))
}
func (l *LanguageChooser) GetLanguage() *Language {
	r := (*Language)(C.gspell_language_chooser_get_language(l.native()))
	return r
}
func (l *LanguageChooser) GetLanguageCode() string {
	r := C.GoString(C.gspell_language_chooser_get_language_code(l.native()))
	return r
}

// SetLanguage sets the selected language.
func (l *LanguageChooser) SetLanguage(language *Language) {
	v1 := (*C.GspellLanguage)(unsafe.Pointer(language.Native()))
	C.gspell_language_chooser_set_language(l.native(), v1)
}
func (l *LanguageChooser) SetLanguageCode(languageCode string) {
	v1 := C.CString(languageCode)
	defer C.free(unsafe.Pointer(v1))
	C.gspell_language_chooser_set_language_code(l.native(), v1)
}

type Navigatorer interface {
	objector
	// Change changes the current word by change_to in the text. word must be the
	// same as returned by the last call to C.gspell_navigator_goto_next().
	//
	// This function doesn't call (*Checker).SetCorrection(). A widget using a
	// Navigator should call (*Checker).SetCorrection() in addition to this
	// function.
	Change(word string, changeTo string)
	// ChangeAll changes all occurrences of word by change_to in the text.
	//
	// This function doesn't call (*Checker).SetCorrection(). A widget using a
	// Navigator should call (*Checker).SetCorrection() in addition to this
	// function.
	ChangeAll(word string, changeTo string)
}

type Navigator struct {
	*glib.Object
}

// native turns the current *Navigator into the native C pointer type.
func (n *Navigator) native() *C.GspellNavigator {
	return (*C.GspellNavigator)(unsafe.Pointer(n.Native()))
}

// Change changes the current word by change_to in the text. word must be the
// same as returned by the last call to C.gspell_navigator_goto_next().
//
// This function doesn't call (*Checker).SetCorrection(). A widget using a
// Navigator should call (*Checker).SetCorrection() in addition to this
// function.
func (n *Navigator) Change(word string, changeTo string) {
	v1 := C.CString(word)
	defer C.free(unsafe.Pointer(v1))
	v2 := C.CString(changeTo)
	defer C.free(unsafe.Pointer(v2))

	C.gspell_navigator_change(n.native(), v1, v2)
}

// ChangeAll changes all occurrences of word by change_to in the text.
//
// This function doesn't call (*Checker).SetCorrection(). A widget using a
// Navigator should call (*Checker).SetCorrection() in addition to this
// function.
func (n *Navigator) ChangeAll(word string, changeTo string) {
	v1 := C.CString(word)
	defer C.free(unsafe.Pointer(v1))
	v2 := C.CString(changeTo)
	defer C.free(unsafe.Pointer(v2))

	C.gspell_navigator_change_all(n.native(), v1, v2)
}

func CheckerErrorQuark() glib.Quark {
	r := glib.Quark(C.gspell_checker_error_quark())
	return r
}
func LanguageGetAvailable() *glib.List {
	r := glib.WrapList(uintptr(unsafe.Pointer(C.gspell_language_get_available())))
	return r
}

// LanguageGetDefault finds the best available language based on the current
// locale.
func LanguageGetDefault() *Language {
	r := (*Language)(C.gspell_language_get_default())
	return r
}
func LanguageLookup(languageCode string) *Language {
	v1 := C.CString(languageCode)
	defer C.free(unsafe.Pointer(v1))
	r := (*Language)(C.gspell_language_lookup(v1))
	return r
}

type Checker struct {
	*glib.Object
}

// wrapChecker wraps the given pointer to *Checker.
func wrapChecker(ptr unsafe.Pointer) *Checker {
	obj := glib.Take(ptr)
	return &Checker{
		Object: obj,
	}
}

func marshalChecker(p uintptr) (interface{}, error) {
	return wrapChecker(unsafe.Pointer(C.g_value_get_object((*C.GValue)(unsafe.Pointer(p))))), nil
}

// CheckerNew creates a new Checker. If language is nil, the default language is
// picked with LanguageGetDefault().
func CheckerNew(language *Language) *Checker {
	v1 := (*C.GspellLanguage)(unsafe.Pointer(language.Native()))
	return wrapChecker(unsafe.Pointer(C.gspell_checker_new(v1)))
}

// native turns the current *Checker into the native C pointer type.
func (c *Checker) native() *C.GspellChecker {
	return (*C.GspellChecker)(unsafe.Pointer(c.Object.Native()))
}

// AddWordToPersonal adds a word to the personal dictionary. It is typically
// saved in the user's home directory.
func (c *Checker) AddWordToPersonal(word string, wordLength int) {
	v1 := C.CString(word)
	defer C.free(unsafe.Pointer(v1))
	v2 := C.gssize(wordLength)

	C.gspell_checker_add_word_to_personal(c.native(), v1, v2)
}

// AddWordToSession adds a word to the session dictionary. Each Checker instance
// has a different session dictionary. The session dictionary is lost when the
// Checker:language property changes or when checker is destroyed or when
// (*Checker).ClearSession() is called.
//
// This function is typically called for an “Ignore All” action.
func (c *Checker) AddWordToSession(word string, wordLength int) {
	v1 := C.CString(word)
	defer C.free(unsafe.Pointer(v1))
	v2 := C.gssize(wordLength)

	C.gspell_checker_add_word_to_session(c.native(), v1, v2)
}

// CheckWord if the Checker:language is nil, i.e. when no dictonaries are
// available, this function returns true to limit the damage.
func (c *Checker) CheckWord(word string, wordLength int) bool {
	v1 := C.CString(word)
	defer C.free(unsafe.Pointer(v1))
	v2 := C.gssize(wordLength)

	r := gobool(C.gspell_checker_check_word(c.native(), v1, v2, nil))
	return r
}

// ClearSession clears the session dictionary.
func (c *Checker) ClearSession() {
	C.gspell_checker_clear_session(c.native())
}

// GetEnchantDict gets the EnchantDict currently used by checker. It permits to
// extend Checker with more features. Note that by doing so, the other classes
// in gspell may no longer work well.
//
// Checker re-creates a new EnchantDict when the Checker:language is changed and
// when the session is cleared.
func (c *Checker) GetEnchantDict() {
	C.gspell_checker_get_enchant_dict(c.native())
}
func (c *Checker) GetLanguage() *Language {
	r := (*Language)(C.gspell_checker_get_language(c.native()))
	return r
}

// GetSuggestions gets the suggestions for word. Free the return value with
// g_slist_free_full(suggestions, g_free).
func (c *Checker) GetSuggestions(word string, wordLength int) *glib.SList {
	v1 := C.CString(word)
	defer C.free(unsafe.Pointer(v1))
	v2 := C.gssize(wordLength)

	r := glib.WrapSList(uintptr(unsafe.Pointer(C.gspell_checker_get_suggestions(c.native(), v1, v2))))
	return r
}

// SetCorrection informs the spell checker that word is replaced/corrected by
// replacement.
func (c *Checker) SetCorrection(word string, wordLength int, replacement string, replacementLength int) {
	v1 := C.CString(word)
	defer C.free(unsafe.Pointer(v1))
	v2 := C.gssize(wordLength)
	v3 := C.CString(replacement)
	defer C.free(unsafe.Pointer(v3))
	v4 := C.gssize(replacementLength)

	C.gspell_checker_set_correction(c.native(), v1, v2, v3, v4)
}

// SetLanguage sets the language to use for the spell checking. If language is
// nil, the default language is picked with LanguageGetDefault().
func (c *Checker) SetLanguage(language *Language) {
	v1 := (*C.GspellLanguage)(unsafe.Pointer(language.Native()))
	C.gspell_checker_set_language(c.native(), v1)
}

type CheckerDialog struct {
	gtk.Dialog
}

// wrapCheckerDialog wraps the given pointer to *CheckerDialog.
func wrapCheckerDialog(ptr unsafe.Pointer) *CheckerDialog {
	obj := glib.Take(ptr)
	return &CheckerDialog{
		Dialog: gtk.Dialog{
			Window: gtk.Window{
				Bin: gtk.Bin{
					Container: gtk.Container{
						Widget: gtk.Widget{
							InitiallyUnowned: glib.InitiallyUnowned{
								Object: obj,
							},
						},
					},
				},
			},
		},
	}
}

func marshalCheckerDialog(p uintptr) (interface{}, error) {
	return wrapCheckerDialog(unsafe.Pointer(C.g_value_get_object((*C.GValue)(unsafe.Pointer(p))))), nil
}

// CheckerDialogNew creates a new CheckerDialog.
func CheckerDialogNew(parent *gtk.Window, navigator *Navigator) *CheckerDialog {
	v1 := (*C.GtkWindow)(unsafe.Pointer(parent.Widget.Native()))
	v2 := (*C.GspellNavigator)(unsafe.Pointer(navigator.Native()))

	return wrapCheckerDialog(unsafe.Pointer(C.gspell_checker_dialog_new(v1, v2)))
}

// native turns the current *CheckerDialog into the native C pointer type.
func (c *CheckerDialog) native() *C.GspellCheckerDialog {
	return (*C.GspellCheckerDialog)(gwidget(&c.Dialog))
}

func (c *CheckerDialog) GetSpellNavigator() *Navigator {
	obj := glib.Take(unsafe.Pointer(C.gspell_checker_dialog_get_spell_navigator(c.native())))
	r := &Navigator{obj}
	return r
}

type Entry struct {
	*glib.Object
}

// wrapEntry wraps the given pointer to *Entry.
func wrapEntry(ptr unsafe.Pointer) *Entry {
	obj := glib.Take(ptr)
	return &Entry{
		Object: obj,
	}
}

func marshalEntry(p uintptr) (interface{}, error) {
	return wrapEntry(unsafe.Pointer(C.g_value_get_object((*C.GValue)(unsafe.Pointer(p))))), nil
}

// native turns the current *Entry into the native C pointer type.
func (e *Entry) native() *C.GspellEntry {
	return (*C.GspellEntry)(unsafe.Pointer(e.Object.Native()))
}

// GetFromGtkEntry returns the Entry of gtk_entry. The returned object is
// guaranteed to be the same for the lifetime of gtk_entry.
func GetFromGtkEntry(gtkEntry *gtk.Entry) *Entry {
	v1 := (*C.GtkEntry)(unsafe.Pointer(gtkEntry.Widget.Native()))
	r := wrapEntry(unsafe.Pointer(C.gspell_entry_get_from_gtk_entry(v1)))
	return r
}

// BasicSetup function is a convenience function that does the following: - Set
// a spell checker. The language chosen is the one returned by
// LanguageGetDefault(). - Set the Entry:inline-spell-checking property to true.
//
// Example: |[ GtkEntry *gtk_entry; GspellEntry *gspell_entry;
//
//    gspell_entry = gspell_entry_get_from_gtk_entry (gtk_entry);
//    gspell_entry_basic_setup (gspell_entry);
//
//
//
//    GtkEntry *gtk_entry;
//    GspellEntry *gspell_entry;
//    GspellChecker *checker;
//    GtkEntryBuffer *gtk_buffer;
//    GspellEntryBuffer *gspell_buffer;
//
//    checker = gspell_checker_new (NULL);
//    gtk_buffer = gtk_entry_get_buffer (gtk_entry);
//    gspell_buffer = gspell_entry_buffer_get_from_gtk_entry_buffer (gtk_buffer);
//    gspell_entry_buffer_set_spell_checker (gspell_buffer, checker);
//    g_object_unref (checker);
//
//    gspell_entry = gspell_entry_get_from_gtk_entry (gtk_entry);
//    gspell_entry_set_inline_spell_checking (gspell_entry, TRUE);
//
func (e *Entry) BasicSetup() {
	C.gspell_entry_basic_setup(e.native())
}
func (e *Entry) GetEntry() *gtk.Entry {
	obj := glib.Take(unsafe.Pointer(C.gspell_entry_get_entry(e.native())))
	r := &gtk.Entry{
		Widget: gtk.Widget{
			InitiallyUnowned: glib.InitiallyUnowned{
				Object: obj,
			},
		},
	}
	return r
}
func (e *Entry) GetInlineSpellChecking() bool {
	r := gobool(C.gspell_entry_get_inline_spell_checking(e.native()))
	return r
}

// SetInlineSpellChecking sets the Entry:inline-spell-checking property.
func (e *Entry) SetInlineSpellChecking(enable bool) {
	v1 := cbool(enable)
	C.gspell_entry_set_inline_spell_checking(e.native(), v1)
}

type EntryBuffer struct {
	*glib.Object
}

// wrapEntryBuffer wraps the given pointer to *EntryBuffer.
func wrapEntryBuffer(ptr unsafe.Pointer) *EntryBuffer {
	obj := glib.Take(ptr)
	return &EntryBuffer{
		Object: obj,
	}
}

func marshalEntryBuffer(p uintptr) (interface{}, error) {
	return wrapEntryBuffer(unsafe.Pointer(C.g_value_get_object((*C.GValue)(unsafe.Pointer(p))))), nil
}

// native turns the current *EntryBuffer into the native C pointer type.
func (e *EntryBuffer) native() *C.GspellEntryBuffer {
	return (*C.GspellEntryBuffer)(unsafe.Pointer(e.Object.Native()))
}

// GetFromGtkEntryBuffer returns the EntryBuffer of gtk_buffer. The returned
// object is guaranteed to be the same for the lifetime of gtk_buffer.
func GetFromGtkEntryBuffer(gtkBuffer *gtk.EntryBuffer) *EntryBuffer {
	v1 := (*C.GtkEntryBuffer)(unsafe.Pointer(gtkBuffer.Native()))
	r := wrapEntryBuffer(unsafe.Pointer(C.gspell_entry_buffer_get_from_gtk_entry_buffer(v1)))
	return r
}

func (e *EntryBuffer) GetBuffer() *gtk.EntryBuffer {
	obj := glib.Take(unsafe.Pointer(C.gspell_entry_buffer_get_buffer(e.native())))
	r := &gtk.EntryBuffer{
		Object: obj,
	}
	return r
}
func (e *EntryBuffer) GetSpellChecker() *Checker {
	r := wrapChecker(unsafe.Pointer(C.gspell_entry_buffer_get_spell_checker(e.native())))
	return r
}

// SetSpellChecker sets a Checker to a EntryBuffer. The gspell_buffer will own a
// reference to spell_checker, so you can release your reference to
// spell_checker if you no longer need it.
func (e *EntryBuffer) SetSpellChecker(spellChecker *Checker) {
	v1 := (*C.GspellChecker)(unsafe.Pointer(spellChecker.Native()))
	C.gspell_entry_buffer_set_spell_checker(e.native(), v1)
}

type LanguageChooserButton struct {
	gtk.Button

	// Interfaces
	LanguageChooserer
	gtk.Actionable
}

// wrapLanguageChooserButton wraps the given pointer to *LanguageChooserButton.
func wrapLanguageChooserButton(ptr unsafe.Pointer) *LanguageChooserButton {
	obj := glib.Take(ptr)
	return &LanguageChooserButton{
		Button: gtk.Button{
			Bin: gtk.Bin{
				Container: gtk.Container{
					Widget: gtk.Widget{
						InitiallyUnowned: glib.InitiallyUnowned{
							Object: obj,
						},
					},
				},
			},
		},
		LanguageChooserer: &LanguageChooser{obj},
		Actionable:        gtk.Actionable{obj},
	}
}

func marshalLanguageChooserButton(p uintptr) (interface{}, error) {
	return wrapLanguageChooserButton(unsafe.Pointer(C.g_value_get_object((*C.GValue)(unsafe.Pointer(p))))), nil
}

// LanguageChooserButtonNew creates a new LanguageChooserButton.
func LanguageChooserButtonNew(currentLanguage *Language) *LanguageChooserButton {
	v1 := (*C.GspellLanguage)(unsafe.Pointer(currentLanguage.Native()))
	return wrapLanguageChooserButton(unsafe.Pointer(C.gspell_language_chooser_button_new(v1)))
}

// native turns the current *LanguageChooserButton into the native C pointer
// type.
func (l *LanguageChooserButton) native() *C.GspellLanguageChooserButton {
	return (*C.GspellLanguageChooserButton)(gwidget(&l.Button))
}

type LanguageChooserDialog struct {
	gtk.Dialog

	// Interfaces
	LanguageChooserer
}

// wrapLanguageChooserDialog wraps the given pointer to *LanguageChooserDialog.
func wrapLanguageChooserDialog(ptr unsafe.Pointer) *LanguageChooserDialog {
	obj := glib.Take(ptr)
	return &LanguageChooserDialog{
		Dialog: gtk.Dialog{
			Window: gtk.Window{
				Bin: gtk.Bin{
					Container: gtk.Container{
						Widget: gtk.Widget{
							InitiallyUnowned: glib.InitiallyUnowned{
								Object: obj,
							},
						},
					},
				},
			},
		},
		LanguageChooserer: &LanguageChooser{obj},
	}
}

func marshalLanguageChooserDialog(p uintptr) (interface{}, error) {
	return wrapLanguageChooserDialog(unsafe.Pointer(C.g_value_get_object((*C.GValue)(unsafe.Pointer(p))))), nil
}

// LanguageChooserDialogNew creates a new LanguageChooserDialog.
func LanguageChooserDialogNew(parent *gtk.Window, currentLanguage *Language, flags gtk.DialogFlags) *LanguageChooserDialog {
	v1 := (*C.GtkWindow)(unsafe.Pointer(parent.Widget.Native()))
	v2 := (*C.GspellLanguage)(unsafe.Pointer(currentLanguage.Native()))
	v3 := C.GtkDialogFlags(flags)

	return wrapLanguageChooserDialog(unsafe.Pointer(C.gspell_language_chooser_dialog_new(v1, v2, v3)))
}

// native turns the current *LanguageChooserDialog into the native C pointer
// type.
func (l *LanguageChooserDialog) native() *C.GspellLanguageChooserDialog {
	return (*C.GspellLanguageChooserDialog)(gwidget(&l.Dialog))
}

type NavigatorTextView struct {
	glib.InitiallyUnowned

	// Interfaces
	Navigatorer
}

// wrapNavigatorTextView wraps the given pointer to *NavigatorTextView.
func wrapNavigatorTextView(ptr unsafe.Pointer) *NavigatorTextView {
	obj := glib.Take(ptr)
	return &NavigatorTextView{
		InitiallyUnowned: glib.InitiallyUnowned{
			Object: obj,
		},
		Navigatorer: &Navigator{obj},
	}
}

func marshalNavigatorTextView(p uintptr) (interface{}, error) {
	return wrapNavigatorTextView(unsafe.Pointer(C.g_value_get_object((*C.GValue)(unsafe.Pointer(p))))), nil
}

// native turns the current *NavigatorTextView into the native C pointer type.
func (n *NavigatorTextView) native() *C.GspellNavigatorTextView {
	return (*C.GspellNavigatorTextView)(unsafe.Pointer(n.InitiallyUnowned.Native()))
}

func New(view *gtk.TextView) *Navigator {
	v1 := (*C.GtkTextView)(unsafe.Pointer(view.Widget.Native()))
	obj := glib.Take(unsafe.Pointer(C.gspell_navigator_text_view_new(v1)))
	r := &Navigator{obj}
	return r
}

func (n *NavigatorTextView) GetView() *gtk.TextView {
	obj := glib.Take(unsafe.Pointer(C.gspell_navigator_text_view_get_view(n.native())))
	r := &gtk.TextView{
		Container: gtk.Container{
			Widget: gtk.Widget{
				InitiallyUnowned: glib.InitiallyUnowned{
					Object: obj,
				},
			},
		},
	}
	return r
}

type TextBuffer struct {
	*glib.Object
}

// wrapTextBuffer wraps the given pointer to *TextBuffer.
func wrapTextBuffer(ptr unsafe.Pointer) *TextBuffer {
	obj := glib.Take(ptr)
	return &TextBuffer{
		Object: obj,
	}
}

func marshalTextBuffer(p uintptr) (interface{}, error) {
	return wrapTextBuffer(unsafe.Pointer(C.g_value_get_object((*C.GValue)(unsafe.Pointer(p))))), nil
}

// native turns the current *TextBuffer into the native C pointer type.
func (t *TextBuffer) native() *C.GspellTextBuffer {
	return (*C.GspellTextBuffer)(unsafe.Pointer(t.Object.Native()))
}

// GetFromGtkTextBuffer returns the TextBuffer of gtk_buffer. The returned
// object is guaranteed to be the same for the lifetime of gtk_buffer.
func GetFromGtkTextBuffer(gtkBuffer *gtk.TextBuffer) *TextBuffer {
	v1 := (*C.GtkTextBuffer)(unsafe.Pointer(gtkBuffer.Native()))
	r := wrapTextBuffer(unsafe.Pointer(C.gspell_text_buffer_get_from_gtk_text_buffer(v1)))
	return r
}

func (t *TextBuffer) GetBuffer() *gtk.TextBuffer {
	obj := glib.Take(unsafe.Pointer(C.gspell_text_buffer_get_buffer(t.native())))
	r := &gtk.TextBuffer{
		Object: obj,
	}
	return r
}
func (t *TextBuffer) GetSpellChecker() *Checker {
	r := wrapChecker(unsafe.Pointer(C.gspell_text_buffer_get_spell_checker(t.native())))
	return r
}

// SetSpellChecker sets a Checker to a TextBuffer. The gspell_buffer will own a
// reference to spell_checker, so you can release your reference to
// spell_checker if you no longer need it.
func (t *TextBuffer) SetSpellChecker(spellChecker *Checker) {
	v1 := (*C.GspellChecker)(unsafe.Pointer(spellChecker.Native()))
	C.gspell_text_buffer_set_spell_checker(t.native(), v1)
}

type TextView struct {
	*glib.Object
}

// wrapTextView wraps the given pointer to *TextView.
func wrapTextView(ptr unsafe.Pointer) *TextView {
	obj := glib.Take(ptr)
	return &TextView{
		Object: obj,
	}
}

func marshalTextView(p uintptr) (interface{}, error) {
	return wrapTextView(unsafe.Pointer(C.g_value_get_object((*C.GValue)(unsafe.Pointer(p))))), nil
}

// native turns the current *TextView into the native C pointer type.
func (t *TextView) native() *C.GspellTextView {
	return (*C.GspellTextView)(unsafe.Pointer(t.Object.Native()))
}

// GetFromGtkTextView returns the TextView of gtk_view. The returned object is
// guaranteed to be the same for the lifetime of gtk_view.
func GetFromGtkTextView(gtkView *gtk.TextView) *TextView {
	v1 := (*C.GtkTextView)(unsafe.Pointer(gtkView.Widget.Native()))
	r := wrapTextView(unsafe.Pointer(C.gspell_text_view_get_from_gtk_text_view(v1)))
	return r
}

// BasicSetup function is a convenience function that does the following: - Set
// a spell checker. The language chosen is the one returned by
// LanguageGetDefault(). - Set the TextView:inline-spell-checking property to
// true. - Set the TextView:enable-language-menu property to true.
//
// Example: |[ GtkTextView *gtk_view; GspellTextView *gspell_view;
//
//    gspell_view = gspell_text_view_get_from_gtk_text_view (gtk_view);
//    gspell_text_view_basic_setup (gspell_view);
//
//
//
//    GtkTextView *gtk_view;
//    GspellTextView *gspell_view;
//    GspellChecker *checker;
//    GtkTextBuffer *gtk_buffer;
//    GspellTextBuffer *gspell_buffer;
//
//    checker = gspell_checker_new (NULL);
//    gtk_buffer = gtk_text_view_get_buffer (gtk_view);
//    gspell_buffer = gspell_text_buffer_get_from_gtk_text_buffer (gtk_buffer);
//    gspell_text_buffer_set_spell_checker (gspell_buffer, checker);
//    g_object_unref (checker);
//
//    gspell_view = gspell_text_view_get_from_gtk_text_view (gtk_view);
//    gspell_text_view_set_inline_spell_checking (gspell_view, TRUE);
//    gspell_text_view_set_enable_language_menu (gspell_view, TRUE);
//
func (t *TextView) BasicSetup() {
	C.gspell_text_view_basic_setup(t.native())
}
func (t *TextView) GetEnableLanguageMenu() bool {
	r := gobool(C.gspell_text_view_get_enable_language_menu(t.native()))
	return r
}
func (t *TextView) GetInlineSpellChecking() bool {
	r := gobool(C.gspell_text_view_get_inline_spell_checking(t.native()))
	return r
}
func (t *TextView) GetView() *gtk.TextView {
	obj := glib.Take(unsafe.Pointer(C.gspell_text_view_get_view(t.native())))
	r := &gtk.TextView{
		Container: gtk.Container{
			Widget: gtk.Widget{
				InitiallyUnowned: glib.InitiallyUnowned{
					Object: obj,
				},
			},
		},
	}
	return r
}

// SetEnableLanguageMenu sets whether to enable the language context menu. If
// enabled, doing a right click on the TextView will show a sub-menu to choose
// the language for the spell checking. If another language is chosen, it
// changes the Checker:language property of the TextBuffer:spell-checker of the
// TextView:buffer of the TextView:view.
func (t *TextView) SetEnableLanguageMenu(enableLanguageMenu bool) {
	v1 := cbool(enableLanguageMenu)
	C.gspell_text_view_set_enable_language_menu(t.native(), v1)
}

// SetInlineSpellChecking enables or disables the inline spell checking.
func (t *TextView) SetInlineSpellChecking(enable bool) {
	v1 := cbool(enable)
	C.gspell_text_view_set_inline_spell_checking(t.native(), v1)
}

type Language C.GspellLanguage

func marshalLanguage(p uintptr) (interface{}, error) {
	return (*C.GspellLanguage)(unsafe.Pointer(C.g_value_get_boxed((*C.GValue)(unsafe.Pointer(p))))), nil
}

// native turns the current *Language into the native C pointer type.
func (l *Language) native() *C.GspellLanguage {
	return (*C.GspellLanguage)(unsafe.Pointer(l))
}

func (l *Language) Native() uintptr {
	return uintptr(unsafe.Pointer(l.native()))
}

// Compare compares alphabetically two languages by their name, as returned by
// C.gspell_language_get_name().
func (l *Language) Compare(languageB *Language) int {
	v1 := (*C.GspellLanguage)(unsafe.Pointer(languageB.Native()))
	r := int(C.gspell_language_compare(l.native(), v1))
	return r
}

// Copy used by language bindings.
func (l *Language) Copy() *Language {
	r := (*Language)(C.gspell_language_copy(l.native()))
	return r
}

// Free used by language bindings.
func (l *Language) Free() {
	C.gspell_language_free(l.native())
}
func (l *Language) GetCode() string {
	r := C.GoString(C.gspell_language_get_code(l.native()))
	return r
}

// GetName returns the language name translated to the current locale. For
// example "French (Belgium)" is returned if the current locale is in English
// and the language code is fr_BE.
func (l *Language) GetName() string {
	r := C.GoString(C.gspell_language_get_name(l.native()))
	return r
}
