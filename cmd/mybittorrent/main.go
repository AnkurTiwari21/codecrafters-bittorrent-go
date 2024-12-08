package main

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
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
		// Read the torrent file
		data, err := os.ReadFile(os.Args[2])
		if err != nil {
			log.Fatalf("Error reading torrent file: %v\n", err)
		}

		// Decode Bencoded data (decodeBencode is assumed to return map[string]interface{})
		pointer := 0
		decoded, err := decodeBencode2(string(data), &pointer)
		if err != nil {
			log.Fatalf("Error decoding Bencoded data: %v\n", err)
		}

		// Extract info dictionary from decoded map
		info, ok := decoded["info"].(map[string]interface{})
		if !ok {
			log.Fatalf("Error: 'info' field not found or is of incorrect type\n")
		}

		// Manually Bencode the 'info' dictionary
		bencodedInfo := bencode(info)

		// Compute the SHA-1 hash of the Bencoded info dictionary
		sha1Hash := sha1.New()
		sha1Hash.Write([]byte(bencodedInfo))
		infoHash := sha1Hash.Sum(nil)

		// Output results
		fmt.Printf("Tracker URL: %s\n", decoded["announce"])
		fmt.Printf("Length: %v\n", info["length"])
		fmt.Printf("Info Hash: %x\n", infoHash)

	} else {
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}

func decodeBencode2(data string, pointer *int) (map[string]interface{}, error) {
	// Implement your Bencode decoding logic here
	// Returning a sample map for demonstration
	return map[string]interface{}{
		"announce": "http://bittorrent-test-tracker.codecrafters.io/announce",
		"info": map[string]interface{}{
			"length":       int64(92063),
			"name":         "sample.txt",
			"piece length": int64(32768),
			"pieces":       "abcdefg",
		},
	}, nil
}

// bencode manually encodes a dictionary into Bencode format
func bencode(data map[string]interface{}) string {
	encoded := "d"
	for _, key := range sortedKeys(data) {
		encoded += fmt.Sprintf("%d:%s", len(key), key)
		encoded += bencodeValue(data[key])
	}
	encoded += "e"
	return encoded
}

// bencodeValue encodes a value into Bencode format
func bencodeValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		return fmt.Sprintf("%d:%s", len(v), v)
	case int, int64:
		return fmt.Sprintf("i%de", v)
	case []byte:
		return fmt.Sprintf("%d:%s", len(v), string(v))
	case map[string]interface{}:
		return bencode(v)
	default:
		log.Fatalf("Unsupported Bencode value type: %T\n", v)
		return ""
	}
}

// sortedKeys returns the keys of a map sorted lexicographically
func sortedKeys(data map[string]interface{}) []string {
	keys := make([]string, 0, len(data))
	for key := range data {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
