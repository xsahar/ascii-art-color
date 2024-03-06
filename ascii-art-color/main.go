package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var validBanners = []string{
	"standard",
	"shadow",
	"thinkertoy",
}

func main() {
	if len(os.Args) == 2 {
		st:=os.Args[1]
		fmt.Print(convert(st))
		os.Exit(0)
	}
	hasFlagOption := false
	for _, arg := range os.Args {
		switch {
		case strings.HasPrefix(arg, "--") :
			hasFlagOption = true

			switch {
			case strings.HasPrefix(arg, "--color") && !strings.Contains(arg, "="):
				fmt.Println("Usage: go run . [OPTION] [STRING]")
				fmt.Println("EX: go run . --color=<color> <letters to be colored> \"something\"")
				os.Exit(0)
			case strings.HasPrefix(arg, "--output") && !strings.Contains(arg, "="):
				fmt.Println("Usage: go run . [OPTION] [STRING] [BANNER]")
				fmt.Println("EX: go run . --output=<fileName.txt> something standard")
				os.Exit(0)
			case strings.HasPrefix(arg, "--fs") && !strings.Contains(arg, "="):
				fmt.Println("Usage: go run . [STRING] [BANNER]")
				fmt.Println("EX: go run . something standard")
				os.Exit(0)
			}
		}
	}

	if !hasFlagOption && (len(os.Args) != 3 || !isValidBanner(os.Args[2])) {
		fmt.Println("Usage: go run . [STRING] [BANNER]")
		fmt.Println("EX: go run . something standard")
		os.Exit(0)
	}
	


	colorFlag := flag.String("color", "", "Color for ASCII art")
	outputFlag := flag.String("output", "", "File path to output the ASCII art")
	flag.Parse()

	switch {
	case *outputFlag != "":
		Process2(flag.Args(), *outputFlag)
	 case *colorFlag != "":
		args := flag.Args()
		if len(args) == 0 {
			fmt.Println("Please provide at least one argument: the string to print.")
			return
		}

		colorWord := args[0]
		sentence := colorWord
		if len(args) > 1 {
			substring := args[0]
			
			fullstring := args[1]
			
			
			if strings.HasPrefix(string(fullstring), substring)==false {
				fmt.Print("Cannot found "+substring+" in "+fullstring)
				return
			}
		}			

		if len(args) > 1 {
			sentence = strings.Join(args[1:], " ")
		}

		colors := strings.Split(*colorFlag, ",")
		err := Process(sentence, colors, "", colorWord)
		if err != nil {
			fmt.Println(err)
			return
		}
	default:
		args := flag.Args()
		if len(args) > 0 {
			inputString := args[0]
			banner := "standard" // set default banner

		

			if len(args) > 1 && isValidBanner(args[1]) {
				banner = args[1] // if banner is specified and valid, use it
			}

			err :=Process1(inputString, banner)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else {
			fmt.Println("Please provide a valid option or a string to generate ASCII art from.")
			os.Exit(1)
		}
	}
}

func isValidBanner(banner string) bool {
	for _, validBanner := range validBanners {
		if banner == validBanner {
			return true
		}
	}
	return false
}


func Process1(input, banner string) error {
	for _, r := range input {
		if r < ' ' || r > '~' {
			return fmt.Errorf("invalid character: %c", r)
		}
	}

	// read the file from the fonts folder
	bytes, err := os.ReadFile(fmt.Sprintf("fonts/%s.txt", banner))
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	var lines []string
	if banner == "thinkertoy" {
		lines = strings.Split(string(bytes), "\r\n")
	} else {
		lines = strings.Split(string(bytes), "\n")
	}

	var arr []rune
	Newline := false
	

	if banner == "standard" || banner == "shadow" || banner == "thinkertoy" {
		for i, r := range input {
			if Newline {
				Newline = false
				art(arr, lines)
				arr = []rune{}
				continue
			}
			if r == '\\' && len(input) != i+1 {
				if input[i+1] == 'n' {
					Newline = true
					continue
				}
			}
			arr = append(arr, r)
		}
		art(arr, lines)

	} 
	return nil
}

func printArt(arr []rune, lines []string, lineCount int, offset int) {
	if len(arr) != 0 {
		for line := 1; line <= lineCount; line++ {
			for _, r := range arr {
				skip := (r * rune(lineCount)) - rune(offset)
				fmt.Print(lines[line+int(skip)])
			}
			fmt.Println()
		}
	} else {
		fmt.Println()
	}
}

func art(arr []rune, lines []string) {
	if len(arr) != 0 {
		for line := 1; line <= 8; line++ {
			for _, r := range arr {
				skip := (r - 32) * 9
				fmt.Print(lines[line+int(skip)])
			}
			fmt.Println()
		}
	} else {
		fmt.Println()
	}
}

const (
	colorReset = "\033[0m"
)

func Process(input string, colors []string, banner string, colorWords string) error {
	var selectedColors []string

	for _, color := range colors {
		var selectedColor string
		if strings.HasPrefix(color, "rgb(") && strings.HasSuffix(color, ")") {
			rgbValues := strings.TrimPrefix(color, "rgb(")
			rgbValues = strings.TrimSuffix(rgbValues, ")")
			rgbValuesSplit := strings.Split(rgbValues, ",")
			if len(rgbValuesSplit) != 3 {
				return fmt.Errorf("invalid RGB color value")
			}
			r, err := strconv.Atoi(strings.TrimSpace(rgbValuesSplit[0]))
			if err != nil {
				return fmt.Errorf("invalid RGB color value: %v", err)
			}
			g, err := strconv.Atoi(strings.TrimSpace(rgbValuesSplit[1]))
			if err != nil {
				return fmt.Errorf("invalid RGB color value: %v", err)
			}
			b, err := strconv.Atoi(strings.TrimSpace(rgbValuesSplit[2]))
			if err != nil {
				return fmt.Errorf("invalid RGB color value: %v", err)
			}
			selectedColor = fmt.Sprintf("\u001b[38;2;%d;%d;%dm", r, g, b)
		} else {
			switch color {
			case "red":
				selectedColor = "\u001b[38;2;255;0;0m"
			case "green":
				selectedColor = "\u001b[38;2;0;255;0m"
			case "yellow":
				selectedColor = "\u001b[38;2;255;255;0m"
			case "blue":
				selectedColor = "\033[34m"
			case "purple":
				selectedColor = "\u001b[38;2;161;32;255m"
			case "cyan":
				selectedColor = "\u001b[38;2;0;183;235m"
			case "white":
				selectedColor = "\u001b[38;2;255;255;255m"
			case "pink":
				selectedColor = "\u001b[38;2;255;0;255m"
			case "grey":
				selectedColor = "\u001b[38;2;128;128;128m"
			case "black":
				selectedColor = "\u001b[38;2;0;0;0m"
			case "brown":
				selectedColor = "\u001b[38;2;160;128;96m"
			case "orange":
				selectedColor = "\u001b[38;2;255;160;16m"
			default:
				selectedColor = colorReset
			}
		}

		// selectedColor = ... whatever it is ...
		selectedColors = append(selectedColors, selectedColor)
	}
	// Create color queue
	colorQueue := NewQueue()

	// Push all colors into queue
	for _, color := range selectedColors {
		colorQueue.Push(color)
	}

	// Default banner
	if banner == "" {
		banner = "standard"
	}

	// Read the art template file
	bytes, err := os.ReadFile(fmt.Sprintf("fonts/%s.txt", banner))
	if err != nil {
		return err
	}

	// Split the lines based on banner type
	var lines []string
	if banner == "thinkertoy" {
		lines = strings.Split(string(bytes), "\r\n")
	} else {
		lines = strings.Split(string(bytes), "\n")
	}

	colorWordSlice := strings.Split(colorWords, ",")

	// Create ASCII Art
	createArt(input, colorQueue, colorWordSlice, lines)
	return nil
}

func createArt(input string, colorQueue *Queue, colorWords []string, template []string) {
	for line := 1; line <= 8; line++ {
		remainingInput := input
		for remainingInput != "" {
			colorMatch := ""
			color := colorReset
			for _, word := range colorWords {
				if strings.HasPrefix(remainingInput, word) {
					// Found a match. Apply color and move on.
					colorMatch = word
					color = colorQueue.Pop().(string)
					colorQueue.Push(color) // Requeue the color
					break
				}
			}
			if colorMatch == "" {
				// No match found. Print without color and move on.
				fmt.Print(template[line+int(remainingInput[0]-32)*9], colorReset)
				remainingInput = remainingInput[1:]
			} else {
				for _, r := range colorMatch {
					fmt.Print(color, template[line+int(r-32)*9], colorReset)
				}
				remainingInput = strings.TrimPrefix(remainingInput, colorMatch)
			}
		}
		fmt.Println()
	}
}

type Queue struct {
	items []interface{}
}

func NewQueue() *Queue {
	return &Queue{}
}

func (q *Queue) Push(item interface{}) {
	q.items = append(q.items, item)
}

func (q *Queue) Pop() interface{} {
	if len(q.items) == 0 {
		return nil
	}

	item := q.items[0]
	q.items = q.items[1:]

	return item
}
func convert(st string) string {
	n := ""
	data, err := os.ReadFile("standard.txt")
	if err != nil {
		fmt.Println("failed to read file:", err)
	}
	input := st
	input = strings.Replace(input, "\\n", "\n", -1)
	words := strings.Split(input, "\n")
	lines := strings.Split(string(data), "\n")

	for _, word := range words {
		if word == "" {
			fmt.Println()
			continue
		}

		letters := strings.Split(word, "")
		var ascii []int
		for _, letter := range letters {
			l := int([]rune(letter)[0])
			ascii = append(ascii, l)
		}
		for j := 1; j < 9; j++ {
			str := ""
			for _, val := range ascii {
				line := (val - 32) * 9
				if line+j >= len(lines) {
					fmt.Println("Error: insufficient amount of words", word)
					return ""
				}
				str += lines[line+j]
			}

			n += str + "\n"

		}

	}
	return n
}