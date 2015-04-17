package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/codegangsta/cli"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type Contact struct {
	Title   string
	Address string
}

func main() {
	app := cli.NewApp()
	app.Name = "blackbook"
	app.Usage = "Keep track of all your ssh information"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "new",
			Value: "",
			Usage: "Creates a new entry into blackbook",
		},
		cli.StringFlag{
			Name:  "del, delete",
			Value: "",
			Usage: "Deletes a contact",
		},
	}
	app.Action = func(c *cli.Context) {
		if len(c.Args()) == 1 {
			title := c.Args()[0]
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Password: ")
			key, err := reader.ReadString('\n')
			if err != nil {
				fmt.Print("Sorry something went wrong")
			}
			c, err := loadContact(title, key)
			if err != nil {
				fmt.Println("Could not find that contact")
			}
			fmt.Print(c.Address)
			exec.Command("/bin/sh", "ssh "+c.Address)
		}
		if c.String("new") != "" {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Address(e.x. john@127.0.0.1): ")
			address, err := reader.ReadString('\n')
			if err != nil {
				fmt.Print("Sorry something went wrong")
			}
			fmt.Print("Password: ")
			key, err := reader.ReadString('\n')
			if err != nil {
				fmt.Print("Sorry something went wrong")
			}
			c := &Contact{Title: c.String("new"), Address: address}
			c.save(key)
			return
		} else if c.String("del") != "" {
			err := os.Remove("contacts/" + c.String("del") + ".txt")
			if err != nil {
				fmt.Println(err)
				return
			}
			return
		}
	}
	app.Run(os.Args)
}

func (c *Contact) save(key string) error {
	filename := c.Title + ".txt"
	cryptoAddress, _ := encrypt([]byte(padKey(key)), []byte(c.Address))
	return ioutil.WriteFile("contacts/"+filename, cryptoAddress, 0600)
}

func loadContact(title string, key string) (*Contact, error) {
	filename := title + ".txt"
	cryptoAddress, err := ioutil.ReadFile("contacts/" + filename)
	if err != nil {
		fmt.Println("Sorry, can't find that file!")
	}
	address, _ := decrypt([]byte(padKey(key)), cryptoAddress)
	return &Contact{Title: title, Address: string(address)}, nil
}

func encrypt(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	b := base64.StdEncoding.EncodeToString(text)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))
	return ciphertext, nil
}

//Decrypt function. Takes in a key and encrypted text. Returns the decrypted text and err
func decrypt(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(text) < aes.BlockSize {
		return nil, errors.New("cipher text too short")
	}

	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)
	data, err := base64.StdEncoding.DecodeString(string(text))
	if err != nil {
		return nil, err
	}
	return data, nil
}

func padKey(key string) string {
	length := len(key)
	if length == 16 || length == 24 || length == 32 {
		return key
	} else if length < 16 {
		padLen := 16 - length
		padding := strings.Repeat("f", padLen)
		return key + padding
	} else if length < 24 {
		padLen := 24 - length
		padding := strings.Repeat("f", padLen)
		return key + padding
	} else if length < 32 {
		padLen := 32 - length
		padding := strings.Repeat("f", padLen)
		return key + padding
	} else {
		return key[:32]
	}
}
