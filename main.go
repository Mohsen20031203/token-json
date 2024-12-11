package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
)

var httpUrls = []string{
	"https://api.vultisig.com/1inch/swap/v6.0/1/tokens",
	"https://api.vultisig.com/1inch/swap/v6.0/43114/tokens",
	"https://api.vultisig.com/1inch/swap/v6.0/8453/tokens",
	"https://api.vultisig.com/1inch/swap/v6.0/81457/tokens",
	"https://api.vultisig.com/1inch/swap/v6.0/42161/tokens",
	"https://api.vultisig.com/1inch/swap/v6.0/137/tokens",
	"https://api.vultisig.com/1inch/swap/v6.0/10/tokens",
	"https://api.vultisig.com/1inch/swap/v6.0/56/tokens",
	"https://api.vultisig.com/1inch/swap/v6.0/25/tokens",
}

type Token struct {
	Address  string   `json:"address"`
	Symbol   string   `json:"symbol"`
	Decimals int      `json:"decimals"`
	Name     string   `json:"name"`
	LogoURI  string   `json:"logoURI"`
	EIP2612  bool     `json:"eip2612"`
	Tags     []string `json:"tags"`
	CmcID    string   `json:"CmcID"`
}

type TokenData struct {
	Tokens map[string]Token `json:"tokens"`
}

func fetchAndProcessURL(url, name string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var tokenData TokenData
	if err := json.Unmarshal(body, &tokenData); err != nil {
		return err
	}

	updatedBody, err := json.MarshalIndent(tokenData, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(name, updatedBody, 0644)
	if err != nil {
		fmt.Print("hello")
	}

	return nil
}

func main() {
	coout := 1
	for _, url := range httpUrls {
		nameFileJson := fmt.Sprintf("%s%d.json", path.Base(url), coout)
		if err := fetchAndProcessURL(url, nameFileJson); err != nil {
			fmt.Printf("Error processing URL %s: %v\n", url, err)
		}
		coout++
	}
}
