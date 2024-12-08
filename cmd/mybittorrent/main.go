package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"unicode"

	"github.com/codecrafters-io/bittorrent-starter-go/cmd/mybittorrent/model"
	// bencode "github.com/jackpal/bencode-go" // Available if you need it!
)

// Ensures gofmt doesn't remove the "os" encoding/json import (feel free to remove this!)
var _ = json.Marshal

// Example:
// - 5:hello -> hello
// - 10:hello12345 -> hello12345
func decodeBencode(bencodedString string, pointer *int) (interface{}, error) {
	if unicode.IsDigit(rune(bencodedString[*pointer])) {
		var firstColonIndex int

		for i := *pointer; i < len(bencodedString)+(*pointer); i++ {
			if bencodedString[i] == ':' {
				firstColonIndex = i
				break
			}
		}

		lengthStr := bencodedString[*pointer:firstColonIndex]
		//3:123
		length, err := strconv.Atoi(lengthStr)
		if err != nil {
			return "", err
		}
		*pointer = firstColonIndex + 1 + length
		return bencodedString[firstColonIndex+1 : firstColonIndex+1+length], nil
	} else if bencodedString[*pointer] == 'i' {
		i := *pointer
		numStr := ""
		i++
		for bencodedString[i]-'0' >= 0 && bencodedString[i]-'0' <= 9 || bencodedString[i] == '-' {
			numStr += string(bencodedString[i])
			i++
		}
		num, err := strconv.Atoi(numStr)
		if err != nil {
			return "", fmt.Errorf("error in conversion..")
		}
		*pointer = i + 1
		return num, nil
	} else if bencodedString[*pointer] == 'l' {
		list := []interface{}{}
		*pointer = *pointer + 1
		for *pointer < len(bencodedString) {
			result, err := decodeBencode(bencodedString, pointer)
			if err != nil {
				return "", fmt.Errorf("error in decoding string | err", err)
			}
			if result != "" {
				list = append(list, result)
			} else {
				return list, nil
			}
		}
		return list, nil
	} else if bencodedString[*pointer] == 'd' {
		dict := map[string]interface{}{}
		*pointer = *pointer + 1
		item := 0
		var key string
		var value interface{}
		for *pointer < len(bencodedString) {
			result, err := decodeBencode(bencodedString, pointer)
			if err != nil {
				return "", fmt.Errorf("error in decoding string | err", err)
			}
			// fmt.Print("returned with ",result)
			// fmt.Print("current pointer ",*pointer)

			if result != "" {
				if item == 0 {
					item = 1
					key = result.(string)
				} else {
					item = 0
					value = result
					dict[key] = value
				}
			} else {
				// fmt.Print("here1")
				return dict, nil
			}
			// fmt.Print(dict)
			// fmt.Print("\n")
		}
		// fmt.Print("here")
		return dict, nil
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
	} else if command == "info" {
		//read the file assigned in command line
		data, err := os.ReadFile(os.Args[2])
		if err != nil {
			fmt.Println("error in opening file | err", err)
			return
		}
		pointer := 0
		decoded, err := decodeBencode(string(data), &pointer)
		if err != nil {
			fmt.Println(err)
			return
		}

		jsonOutput, _ := json.Marshal(decoded)
		var FileData model.File
		err = json.Unmarshal(jsonOutput, &FileData)
		if err != nil {
			fmt.Errorf("error is ", err)
			return
		}
		fmt.Println("Tracker URL: ",FileData.Announce)
		fmt.Println("Length: ",FileData.Info.Length)
	} else {
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}

// lli798e6:bananaee

// 1. [] lli798e6:bananaee
// 2. [798,] li798e6:bananaee

// 3. returns 798 and 6:bananaee
