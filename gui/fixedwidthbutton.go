// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"github.com/g3n/engine/window"
)

/***************************************

 FixedWidthButton Panel
 +-------------------------------+
 |  Image/Icon      Label        |
 |  +----------+   +----------+  |
 |  |          |   |          |  |
 |  |          |   |          |  |
 |  +----------+   +----------+  |
 +-------------------------------+

****************************************/

// FixedWidthButton represents a button GUI element which does not resize
type FixedWidthButton struct {
	Panel                   // Embedded Panel
	Label     *Label        // Label panel
	image     *Image        // pointer to button image (may be nil)
	icon      *Label        // pointer to button icon (may be nil
	styles    *ButtonStyles // pointer to current button styles
	mouseOver bool          // true if mouse is over button
	pressed   bool          // true if button is pressed
}

// NewFixedWidthButton creates and returns a pointer to a new fixed-width button widget
// with the specified text for the button label.
func NewFixedWidthButton(text string, width, height float32) *FixedWidthButton {

	b := new(FixedWidthButton)
	b.styles = &StyleDefault().Button

	// Initializes the button panel
	b.Panel.Initialize(b, width, height)

	// Subscribe to panel events
	b.Subscribe(OnKeyDown, b.onKey)
	b.Subscribe(OnKeyUp, b.onKey)
	b.Subscribe(OnMouseUp, b.onMouse)
	b.Subscribe(OnMouseDown, b.onMouse)
	b.Subscribe(OnMouseUpOut, b.onMouse)
	b.Subscribe(OnCursor, b.onCursor)
	b.Subscribe(OnCursorEnter, b.onCursor)
	b.Subscribe(OnCursorLeave, b.onCursor)
	b.Subscribe(OnEnable, func(name string, ev interface{}) { b.update() })
	b.Subscribe(OnResize, func(name string, ev interface{}) { b.recalc() })

	// Creates label
	b.Label = NewLabel(text)
	b.Label.Subscribe(OnResize, func(name string, ev interface{}) { b.recalc() })
	b.Panel.Add(b.Label)

	b.recalc() // recalc first then update!
	b.update()
	return b
}

// SetIcon sets the button icon from the default Icon font.
// If there is currently a selected image, it is removed
func (b *FixedWidthButton) SetIcon(icode string) {

	ico := NewIcon(icode)
	if b.image != nil {
		b.Panel.Remove(b.image)
		b.image = nil
	}
	if b.icon != nil {
		b.Panel.Remove(b.icon)
	}
	b.icon = ico
	b.icon.SetFontSize(b.Label.FontSize() * 1.4)
	b.Panel.Add(b.icon)

	b.recalc()
	b.update()
}

// SetImage sets the button left image from the specified filename
// If there is currently a selected icon, it is removed
func (b *FixedWidthButton) SetImage(imgfile string) error {

	img, err := NewImage(imgfile)
	if err != nil {
		return err
	}
	if b.image != nil {
		b.Panel.Remove(b.image)
	}
	b.image = img
	b.Panel.Add(b.image)
	b.recalc()
	return nil
}

// SetStyles set the button styles overriding the default style
func (b *FixedWidthButton) SetStyles(bs *ButtonStyles) {

	b.styles = bs
	b.update()
}

// onCursor process subscribed cursor events
func (b *FixedWidthButton) onCursor(evname string, ev interface{}) {

	switch evname {
	case OnCursorEnter:
		b.mouseOver = true
		b.update()
	case OnCursorLeave:
		b.mouseOver = false
		b.update()
	}
}

// onMouseEvent process subscribed mouse events
func (b *FixedWidthButton) onMouse(evname string, ev interface{}) {

	if !b.Enabled() {
		return
	}

	switch evname {
	case OnMouseDown:
		Manager().SetKeyFocus(b)
		b.pressed = true
		b.update()
	case OnMouseUpOut:
		fallthrough
	case OnMouseUp:
		if b.pressed && b.mouseOver {
			b.Dispatch(OnClick, nil)
		}
		b.pressed = false
		b.update()
	default:
		return
	}
}

// onKey processes subscribed key events
func (b *FixedWidthButton) onKey(evname string, ev interface{}) {

	kev := ev.(*window.KeyEvent)
	if kev.Key != window.KeyEnter {
		return
	}
	switch evname {
	case OnKeyDown:
		b.pressed = true
		b.update()
		b.Dispatch(OnClick, nil)
	case OnKeyUp:
		b.pressed = false
		b.update()
	}
}

// update updates the button visual state
func (b *FixedWidthButton) update() {

	if !b.Enabled() {
		b.applyStyle(&b.styles.Disabled)
		return
	}
	if b.pressed && b.mouseOver {
		b.applyStyle(&b.styles.Pressed)
		return
	}
	if b.mouseOver {
		b.applyStyle(&b.styles.Over)
		return
	}
	b.applyStyle(&b.styles.Normal)
}

// applyStyle applies the specified button style
func (b *FixedWidthButton) applyStyle(bs *ButtonStyle) {

	b.Panel.ApplyStyle(&bs.PanelStyle)
	if b.icon != nil {
		b.icon.SetColor4(&bs.FgColor)
	}
	b.Label.SetColor4(&bs.FgColor)
}

// recalc recalculates all dimensions and position from inside out
func (b *FixedWidthButton) recalc() {

	// Image or icon width
	imgWidth := float32(0)
	spacing := float32(4)
	if b.image != nil {
		imgWidth = b.image.Width()
	} else if b.icon != nil {
		imgWidth = b.icon.Width()
	}
	if imgWidth == 0 {
		spacing = 0
	}

	// Centralize horizontally
	px := (b.width - b.Label.width) / 2.0

	// Set label position
	ly := (b.height - b.Label.Height()) / 2
	b.Label.SetPosition(px+imgWidth+spacing, ly)

	// Image/icon position
	if b.image != nil {
		iy := (b.height - b.image.height) / 2
		b.image.SetPosition(px, iy)
	} else if b.icon != nil {
		b.icon.SetPosition(px, ly)
	}
}
