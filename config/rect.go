package config

import (
	"image"
)

type Rect struct {
	X, Y, Width, Height float64
}

type Line struct {
	X1, Y1, X2, Y2 float64
}

const (
	INSIDE = 0
	LEFT   = 1
	RIGHT  = 2
	BOTTOM = 4
	TOP    = 8
)

func NewLine(x1, y1, x2, y2 float64) Line {
	return Line{
		X1: x1,
		Y1: y1,
		X2: x2,
		Y2: y2,
	}
}
func NewRectangle(x, y, width, height float64) Rect {
	return Rect{
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
	}
}

func NewRect(x, y, width, height float64) image.Rectangle {
	return image.Rectangle{
		Min: image.Point{
			X: int(x),
			Y: int(y),
		},
		Max: image.Point{
			X: int(x + width),
			Y: int(y + height),
		},
	}
}

func (r Rect) MaxX() float64 {
	return r.X + r.Width
}

func (r Rect) MaxY() float64 {
	return r.Y + r.Height
}

func (r Rect) Intersects(other Rect) bool {
	return r.X <= other.MaxX() &&
		other.X <= r.MaxX() &&
		r.Y <= other.MaxY() &&
		other.Y <= r.MaxY()
}

func IntersectRect(a, b image.Rectangle) bool {
	interseption := a.Intersect(b)
	if interseption.Min.X == 0 && interseption.Min.Y == 0 && interseption.Max.X == 0 && interseption.Max.Y == 0 {
		return false
	}
	return true
}

// code by Phind (Cohen-Sutherland algorithm) START
func computeCode(x, y float64, r image.Rectangle) int {
	code := INSIDE
	if x < float64(r.Min.X) {
		code |= LEFT
	} else if x > float64(r.Max.X) {
		code |= RIGHT
	}
	if y < float64(r.Min.Y) {
		code |= BOTTOM
	} else if y > float64(r.Max.Y) {
		code |= TOP
	}
	return code
}

func cohenSutherlandClip(x1, y1, x2, y2 float64, r image.Rectangle) (float64, float64, float64, float64) {
	code1 := computeCode(x1, y1, r)
	code2 := computeCode(x2, y2, r)
	for {
		if code1&code2 != 0 || code1 == 0 && code2 == 0 {
			// Both endpoints are outside or inside the rectangle, in same region
			return 0, 0, 0, 0
		} else {
			return 1, 1, 1, 1
		}
	}
}

func IntersectLine(l Line, r image.Rectangle) bool {
	x1, y1, x2, y2 := l.X1, l.Y1, l.X2, l.Y2
	x1, y1, x2, y2 = cohenSutherlandClip(x1, y1, x2, y2, r)
	if x1 == 0 && y1 == 0 && x2 == 0 && y2 == 0 {
		return false
	}
	return true
}

// code by Phind END
