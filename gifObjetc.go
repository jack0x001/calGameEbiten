package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"image/gif"
	"os"
	"time"
)

type GifObject struct {
	Anim        *gif.GIF
	FrameNum    int
	FrameTime   time.Time
	FrameImage  *ebiten.Image
	isPaused    bool  //暂停
	CustomDelay int64 //自定义GIF两个图片帧之间的时间间隔, 如果为0, 则使用GIF的自带帧间隔
}

func (obj *GifObject) Pause() {
	obj.isPaused = true
}

func (obj *GifObject) ContinuePlay() {
	obj.isPaused = false
}

func (obj *GifObject) Update() error {
	//如果传入的delay为负数, 则使用GIF自带的delay
	delay := int64(obj.Anim.Delay[obj.FrameNum])
	if obj.CustomDelay > 0 {
		delay = obj.CustomDelay
	}
	currentTime := time.Now()
	elapsed := currentTime.Sub(obj.FrameTime)

	if obj.FrameImage != nil && elapsed.Milliseconds() < delay*10 {
		return nil
	}

	// frame rectangle
	var bounds = obj.Anim.Image[obj.FrameNum].Bounds()
	// base picture size
	imgW, imgH := obj.Anim.Config.Width, obj.Anim.Config.Height
	// frame image
	srcImg := ebiten.NewImageFromImage(obj.Anim.Image[obj.FrameNum])
	//src_w, src_h := src_img.Size()
	imgOpts := &ebiten.DrawImageOptions{}
	// shifting by rectangle properties
	imgOpts.GeoM.Translate(float64(bounds.Min.X), float64(bounds.Min.Y))

	// image from reference size
	obj.FrameImage = ebiten.NewImage(imgW, imgH)
	//img.Fill(color.RGBA{0xF0, 0x10, 0x0F, 0xFF})

	obj.FrameImage.DrawImage(srcImg, imgOpts)
	obj.FrameTime = currentTime
	if !obj.isPaused {
		obj.FrameNum += 1
		if obj.FrameNum >= len(obj.Anim.Image) {
			obj.FrameNum = 0
		}
	}

	return nil
}

func (obj *GifObject) Draw(screen *ebiten.Image, where image.Rectangle) {
	if obj.FrameImage != nil {
		imgW, imgH := obj.Anim.Config.Width, obj.Anim.Config.Height
		scrW, _ := where.Dx(), where.Dy()
		opts := &ebiten.DrawImageOptions{}

		//等比例缩放
		ratioOfImgWidthAndHeight := float64(imgW) / float64(imgH)
		scaleOfWidth := float64(scrW) / float64(imgW)
		scaleOfHeight := scaleOfWidth * ratioOfImgWidthAndHeight
		opts.GeoM.Scale(scaleOfWidth, scaleOfHeight)

		//pos to place image
		opts.GeoM.Translate(float64(where.Min.X), float64(where.Min.Y))
		////color white
		//opts.ColorM.Scale(1, 1, 1, 1)

		screen.DrawImage(obj.FrameImage, opts)
	}
}

func MakeGifObject(filePath string) (obj *GifObject, err error) {
	obj = new(GifObject)

	fileReader, err := os.Open(filePath)
	if err != nil {
		return
	}
	obj.Anim, err = gif.DecodeAll(fileReader)
	if err != nil {
		return
	}

	return
}
