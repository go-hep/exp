// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package roots

import (
	"math"
	"testing"

	"gonum.org/v1/gonum/cmplxs"
)

func TestPolynomial(t *testing.T) {
	for _, test := range []struct {
		ps   []float64
		dst  []complex128
		want []complex128
	}{
		{
			ps:   []float64{2, 1},
			want: []complex128{complex(-0.5, 0)},
		},
		{
			ps:   []float64{1, 2},
			want: []complex128{complex(-2, 0)},
		},
		{
			ps:   []float64{1, 2},
			dst:  make([]complex128, 1),
			want: []complex128{complex(-2, 0)},
		},
		{
			ps:   []float64{2, 0},
			want: []complex128{0},
		},
		{
			ps:   []float64{0, 2},
			want: []complex128{},
		},
		// -- quadratic --
		{
			ps: []float64{3, 2, 1},
			want: []complex128{
				complex(-0.3333333333333333, -0.4714045207910317),
				complex(-0.3333333333333333, +0.4714045207910317),
			},
		},
		// -- cubic --
		{
			// p == q == 0
			ps: []float64{1, 1, 1 / 3., 1 / 27.},
			want: []complex128{
				complex(-1/3., 0),
				complex(-1/3., 0),
				complex(-1/3., 0),
			},
		},
		{
			// Δ == 0
			ps: []float64{1, 0, -1, 2 / 3. * math.Sqrt(1/3.)},
			want: []complex128{
				complex(-1.1547005383792515, 0),
				complex(0.5773502691896257, 0),
				complex(0.5773502691896257, 0),
			},
		},
		{
			// Δ > 0
			ps: []float64{4, 3, 2, 1},
			want: []complex128{
				complex(-0.605829586188268, 0),
				complex(-0.072085206905866, -0.6383267351483765),
				complex(-0.072085206905866, +0.6383267351483765),
			},
		},
		{
			// Δ < 0
			ps: []float64{1, -6, 11, -6},
			want: []complex128{
				complex(1, 0),
				complex(2, 0),
				complex(3, 0),
			},
		},
		// -- quartic --
		{
			// x^4 + ax^3 =0, a>0
			ps: []float64{1, 5, 0, 0, 0},
			want: []complex128{
				complex(-5, 0),
				complex(0, 0),
				complex(0, 0),
				complex(0, 0),
			},
		},
		{
			// x^4 + ax^3 =0, a<0
			ps: []float64{1, -5, 0, 0, 0},
			want: []complex128{
				complex(0, 0),
				complex(0, 0),
				complex(0, 0),
				complex(5, 0),
			},
		},
		{
			//  x^4 + d = 0, d>0
			ps: []float64{1, 0, 0, 0, 9},
			want: []complex128{
				complex(-1.2247448713915892, -1.2247448713915890),
				complex(-1.2247448713915890, +1.2247448713915892),
				complex(+1.2247448713915890, -1.2247448713915892),
				complex(+1.2247448713915892, +1.2247448713915890),
			},
		},
		{
			//  x^4 + d = 0, d<0
			ps: []float64{1, 0, 0, 0, -9},
			want: []complex128{
				complex(-1.7320508075688772, -0),
				complex(-0, -1.7320508075688772),
				complex(0, 1.7320508075688772),
				complex(1.7320508075688772, 0),
			},
		},
		{
			//  x^4 + ax^3 + bx^2 = 0
			ps: []float64{1, 25, 4, 0, 0},
			want: []complex128{
				complex(-24.838962679253065, 0),
				complex(-0.16103732074693405, 0),
				complex(0, 0),
				complex(0, 0),
			},
		},
		{
			//  x^4 + ax^3 + bx^2 + cx + d = 0, R != 0
			ps: []float64{5, 4, 3, 2, 1},
			want: []complex128{
				complex(-0.5378322749029899, -0.358284686345128),
				complex(-0.5378322749029899, +0.358284686345128),
				complex(+0.13783227490298988, -0.6781543891053364),
				complex(+0.13783227490298988, +0.6781543891053364),
			},
		},
		{
			//  x^4 + ax^3 + bx^2 + cx + d = 0, R == 0
			ps: []float64{
				2.2206846808021337, 7.643281053997895, 8.831759446092846,
				3.880673545129404, 0.5724551380144077,
			},
			// FIXME(sbinet): using gsl_poly_complex_solve, one gets:
			// z0 = -1.342828035596519198 +0.000000000000000000
			// z1 = -1.342830052099751592 +0.000000000000000000
			// z2 = -0.378099893086095040 -0.000001007044401213
			// z3 = -0.378099893086095040 +0.000001007044401213
			want: []complex128{
				complex(-1.3428290438500117, -1.345173146522382e-06),
				complex(-1.3428290438500117, +1.345173146522382e-06),
				complex(-0.3780998930842191, -1.345173146522382e-06),
				complex(-0.3780998930842191, +1.345173146522382e-06),
			},
		},
	} {
		t.Run("", func(t *testing.T) {
			const tol = 1e-15
			ps := make([]float64, len(test.ps))
			copy(ps, test.ps)
			got := Polynomial(test.dst, ps)
			if got, want := got, test.want; !cmplxs.EqualApprox(got, want, tol) {
				t.Fatalf(
					"invalid polynomial roots from %g:\ngot= %g\nwant=%g",
					test.ps, got, want,
				)
			}
			if a, b := math.Cbrt(-0.5*-3), -math.Cbrt(0.5*-3); a != b {
				t.Fatalf("a=%g, b=%g", a, b)
			}
		})
	}
}
