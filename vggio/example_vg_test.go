// Copyright 2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vggio_test

import (
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/unit"
	"go-hep.org/x/exp/vggio"
	"go-hep.org/x/hep/hplot"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

func ExampleCanvas() {
	const (
		w = 20 * vg.Centimeter
		h = 15 * vg.Centimeter
	)
	go func(w, h vg.Length) {
		win := app.NewWindow(app.Title("vg-gio"),
			app.Size(unit.Px(float32(w.Points())), unit.Px(float32(h.Points()))),
		)
		done := time.NewTimer(2 * time.Second)
		defer done.Stop()
		for {
			select {
			case e := <-win.Events():
				switch e := e.(type) {
				case system.FrameEvent:
					p := hplot.New()
					p.Title.Text = "My title"
					p.X.Label.Text = "X"
					p.Y.Label.Text = "Y"

					cnv := vggio.New(e, w, h)
					p.Draw(draw.New(cnv))
					cnv.Paint(e)

				case key.Event:
					switch e.Name {
					case "Q", key.NameEscape:
						os.Exit(0)
					}
				}
			case <-done.C:
				os.Exit(0)
			}
		}
	}(w, h)

	app.Main()

	// Output:
}
