// Copyright 2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package vggio provides a vg.Canvas implementation backed by Gioui.
package vggio // import "go-hep.org/x/exp/vggio"

import (
	"image/color"

	"gioui.org/f32"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"

	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/vgimg"
)

// Canvas implements the vg.Canvas interface,
// drawing to an image.Image using vgimg and painting that image
// into a Gioui context.
type Canvas struct {
	*vgimg.Canvas
	gtx layout.Context
}

func (c *Canvas) pt32(p vg.Point) f32.Point {
	_, h := c.Size()
	dpi := c.DPI()
	return f32.Point{
		X: float32(p.X.Dots(dpi)),
		Y: float32(h.Dots(dpi) - p.Y.Dots(dpi)),
	}
}

// New returns a new image canvas with the provided width and height.
func New(e system.FrameEvent, w, h vg.Length) *Canvas {
	c := &Canvas{
		gtx: layout.NewContext(new(op.Ops), e),
		Canvas: vgimg.NewWith(
			vgimg.UseDPI(vgimg.DefaultDPI),
			vgimg.UseWH(w, h),
			vgimg.UseBackgroundColor(color.White),
		),
	}
	return c
}

// Paint paints the canvas' content on the screen.
func (c *Canvas) Paint(e system.FrameEvent) {
	w, h := c.Size()
	box := vg.Rectangle{Max: vg.Point{X: w, Y: h}}
	img := c.Canvas.Image()
	ops := c.gtx.Ops
	min := c.pt32(box.Min)
	max := c.pt32(box.Max)
	r32 := f32.Rect(min.X, min.Y, max.X, max.Y)

	paint.NewImageOp(img).Add(ops)
	paint.PaintOp{Rect: r32}.Add(ops)

	e.Frame(ops)
}
