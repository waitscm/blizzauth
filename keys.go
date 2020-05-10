package blizzauth

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
)

type keys struct {
	id     string
	secret string
	name   string
}

// KeyDir Directory in the user's $HOME
const KeyDir = ".blizzard"

func newKeys(name string) *keys {

	id, err := readKeyFromFile(getIDFilePath(name))
	if err != nil {
		log.Println(err)
		return nil
	}

	secret, err := readKeyFromFile(getSecretFilePath(name))
	if err != nil {
		log.Println(err)
		return nil
	}

	return &keys{
		id:     id,
		secret: secret,
		name:   name,
	}
}

func getIDFilePath(name string) string {
	return fmt.Sprintf("%v/%v/%v.id", os.Getenv("HOME"), KeyDir, name)
}

func getSecretFilePath(name string) string {
	return fmt.Sprintf("%v/%v/%v.secret", os.Getenv("HOME"), KeyDir, name)
}

func readKeyFromFile(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		log.Println(err)
		return "", nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			return line, nil
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return "", errors.New("No key found in " + filename)
}
