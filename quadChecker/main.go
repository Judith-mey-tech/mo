package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	inputBytes, err := os.ReadFile("/dev/stdin")
	if err != nil {
		fmt.Println("Not a quad function")
		return
	}
	input := string(inputBytes)

	quadNames := []string{"quadA", "quadB", "quadC", "quadD", "quadE"}
	matches := []string{}

	for _, quad := range quadNames {
		// Проверяем наличие исполняемого файла
		if _, err := os.Stat("./" + quad); err != nil {
			continue
		}

		// Пробуем все разумные размеры: от 1 до 100 по ширине и высоте
		for height := 1; height <= 100; height++ {
			for width := 1; width <= 100; width++ {
				cmd := exec.Command("./"+quad, fmt.Sprint(width), fmt.Sprint(height))
				output, err := cmd.Output()
				if err != nil {
					continue
				}
				if string(output) == input {
					result := fmt.Sprintf("[%s] [%d] [%d]", quad, width, height)
					matches = append(matches, result)
					goto NextQuad
				}
			}
		}
	NextQuad:
	}

	if len(matches) == 0 {
		fmt.Println("Not a quad function")
		return
	}

	// Сортировка по алфавиту — без sort, руками (bubble sort)
	for i := 0; i < len(matches); i++ {
		for j := i + 1; j < len(matches); j++ {
			if matches[i] > matches[j] {
				matches[i], matches[j] = matches[j], matches[i]
			}
		}
	}

	fmt.Println(strings.Join(matches, " || "))
}