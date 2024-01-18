package config

import (
	"image"
	"math"
)

type Rect struct {
	X, Y, Width, Height float64
}

type Line struct {
	X1, Y1, X2, Y2 float64
}

type Circle struct {
	X, Y, Radius float64
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

func IntersectLine(l Line, r image.Rectangle) bool {
	x1, y1, x2, y2 := l.X1, l.Y1, l.X2, l.Y2
	res := LineTorectIntersection(x1, y1, x2, y2, r)
	return res
}

func IntersectCircle(r image.Rectangle, c Circle) bool {
	testX := c.X
	testY := c.Y

	// which edge is closest?
	if c.X < float64(r.Min.X) {
		testX = float64(r.Min.X) // test left edge
	} else if c.X > float64(r.Max.X) {
		testX = float64(r.Max.X) // right edge
	}
	if c.Y < float64(r.Min.Y) {
		testY = float64(r.Min.Y) // top edge
	} else if c.Y > float64(r.Max.Y) {
		testY = float64(r.Max.Y) // bottom edge
	}

	// get distance from closest edges
	distX := c.X - testX
	distY := c.Y - testY
	distance := math.Sqrt((distX * distX) + (distY * distY))

	// if the distance is less than the radius, collision!
	return distance <= c.Radius
}

// from https://jeffreythompson.org/collision-detection/line-rect.php
func lineToLineIntersection(x1 float64, y1 float64, x2 float64, y2 float64, x3 float64, y3 float64, x4 float64, y4 float64) bool {
	// calculate the direction of the lines
	uA := ((x4-x3)*(y1-y3) - (y4-y3)*(x1-x3)) / ((y4-y3)*(x2-x1) - (x4-x3)*(y2-y1))
	uB := ((x2-x1)*(y1-y3) - (y2-y1)*(x1-x3)) / ((y4-y3)*(x2-x1) - (x4-x3)*(y2-y1))

	// if uA and uB are between 0-1, lines are colliding
	if uA >= 0 && uA <= 1 && uB >= 0 && uB <= 1 {
		return true
	}
	return false
}

func LineTorectIntersection(x1, y1, x2, y2 float64, r image.Rectangle) bool {
	left := lineToLineIntersection(x1, y1, x2, y2, float64(r.Min.X), float64(r.Min.Y), float64(r.Min.X), float64(r.Max.Y))
	right := lineToLineIntersection(x1, y1, x2, y2, float64(r.Max.X), float64(r.Min.Y), float64(r.Max.X), float64(r.Max.Y))
	top := lineToLineIntersection(x1, y1, x2, y2, float64(r.Min.X), float64(r.Min.Y), float64(r.Max.X), float64(r.Min.Y))
	bottom := lineToLineIntersection(x1, y1, x2, y2, float64(r.Min.X), float64(r.Max.Y), float64(r.Max.X), float64(r.Max.Y))
	if left || right || top || bottom {
		return true
	}
	return false
}

// https://en.wikipedia.org/wiki/Cohen%E2%80%93Sutherland_algorithm(Cohen-Sutherland algorithm)
// func computeCode(x, y float64, r image.Rectangle) int {
// 	code := INSIDE
// 	if x < float64(r.Min.X) {
// 		code |= LEFT
// 	} else if x > float64(r.Max.X) {
// 		code |= RIGHT
// 	}
// 	if y < float64(r.Min.Y) {
// 		code |= BOTTOM
// 	} else if y > float64(r.Max.Y) {
// 		code |= TOP
// 	}
// 	return code
// }

// func cohenSutherlandClip(x1, y1, x2, y2 float64, r image.Rectangle) bool {
// 	code1 := computeCode(x1, y1, r)
// 	code2 := computeCode(x2, y2, r)
// 	var code int
// 	var x, y float64
// 	for {
// 		if code1&code2 == 0 {
// 			// Both endpoints are inside the rectangle
// 			return false
// 		}
// 		if code1&code2 != 0 {
// 			// Both endpoints are outside the rectangle, in same region
// 			return false
// 		}
// 		if code1 != 0 {
// 			code = code1
// 		} else {
// 			code = code2
// 		}

// 		if (LEFT & code) != 0 {
// 			x = float64(r.Min.X)
// 			y = y1 + (y2-y1)*(float64(r.Min.X)-x1)/(x2-x1)
// 		} else if (RIGHT & code) != 0 {
// 			x = float64(r.Max.X)
// 			y = y1 + (y2-y1)*(float64(r.Max.X)-x1)/(x2-x1)
// 		} else if (BOTTOM & code) != 0 {
// 			y = float64(r.Min.Y)
// 			x = x1 + (x2-x1)*(float64(r.Min.Y)-y1)/(y2-y1)
// 		} else if (TOP & code) != 0 {
// 			y = float64(r.Max.Y)
// 			x = x1 + (x2-x1)*(float64(r.Max.Y)-y1)/(y2-y1)
// 		}
// 		if code == code1 {
// 			x1 = x
// 			y1 = y
// 			code1 = computeCode(x, y, r)
// 		} else {
// 			x2 = x
// 			y2 = y
// 			code2 = computeCode(x, y, r)
// 		}
// 		if x != 0 || y != 0 {
// 			return true
// 		}
// 	}
// }
