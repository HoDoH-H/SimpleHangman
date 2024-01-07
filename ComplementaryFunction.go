package SimpleHangman

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
)

type Data struct {
	Life                      int
	AncientLetter             []string
	LetterFound               []string
	WordFound                 bool
	Word                      string
	FormatizedWord            string
	GameOver                  bool
	LetterTriedFormatizedText string
}

// Load the save
func LoadGame(save *Data) {
	data, _ := os.ReadFile("save/save.json")
	err := json.Unmarshal(data, save)
	if err != nil {
		log.Fatal(err)
	}
}

// Save your datas
func SaveGame(save *Data) {
	data, err := json.Marshal(save)
	if err != nil {
		fmt.Println(err)
	}
	err = os.WriteFile("save/save.json", data, 0644)
	if err != nil {
		os.Create("save/save.json")
		os.WriteFile("save/save.json", data, 0644)
	}
}

// Check if a string with a lengh greater than 1 is equal to the word we are searching
func CheckWord(try string, save *Data) bool {
	if len(try) == len(save.Word) {
		for e := range try {
			if string(try[e]) != string(save.Word[e]) {
				return false
			}
		}
		return true
	} else {
		return false
	}
}

// Check if a string with a lengh of 1 is inside the word we are searching
func CheckLetter(letter, word string) bool {
	for e := range word {
		if letter == string(word[e]) {
			return true
		}
	}
	return false
}

// Formatize the string (delete numbers, special characters, etc)
func FormatAns(ans string) string {
	realAns := ""
	if len(string(ans)) == 1 {
		if ans >= "a" && ans <= "z" {
			realAns = ans
		} else if ans >= "A" && ans <= "Z" {
			for i := range ans {
				realAns = string(rune(ans[i]) + 32)
			}
		}
	} else {
		for e := range ans {
			if string(ans[e]) >= "a" && string(ans[e]) <= "z" {
				realAns += string(ans[e])
			} else if string(ans[e]) >= "A" && string(ans[e]) <= "Z" {
				realAns += string(rune(ans[e]) + 32)
			}
		}
	}
	return realAns
}

// Get the word at a certain position inside the dictionary
func GetWord(data *Data, fileName string) {
	lineCount := 0
	file, _ := os.Open(fileName)
	fileScanner := bufio.NewScanner(file)
	for fileScanner.Scan() {
		lineCount++
	}
	file.Close()
	selectedNumber := rand.Intn(lineCount)
	file, _ = os.Open(fileName)
	fileScanner = bufio.NewScanner(file)

	line := 0
	for fileScanner.Scan() {
		if line == selectedNumber {
			data.Word = fileScanner.Text()
			return
		}
		line++
	}
	data.Word = ""
}

// Check if the player found the word
func IsWordDiscovered(save *Data) {
	for e := range save.Word {
		letterIsFound := false
		for ee := range save.LetterFound {
			if string(save.Word[e]) == string(save.LetterFound[ee]) {
				letterIsFound = true
			}
		}
		if !letterIsFound {
			save.WordFound = false
			return
		}
	}
	save.WordFound = true
}

// Check if the player already tried the letter
func IsLetterAlreadyTried(letter string, save *Data) bool {
	for i := range save.AncientLetter {
		if letter == save.AncientLetter[i] {
			return true
		}
	}
	return false
}

// Used to reduce the number of fmt.Println("") we have in the code, for better readability
func LineJump(n int) {
	for i := 0; i < n; i++ {
		fmt.Println("")
	}
}

// Print hangman drawing
func PrintDraw(save *Data) string {
	file, _ := os.Open("textFiles/hangman.txt")
	fileScanner := bufio.NewScanner(file)
	line := 1
	visual := ""
	if save.Life == 0 {
		return visual
	} else {
		for fileScanner.Scan() {
			if fileScanner.Text() == "" || fileScanner.Text() == " " {
				line++
			} else if line == save.Life {
				visual += "\n"
				visual += fileScanner.Text()
			}
		}
		return visual
	}
}

// Add all the new letter into the letterfound tab and anciantletter tab
func SplitWordToFindLetter(try string, save *Data) {
	validLetter := []string{}
	for i := range try {
		if CheckLetter(string(try[i]), save.Word) {
			validLetter = append(validLetter, string(try[i]))
		}
		newLetter := true
		for e := range save.AncientLetter {
			if string(try[i]) == save.AncientLetter[e] {
				newLetter = false
			}
		}
		if newLetter {
			save.AncientLetter = append(save.AncientLetter, string(try[i]))
		}
	}
	for i := range validLetter {
		newLetter := true
		for e := range save.LetterFound {
			if string(validLetter[i]) == save.LetterFound[e] {
				newLetter = false
			}
		}
		if newLetter {
			save.LetterFound = append(save.LetterFound, string(validLetter[i]))
		}
	}
}

// Just update the life left
func UpdateLife(try string, save *Data) {
	if len(try) == 1 {
		for e := range save.Word {
			if try == string(save.Word[e]) {
				return
			}
		}
		save.Life--
	} else if len(try) > 1 {
		if len(try) == len(save.Word) {
			for e := range try {
				isGood := false
				for ee := range save.Word {
					if try[e] == save.Word[ee] {
						isGood = true
					}
				}
				if !isGood {
					save.Life--
				}
			}
		} else {
			if len(try) < len(save.Word) {
				save.Life += len(save.Word) - len(try)
				for e := range try {
					isGood := false
					for ee := range save.Word {
						if try[e] == save.Word[ee] {
							isGood = true
						}
					}
					if !isGood {
						save.Life--
					}
				}
			} else {
				save.Life += len(try) - len(save.Word)
				for e := range save.Word {
					isGood := false
					for ee := range try {
						if try[ee] == save.Word[e] {
							isGood = true
						}
					}
					if !isGood {
						save.Life--
					}
				}
			}
		}
	}
}

// Print the word with '_' when a letter isn't found
func VisualWord(save *Data) {
	visual := ""
	for e := range save.Word {
		letterIsFound := false
		for ee := range save.LetterFound {
			if string(save.Word[e]) == string(save.LetterFound[ee]) {
				visual += string(rune(save.Word[e] - 32))
				letterIsFound = true
			}
		}
		if !letterIsFound {
			visual += "_"
		}
	}
	save.FormatizedWord = visual
}

// Get a formatize string of all the letters we found
func VisualLetterFound(save *Data) string {
	visual := ""
	for i := range save.LetterFound {
		visual += save.LetterFound[i]
		if i < len(save.LetterFound) {
			visual += " | "
		}
	}
	return visual
}

// Get a formatize string of all the letters we tried
func VisualLetterTried(save *Data) string {
	visual := ""
	for i := range save.AncientLetter {
		visual += save.AncientLetter[i]
		if i < len(save.AncientLetter) {
			visual += " | "
		}
	}
	return visual
}
