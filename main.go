package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"slices"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Swallow struct {
	x, y int
}

const (
	xSize = 32
	ySize = 24
)

var (
	gameOver  = false
	bkg       = color.Black
	lightgrey = color.RGBA{0xc2, 0xc5, 0xc6, 0xff}
	yellow    = color.RGBA{0xff, 0xb3, 0x5d, 0xff}
	grey      = color.RGBA{0x77, 0x7c, 0x7e, 0xff}
	green     = color.RGBA{0x60, 0xa6, 0x65, 0xff}
	// foodBeingSwallowed = make([]Swallow, 0)
	mysleepMillis = 100
	keys          []ebiten.Key
	leftPressed   = false
	rightPressed  = false
	upPressed     = false
	downPressed   = false
	leftNano      = int64(0)
	rightNano     = int64(0)
	upNano        = int64(0)
	downNano      = int64(0)
	boardTile1    = ebiten.NewImage(8, 8)
	boardTile2    = ebiten.NewImage(8, 8)
	snakeTile     = ebiten.NewImage(10, 10)
	foodTile      = ebiten.NewImage(6, 6)
	snakeMatrix   [xSize][ySize]bool
)

var mySnake = &List{}

type Node struct {
	next *Node
	prev *Node
	x, y int
}

type List struct {
	head         *Node
	tail         *Node
	dir          string
	xSize, ySize int
	xFood, yFood int
	len          int
}

func (l *List) init(x, y int) {
	boardTile1.Fill(lightgrey)
	boardTile2.Fill(grey)
	snakeTile.Fill(green)
	foodTile.Fill(yellow)
	l.head = &Node{x: x, y: y}
	l.tail = &Node{x: x - 1, y: y}
	snakeMatrix[x][y] = true
	snakeMatrix[x-1][y] = true
	l.head.next = l.tail
	l.head.prev = nil
	l.tail.next = nil
	l.tail.prev = l.head
	l.len = 2
	l.dir = "right"
	l.xSize = xSize
	l.ySize = ySize
	l.generateFood()
	// fmt.Println(snakeMatrix[0][0], 101)
}

func (l *List) next(n *Node) *Node {
	return n.next
}

func (l *List) nextStep() {

	snakeMatrix[l.tail.x][l.tail.y] = false
	newTail := l.tail.prev
	l.tail.prev = nil
	newTail.next = nil
	l.tail = newTail

	newHeadX := l.head.x
	newHeadY := l.head.y
	switch l.dir {
	case "down":
		newHeadY = (l.head.y + 1) % l.ySize
	case "up":
		newHeadY = (l.head.y - 1 + l.ySize) % l.ySize
	case "right":
		newHeadX = (l.head.x + 1) % l.xSize
	case "left":
		newHeadX = (l.head.x - 1 + l.xSize) % l.xSize
	}
	newHead := &Node{x: newHeadX, y: newHeadY}
	newHead.next = l.head
	l.head.prev = newHead
	l.head = newHead

	if snakeMatrix[l.head.x][l.head.y] {
		l.gameOver()
	}
	snakeMatrix[newHeadX][newHeadY] = true

}

func (l *List) gameOver() {
	gameOver = true
	fmt.Println("Game Over!")
}
func (l *List) changeDir(direction string) {
	if l.dir == "up" && direction == "down" {
		return
	}
	if l.dir == "down" && direction == "up" {
		return
	}
	if l.dir == "right" && direction == "left" {
		return
	}
	if l.dir == "left" && direction == "right" {
		return
	}
	l.dir = direction
}

func checkDirection() (string, bool) {
	keys = inpututil.AppendPressedKeys(keys[:0])
	keyString := []string{}
	for _, key := range keys {
		keyString = append(keyString, key.String())
	}

	if slices.Contains(keyString, "ArrowUp") {
		if !upPressed {
			upPressed = true
			upNano = time.Now().UnixNano()
		}
	} else {
		upPressed = false
		upNano = int64(0)
	}
	// ###################################################
	if slices.Contains(keyString, "ArrowDown") {
		if !downPressed {
			downPressed = true
			downNano = time.Now().UnixNano()
		}
	} else {
		downPressed = false
		downNano = int64(0)
	}
	// ####################################################
	if slices.Contains(keyString, "ArrowLeft") {
		if !leftPressed {
			leftPressed = true
			leftNano = time.Now().UnixNano()
		}
	} else {
		leftPressed = false
		leftNano = int64(0)
	}
	// #####################################
	if slices.Contains(keyString, "ArrowRight") {
		if !rightPressed {
			rightPressed = true
			rightNano = time.Now().UnixNano()
		}
	} else {
		rightPressed = false
		rightNano = int64(0)
	}

	maxNano := int64(0)
	lastPressedKey := ""
	returnBool := false
	if upNano > maxNano {
		maxNano = upNano
		lastPressedKey = "up"
		returnBool = true
	}
	if downNano > maxNano {
		maxNano = downNano
		lastPressedKey = "down"
		returnBool = true
	}

	if rightNano > maxNano {
		maxNano = rightNano
		lastPressedKey = "right"
		returnBool = true
	}

	if leftNano > maxNano {
		maxNano = leftNano
		lastPressedKey = "left"
		returnBool = true
	}
	return lastPressedKey, returnBool
}

func (l *List) generateFood() {
	xFood := rand.Intn(l.xSize)
	yFood := rand.Intn(l.ySize)
	fmt.Println(l.len)
	l.xFood, l.yFood = xFood, yFood
}

func (l *List) eatNGrow() {
	newElement := &Node{x: l.tail.x, y: l.tail.y}
	l.tail.next = newElement
	newElement.prev = l.tail
	l.tail = newElement
	l.len = l.len + 1
	mysleepMillis = max(10, mysleepMillis-2)
}

// Game implements ebiten.Game interface.
type Game struct{}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
	time.Sleep(time.Duration(mysleepMillis) * time.Millisecond)
	if mySnake.head.x == mySnake.xFood && mySnake.head.y == mySnake.yFood {
		mySnake.eatNGrow()
		mySnake.generateFood()
	}
	if !gameOver {
		mySnake.nextStep()
	}
	return nil
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(bkg)
	for i := 0; i < 32; i++ {
		for j := 0; j < 24; j++ {
			op := &ebiten.DrawImageOptions{}
			if snakeMatrix[i][j] {
				op.GeoM.Translate(float64(i)*10-1, float64(j)*10-1)
				screen.DrawImage(snakeTile, op)
			} else if (i%2 == 0) != (j%2 == 0) {
				op.GeoM.Translate(float64(i)*10, float64(j)*10)
				screen.DrawImage(boardTile1, op)
			} else {
				op.GeoM.Translate(float64(i)*10, float64(j)*10)
				screen.DrawImage(boardTile2, op)
			}
		}
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(mySnake.xFood)*10+1, float64(mySnake.yFood)*10+1)
	screen.DrawImage(foodTile, op)

}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	game := &Game{}
	mySnake.init(15, 15)
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Snake")
	go func() {
		for {
			if direction, ok := checkDirection(); ok {
				mySnake.changeDir(direction)
			}
		}
	}()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
