package ctk

import (
	"regexp"
	"strings"

	"github.com/kckrinke/go-cdk"
	"github.com/kckrinke/go-cdk/utils"
)

// CDK type-tag for Label objects
const TypeLabel cdk.CTypeTag = "ctk-label"

func init() {
	_ = cdk.TypesManager.AddType(TypeLabel, func() interface{} { return MakeLabel() })
	ctkBuilderTranslators[TypeLabel] = func(builder Builder, widget Widget, name, value string) error {
		switch strings.ToLower(name) {
		case "wrap":
			isTrue := utils.IsTrue(value)
			if err := widget.SetBoolProperty(PropertyWrap, isTrue); err != nil {
				return err
			}
			if isTrue {
				if wmi, err := widget.GetStructProperty(PropertyWrapMode); err == nil {
					if wm, ok := wmi.(cdk.WrapMode); ok {
						if wm == cdk.WRAP_NONE {
							if err := widget.SetStructProperty(PropertyWrapMode, cdk.WRAP_WORD); err != nil {
								widget.LogErr(err)
							}
						}
					}
				}
			}
			return nil
		}
		return ErrFallthrough
	}
}

// Label Hierarchy:
//	Object
//	  +- Widget
//	    +- Misc
//	      +- Label
//	        +- AccelLabel
//	        +- TipsQuery
type Label interface {
	Misc
	Alignable
	Buildable

	Init() (already bool)
	Build(builder Builder, element *CBuilderElement) error
	SetText(text string)
	SetAttributes(attrs cdk.Style)
	SetMarkup(text string) (parseError error)
	SetMarkupWithMnemonic(str string) (err error)
	SetJustify(justify cdk.Justification)
	SetEllipsize(mode bool)
	SetWidthChars(nChars int)
	SetMaxWidthChars(nChars int)
	SetLineWrap(wrap bool)
	SetLineWrapMode(wrapMode cdk.WrapMode)
	GetMnemonicKeyVal() (value rune)
	GetSelectable() (value bool)
	GetText() (value string)
	SelectRegion(startOffset int, endOffset int)
	SetMnemonicWidget(widget Widget)
	SetSelectable(setting bool)
	SetTextWithMnemonic(str string)
	GetAttributes() (value cdk.Style)
	GetJustify() (value cdk.Justification)
	GetEllipsize() (value bool)
	GetWidthChars() (value int)
	GetMaxWidthChars() (value int)
	GetLabel() (value string)
	GetLineWrap() (value bool)
	GetLineWrapMode() (value cdk.WrapMode)
	GetMnemonicWidget() (value Widget)
	GetSelectionBounds() (start int, end int, nonEmpty bool)
	GetUseMarkup() (value bool)
	GetUseUnderline() (value bool)
	GetSingleLineMode() (value bool)
	SetLabel(str string)
	SetUseMarkup(setting bool)
	SetUseUnderline(setting bool)
	SetSingleLineMode(singleLineMode bool)
	GetCurrentUri() (value string)
	SetTrackVisitedLinks(trackLinks bool)
	GetTrackVisitedLinks() (value bool)
	SetTheme(theme cdk.Theme)
	GetClearText() (text string)
	GetPlainText() (text string)
	GetCleanText() (text string)
	GetPlainTextInfo() (maxWidth, lineCount int)
	GetPlainTextInfoAtWidth(width int) (maxWidth, lineCount int)
	GetSizeRequest() (width, height int)
	Resize() cdk.EventFlag
	Draw(canvas cdk.Canvas) cdk.EventFlag
	Invalidate() cdk.EventFlag
}

// The CLabel structure implements the Label interface and is
// exported to facilitate type embedding with custom implementations. No member
// variables are exported as the interface methods are the only intended means
// of interacting with Label objects
type CLabel struct {
	CMisc

	text    string
	tbuffer cdk.TextBuffer
	tbStyle cdk.Style
	canvas  *cdk.CCanvas
}

// Default constructor for Label objects
func MakeLabel() *CLabel {
	return NewLabel("")
}

func NewLabel(plain string) *CLabel {
	l := new(CLabel)
	l.Init()
	l.SetText(plain)
	return l
}

// Creates a new Label, containing the text in str . If characters in str
// are preceded by an underscore, they are underlined. If you need a literal
// underscore character in a label, use '__' (two underscores). The first
// underlined character represents a keyboard accelerator called a mnemonic.
// The mnemonic key can be used to activate another widget, chosen
// automatically, or explicitly using SetMnemonicWidget. If
// SetMnemonicWidget is not called, then the first activatable
// ancestor of the Label will be chosen as the mnemonic widget. For
// instance, if the label is inside a button or menu item, the button or menu
// item will automatically become the mnemonic widget and be activated by the
// mnemonic.
// Parameters:
// 	str	The text of the label, with an underscore in front of the
// mnemonic character
// Returns:
// 	the new Label
func NewLabelWithMnemonic(str string) (value *CLabel) {
	l := new(CLabel)
	l.Init()
	l.SetTextWithMnemonic(str)
	return l
}

func NewLabelWithMarkup(markup string) (label *CLabel, err error) {
	label = new(CLabel)
	label.Init()
	err = label.SetMarkup(markup)
	return
}

// Label object initialization. This must be called at least once to setup
// the necessary defaults and allocate any memory structures. Calling this more
// than once is safe though unnecessary. Only the first call will result in any
// effect upon the Label instance
func (l *CLabel) Init() (already bool) {
	if l.InitTypeItem(TypeLabel, l) {
		return true
	}
	l.CMisc.Init()
	l.flags = NULL_WIDGET_FLAG
	l.SetFlags(PARENT_SENSITIVE)
	l.SetFlags(APP_PAINTABLE)
	_ = l.InstallProperty(PropertyAttributes, cdk.StructProperty, true, nil)
	_ = l.InstallProperty(PropertyCursorPosition, cdk.IntProperty, false, 0)
	_ = l.InstallProperty(PropertyEllipsize, cdk.BoolProperty, true, false)
	_ = l.InstallProperty(PropertyJustify, cdk.StructProperty, true, cdk.JUSTIFY_LEFT)
	_ = l.InstallProperty(PropertyLabel, cdk.StringProperty, true, "")
	_ = l.InstallProperty(PropertyMaxWidthChars, cdk.IntProperty, true, -1)
	_ = l.InstallProperty(PropertyMnemonicKeyVal, cdk.IntProperty, false, rune(0))
	_ = l.InstallProperty(PropertyMnemonicWidget, cdk.StructProperty, true, nil)
	_ = l.InstallProperty(PropertySelectable, cdk.BoolProperty, true, false)
	_ = l.InstallProperty(PropertySelectionBound, cdk.IntProperty, false, 0)
	_ = l.InstallProperty(PropertySingleLineMode, cdk.BoolProperty, true, false)
	_ = l.InstallProperty(PropertyTrackVisitedLinks, cdk.BoolProperty, true, true)
	_ = l.InstallProperty(PropertyUseMarkup, cdk.BoolProperty, true, false)
	_ = l.InstallProperty(PropertyUseUnderline, cdk.BoolProperty, true, false)
	_ = l.InstallProperty(PropertyWidthChars, cdk.IntProperty, true, -1)
	_ = l.InstallProperty(PropertyWrap, cdk.BoolProperty, true, false)
	_ = l.InstallProperty(PropertyWrapMode, cdk.StructProperty, true, cdk.WRAP_WORD)
	l.text = ""
	l.tbuffer = nil
	// _ = l.SetBoolProperty(PropertyDebug, true)
	l.canvas = cdk.NewCanvas(cdk.Point2I{}, cdk.Rectangle{}, l.GetTheme().Content.Normal)
	l.Invalidate()
	return false
}

func (l *CLabel) Build(builder Builder, element *CBuilderElement) error {
	l.Freeze()
	defer l.Thaw()
	if name, ok := element.Attributes["id"]; ok {
		l.SetName(name)
	}
	for k, v := range element.Properties {
		switch cdk.Property(k) {
		case PropertyLabel:
			l.SetLabel(v)
		default:
			element.ApplyProperty(k, v)
		}
	}
	element.ApplySignals()
	return nil
}

// Sets the text within the Label widget. It overwrites any text that was
// there before. This will also clear any previously set mnemonic
// accelerators.
// Parameters:
// 	str	The text you want to set
func (l *CLabel) SetText(text string) {
	l.Lock()
	defer l.Unlock()
	l.SetUseMarkup(false)
	l.text = text
	l.tbuffer = cdk.NewTextBuffer(text, l.GetTheme().Content.Normal, l.GetUseUnderline())
	l.Invalidate()
}

// Sets a PangoAttrList; the attributes in the list are applied to the label
// text.
// Parameters:
// 	attrs	a PangoAttrList
func (l *CLabel) SetAttributes(attrs cdk.Style) {
	if err := l.SetStructProperty(PropertyAttributes, attrs); err != nil {
		l.LogErr(err)
	}
}

// Parses str which is marked up with the Pango text markup language, setting
// the label's text and attribute list based on the parse results. If the str
// is external data, you may need to escape it with g_markup_escape_text or
// g_markup_printf_escaped:
// Parameters:
// 	str	a markup string (see Pango markup format)
func (l *CLabel) SetMarkup(text string) (parseError error) {
	l.Lock()
	defer l.Unlock()
	var m cdk.Tango
	if m, parseError = cdk.NewMarkup(text, l.GetTheme().Content.Normal); parseError != nil {
		return parseError
	}
	l.SetUseMarkup(true)
	l.text = text
	l.tbuffer = m.TextBuffer(l.GetUseUnderline())
	l.Invalidate()
	return nil
}

// Parses str which is marked up with the Pango text markup language, setting
// the label's text and attribute list based on the parse results. If
// characters in str are preceded by an underscore, they are underlined
// indicating that they represent a keyboard accelerator called a mnemonic.
// The mnemonic key can be used to activate another widget, chosen
// automatically, or explicitly using SetMnemonicWidget.
// Parameters:
// 	str	a markup string (see Pango markup format)
func (l *CLabel) SetMarkupWithMnemonic(str string) (err error) {
	l.SetUseUnderline(true)
	err = l.SetMarkup(str)
	return
}

// Sets the alignment of the lines in the text of the label relative to each
// other. GTK_JUSTIFY_LEFT is the default value when the widget is first
// created with New. If you instead want to set the alignment of
// the label as a whole, use MiscSetAlignment instead.
// SetJustify has no effect on labels containing only a single
// line.
// Parameters:
// 	jtype	a Justification
func (l *CLabel) SetJustify(justify cdk.Justification) {
	if err := l.SetStructProperty(PropertyJustify, justify); err != nil {
		l.LogErr(err)
	}
}

// Sets the mode used to ellipsize (add an ellipsis: "...") to the text if
// there is not enough space to render the entire string.
// Parameters:
// 	mode	bool
func (l *CLabel) SetEllipsize(mode bool) {
	if err := l.SetBoolProperty(PropertyEllipsize, mode); err != nil {
		l.LogErr(err)
	}
}

// Sets the desired width in characters of label to n_chars .
// Parameters:
// 	nChars	the new desired width, in characters.
func (l *CLabel) SetWidthChars(nChars int) {
	if err := l.SetIntProperty(PropertyWidthChars, nChars); err != nil {
		l.LogErr(err)
	}
}

// Sets the desired maximum width in characters of label to n_chars .
// Parameters:
// 	nChars	the new desired maximum width, in characters.
func (l *CLabel) SetMaxWidthChars(nChars int) {
	if err := l.SetIntProperty(PropertyMaxWidthChars, nChars); err != nil {
		l.LogErr(err)
	}
}

// Toggles line wrapping within the Label widget. TRUE makes it break
// lines if text exceeds the widget's size. FALSE lets the text get cut off
// by the edge of the widget if it exceeds the widget size. Note that setting
// line wrapping to TRUE does not make the label wrap at its parent
// container's width, because CTK widgets conceptually can't make their
// requisition depend on the parent container's size. For a label that wraps
// at a specific position, set the label's width using
// SetSizeRequest.
// Parameters:
// 	wrap	the setting
func (l *CLabel) SetLineWrap(wrap bool) {
	if err := l.SetBoolProperty(PropertyWrap, wrap); err != nil {
		l.LogErr(err)
	}
}

// If line wrapping is on (see SetLineWrap) this controls how
// the line wrapping is done. The default is PANGO_WRAP_WORD which means wrap
// on word boundaries.
// Parameters:
// 	wrapMode	the line wrapping mode
func (l *CLabel) SetLineWrapMode(wrapMode cdk.WrapMode) {
	if err := l.SetStructProperty(PropertyWrapMode, wrapMode); err != nil {
		l.LogErr(err)
	}
}

// If the label has been set so that it has an mnemonic key this function
// returns the keyval used for the mnemonic accelerator. If there is no
// mnemonic set up it returns GDK_VoidSymbol.
// Returns:
// 	GDK keyval usable for accelerators, or GDK_VoidSymbol
func (l *CLabel) GetMnemonicKeyVal() (value rune) {
	if l.GetUseUnderline() {
		label := l.GetClearText()
		if rxLabelMnemonic.MatchString(label) {
			m := rxLabelMnemonic.FindStringSubmatch(label)
			if len(m) > 1 {
				return rune(strings.ToLower(m[1])[0])
			}
		}
	}
	return
}

// Gets the value set by SetSelectable.
// Returns:
// 	TRUE if the user can copy text from the label
func (l *CLabel) GetSelectable() (value bool) {
	var err error
	if value, err = l.GetBoolProperty(PropertySelectable); err != nil {
		l.LogErr(err)
	}
	return
}

// Fetches the text from a label widget, as displayed on the screen. This
// does not include any embedded underlines indicating mnemonics or Tango
// markup. (See GetLabel)
// Returns:
// 	the text in the label widget.
func (l *CLabel) GetText() (value string) {
	value = l.GetCleanText()
	return
}

// Selects a range of characters in the label, if the label is selectable.
// See SetSelectable. If the label is not selectable, this
// function has no effect. If start_offset or end_offset are -1, then the end
// of the label will be substituted.
// Parameters:
// 	startOffset	start offset (in characters not bytes)
// 	endOffset	end offset (in characters not bytes)
func (l *CLabel) SelectRegion(startOffset int, endOffset int) {}

// If the label has been set so that it has an mnemonic key (using i.e.
// SetMarkupWithMnemonic, SetTextWithMnemonic,
// NewWithMnemonic or the "use_underline" property) the label
// can be associated with a widget that is the target of the mnemonic. When
// the label is inside a widget (like a Button or a Notebook tab) it is
// automatically associated with the correct widget, but sometimes (i.e. when
// the target is a Entry next to the label) you need to set it explicitly
// using this function. The target widget will be accelerated by emitting the
// Widget::mnemonic-activate signal on it. The default handler for this
// signal will activate the widget if there are no mnemonic collisions and
// toggle focus between the colliding widgets otherwise.
// Parameters:
// 	widget	the target Widget.
func (l *CLabel) SetMnemonicWidget(widget Widget) {
	if err := l.SetStructProperty(PropertyMnemonicWidget, widget); err != nil {
		l.LogErr(err)
	} else {
		l.Invalidate()
	}
}

// Selectable labels allow the user to select text from the label, for
// copy-and-paste.
// Parameters:
// 	setting	TRUE to allow selecting text in the label
func (l *CLabel) SetSelectable(setting bool) {
	if err := l.SetBoolProperty(PropertySelectable, setting); err != nil {
		l.LogErr(err)
	}
}

// Sets the label's text from the string str . If characters in str are
// preceded by an underscore, they are underlined indicating that they
// represent a keyboard accelerator called a mnemonic. The mnemonic key can
// be used to activate another widget, chosen automatically, or explicitly
// using SetMnemonicWidget.
// Parameters:
// 	str	a string
func (l *CLabel) SetTextWithMnemonic(str string) {
	l.SetUseUnderline(true)
	l.SetText(str)
}

// Gets the attribute list that was set on the label using
// SetAttributes, if any. This function does not reflect
// attributes that come from the labels markup (see SetMarkup).
// If you want to get the effective attributes for the label, use
// pango_layout_get_attribute (GetLayout (label)).
// Returns:
// 	the attribute list, or NULL if none was set.
// 	[transfer none]
func (l *CLabel) GetAttributes() (value cdk.Style) {
	var ok bool
	if v, err := l.GetStructProperty(PropertyAttributes); err != nil {
		l.LogErr(err)
	} else if value, ok = v.(cdk.Style); !ok {
		l.LogError("value stored in PropertyAttributes is not of cdk.Style type: %v (%T)", v, v)
	}
	return
}

// Returns the justification of the label. See SetJustify.
// Returns:
// 	Justification
func (l *CLabel) GetJustify() (value cdk.Justification) {
	var ok bool
	if v, err := l.GetStructProperty(PropertyJustify); err != nil {
		l.LogErr(err)
	} else if value, ok = v.(cdk.Justification); !ok {
		l.LogError("value stored in PropertyJustify is not of cdk.Justification type: %v (%T)", v, v)
	}
	return
}

// Returns the ellipsizing position of the label. See
// SetEllipsize.
// Returns:
// 	PangoEllipsizeMode
func (l *CLabel) GetEllipsize() (value bool) {
	var err error
	if value, err = l.GetBoolProperty(PropertyEllipsize); err != nil {
		l.LogErr(err)
	}
	return
}

// Retrieves the desired width of label , in characters. See
// SetWidthChars.
// Returns:
// 	the width of the label in characters.
func (l *CLabel) GetWidthChars() (value int) {
	var err error
	if value, err = l.GetIntProperty(PropertyWidthChars); err != nil {
		l.LogErr(err)
	}
	return
}

// Retrieves the desired maximum width of label , in characters. See
// SetWidthChars.
// Returns:
// 	the maximum width of the label in characters.
func (l *CLabel) GetMaxWidthChars() (value int) {
	var err error
	if value, err = l.GetIntProperty(PropertyMaxWidthChars); err != nil {
		l.LogErr(err)
	}
	return
}

// Fetches the text from a label widget including any embedded underlines
// indicating mnemonics and Pango markup. (See GetText).
// Returns:
// 	the text of the label widget. This string is owned by the
// 	widget and must not be modified or freed.
func (l *CLabel) GetLabel() (value string) {
	var err error
	if value, err = l.GetStringProperty(PropertyLabel); err != nil {
		l.LogErr(err)
	}
	return
}

// Returns whether lines in the label are automatically wrapped. See
// SetLineWrap.
// Returns:
// 	TRUE if the lines of the label are automatically wrapped.
func (l *CLabel) GetLineWrap() (value bool) {
	var err error
	if value, err = l.GetBoolProperty(PropertyWrap); err != nil {
		l.LogErr(err)
	}
	return
}

// Returns line wrap mode used by the label. See SetLineWrapMode.
// Returns:
// 	the current cdk.WrapMode
func (l *CLabel) GetLineWrapMode() (value cdk.WrapMode) {
	// if !l.GetLineWrap() {
	// 	return cdk.WRAP_NONE
	// }
	var ok bool
	if v, err := l.GetStructProperty(PropertyWrapMode); err != nil {
		l.LogErr(err)
	} else if value, ok = v.(cdk.WrapMode); !ok {
		l.LogError("value stored in PropertyWrap is not of cdk.WrapMode type: %v (%T)", v, v)
	}
	return
}

// Retrieves the target of the mnemonic (keyboard shortcut) of this label.
// See SetMnemonicWidget.
// Returns:
// 	the target of the label's mnemonic, or NULL if none has been
// 	set and the default algorithm will be used.
// 	[transfer none]
func (l *CLabel) GetMnemonicWidget() (value Widget) {
	if v, err := l.GetStructProperty(PropertyMnemonicWidget); err == nil {
		value, _ = v.(Widget)
	} else {
		l.LogErr(err)
	}

	return
}

// Gets the selected range of characters in the label, returning TRUE if
// there's a selection.
// Parameters:
// 	start	return location for start of selection, as a character offset.
// 	end	return location for end of selection, as a character offset.
// Returns:
// 	TRUE if selection is non-empty
func (l *CLabel) GetSelectionBounds() (start int, end int, nonEmpty bool) {
	return 0, 0, false
}

// Returns whether the label's text is interpreted as marked up with the
// Pango text markup language. See SetUseMarkup.
// Returns:
// 	TRUE if the label's text will be parsed for markup.
func (l *CLabel) GetUseMarkup() (value bool) {
	var err error
	if value, err = l.GetBoolProperty(PropertyUseMarkup); err != nil {
		l.LogErr(err)
	}
	return
}

// Returns whether an embedded underline in the label indicates a mnemonic.
// See SetUseUnderline.
// Returns:
// 	TRUE whether an embedded underline in the label indicates the
// 	mnemonic accelerator keys.
func (l *CLabel) GetUseUnderline() (value bool) {
	var err error
	if value, err = l.GetBoolProperty(PropertyUseUnderline); err != nil {
		l.LogErr(err)
	}
	return
}

// Returns whether the label is in single line mode.
// Returns:
// 	TRUE when the label is in single line mode.
func (l *CLabel) GetSingleLineMode() (value bool) {
	var err error
	if value, err = l.GetBoolProperty(PropertySingleLineMode); err != nil {
		l.LogErr(err)
	}
	return
}

// Sets the text of the label. The label is interpreted as including embedded
// underlines and/or Pango markup depending on the values of the
// “use-underline”" and “use-markup” properties.
// Parameters:
// 	str	the new text to set for the label
func (l *CLabel) SetLabel(str string) {
	if err := l.SetStringProperty(PropertyLabel, str); err != nil {
		l.LogErr(err)
	} else {
		if l.GetUseMarkup() {
			if err := l.SetMarkup(str); err != nil {
				l.LogErr(err)
			}
		} else {
			l.SetText(str)
		}
	}
}

// Sets whether the text of the label contains markup in Pango's text markup
// language. See SetMarkup.
// Parameters:
// 	setting	TRUE if the label's text should be parsed for markup.
func (l *CLabel) SetUseMarkup(setting bool) {
	if err := l.SetBoolProperty(PropertyUseMarkup, setting); err != nil {
		l.LogErr(err)
	} else {
		l.Invalidate()
	}
}

// If true, an underline in the text indicates the next character should be
// used for the mnemonic accelerator key.
// Parameters:
// 	setting	TRUE if underlines in the text indicate mnemonics
func (l *CLabel) SetUseUnderline(setting bool) {
	if err := l.SetBoolProperty(PropertyUseUnderline, setting); err != nil {
		l.LogErr(err)
	} else {
		l.Invalidate()
	}
}

// Sets whether the label is in single line mode.
// Parameters:
// 	singleLineMode	TRUE if the label should be in single line mode
func (l *CLabel) SetSingleLineMode(singleLineMode bool) {
	if err := l.SetBoolProperty(PropertySingleLineMode, singleLineMode); err != nil {
		l.LogErr(err)
	} else {
		l.Invalidate()
	}
}

// Returns the URI for the currently active link in the label. The active
// link is the one under the mouse pointer or, in a selectable label, the
// link in which the text cursor is currently positioned. This function is
// intended for use in a “activate-link” handler or for use in a
// “query-tooltip” handler.
// Returns:
// 	the currently active URI. The string is owned by CTK and must
// 	not be freed or modified.
func (l *CLabel) GetCurrentUri() (value string) {
	return ""
}

// Sets whether the label should keep track of clicked links (and use a
// different color for them).
// Parameters:
// 	trackLinks	TRUE to track visited links
func (l *CLabel) SetTrackVisitedLinks(trackLinks bool) {
	if err := l.SetBoolProperty(PropertyTrackVisitedLinks, trackLinks); err != nil {
		l.LogErr(err)
	}
}

// Returns whether the label is currently keeping track of clicked links.
// Returns:
// 	TRUE if clicked links are remembered
func (l *CLabel) GetTrackVisitedLinks() (value bool) {
	var err error
	if value, err = l.GetBoolProperty(PropertyTrackVisitedLinks); err != nil {
		l.LogErr(err)
	}
	return
}

// Set the Theme for the Widget instance. This will also refresh the requested
// theme. A request theme is a transient theme, based on the actually set theme
// and adjusted for focus. If the given theme is equivalent to the current theme
// then no action is taken. After verifying that the given theme is different,
// this method emits a set-theme signal and if the listeners return EVENT_PASS,
// the changes are applied and the Widget.Invalidate() method is called
func (l *CLabel) SetTheme(theme cdk.Theme) {
	if theme.String() != l.GetTheme().String() {
		if f := l.Emit(SignalSetTheme, l, theme); f == cdk.EVENT_PASS {
			l.CObject.SetTheme(theme)
			l.Invalidate()
		}
	}
}

var (
	rxLabelPlainText = regexp.MustCompile(`(?msi)(_)([a-z])`)
	rxLabelMnemonic  = regexp.MustCompile(`(?msi)_([a-z])`)
)

func (l *CLabel) GetClearText() (text string) {
	if l.tbuffer == nil {
		return ""
	}
	text = l.tbuffer.ClearText(l.GetLineWrapMode(), l.GetEllipsize(), l.GetJustify(), l.GetMaxWidthChars())
	if l.GetSingleLineMode() {
		if strings.Contains(text, "\n") {
			if idx := strings.Index(text, "\n"); idx >= 0 {
				text = text[:idx]
			}
		}
	}
	return
}

func (l *CLabel) GetPlainText() (text string) {
	if l.tbuffer == nil {
		return ""
	}
	text = l.tbuffer.PlainText(l.GetLineWrapMode(), l.GetEllipsize(), l.GetJustify(), l.GetMaxWidthChars())
	if l.GetSingleLineMode() {
		if strings.Contains(text, "\n") {
			if idx := strings.Index(text, "\n"); idx >= 0 {
				text = text[:idx]
			}
		}
	}
	return
}

func (l *CLabel) GetCleanText() (text string) {
	text = rxLabelPlainText.ReplaceAllString(l.GetPlainText(), "$2")
	return
}

func (l *CLabel) GetPlainTextInfo() (maxWidth, lineCount int) {
	if l.tbuffer == nil {
		return -1, -1
	}
	wrapMode := l.GetLineWrapMode()
	maxWidth, lineCount = l.tbuffer.PlainTextInfo(wrapMode, l.GetEllipsize(), l.GetJustify(), l.GetMaxWidthChars())
	return
}

func (l *CLabel) GetPlainTextInfoAtWidth(width int) (maxWidth, lineCount int) {
	if l.tbuffer == nil {
		return -1, -1
	}
	maxWidth, lineCount = l.tbuffer.PlainTextInfo(l.GetLineWrapMode(), l.GetEllipsize(), l.GetJustify(), width)
	return
}

// TODO: label size request handling is wonky
// TODO: scrolled viewport resize issues

func (l *CLabel) GetSizeRequest() (width, height int) {
	size := cdk.NewRectangle(l.CWidget.GetSizeRequest())
	if size.W <= -1 {
		if wc := l.GetWidthChars(); wc <= -1 {
			if mwc := l.GetMaxWidthChars(); mwc <= -1 {
				alloc := l.GetAllocation()
				if alloc.W > 0 {
					size.W, _ = l.GetPlainTextInfoAtWidth(alloc.W)
				} else {
					size.W, _ = l.GetPlainTextInfo()
				}
			} else {
				size.W = mwc
			}
		} else {
			size.W = wc
		}
	}
	if size.H <= -1 {
		_, size.H = l.GetPlainTextInfoAtWidth(size.W)
	}
	// add padding
	xPadding, yPadding := l.GetPadding()
	size.W += xPadding * 2
	size.H += yPadding * 2
	// min size of 3 according to GTK
	return size.W, size.H
}

func (l *CLabel) Resize() cdk.EventFlag {
	l.Lock()
	defer l.Unlock()
	alloc := l.GetAllocation()
	if !l.IsVisible() || alloc.W <= 0 || alloc.H <= 0 {
		l.LogTrace("Label.Resize(): not visible, zero width or zero height")
		return cdk.EVENT_PASS
	}
	var pos cdk.Point2I
	size := cdk.NewRectangle(alloc.W, alloc.H)

	req := l.CWidget.SizeRequest()
	if req.W <= -1 {
		size.W, size.H = l.GetPlainTextInfoAtWidth(alloc.W)
	} else {
		size.W, size.H = l.GetPlainTextInfoAtWidth(req.W)
	}
	if req.H > -1 {
		size.H = req.H
	}
	size.Clamp(0, 0, alloc.W, alloc.H)

	xAlign, yAlign := l.GetAlignment()

	if size.W < alloc.W {
		if size.W < alloc.W {
			delta := alloc.W - size.W
			pos.X += int(float64(delta) * xAlign)
		}
	}

	if size.H < alloc.H {
		if size.H < alloc.H {
			delta := alloc.H - size.H
			pos.Y += int(float64(delta) * yAlign)
		}
	}

	l.canvas.SetOrigin(pos)
	l.canvas.Resize(*size, l.getStyleRequest())
	l.Invalidate()
	l.Emit(SignalResize, l)
	return cdk.EVENT_STOP
}

func (l *CLabel) Draw(canvas cdk.Canvas) cdk.EventFlag {
	l.Lock()
	defer l.Unlock()
	alloc := l.GetAllocation()
	if !l.IsVisible() || alloc.W <= 0 || alloc.H <= 0 {
		l.LogTrace("Label.Draw(): not visible, zero width or zero height")
		return cdk.EVENT_PASS
	}

	if l.tbuffer != nil {
		// if l.GetTheme().String() != l.GetThemeRequest().String() {
		// 	l.Invalidate()
		// }
		l.tbuffer.Draw(l.canvas, l.GetSingleLineMode(), l.GetLineWrapMode(), l.GetEllipsize(), l.GetJustify(), cdk.ALIGN_TOP)
		if err := canvas.Composite(l.canvas); err != nil {
			l.LogError("composite error: %v", err)
		}
	}

	if debug, _ := l.GetBoolProperty(cdk.PropertyDebug); debug {
		canvas.DebugBox(cdk.ColorSilver, l.ObjectInfo())
	}
	return cdk.EVENT_STOP
}

func (l *CLabel) getMaxCharsRequest() (maxWidth int) {
	alloc := l.GetAllocation()
	maxWidth = l.GetMaxWidthChars()
	if maxWidth <= -1 {
		w, _ := l.GetSizeRequest()
		if w > -1 {
			maxWidth = w
		} else {
			maxWidth = alloc.W
		}
	}
	return
}

func (l *CLabel) Invalidate() cdk.EventFlag {
	style := l.getStyleRequest()
	_ = l.refreshBufferWithStyle(style)
	theme := l.GetThemeRequest()
	theme.Content.FillRune = rune(0)
	l.canvas.Fill(theme)
	l.refreshMnemonics()
	return cdk.EVENT_STOP
}

func (l *CLabel) getStyleRequest() (style cdk.Style) {
	style = l.GetThemeRequest().Content.Normal
	return
}

func (l *CLabel) refreshBufferWithStyle(style cdk.Style) error {
	if l.tbStyle.String() != style.String() {
		l.tbStyle = style
		if l.GetUseMarkup() {
			if m, err := cdk.NewMarkup(l.text, style); err != nil {
				return err
			} else {
				l.tbuffer = m.TextBuffer(l.GetUseUnderline())
			}
		} else if l.tbuffer != nil {
			l.tbuffer = cdk.NewTextBuffer(l.text, style, l.GetUseUnderline())
		}
	}
	return nil
}

func (l *CLabel) refreshMnemonics() {
	if w := l.GetWindow(); w != nil {
		if widget := l.GetMnemonicWidget(); widget != nil {
			w.RemoveWidgetMnemonics(widget)
		} else {
			if parent := l.GetParent(); parent != nil {
				w.RemoveWidgetMnemonics(parent)
			}
		}
	}
	if l.GetUseUnderline() {
		if w := l.GetWindow(); w != nil {
			if keyval := l.GetMnemonicKeyVal(); keyval > 0 {
				if widget := l.GetMnemonicWidget(); widget != nil {
					w.AddMnemonic(keyval, widget)
				} else {
					if parent := l.GetParent(); parent != nil {
						if pw, ok := parent.(Widget); ok && pw.IsSensitive() && pw.IsVisible() {
							w.AddMnemonic(keyval, parent)
						}
					}
				}
			}
		}
	}
}

// A list of style attributes to apply to the text of the label.
// Flags: Read / Write
const PropertyAttributes cdk.Property = "attributes"

// The current position of the insertion cursor in chars.
// Flags: Read
// Allowed values: >= 0
// Default value: 0
const PropertyCursorPosition cdk.Property = "cursor-position"

// The preferred place to ellipsize the string, if the label does not have
// enough room to display the entire string, specified as a bool.
// Flags: Read / Write
// Default value: false
const PropertyEllipsize cdk.Property = "ellipsize"

// The alignment of the lines in the text of the label relative to each
// other. This does NOT affect the alignment of the label within its
// allocation. See Misc::xAlign for that.
// Flags: Read / Write
// Default value: cdk.JUSTIFY_LEFT
const PropertyJustify cdk.Property = "justify"

// The text of the label.
// Flags: Read / Write
// Default value: ""
const PropertyLabel cdk.Property = "label"

// The desired maximum width of the label, in characters. If this property is
// set to -1, the width will be calculated automatically, otherwise the label
// will request space for no more than the requested number of characters. If
// the “width-chars” property is set to a positive value, then the
// "max-width-chars" property is ignored.
// Flags: Read / Write
// Allowed values: >= -1
// Default value: -1
const PropertyMaxWidthChars cdk.Property = "max-width-chars"

// The mnemonic accelerator key for this label.
// Flags: Read
// Default value: 16777215
const PropertyMnemonicKeyVal cdk.Property = "mnemonic-key-val"

// The widget to be activated when the label's mnemonic key is pressed.
// Flags: Read / Write
const PropertyMnemonicWidget cdk.Property = "mnemonic-widget"

// Whether the label text can be selected with the mouse.
// Flags: Read / Write
// Default value: FALSE
const PropertySelectable cdk.Property = "selectable"

// The position of the opposite end of the selection from the cursor in
// chars.
// Flags: Read
// Allowed values: >= 0
// Default value: 0
const PropertySelectionBound cdk.Property = "selection-bound"

// Whether the label is in single line mode. In single line mode, the height
// of the label does not depend on the actual text, it is always set to
// ascent + descent of the font. This can be an advantage in situations where
// resizing the label because of text changes would be distracting, e.g. in a
// statusbar.
// Flags: Read / Write
// Default value: FALSE
const PropertySingleLineMode cdk.Property = "single-line-mode"

// Set this property to TRUE to make the label track which links have been
// clicked. It will then apply the ::visited-link-color color, instead of
// ::link-color.
// Flags: Read / Write
// Default value: TRUE
const PropertyTrackVisitedLinks cdk.Property = "track-visited-links"

// The text of the label includes XML markup. See pango_parse_markup.
// Flags: Read / Write
// Default value: FALSE
const PropertyUseMarkup cdk.Property = "use-markup"

// If set, an underline in the text indicates the next character should be
// used for the mnemonic accelerator key.
// Flags: Read / Write
// Default value: FALSE
const PropertyLabelUseUnderline cdk.Property = "use-underline"

// The desired width of the label, in characters. If this property is set to
// -1, the width will be calculated automatically, otherwise the label will
// request either 3 characters or the property value, whichever is greater.
// If the "width-chars" property is set to a positive value, then the
// “max-width-chars” property is ignored.
// Flags: Read / Write
// Allowed values: >= -1
// Default value: -1
const PropertyWidthChars cdk.Property = "width-chars"

// If set, wrap lines if the text becomes too wide.
// Flags: Read / Write
// Default value: FALSE
const PropertyWrap cdk.Property = "wrap"

// If line wrapping is on (see the “wrap” property) this controls how the
// line wrapping is done. The default is PANGO_WRAP_WORD, which means wrap on
// word boundaries.
// Flags: Read / Write
// Default value: PANGO_WRAP_WORD
const PropertyWrapMode cdk.Property = "wrap-mode"

// A keybinding signal which gets emitted when the user activates a link in
// the label. Applications may also emit the signal with
// g_signal_emit_by_name if they need to control activation of URIs
// programmatically. The default bindings for this signal are all forms of
// the Enter key.
const SignalActivateCurrentLink cdk.Signal = "activate-current-link"

// The signal which gets emitted to activate a URI. Applications may connect
// to it to override the default behaviour, which is to call ShowUri.
const SignalActivateLink cdk.Signal = "activate-link"

// The ::copy-clipboard signal is a which gets emitted to copy the selection
// to the clipboard. The default binding for this signal is Ctrl-c.
const SignalCopyClipboard cdk.Signal = "copy-clipboard"

// The ::move-cursor signal is a which gets emitted when the user initiates a
// cursor movement. If the cursor is not visible in entry , this signal
// causes the viewport to be moved instead. Applications should not connect
// to it, but may emit it with g_signal_emit_by_name if they need to
// control the cursor programmatically. The default bindings for this signal
// come in two variants, the variant with the Shift modifier extends the
// selection, the variant without the Shift modifer does not. There are too
// many key combinations to list them all here.
// Listener function arguments:
// 	step MovementStep	the granularity of the move, as a GtkMovementStep
// 	count int	the number of step units to move
// 	extendSelection bool	TRUE if the move should extend the selection
const SignalMoveCursor cdk.Signal = "move-cursor"

// The ::populate-popup signal gets emitted before showing the context menu
// of the label. Note that only selectable labels have context menus. If you
// need to add items to the context menu, connect to this signal and append
// your menuitems to the menu .
// Listener function arguments:
// 	menu Menu	the menu that is being populated
const SignalPopulatePopup cdk.Signal = "populate-popup"
