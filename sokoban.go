package main

import (
	"strings"
	"github.com/gdamore/tcell"
	"regexp"
	"errors"
	"bufio"
	"os"
	"log"
		"strconv"
	)

var mapState [][]rune
var gameWon, playing = false, false
var gameWindow tcell.Screen
var playerPos [2]int
var moves = 0

func logErrorIfNeeded(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func convertMapToString() (mapString string) {
	for _, row := range mapState {
		mapString += string(row)
	}
	return
}

func isGameWon() bool {
	mapString := convertMapToString()
	won := !strings.Contains(mapString, "+") && !strings.Contains(mapString, ".")
	file, fileErr := os.Open("/Users/james/projects/go-projects/sokoban/log.log")
	logErrorIfNeeded(fileErr)
	writer := bufio.NewWriter(file)
	writer.WriteString("won: ")
	writer.WriteString(strconv.FormatBool(won))
	writer.WriteString("\n")
	writer.Flush()
	file.Close()
	return won
}

func determineFuturePosition(direction rune, predictingCrate bool) (pos [2]int) {
	predictionSize := 1
	if predictingCrate {
		predictionSize = 2
	}

	switch direction {
	case 'w': pos = [2]int{playerPos[0], playerPos[1] - predictionSize}
	case 'a': pos = [2]int{playerPos[0] - predictionSize, playerPos[1]}
	case 's': pos = [2]int{playerPos[0], playerPos[1] + predictionSize}
	case 'd': pos = [2]int{playerPos[0] + predictionSize, playerPos[1]}
	}

	return pos
}

func processMovement(input rune) {
	futurePlayerPos := determineFuturePosition(input, false)

	futureChar := mapState[futurePlayerPos[0]][futurePlayerPos[1]]
	switch futureChar {
	case ' ': mapState[futurePlayerPos[0]][futurePlayerPos[1]] = '@'
	case '.': mapState[futurePlayerPos[0]][futurePlayerPos[1]] = '+'
	case 'o', '*':
		{
			futureCratePos := determineFuturePosition(input, true)
			futureCratePosChar := mapState[futureCratePos[0]][futureCratePos[1]]
			switch futureCratePosChar{
			case ' ': mapState[futureCratePos[0]][futureCratePos[1]] = 'o'
			case '.': mapState[futureCratePos[0]][futureCratePos[1]] = '*'
			default: return
			}
			switch futureChar {
			case 'o': mapState[futurePlayerPos[0]][futurePlayerPos[1]] = '@'
			case '*': mapState[futurePlayerPos[0]][futurePlayerPos[1]] = '+'
			}

		}
	default: return
	}

	playerPosChar := mapState[playerPos[0]][playerPos[1]]
	switch playerPosChar {
	case '@': mapState[playerPos[0]][playerPos[1]] = ' '
	case '.': mapState[playerPos[0]][playerPos[1]] = '+'
	case '+': mapState[playerPos[0]][playerPos[1]] = '.'
	}

	playerPos = futurePlayerPos
	moves++
}

func processEvents() {
	if event := gameWindow.PollEvent(); event != nil {
		switch t := event.(type) {
		case *tcell.EventError, *tcell.EventInterrupt: gameWindow.Fini()
		case *tcell.EventKey:
			switch t.Key() {
			case tcell.KeyCtrlC: gameWindow.Fini()
			case tcell.KeyRune:
				switch t.Rune() {
				case 'w', 'a', 's', 'd':
					processMovement(t.Rune())
				}
			}
		}
	}

}

func play() {
	playing = true
	for playing {
		if !isGameWon() {
			processEvents()
			drawMap()
			file, fileErr := os.Open("/Users/james/projects/go-projects/sokoban/log.log")
			logErrorIfNeeded(fileErr)
			file.WriteString("main loop\n")
			file.Sync()
			file.Close()
		} else {
			file, fileErr := os.Open("log.log")
			logErrorIfNeeded(fileErr)
			writer := bufio.NewWriter(file)
			writer.WriteString("in else!\n")
			writer.Flush()
			file.Close()
			if event := gameWindow.PollEvent(); event != nil {
				switch t := event.(type) {
				case *tcell.EventError, *tcell.EventInterrupt: gameWindow.Fini()
				case *tcell.EventKey:
					switch t.Key() {
					case tcell.KeyCtrlC:
						gameWindow.Fini()
					}
				default: continue
				}
			}
			mapState[40][50] = 'C'
			mapState[41][50] = 'o'
			mapState[42][50] = 'n'
			mapState[43][50] = 'g'
			mapState[44][50] = 'r'
			mapState[45][50] = 'a'
			mapState[46][50] = 't'
			mapState[47][50] = 'u'
			mapState[48][50] = 'l'
			mapState[49][50] = 'a'
			mapState[50][50] = 't'
			mapState[51][50] = 'i'
			mapState[52][50] = 'o'
			mapState[53][50] = 'n'
			mapState[54][50] = 's'
			mapState[55][50] = '!'
			drawMap()
		}
	}
}

func processMenuInput(input string) (string, error) {
	if ok, _ := regexp.MatchString("", input); ok {
		return "maps/level_" + input + ".map", nil
	} else {
		return "", errors.New("invalid input given")
	}
}

func drawMap() {
	gameWindow.Clear()
	for rowI, row := range mapState {
		for cellI, cell := range row {
			gameWindow.SetContent(rowI, cellI, rune(cell), nil, tcell.StyleDefault)
		}
	}
	gameWindow.Show()
}

func populateMapState() {
	moves = 0
	mapState = make([][]rune, 100)
	for i := range mapState {
		mapState[i] = make([]rune, 100)
		for j := range mapState[i] {
			mapState[i][j] = ' '
		}
	}
	maxWidth, curLine := 0, 0

	path, menuErr := processMenuInput("1")
	logErrorIfNeeded(menuErr)

	file, fileErr := os.Open(path)
	logErrorIfNeeded(fileErr)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
		for i, char := range line {
			if char == 0 {
				continue
			}
			mapState[i][curLine] = char
			if char == '@' {
				playerPos = [2]int{i, curLine}
			}
		}
		curLine++
	}
	//resizeMapState(maxWidth, curLine)
}

func startup() {
	var err error
	gameWindow, err = tcell.NewScreen()
	if err != nil {
		panic(err)
	}
	gameWindow.Init()
	populateMapState()
	gameWindow.Show()
}

func main() {
	startup()
	play()
}