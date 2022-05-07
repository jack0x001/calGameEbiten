package main

import (
	"bufio"
	"bytes"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	resources "github.com/hajimehoshi/ebiten/v2/examples/resources/images/flappy"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image/gif"
	"io/ioutil"
	"os"
	"time"

	"image"
	"log"
)

var (
	titleArcadeFont font.Face
	arcadeFont      font.Face
	smallArcadeFont font.Face
	marioImage      *ebiten.Image
	coinImage       *ebiten.Image
	// the bottom image is used for the background
	tilesImage *ebiten.Image

	audioContext  *audio.Context
	okPlayer      *audio.Player
	errorPlayer   *audio.Player
	happyPlayer   *audio.Player
	welcomePlayer *audio.Player
	comeOnPlayer  *audio.Player
)

func init() {
	ttf, err := opentype.Parse(fonts.PressStart2P_ttf)
	if err != nil {
		log.Fatal(err)
	}
	const dpi = 72
	titleArcadeFont, err = opentype.NewFace(ttf, &opentype.FaceOptions{
		Size:    titleFontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	arcadeFont, err = opentype.NewFace(ttf, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	smallArcadeFont, err = opentype.NewFace(ttf, &opentype.FaceOptions{
		Size:    smallFontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	marioImage = createImageFromFile("./res/mario.png")
	coinImage = createImageFromFile("./res/coin.png")

	img, _, err := image.Decode(bytes.NewReader(resources.Tiles_png))
	if err != nil {
		log.Fatal(err)
	}
	tilesImage = ebiten.NewImageFromImage(img)

	errorPlayer = createWavePlayer("./res/error.wav")
	okPlayer = createWavePlayer("./res/ok.wav")
	happyPlayer = createWavePlayer("./res/happy.wav")
	welcomePlayer = createWavePlayer("./res/welcome.wav")
	comeOnPlayer = createWavePlayer("./res/comeOn.wav")
}

func createImageFromFile(file string) *ebiten.Image {
	fp, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer func(fp *os.File) {
		err := fp.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(fp)
	bytesArray, err := ioutil.ReadAll(fp)
	if err != nil {
		log.Fatal(err)
	}
	img, _, err := image.Decode(bytes.NewReader(bytesArray))
	if err != nil {
		log.Fatal(err)
	}
	return ebiten.NewImageFromImage(img)
}

func createWavePlayer(file string) *audio.Player {

	if audioContext == nil {
		audioContext = audio.NewContext(48000)
	}

	fp, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}

	// os.File.Read(), io.ReadFull(), and
	// io.ReadAtLeast() all work with a fixed
	// byte slice that you make before you read
	// ioutil.ReadAll() will read every byte
	// from the reader (in this case a file),
	// and return a slice of unknown slice
	data, err := ioutil.ReadAll(fp)
	if err != nil {
		log.Fatal(err)
	}
	decoder, err := wav.DecodeWithSampleRate(48000, bytes.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}
	player, err := audioContext.NewPlayer(decoder)
	if err != nil {
		log.Fatal(err)
	}

	return player
}

func PlaySoundOK() {
	err := okPlayer.Rewind()
	if err != nil {
		return
	}
	okPlayer.Play()
}

func PlaySoundComeOn() {
	err := comeOnPlayer.Rewind()
	if err != nil {
		return
	}
	comeOnPlayer.Play()
}

func PlaySoundError() {
	err := errorPlayer.Rewind()
	if err != nil {
		return
	}
	errorPlayer.Play()
}

func PlaySoundWelcome() {
	time.Sleep(2 * time.Second)

	err := welcomePlayer.Rewind()
	if err != nil {
		return
	}
	welcomePlayer.Play()
}

func PlaySoundHappyEnding() {
	err := happyPlayer.Rewind()
	if err != nil {
		return
	}
	happyPlayer.Play()
}

func decodeGif(file string) *gif.GIF {
	inputFile, err := os.Open(file)
	defer func(inputFile *os.File) {
		err := inputFile.Close()
		if err != nil {
			log.Println(err)
		}
	}(inputFile)
	if err != nil {
		log.Println(err)
		return nil
	}

	r := bufio.NewReader(inputFile)

	g, err := gif.DecodeAll(r)
	if err != nil {
		log.Println(err)
		return nil
	}
	return g
}
