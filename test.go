package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

func main() {
	// Считываем весь ввод из стандартного потока
	data, err := io.ReadAll(os.Stdin)
	if err != nil || len(data) == 0 {
		fmt.Println("Not a quad function.")
		return
	}
	s := string(data)

	// Убираем лишний символ новой строки в конце, если он есть
	s = strings.TrimSuffix(s, "\n")

	lines := strings.Split(s, "\n")
	height := len(lines)
	if height == 0 {
		fmt.Println("Not a quad function.")
		return
	}

	width := len(lines[0])
	if width == 0 {
		fmt.Println("Not a quad function.")
		return
	}

	// Проверка на корректную ширину всех строк
	for _, line := range lines {
		if len(line) != width {
			fmt.Println("Not a quad function.")
			return
		}
	}

	var matches []string

	// Проверяем случай 1x1
	if width == 1 && height == 1 {
		switch s {
		case "o":
			matches = append(matches, fmt.Sprintf("quadA [%d] [%d]", width, height))
		case "/":
			matches = append(matches, fmt.Sprintf("quadB [%d] [%d]", width, height))
		case "A":
			matches = append(matches, fmt.Sprintf("quadC [%d] [%d]", width, height))
			matches = append(matches, fmt.Sprintf("quadD [%d] [%d]", width, height))
			matches = append(matches, fmt.Sprintf("quadE [%d] [%d]", width, height))
		}
	} else {
		// Проверяем общие случаи (больше 1x1)
		checkAllQuads(lines, width, height, &matches)
	}

	if len(matches) > 0 {
		sort.Strings(matches)
		fmt.Println(strings.Join(matches, " || "))
	} else {
		fmt.Println("Not a quad function.")
	}
}

// checkAllQuads выполняет все проверки, кроме случая 1x1.
func checkAllQuads(lines []string, width, height int, matches *[]string) {
	// quadA
	if lines[0][0] == 'o' && lines[0][width-1] == 'o' &&
		lines[height-1][0] == 'o' && lines[height-1][width-1] == 'o' &&
		checkHorizontalBorder(lines[0], '-', width) &&
		checkHorizontalBorder(lines[height-1], '-', width) &&
		checkVerticalBorders(lines, '|', width, height) &&
		checkInnerSpace(lines, width, height) {
		*matches = append(*matches, fmt.Sprintf("quadA [%d] [%d]", width, height))
	}

	// quadB
	if lines[0][0] == '/' && lines[0][width-1] == '\\' &&
		lines[height-1][0] == '\\' && lines[height-1][width-1] == '/' &&
		checkHorizontalBorder(lines[0], '*', width) &&
		checkHorizontalBorder(lines[height-1], '*', width) &&
		checkVerticalBorders(lines, '*', width, height) &&
		checkInnerSpace(lines, width, height) {
		*matches = append(*matches, fmt.Sprintf("quadB [%d] [%d]", width, height))
	}

	// quadC
	if lines[0][0] == 'A' && lines[0][width-1] == 'A' &&
		lines[height-1][0] == 'C' && lines[height-1][width-1] == 'C' &&
		checkHorizontalBorder(lines[0], 'B', width) &&
		checkHorizontalBorder(lines[height-1], 'B', width) &&
		checkVerticalBorders(lines, 'B', width, height) &&
		checkInnerSpace(lines, width, height) {
		*matches = append(*matches, fmt.Sprintf("quadC [%d] [%d]", width, height))
	}

	// quadD
	if lines[0][0] == 'A' && lines[0][width-1] == 'C' &&
		lines[height-1][0] == 'A' && lines[height-1][width-1] == 'C' &&
		checkHorizontalBorder(lines[0], 'B', width) &&
		checkHorizontalBorder(lines[height-1], 'B', width) &&
		checkVerticalBorders(lines, 'B', width, height) &&
		checkInnerSpace(lines, width, height) {
		*matches = append(*matches, fmt.Sprintf("quadD [%d] [%d]", width, height))
	}

	// quadE
	if lines[0][0] == 'A' && lines[0][width-1] == 'C' &&
		lines[height-1][0] == 'C' && lines[height-1][width-1] == 'A' &&
		checkHorizontalBorder(lines[0], 'B', width) &&
		checkHorizontalBorder(lines[height-1], 'B', width) &&
		checkVerticalBorders(lines, 'B', width, height) &&
		checkInnerSpace(lines, width, height) {
		*matches = append(*matches, fmt.Sprintf("quadE [%d] [%d]", width, height))
	}
}

// checkHorizontalBorder проверяет, что горизонтальные границы состоят из нужного символа.
func checkHorizontalBorder(line string, char rune, width int) bool {
	if width <= 2 {
		return true
	}
	for i := 1; i < width-1; i++ {
		if rune(line[i]) != char {
			return false
		}
	}
	return true
}

// checkVerticalBorders проверяет, что вертикальные границы состоят из нужного символа.
func checkVerticalBorders(lines []string, char rune, width, height int) bool {
	if height <= 2 {
		return true
	}
	for i := 1; i < height-1; i++ {
		if rune(lines[i][0]) != char || rune(lines[i][width-1]) != char {
			return false
		}
	}
	return true
}

// checkInnerSpace проверяет, что внутреннее пространство состоит из пробелов.
func checkInnerSpace(lines []string, width, height int) bool {
	if width <= 2 || height <= 2 {
		return true
	}
	for i := 1; i < height-1; i++ {
		for j := 1; j < width-1; j++ {
			if lines[i][j] != ' ' {
				return false
			}
		}
	}
	return true
}