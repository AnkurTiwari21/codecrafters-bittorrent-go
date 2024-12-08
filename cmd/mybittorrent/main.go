package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"unicode"
	// bencode "github.com/jackpal/bencode-go" // Available if you need it!
)

// Ensures gofmt doesn't remove the "os" encoding/json import (feel free to remove this!)
var _ = json.Marshal

// Example:
// - 5:hello -> hello
// - 10:hello12345 -> hello12345
func decodeBencode(bencodedString string, pointer *int) (interface{}, error) {
	if unicode.IsDigit(rune(bencodedString[*pointer])) {
		// fmt.Print("in string")
		// fmt.Print("pointer is ")
		// fmt.Print(*pointer)
		var firstColonIndex int

		for i := *pointer; i < len(bencodedString)+(*pointer); i++ {
			if bencodedString[i] == ':' {
				firstColonIndex = i
				// fmt.Print("first colon ")
				// fmt.Print(firstColonIndex)
				break
			}
		}

		lengthStr := bencodedString[*pointer:firstColonIndex]

		length, err := strconv.Atoi(lengthStr)
		if err != nil {
			return "", err
		}
		*pointer = firstColonIndex + 1 + length
		return bencodedString[firstColonIndex+1 : firstColonIndex+1+length], nil
	} else if bencodedString[*pointer] == 'i' {
		// fmt.Print("in integer")
		// fmt.Print("pointer is ")
		// fmt.Print(*pointer)
		i := *pointer
		numStr := ""
		i++
		for bencodedString[i]-'0' >= 0 && bencodedString[i]-'0' <= 9 || bencodedString[i] == '-' {
			numStr += string(bencodedString[i])
			i++
		}
		i--
		num, err := strconv.Atoi(numStr)
		if err != nil {
			return "", fmt.Errorf("error in conversion..")
		}
		*pointer = i + 1
		return num, nil
	} else if bencodedString[*pointer] == 'l' {
		//it is a list
		// fmt.Print("in list")
		// fmt.Print("pointer is ")
		// fmt.Print(*pointer)
		list := []interface{}{}
		*pointer = *pointer + 1
		for *pointer < len(bencodedString) {
			result, err := decodeBencode(bencodedString, pointer)
			// fmt.Print("\n returned with ")
			// fmt.Print(result)
			// fmt.Print("\n current index is ")
			// fmt.Print(*pointer)
			if err != nil {
				return "", fmt.Errorf("error in decoding string | err", err)
			}
			if result != "" {
				list = append(list, result)
			}
		}
		return list, nil
	} else if bencodedString[*pointer] == 'e' {
		*pointer = *pointer + 1
		return "", nil
	} else {
		return "", fmt.Errorf("Only strings are supported at the moment")
	}
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Fprintln(os.Stderr, "Logs from your program will appear here!")

	command := os.Args[1]

	if command == "decode" {
		// Uncomment this block to pass the first stage
		//
		bencodedValue := os.Args[2]
		pointer := 0
		decoded, err := decodeBencode(bencodedValue, &pointer)
		if err != nil {
			fmt.Println(err)
			return
		}

		jsonOutput, _ := json.Marshal(decoded)
		fmt.Println(string(jsonOutput))
	} else {
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}

// lli798e6:bananaee

// 1. [] lli798e6:bananaee
// 2. [798,] li798e6:bananaee

// 3. returns 798 and 6:bananaee
