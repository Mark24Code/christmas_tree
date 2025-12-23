package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

const (
	minWidth  = 50
	minHeight = 25
	width     = 80
	height    = 24
)

type Snowflake struct {
	x, y  float64
	speed float64
	char  rune
}

type Scene int

const (
	SceneSnow Scene = iota
	SceneSanta
	SceneTree
)

type App struct {
	screen     tcell.Screen
	name       string
	snowflakes []Snowflake
	frame      int
	scene      Scene
	sceneFrame int
}

func main() {
	// Custom usage/help message
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "🎄 Christmas Magic - A Gorgeous Terminal Animation\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		fmt.Fprintf(os.Stderr, "  -name string\n")
		fmt.Fprintf(os.Stderr, "        Your name to display at the bottom (default \"Merry Christmas\")\n")
		fmt.Fprintf(os.Stderr, "  -h, --help\n")
		fmt.Fprintf(os.Stderr, "        Show this help message\n\n")
		fmt.Fprintf(os.Stderr, "Controls:\n")
		fmt.Fprintf(os.Stderr, "  ESC / Ctrl+C / q  Exit the animation\n\n")
		fmt.Fprintf(os.Stderr, "Author:\n")
		fmt.Fprintf(os.Stderr, "  Name:  Mark24Code\n")
		fmt.Fprintf(os.Stderr, "  Email: mark.zhangyoung@gmail.com\n")
		fmt.Fprintf(os.Stderr, "  Repo:  https://github.com/Mark24Code/christmas_tree\n\n")
	}

	name := flag.String("name", "Merry Christmas", "Your name to display")
	flag.Parse()

	s, err := tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	if err := s.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	w, h := s.Size()
	if w < minWidth || h < minHeight {
		s.Fini()
		fmt.Fprintf(os.Stderr, "Terminal too small! Need %dx%d, got %dx%d\n",
			minWidth, minHeight, w, h)
		os.Exit(1)
	}

	app := &App{
		screen: s,
		name:   *name,
		scene:  SceneSnow,
	}
	app.initSnow()
	app.run()
}

func (a *App) initSnow() {
	rand.Seed(time.Now().UnixNano())
	a.snowflakes = make([]Snowflake, 100)
	for i := range a.snowflakes {
		a.snowflakes[i] = Snowflake{
			x:     rand.Float64() * width,
			y:     rand.Float64() * height,
			speed: 0.2 + rand.Float64()*0.3,
			char:  []rune{'*', '·', '•', '❄', '❅', '❆'}[rand.Intn(6)],
		}
	}
}

func (a *App) run() {
	defer a.screen.Fini()

	quit := make(chan struct{})
	go func() {
		for {
			ev := a.screen.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC || ev.Rune() == 'q' {
					close(quit)
					return
				}
			case *tcell.EventResize:
				a.screen.Sync()
			}
		}
	}()

	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-quit:
			return
		case <-ticker.C:
			a.update()
			a.draw()
			a.frame++
			a.sceneFrame++

			// Scene transitions
			switch a.scene {
			case SceneSnow:
				if a.sceneFrame > 120 { // 6 seconds
					a.scene = SceneSanta
					a.sceneFrame = 0
				}
			case SceneSanta:
				if a.sceneFrame > 120 { // 6 seconds
					a.scene = SceneTree
					a.sceneFrame = 0
				}
			case SceneTree:
				// Stay on tree scene, don't exit automatically
				// User can press ESC, Ctrl+C, or 'q' to exit
			}
		}
	}
}

func (a *App) update() {
	// Update snowflakes
	for i := range a.snowflakes {
		a.snowflakes[i].y += a.snowflakes[i].speed
		a.snowflakes[i].x += math.Sin(a.snowflakes[i].y*0.1+float64(i)) * 0.3

		if a.snowflakes[i].y > height {
			a.snowflakes[i].y = 0
			a.snowflakes[i].x = rand.Float64() * width
		}
	}
}

func (a *App) draw() {
	a.screen.Clear()

	switch a.scene {
	case SceneSnow:
		a.drawSceneSnow()
	case SceneSanta:
		a.drawSceneSanta()
	case SceneTree:
		a.drawSceneTree()
	}

	a.screen.Show()
}

func (a *App) drawSceneSnow() {
	// White fade-in background
	fade := float64(a.sceneFrame) / 30.0
	if fade > 1.0 {
		fade = 1.0
	}

	// Background
	bgStyle := tcell.StyleDefault.Background(tcell.ColorWhite)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if rand.Float64() < fade {
				a.screen.SetContent(x, y, ' ', nil, bgStyle)
			}
		}
	}

	// Draw snowflakes
	snowStyle := tcell.StyleDefault.Foreground(tcell.ColorSilver).Background(tcell.ColorWhite)
	for _, sf := range a.snowflakes {
		if sf.y > 0 && int(sf.y) < height && int(sf.x) >= 0 && int(sf.x) < width {
			a.screen.SetContent(int(sf.x), int(sf.y), sf.char, nil, snowStyle)
		}
	}

	// Title
	if a.sceneFrame > 40 {
		title := "❄  Mark24Code presents  ❄"
		titleStyle := tcell.StyleDefault.Foreground(tcell.ColorNavy).Background(tcell.ColorWhite).Bold(true)
		startX := (width - len(title)) / 2
		for i, ch := range title {
			a.screen.SetContent(startX+i, height/2-1, ch, nil, titleStyle)
		}
	}
}

func (a *App) drawSceneSanta() {
	// Night sky background
	skyStyle := tcell.StyleDefault.Background(tcell.ColorNavy)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			a.screen.SetContent(x, y, ' ', nil, skyStyle)
		}
	}

	// Stars twinkling
	starStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorNavy)
	for i := 0; i < 40; i++ {
		x := (i*7 + a.frame/3) % width
		y := (i * 11) % (height - 10)
		if (a.frame/10+i)%3 == 0 {
			a.screen.SetContent(x, y, '✦', nil, starStyle)
		} else {
			a.screen.SetContent(x, y, '·', nil, starStyle)
		}
	}

	// Ground with snow
	groundY := height - 8
	snowGroundStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorNavy)
	for x := 0; x < width; x++ {
		for dy := 0; dy < 3; dy++ {
			y := groundY + dy
			if y < height {
				snowChar := '~'
				if dy == 0 && (x+a.frame/5)%7 == 0 {
					snowChar = '❅'
				}
				a.screen.SetContent(x, y, snowChar, nil, snowGroundStyle)
			}
		}
	}

	// Draw houses on the ground
	houses := []int{10, 30, 55, 70} // X positions for houses
	houseStyle := tcell.StyleDefault.Foreground(tcell.ColorOlive).Background(tcell.ColorNavy)
	roofStyle := tcell.StyleDefault.Foreground(tcell.ColorMaroon).Background(tcell.ColorNavy)
	windowStyle := tcell.StyleDefault.Foreground(tcell.ColorYellow).Background(tcell.ColorNavy).Bold(true)

	for _, houseX := range houses {
		// Roof
		roofLines := []string{
			"  /\\  ",
			" /  \\ ",
			"/____\\",
		}
		for i, line := range roofLines {
			y := groundY - 6 + i
			if y >= 0 && y < height {
				for j, ch := range line {
					x := houseX + j
					if x >= 0 && x < width {
						a.screen.SetContent(x, y, ch, nil, roofStyle)
					}
				}
			}
		}

		// House body
		houseBody := []string{
			"|    |",
			"| [] |",
			"|    |",
		}
		for i, line := range houseBody {
			y := groundY - 3 + i
			if y >= 0 && y < height {
				for j, ch := range line {
					x := houseX + j
					if x >= 0 && x < width {
						style := houseStyle
						// Windows that twinkle
						if ch == '[' || ch == ']' {
							if (a.frame/15+houseX)%2 == 0 {
								style = windowStyle
							}
						}
						a.screen.SetContent(x, y, ch, nil, style)
					}
				}
			}
		}
	}

	// Santa and reindeer moving across the sky from right to left using emojis
	santaX := width - int(float64(a.sceneFrame)*0.8)
	santaY := 5 + int(math.Sin(float64(a.sceneFrame)*0.1)*1.5)

	emojiStyle := tcell.StyleDefault.Background(tcell.ColorNavy)

	// Reindeer team (4 reindeer) - now leading on the left side
	reindeerEmoji := "🦌"
	for i := 0; i < 4; i++ {
		x := santaX - 12 + i*3
		y := santaY
		if x >= 0 && x < width && y >= 0 && y < height {
			for j, ch := range reindeerEmoji {
				if x+j < width {
					a.screen.SetContent(x+j, y, ch, nil, emojiStyle)
				}
			}
		}
	}

	// Sleigh emoji
	sleighX := santaX + 2
	sleighEmoji := "🛷"
	if sleighX >= 0 && sleighX < width-2 && santaY >= 0 && santaY < height {
		for i, ch := range sleighEmoji {
			if sleighX+i < width {
				a.screen.SetContent(sleighX+i, santaY, ch, nil, emojiStyle)
			}
		}
	}

	// Santa emoji
	santaEmojiX := santaX + 4
	santaEmoji := "🎅"
	if santaEmojiX >= 0 && santaEmojiX < width-2 && santaY >= 0 && santaY < height {
		for i, ch := range santaEmoji {
			if santaEmojiX+i < width {
				a.screen.SetContent(santaEmojiX+i, santaY, ch, nil, emojiStyle)
			}
		}
	}

	// Magic sparkles trail - now trailing to the right
	trailStyle := tcell.StyleDefault.Foreground(tcell.ColorYellow).Background(tcell.ColorNavy)
	for i := 0; i < 10; i++ {
		x := santaX + 8 + i*2 // Trail to the right side
		y := santaY + (i % 2)
		if x >= 0 && x < width && y >= 0 && y < height {
			sparkle := []rune{'✦', '✧', '⋆', '*', '·'}[i%5]
			brightness := 10 - i
			style := trailStyle
			if brightness > 7 {
				style = style.Bold(true)
			} else if brightness < 4 {
				style = tcell.StyleDefault.Foreground(tcell.ColorOlive).Background(tcell.ColorNavy)
			}
			a.screen.SetContent(x, y, sparkle, nil, style)
		}
	}

	// Falling snow
	for _, sf := range a.snowflakes[:25] {
		if int(sf.y) < groundY && int(sf.x) >= 0 && int(sf.x) < width {
			style := tcell.StyleDefault.Foreground(tcell.ColorSilver).Background(tcell.ColorNavy)
			a.screen.SetContent(int(sf.x), int(sf.y), sf.char, nil, style)
		}
	}
}

func (a *App) drawSceneTree() {
	// Dark background
	bgStyle := tcell.StyleDefault.Background(tcell.ColorBlack)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			a.screen.SetContent(x, y, ' ', nil, bgStyle)
		}
	}

	// Snow on ground
	snowStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack)
	for x := 0; x < width; x++ {
		for dy := 0; dy < 2; dy++ {
			y := height - 3 + dy
			if y < height {
				a.screen.SetContent(x, y, '~', nil, snowStyle)
			}
		}
	}

	// Falling snow - more prominent for tree scene
	snowColors := []tcell.Color{tcell.ColorWhite, tcell.ColorSilver, tcell.ColorWhite}
	for i, sf := range a.snowflakes {
		if int(sf.y) < height-3 && int(sf.x) >= 0 && int(sf.x) < width {
			color := snowColors[i%len(snowColors)]
			style := tcell.StyleDefault.Foreground(color).Background(tcell.ColorBlack)
			// Make some snowflakes bold for emphasis
			if i%3 == 0 {
				style = style.Bold(true)
			}
			a.screen.SetContent(int(sf.x), int(sf.y), sf.char, nil, style)
		}
	}

	centerX := width / 2
	treeStartY := 4

	// Draw complete gorgeous tree immediately
	// Leave space for 3-layer trunk plus ground snow
	treeHeight := height - 10 // More space for trunk and snow
	a.drawTree2D(centerX, treeStartY, treeHeight)

	// Star on top (always visible)
	starStyle := tcell.StyleDefault.Foreground(tcell.ColorYellow).Background(tcell.ColorBlack).Bold(true)
	if (a.frame/5)%2 == 0 {
		starStyle = starStyle.Foreground(tcell.ColorGold)
	}
	a.screen.SetContent(centerX, treeStartY-1, '★', nil, starStyle)

	// Lights animation (always visible)
	a.drawTreeLights(centerX, treeStartY, treeHeight)

	// Colorful text at bottom with 1-second color changes
	colorfulColors := []tcell.Color{
		tcell.ColorRed,
		tcell.ColorYellow,
		tcell.ColorGreen,
		tcell.ColorBlue,
		tcell.ColorPurple,
		tcell.ColorFuchsia,
		tcell.ColorOrange,
		tcell.ColorPink,
	}

	// Change color every 20 frames (1 second)
	colorIndex := (a.frame / 20) % len(colorfulColors)

	if a.name == "Merry Christmas" {
		// Only show "Merry Christmas" if using default name
		text := "Merry Christmas"
		startX := centerX - len(text)/2
		for i, ch := range text {
			color := colorfulColors[(colorIndex+i)%len(colorfulColors)]
			style := tcell.StyleDefault.Foreground(color).Background(tcell.ColorBlack).Bold(true)
			a.screen.SetContent(startX+i, height-1, ch, nil, style)
		}
	} else {
		// Show two lines: "Merry Christmas" and "@{name}"
		line1 := "Merry Christmas"
		line2 := "@" + a.name

		// First line - centered
		startX1 := centerX - len(line1)/2
		for i, ch := range line1 {
			color := colorfulColors[(colorIndex+i)%len(colorfulColors)]
			style := tcell.StyleDefault.Foreground(color).Background(tcell.ColorBlack).Bold(true)
			a.screen.SetContent(startX1+i, height-2, ch, nil, style)
		}

		// Second line - centered
		startX2 := centerX - len(line2)/2
		for i, ch := range line2 {
			color := colorfulColors[(colorIndex+i+3)%len(colorfulColors)]
			style := tcell.StyleDefault.Foreground(color).Background(tcell.ColorBlack).Bold(true)
			a.screen.SetContent(startX2+i, height-1, ch, nil, style)
		}
	}
}

func (a *App) drawTree2D(centerX, startY, treeHeight int) {
	greenStyle := tcell.StyleDefault.Foreground(tcell.ColorGreen).Background(tcell.ColorBlack)
	whiteStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack)
	brownStyle := tcell.StyleDefault.Foreground(tcell.ColorOlive).Background(tcell.ColorBlack)

	// Simple triangular tree - like the Ruby version
	// Tree crown with * characters
	for row := 0; row < treeHeight; row++ {
		y := startY + row
		// Width increases as we go down: row 0 = 1 char, row 1 = 3 chars, row 2 = 5 chars, etc.
		rowWidth := 2*row + 1

		for col := 0; col < rowWidth; col++ {
			x := centerX - row + col
			if x >= 0 && x < width && y >= 0 && y < height {
				char := '*'
				style := greenStyle

				// Add some white snowflakes (30% of positions)
				// Use deterministic pseudo-random based on position and time
				// Change pattern every 1 second (20 frames) to match text animation
				colorChangeFrame := a.frame / 20
				positionSeed := (row*89 + col*67 + colorChangeFrame*43) % 10

				if positionSeed < 3 { // 30% chance
					char = '*'
					style = whiteStyle
				}

				a.screen.SetContent(x, y, char, nil, style)
			}
		}
	}

	// Simple trunk (like mWm in Ruby version) - 2 layers for visual
	trunkY := startY + treeHeight
	trunkHeight := 2
	for i := 0; i < trunkHeight; i++ {
		y := trunkY + i
		if y >= 0 && y < height {
			// Draw "mWm" pattern
			a.screen.SetContent(centerX-1, y, 'm', nil, brownStyle)
			a.screen.SetContent(centerX, y, 'W', nil, brownStyle)
			a.screen.SetContent(centerX+1, y, 'm', nil, brownStyle)
		}
	}
}

func (a *App) drawTreeLights(centerX, startY, treeHeight int) {
	colors := []tcell.Color{
		tcell.ColorRed,
		tcell.ColorBlue,
		tcell.ColorYellow,
		tcell.ColorPurple,
		tcell.ColorTeal,
		tcell.ColorFuchsia,
		tcell.ColorOrange,
		tcell.ColorPink,
	}

	// Store light positions from previous row to check vertical spacing
	type LightPos struct {
		row, col int
	}
	var allLights []LightPos

	// Draw lights on the tree with sparse distribution in both x and y directions
	// Some lights are always on, others blink with 1-second intervals
	for row := 0; row < treeHeight; row++ {
		y := startY + row
		rowWidth := 2*row + 1

		// Skip some rows entirely for better y-axis distribution
		rowSeed := row * 73
		skipRow := rowSeed%3 == 0 // Skip about 33% of rows

		if skipRow {
			continue
		}

		// Track previous column that had a light in this row to avoid horizontal adjacency
		lastLightCol := -3 // Start at -3 so first position can have a light

		// Go through each position in the row
		for col := 0; col < rowWidth; col++ {
			x := centerX - row + col

			if x >= 0 && x < width && y >= 0 && y < height {
				// Use deterministic pseudo-random based on position
				positionSeed := row*97 + col*53

				// Base chance of having a light (20% for sparse distribution)
				hasLight := positionSeed%10 < 2

				// Don't place light if it would be adjacent to previous light in same row
				if hasLight && (col - lastLightCol) <= 2 {
					hasLight = false
				}

				// Check if too close to lights in previous rows (vertical spacing)
				if hasLight {
					for _, prevLight := range allLights {
						rowDiff := row - prevLight.row
						colDiff := col - prevLight.col
						if colDiff < 0 {
							colDiff = -colDiff
						}

						// Don't place light if within 2 positions vertically and 2 positions horizontally
						if rowDiff <= 2 && colDiff <= 2 {
							hasLight = false
							break
						}
					}
				}

				if hasLight {
					lastLightCol = col
					allLights = append(allLights, LightPos{row: row, col: col})

					// Use row and column to determine light behavior
					lightSeed := row*31 + col*17

					// 40% of lights are always on, 60% blink
					isAlwaysOn := lightSeed%10 < 4

					// Color changes every 1 second (20 frames), matching text animation
					colorChangeFrame := (a.frame / 20) % len(colors)
					colorOffset := lightSeed % len(colors)
					color := colors[(colorChangeFrame+colorOffset)%len(colors)]

					var shouldShow bool
					if isAlwaysOn {
						// Always on lights
						shouldShow = true
					} else {
						// Blinking lights - blink every 1 second (20 frames)
						// Different lights blink at different offsets
						blinkOffset := lightSeed % 20
						blinkPhase := (a.frame + blinkOffset) % 40 // 2-second cycle

						// On for 1 second, off for 1 second
						shouldShow = blinkPhase < 20
					}

					if shouldShow {
						style := tcell.StyleDefault.
							Foreground(color).
							Background(tcell.ColorBlack).
							Bold(true) // All lights are bold for visibility

						// Use 'o' character for lights
						a.screen.SetContent(x, y, 'o', nil, style)
					}
				}
			}
		}
	}
}
