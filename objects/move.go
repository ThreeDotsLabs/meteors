package objects

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

func RotateAndTranslateObject(angle float64, object *ebiten.Image, screen *ebiten.Image, x, y float64) {
	w := object.Bounds().Dx()
	h := object.Bounds().Dy()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
	op.GeoM.Rotate(angle)
	op.GeoM.Translate(x+float64(w)/2, y+float64(h)/2)
	//vector.StrokeRect(screen, float32(x), float32(y), float32(w), float32(h), 2, color.RGBA{255, 255, 255, 255}, false)
	screen.DrawImage(object, op)
}

func RotateAndTranslateAnimation(angle float64, object *ebiten.Image, screen *ebiten.Image, x, y float64) {
	// w := object.Bounds().Dx()
	// h := object.Bounds().Dy()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Rotate(angle)
	op.GeoM.Translate(x, y)
	//vector.StrokeRect(screen, float32(x), float32(y), float32(w), float32(h), 2, color.RGBA{255, 255, 255, 255}, false)
	screen.DrawImage(object, op)
}

func ScaleImg(img *ebiten.Image, scale float64) *ebiten.Image {
	w := int(float64(img.Bounds().Dx()) * scale)
	h := int(float64(img.Bounds().Dy()) * scale)
	if w > 0 && h > 0 {
		scaledImage := ebiten.NewImage(w, h)
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Scale(scale, scale)
		scaledImage.DrawImage(img, opts)
		return scaledImage
	}

	return nil
}

func RandInt(min, max int) int {
	return min + rand.Intn(max-min)
}
