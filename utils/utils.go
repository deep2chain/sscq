package utils

import (
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"strings"

	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/types"
)

// LoadGenesisDoc reads and unmarshals GenesisDoc from the given file.
func LoadGenesisDoc(cdc *amino.Codec, genFile string) (genDoc types.GenesisDoc, err error) {
	genContents, err := ioutil.ReadFile(genFile)
	if err != nil {
		return genDoc, err
	}

	if err := cdc.UnmarshalJSON(genContents, &genDoc); err != nil {
		return genDoc, err
	}

	return genDoc, err
}

//
func WriteString(filepath string, msg string) error {
	err := ioutil.WriteFile(filepath, []byte(msg), 0644)
	return err
}

//
func ReadString(filepath string, lineNum int) (line string, lastLine int, err error) {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
		return line, lastLine, err
	}
	defer file.Close()

	sc := bufio.NewScanner(file)
	for sc.Scan() {
		lastLine++
		if lastLine == lineNum {
			// you can return sc.Bytes() if you need output in []bytes
			return sc.Text(), lastLine, sc.Err()
		}
	}
	return line, lastLine, io.EOF
}

// junying-todo, 2020-01-15
// this is used to export accounts from accounts.text into genesis.json
// in: accounts text file path
// out: accounts, balances
func ReadAccounts(filepath string) (accounts []string, balances []int, err error) {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
		return nil, nil, err
	}
	defer file.Close()

	sc := bufio.NewScanner(file)
	linenum := 0
	for sc.Scan() {
		linenum++
		line := sc.Text()
		items := strings.Split(line, "	")
		// fmt.Print(items[0], ",", items[1], "\n")
		accounts = append(accounts, items[0])
		balance, err := strconv.Atoi(items[1])
		if err != nil {
			return accounts, balances, err
		}
		// fmt.Print(balance, "\n")
		balances = append(balances, balance)
	}
	return accounts, balances, sc.Err()
}
