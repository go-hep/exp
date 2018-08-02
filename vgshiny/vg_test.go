// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vgshiny

import (
	"image"
	"image/color"
	"image/draw"
	"testing"

	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/paint"
)

func TestCanvas(t *testing.T) {
	c, err := New(fakeScreen{}, 640, 480)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Release()

	done := make(chan int)
	go func() {
		c.Run(nil)
		done <- 1
	}()

	c.Send(paint.Event{})
	c.Send(key.Event{Code: key.CodeEscape, Direction: key.DirPress})
	<-done
}

type fakeScreen struct{}

func (fakeScreen) NewBuffer(size image.Point) (screen.Buffer, error) {
	return &fakeBuffer{new(image.RGBA)}, nil
}

func (fakeScreen) NewTexture(size image.Point) (screen.Texture, error) {
	return nil, nil
}

func (fakeScreen) NewWindow(opts *screen.NewWindowOptions) (screen.Window, error) {
	return newWindow(), nil
}

type fakeBuffer struct {
	img *image.RGBA
}

func newBuffer(size image.Point) *fakeBuffer {
	fb := fakeBuffer{
		img: image.NewRGBA(image.Rectangle{
			Max: size,
		}),
	}
	return &fb
}

// Release releases the Buffer's resources, after all pending uploads and
// draws resolve.
//
// The behavior of the Buffer after Release, whether calling its methods or
// passing it as an argument, is undefined.
func (*fakeBuffer) Release() {}

// Size returns the size of the Buffer's image.
func (fb *fakeBuffer) Size() image.Point {
	return fb.img.Rect.Max
}

// Bounds returns the bounds of the Buffer's image. It is equal to
// image.Rectangle{Max: b.Size()}.
func (fb *fakeBuffer) Bounds() image.Rectangle {
	return fb.img.Rect.Bounds()
}

// RGBA returns the pixel buffer as an *image.RGBA.
//
// Its contents should not be accessed while the Buffer is uploading.
//
// The contents of the returned *image.RGBA's Pix field (of type []byte)
// can be modified at other times, but that Pix slice itself (i.e. its
// underlying pointer, length and capacity) should not be modified at any
// time.
//
// The following is valid:
//	m := buffer.RGBA()
//	if len(m.Pix) >= 4 {
//		m.Pix[0] = 0xff
//		m.Pix[1] = 0x00
//		m.Pix[2] = 0x00
//		m.Pix[3] = 0xff
//	}
// or, equivalently:
//	m := buffer.RGBA()
//	m.SetRGBA(m.Rect.Min.X, m.Rect.Min.Y, color.RGBA{0xff, 0x00, 0x00, 0xff})
// and using the standard library's image/draw package is also valid:
//	dst := buffer.RGBA()
//	draw.Draw(dst, dst.Bounds(), etc)
// but the following is invalid:
//	m := buffer.RGBA()
//	m.Pix = anotherByteSlice
// and so is this:
//	*buffer.RGBA() = anotherImageRGBA
func (fb *fakeBuffer) RGBA() *image.RGBA {
	return fb.img
}

type fakeWindow struct {
	screen.Window

	evt chan interface{}
}

func newWindow() *fakeWindow {
	return &fakeWindow{
		evt: make(chan interface{}),
	}
}

func (win *fakeWindow) Release()               { close(win.evt) }
func (win *fakeWindow) Send(evt interface{})   { win.evt <- evt }
func (win *fakeWindow) NextEvent() interface{} { return <-win.evt }

func (win *fakeWindow) Fill(dr image.Rectangle, src color.Color, op draw.Op) {}

func (win *fakeWindow) Upload(dp image.Point, src screen.Buffer, sr image.Rectangle) {}

// Publish flushes any pending Upload and Draw calls to the window, and
// swaps the back buffer to the front.
func (win *fakeWindow) Publish() screen.PublishResult {
	return screen.PublishResult{}
}

var (
	_ screen.Screen = (*fakeScreen)(nil)
	_ screen.Buffer = (*fakeBuffer)(nil)
	_ screen.Window = (*fakeWindow)(nil)
)
