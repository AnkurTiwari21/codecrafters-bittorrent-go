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
func decodeBencode(bencodedString string) (interface{}, error) {
	if unicode.IsDigit(rune(bencodedString[0])) {
		var firstColonIndex int

		for i := 0; i < len(bencodedString); i++ {
			if bencodedString[i] == ':' {
				firstColonIndex = i
				break
			}
		}

		lengthStr := bencodedString[:firstColonIndex]

		length, err := strconv.Atoi(lengthStr)
		if err != nil {
			return "", err
		}

		return bencodedString[firstColonIndex+1 : firstColonIndex+1+length], nil
	} else if bencodedString[0] == 'i' {
		numString := bencodedString[1 : len(bencodedString)-1]
		num, err := strconv.Atoi(numString)
		if err != nil {
			return "", fmt.Errorf("error in conversion")
		}
		return num, nil
	} else if bencodedString[0] == 'l' {
		//it is a list
		list := []interface{}{}
		// fmt.Print("hi")
		for i := 1; i < len(bencodedString); i++ {
			if bencodedString[i] == 'i' {
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
				list = append(list, num)
				// fmt.Print("appending")
				// fmt.Print(num)
			}else if(unicode.IsDigit(rune(bencodedString[i]))){
				str := ""
				lenStr := ""
				for unicode.IsDigit(rune(bencodedString[i])){
					lenStr += string(bencodedString[i])
					i++;
				}
				i++; //skip in colon
				length,err:=strconv.Atoi(lenStr)
				if err!= nil {
					return "", fmt.Errorf("error in converting for string..")
				}
				for k:=0;k<length;k++{
					str += string(bencodedString[i])
					i++;
				}
				list = append(list, str)
				// fmt.Print("appending")
				// fmt.Print(str)
				i--;
			}
		}
		return list, nil
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

		decoded, err := decodeBencode(bencodedValue)
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
