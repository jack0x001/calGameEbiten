package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"image"
	"image/color"
	_ "image/png"
	"strconv"
	"strings"
)

const (
	screenWidth       = 640
	screenHeight      = 480
	tileSize          = 32
	titleFontSize     = fontSize * 1.5
	fontSize          = 24
	smallFontSize     = fontSize / 2
	maxAnswer         = 20 //问题涉及到的最大数字, 比如20以内的加减法
	maxQuestionNumber = 30 //题目数量
	comeOnTimeDelay   = 60 //播放come on 进行催促的时间间隔, 秒
)

// Mode 模式(状态)
type Mode int

const (
	ModeTitle Mode = iota
	ModeGame
	ModeGameOver
)

//Game 定义了Game结构。
//Game实现了ebiten.Game接口。
//Ebiten.Game拥有一个Ebiten游戏的必要功能。更新、绘图和布局.
type Game struct {
	mode Mode

	x16  int
	y16  int
	vy16 int

	// Camera
	cameraX int
	cameraY int

	pressed []string

	score int

	question *Question

	dancingGopher *GifObject
	flagGif       *GifObject
	walkingMario  *GifObject

	updateCount int
}

func (game *Game) init() {

	game.x16 = 0
	game.y16 = 380 * 16

	game.cameraX = -240
	game.cameraY = 0

	game.score = 0
	game.question = NewQuestion()

	dance, err := MakeGifObject("./res/dancingGopher.gif")
	if err != nil {
		panic(err)
	}
	game.dancingGopher = dance

	flag, err := MakeGifObject("./res/flag.gif")
	if err != nil {
		panic(err)
	}
	game.flagGif = flag

	mario, err := MakeGifObject("./res/walkingMario.gif")
	if err != nil {
		panic(err)
	}
	game.walkingMario = mario

}

//Draw 每一帧都会被调用。
//帧是渲染的一个时间单位，这取决于显示器的刷新率。如果显示器的刷新率是60[Hz]，Draw每秒被调用60次。
//Draw需要一个参数screen，它是一个指向ebiten.Image的指针。
//在Ebiten中，所有图像如从图像文件创建的图像、屏幕外的图像（临时渲染目标）和屏幕都表示为ebiten.Image对象。
//screen参数是渲染的最终目的地。窗口显示每一帧屏幕的最终状态。
func (game *Game) Draw(screen *ebiten.Image) {

	switch game.mode {
	case ModeTitle:
		game.drawTitle(screen)
		game.drawBottomTileImages(screen)
	case ModeGame:
		game.drawGame(screen)
		game.drawBottomTileImages(screen)
		game.drawWalkingMarioAndFlag(screen)
	case ModeGameOver:
		game.drawGameOver(screen)
	default:
		panic("unexpected game mode")
	}
}

func (game *Game) drawTitle(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 0x80, G: 0xa0, B: 0xc0, A: 0xff})

	titleTexts := []string{"PLAY MATH"}
	texts := []string{"", "", "", "", "", "", "", "PRESS SPACE KEY", "", "", "", ""}
	for i, l := range titleTexts {
		x := (screenWidth - len(l)*titleFontSize) / 2
		text.Draw(screen, l, titleArcadeFont, x, (i+4)*titleFontSize, color.White)
	}
	for i, l := range texts {
		x := (screenWidth - len(l)*fontSize) / 2
		text.Draw(screen, l, arcadeFont, x, (i+4)*fontSize, color.White)
	}

	game.drawMario(screen)
}

//绘制底部图片
func (game *Game) drawBottomTileImages(screen *ebiten.Image) {
	const (
		nx = screenWidth / tileSize
		ny = screenHeight / tileSize
	)

	op := &ebiten.DrawImageOptions{}
	for i := -2; i < nx+1; i++ {
		// ground
		op.GeoM.Reset()
		op.GeoM.Translate(float64(i*tileSize-floorMod(game.cameraX, tileSize)),
			float64((ny-1)*tileSize-floorMod(game.cameraY, tileSize)))
		screen.DrawImage(tilesImage.SubImage(image.Rect(0, 0, tileSize, tileSize)).(*ebiten.Image), op)

	}
}

//绘制卡通人物
func (game *Game) drawMario(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}

	op.GeoM.Translate(450, 100)
	screen.DrawImage(marioImage, op)
}

var targetMarioPos = image.Rectangle{}
var currentMarioPos = image.Rectangle{}

func (game *Game) drawWalkingMarioAndFlag(screen *ebiten.Image) {

	marioScale := 4
	flagScale := 5

	//mario
	stepLength := float64(screenWidth-game.flagGif.Anim.Config.Width/flagScale) / float64(maxQuestionNumber+1) //每一步的长度 (留一部分画🚩)
	minX := game.score * int(stepLength)
	minY := screenHeight - game.walkingMario.Anim.Config.Height/marioScale - tileSize //确保站在草地上方
	maxX := minX + game.walkingMario.Anim.Config.Width/marioScale
	maxY := screenHeight - tileSize
	targetMarioPos = image.Rectangle{Min: image.Point{X: minX, Y: minY}, Max: image.Point{X: maxX, Y: maxY}}
	currentMinX := currentMarioPos.Min.X + 1
	if currentMinX > targetMarioPos.Min.X {
		currentMinX = targetMarioPos.Min.X
		game.walkingMario.Pause()
	} else {
		game.walkingMario.CustomDelay = 5
		game.walkingMario.ContinuePlay()
	}

	currentMaxX := currentMarioPos.Max.X + 1
	if currentMaxX > targetMarioPos.Max.X {
		currentMaxX = targetMarioPos.Max.X
	}
	currentMarioPos = image.Rectangle{
		Min: image.Point{X: currentMinX, Y: targetMarioPos.Min.Y},
		Max: image.Point{X: currentMaxX, Y: targetMarioPos.Max.Y},
	}
	game.walkingMario.Draw(screen, currentMarioPos)

	//flag
	minX = screenWidth - game.flagGif.Anim.Config.Width/flagScale
	minY = screenHeight - game.flagGif.Anim.Config.Height/flagScale - 20
	maxX = screenWidth
	maxY = screenHeight //确保站在草地上方
	pos := image.Rectangle{Min: image.Point{X: minX, Y: minY}, Max: image.Point{X: maxX, Y: maxY}}
	game.flagGif.Draw(screen, pos)
}

func (game *Game) drawFlag(screen *ebiten.Image) {

}

var lastRandomX = 0
var lastRandomY = 0

func (game *Game) drawGameOver(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 0x80, G: 0xa0, B: 0xc0, A: 0xff})

	texts := []string{"YEAH!"}
	for i, l := range texts {
		x := (screenWidth - len(l)*titleFontSize) / 2
		text.Draw(screen, l, titleArcadeFont, x, (i+4)*titleFontSize, color.White)
	}

	game.drawDancingGopher(screen)
}

func (game *Game) drawDancingGopher(screen *ebiten.Image) {
	halfScreenWidth := screenWidth / 2
	halfScreenHeight := screenHeight / 2

	if game.updateCount%60 == 0 {
		lastRandomX = randInt(0-halfScreenWidth/2, halfScreenWidth/2)
		lastRandomY = randInt(0-halfScreenHeight/2, halfScreenHeight/2)
	}

	centerX := halfScreenWidth + lastRandomX
	centerY := halfScreenHeight + lastRandomY
	halfWidth := game.dancingGopher.Anim.Config.Width / 2
	halfHeight := game.dancingGopher.Anim.Config.Height / 2
	minX := centerX - halfWidth
	maxX := centerX + halfWidth
	minY := centerY - halfHeight
	maxY := centerY + halfHeight

	game.dancingGopher.Draw(screen, image.Rectangle{
		Min: image.Point{X: minX, Y: minY},
		Max: image.Point{X: maxX, Y: maxY},
	})

	game.dancingGopher.Draw(screen, image.Rectangle{
		Min: image.Point{X: minX + 2*halfWidth, Y: minY + halfHeight},
		Max: image.Point{X: maxX + halfWidth, Y: maxY},
	})
}

func (game *Game) drawGame(screen *ebiten.Image) {

	screen.Fill(color.RGBA{R: 0x80, G: 0xa0, B: 0xc0, A: 0xff})

	//draw question
	questionsTexts := []string{game.question.String()}
	for i, l := range questionsTexts {
		x := (screenWidth - len(l)*titleFontSize) / 2
		text.Draw(screen, l, titleArcadeFont, x, (i+4)*titleFontSize, color.White)
	}

	//draw input
	keyString := ""
	for _, p := range game.pressed {
		keyString += p
	}
	texts := []string{"", "", "", "", "", "", "", keyString, "", "", "", ""}

	for i, l := range texts {
		x := (screenWidth - len(l)*fontSize) / 2
		text.Draw(screen, l, arcadeFont, x, (i+4)*fontSize, color.White)
	}

	//draw score
	scoreStr := fmt.Sprintf("%04d", game.score)
	scoreX := screenWidth - len(scoreStr)*fontSize
	scoreY := fontSize
	text.Draw(screen, scoreStr, arcadeFont, scoreX, scoreY, color.White)

	//draw coin
	op := &ebiten.DrawImageOptions{}
	coinScale := 0.25
	coinWidth, _ := coinImage.Size()
	op.GeoM.Scale(coinScale, coinScale)
	op.GeoM.Translate(float64(scoreX)-float64(coinWidth)*coinScale, 0.0)
	screen.DrawImage(coinImage, op)

}

func (game *Game) ReTiming() {
	game.updateCount = 0
}

//Update 函数，每一个Tick都会被调用。
//Tick是逻辑更新的一个时间单位。默认值是1/60[s]，
//那么Update默认每秒被调用60次（即一个Ebiten游戏以每秒60次的速度工作）。
//Update更新游戏的逻辑状态。
//Update返回一个错误值。一般来说，当更新函数返回一个非零的错误时，Ebiten游戏就暂停了。
func (game *Game) Update() error {

	if game.mode == ModeTitle && inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		game.mode = ModeGame
	}

	//get max number of type int
	maxInt := int(^uint(0) >> 1)
	if game.updateCount >= maxInt-10 {
		game.updateCount = 0
	}
	game.updateCount++

	switch game.mode {
	case ModeGame:

		err := game.flagGif.Update()
		if err != nil {
			return err
		}
		err = game.walkingMario.Update()
		if err != nil {
			return err
		}

		if game.updateCount > 60*comeOnTimeDelay {
			game.ReTiming()
			PlaySoundComeOn()
		}

		//TODO: 点击键盘时 加入声音

		for i := 0; i <= 9; i++ {
			if inpututil.IsKeyJustPressed(ebiten.Key(i) + ebiten.KeyDigit0) {
				game.pressed = append(game.pressed, string(rune(i+'0')))
			}
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
			if len(game.pressed) > 0 {
				game.pressed = game.pressed[:len(game.pressed)-1]
			}
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			if len(game.pressed) > 0 {
				//如果输入的数字相加等于answer，则加分
				if answer, _ := strconv.Atoi(strings.Join(game.pressed, "")); answer == game.question.Answer() {
					PlaySoundOK()
					game.ReTiming()
					game.score += 1
					if game.score == maxQuestionNumber {
						game.mode = ModeGameOver
						PlaySoundHappyEnding()
					} else {
						game.question = NewQuestion()
					}

				} else {
					PlaySoundError()
				}
				game.pressed = []string{}
			}
		}
	case ModeGameOver:
		err := game.dancingGopher.Update()
		if err != nil {
			return err
		}
	default:
		break
	}

	return nil
}

//Layout 函数。
//Layout接受一个外部尺寸，也就是桌面上的窗口尺寸，并返回游戏的逻辑屏幕尺寸。
//这段代码忽略了参数并返回固定值。这意味着游戏的屏幕尺寸总是相同的，无论窗口的尺寸是多少。
//当窗口可以调整大小的时候Layout将更有意义。
func (game *Game) Layout(int, int) (int, int) {
	return screenWidth, screenHeight
}
func NewGame() *Game {
	game := Game{}
	game.init()
	return &game
}
