package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"log"
	"math"

	_ "embed"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	ScreenWidth  = 1000
	ScreenHeight = 600
	Size         = 150
)

var (
	//go:embed sample_image.png
	sampleImage     []byte
	mplusFaceSource *text.GoTextFaceSource
	whiteImage      = ebiten.NewImage(3, 3)
)

func init() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}
	mplusFaceSource = s
}

func main() {
	game, err := NewGame()
	if err != nil {
		log.Fatal(err)
	}
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Ebitengine Image Problem")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

type Game struct {
	onceText    string
	onceImage   *ebiten.Image
	onceLayer   [4]*ebiten.Image
	updateText  string
	updateImage *ebiten.Image
	updateLayer [4]*ebiten.Image
}

func NewGame() (*Game, error) {
	g := &Game{}
	g.onceText = "Hello, 世界!: "
	img, _, _ := image.Decode(bytes.NewReader(sampleImage))
	g.onceImage = ebiten.NewImageFromImage(img)
	for i := 0; i < 4; i++ {
		g.onceLayer[i] = ebiten.NewImage(5*Size, (4-i)*Size)
	}
	return g, nil
}

func (g *Game) Update() error {
	g.updateText = "Hello, 世界!: "
	img, _, _ := image.Decode(bytes.NewReader(sampleImage))
	g.updateImage = ebiten.NewImageFromImage(img)
	for i := 0; i < 4; i++ {
		g.updateLayer[i] = ebiten.NewImage(5*Size, (4-i)*Size)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x80, 0xa0, 0xc0, 0xff})
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, Size)
	for i := 3; i >= 0; i-- {
		layerStr := fmt.Sprintf("%d層", i+1)
		g.drawTxt(g.onceLayer[i], g.onceText+layerStr, 0, 0)
		g.drawImg(g.onceLayer[i], g.onceImage, 0, Size/3)
		g.drawVec(g.onceLayer[i], Size*2/3, Size/3)
		g.drawTxt(g.updateLayer[i], g.updateText+layerStr, 0, 0)
		g.drawImg(g.updateLayer[i], g.updateImage, 0, Size/3)
		g.drawVec(g.updateLayer[i], Size*2/3, Size/3)
		if i == 0 {
			ocOp := &ebiten.DrawImageOptions{}
			ocOp.GeoM.Translate(0, 0)
			screen.DrawImage(g.onceLayer[i], ocOp)
			upOp := &ebiten.DrawImageOptions{}
			upOp.GeoM.Translate(3*Size, 0)
			screen.DrawImage(g.updateLayer[i], upOp)
			break
		}
		g.onceLayer[i-1].DrawImage(g.onceLayer[i], op)
		g.updateLayer[i-1].DrawImage(g.updateLayer[i], op)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func (g *Game) drawTxt(screen *ebiten.Image, txt string, x, y int) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	op.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, txt, &text.GoTextFace{
		Source: mplusFaceSource,
		Size:   Size / 3,
	}, op)
}

func (g *Game) drawImg(screen, img *ebiten.Image, x, y int) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(Size*2/3)/float64(img.Bounds().Dx()), float64(Size*2/3)/float64(img.Bounds().Dy()))
	op.GeoM.Translate(float64(x), float64(y))
	op.Filter = ebiten.FilterLinear
	screen.DrawImage(img, op)
}

func (g *Game) drawVec(screen *ebiten.Image, x, y int) {
	r := Size / 3
	rf, xf, yf := float32(r), float32(x), float32(y)
	var path vector.Path
	path.Arc(xf+rf, yf+rf, rf, 0, 2*math.Pi, vector.Clockwise)
	path.Close()
	vs, is := path.AppendVerticesAndIndicesForFilling(nil, nil)
	op := &ebiten.DrawTrianglesOptions{}
	op.AntiAlias = true
	op.FillRule = ebiten.NonZero
	whiteImage.Fill(color.White)
	screen.DrawTriangles(vs, is, whiteImage, op)
}
