// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package roots

import (
	"fmt"
	"math"
	"math/cmplx"
	"sort"

	"gonum.org/v1/gonum/mat"
)

// Polynomial returns the roots of an n-th degree polynomial of the form:
//
//	p[0] * x**n + p[1] * x**(n-1) + ... + p[n-1]*x + p[n]
func Polynomial(dst []complex128, ps []float64) []complex128 {
	n := len(ps) - 1
	if dst == nil {
		dst = make([]complex128, n)
	}
	if len(dst) != n {
		panic("roots: length mismatch")
	}

	for _, v := range ps {
		if v != 0 {
			break
		}
		n--
	}
	dst = dst[:n]

	ps = ps[len(ps)-1-n:]

	switch n {
	case 0:
		return dst

	case 1:
		a := ps[0]
		b := ps[1]
		if a == 0 {
			return dst
		}
		r := -b / a
		dst[0] = complex(r, 0)
		return dst

	case 2:
		dst[0], dst[1] = Poly2(ps[0], ps[1], ps[2])
		return dst

	case 3:
		dst[0], dst[1], dst[2] = Poly3(ps[0], ps[1], ps[2], ps[3])
		return dst

	case 4:
		dst[0], dst[1], dst[2], dst[3] = Poly4(ps[0], ps[1], ps[2], ps[3], ps[4])
		return dst

	default:
		// use companion matrix for n>4

		// FIXME: balance Hessenberg matrix (the companion matrix) as
		// laid out in Numerical Recipes (3rd Ed.), 11.6,11.7,11.8
		// This is what GSL is doing according to:
		//   https://www.gnu.org/software/gsl/doc/html/poly.html
		norm := 1 / ps[0]
		comp := mat.NewDense(n, n, nil)
		row := comp.RawRowView(0)
		copy(row, ps[1:])
		for i := range row {
			row[i] *= -norm
		}
		for i := 1; i < n; i++ {
			row := comp.RawRowView(i)
			row[i-1] = 1
		}
		var eigen mat.Eigen
		if eigen.Factorize(comp, mat.EigenNone) {
			eigen.Values(dst)
		}

		sortCmplx(dst)
		return dst
	}
}

func sortCmplx(xs []complex128) {
	sort.Slice(xs, func(i, j int) bool {
		return min128(xs[i], xs[j])
	})
}

func min128(zi, zj complex128) bool {
	ri := real(zi)
	rj := real(zj)
	if ri == rj {
		ii := imag(zi)
		ij := imag(zj)
		return ii < ij
	}
	return ri < rj
}

// Poly2 returns the roots of the following 2nd degree polynomial:
//
//	p0 x^2 + p1 x + p2 = 0
func Poly2(p0, p1, p2 float64) (z0, z1 complex128) {
	// reduce to: x^2 + a.x + b = 0
	norm := 1 / p0
	b := p2 * norm
	a := p1 * norm

	ha := -0.5 * a
	delta := cmplx.Sqrt(complex(ha*ha-b, 0))
	z0 = complex(ha, 0) + delta
	z1 = complex(ha, 0) - delta
	if min128(z0, z1) {
		return z0, z1
	}
	return z1, z0
}

// Poly3 returns the roots of the following 3rd degree polynomial:
//
//	p0 x^3 + p1 x^2 + p2 x + p3 = 0
func Poly3(p0, p1, p2, p3 float64) (z0, z1, z2 complex128) {
	// use Cardano/Tartaglia/Vieta formulae.
	// see:
	//  https://en.wikipedia.org/wiki/Cubic_equation#Cardano's_formula
	//  https://trans4mind.com/personal_development/mathematics/polynomials/cubicAlgebra.htm

	// reduce to monic form.
	//  x^3 + ax^2 + bx + c = 0
	ip := 1 / p0
	a := p1 * ip
	b := p2 * ip
	c := p3 * ip

	const (
		k3  = 1 / 3.0
		k27 = 1 / 27.0
	)
	a2 := a * a
	p := k3 * (3*b - a2)
	q := k27 * (2*a2*a - 9*a*b + 27*c)

	if p == 0 && q == 0 {
		x := complex(-a*k3, 0)
		return x, x, x
	}

	Δ := 0.25*(q*q) + (p*p*p)*k27

	const ε = 1e-15
	switch {
	case Δ == 0, math.Abs(Δ) < ε:
		cbrt := math.Cbrt(0.5 * q)
		x0 := complex(-2*cbrt-a*k3, 0)
		x1 := complex(+1*cbrt-a*k3, 0)

		if real(x0) < real(x1) {
			return x0, x1, x1
		}
		return x1, x1, x0

	case Δ > ε:
		sq := math.Sqrt(Δ)
		hq := 0.5 * q
		u1 := math.Cbrt(-hq + sq)
		v1 := math.Cbrt(+hq + sq)
		re := -0.5*(u1-v1) - a*k3
		im := +0.5 * (u1 + v1) * math.Sqrt(3)

		zs := [3]complex128{
			complex(u1-v1-a*k3, 0),
			complex(re, +im),
			complex(re, -im),
		}
		sortCmplx(zs[:])
		return zs[0], zs[1], zs[2]

	case Δ < ε:
		pp := -p / 3
		r := math.Sqrt(pp * pp * pp)
		θ := math.Acos(-0.5 * q / r)
		math.Cos(θ * k3)

		r3 := 2 * math.Cbrt(r)
		zs := [3]complex128{
			complex(r3*math.Cos(θ*k3)-a*k3, 0),
			complex(r3*math.Cos((θ+2*math.Pi)*k3)-a*k3, 0),
			complex(r3*math.Cos((θ+4*math.Pi)*k3)-a*k3, 0),
		}
		sortCmplx(zs[:])
		return zs[0], zs[1], zs[2]
	}
	panic(fmt.Errorf("impossible delta=%g, p0=%g p1=%g p2=%g p3=%g",
		Δ, p0, p1, p2, p3,
	))
}

// Poly4 returns the roots of the following 4th degree polynomial:
//
//	p0 x^4 + p1 x^3 + p2 x^2 + p3 x + p4 = 0
func Poly4(p0, p1, p2, p3, p4 float64) (z0, z1, z2, z3 complex128) {
	// https://en.wikipedia.org/wiki/Quartic_function
	// https://doi.org/10.1016/j.cam.2010.04.015

	// reduce to monic form.
	//  x^4 + ax^3 + bx^2 + cx + d = 0
	ip := 1 / p0
	a := p1 * ip
	b := p2 * ip
	c := p3 * ip
	d := p4 * ip

	// handle degenerate cases.
	if b == 0 && c == 0 {
		// x^4 + ax^3 + d = 0
		switch {
		case d == 0:
			// x^4 + ax^3 = 0
			if a > 0 {
				z0 = complex(-a, 0)
			} else {
				z3 = complex(-a, 0)
			}
			return z0, z1, z2, z3
		case a == 0:
			//  x^4 + d = 0
			if d > 0 {
				z3 = cmplx.Sqrt(complex(math.Sqrt(d), 0) * complex(0, 1))
				z2 = complex(0, -1) * z3
				z1 = complex(-real(z2), -imag(z2))
				z0 = complex(-real(z3), -imag(z3))
			} else {
				z3 = cmplx.Sqrt(complex(math.Sqrt(-d), 0))
				z2 = z3 * complex(0, 1)
				z1 = complex(-real(z2), -imag(z2))
				z0 = complex(-real(z3), -imag(z3))
			}
			return z0, z1, z2, z3
		}
	}

	if c == 0 && d == 0 {
		// x^4 + ax^3 + bx^2 = 0
		var zs [4]complex128
		zs[0], zs[1] = Poly2(1, a, b)
		sortCmplx(zs[:])
		return zs[0], zs[1], zs[2], zs[3]
	}

	// construct and solve cubic resolvent.
	//  y^3 + A y^2 + B y + C = 0
	a2 := a * a
	ya := -b
	yb := c*a - 4*d
	yc := 4*b*d - c*c - a2*d
	y1, y2, y3 := Poly3(1, ya, yb, yc)
	var y float64
	switch {
	case imag(y1) == 0:
		y = real(y1)
	case imag(y2) == 0:
		y = real(y2)
	case imag(y3) == 0:
		y = real(y3)
	}

	var (
		R = cmplx.Sqrt(complex(0.25*a2-b+y, 0))

		f = complex(0.75*a2-2*b, 0)
		g = 2 * cmplx.Sqrt(complex(y*y-4*d, 0))

		D complex128
		E complex128
	)

	switch {
	case R == 0:
		D = cmplx.Sqrt(f + g)
		E = cmplx.Sqrt(f - g)
	default:
		h := complex(0.25*(4*a*b-8*c-a*a2), 0) / R
		r := R * R
		D = cmplx.Sqrt(f - r + h)
		E = cmplx.Sqrt(f - r - h)
	}

	rr := 0.5 * R
	dd := 0.5 * D
	ee := 0.5 * E
	za := 0.25 * complex(a, 0)

	var zs [4]complex128

	zs[0] = -za + rr + dd
	zs[1] = -za + rr - dd
	zs[2] = -za - rr + ee
	zs[3] = -za - rr - ee

	sortCmplx(zs[:])
	return zs[0], zs[1], zs[2], zs[3]
}
