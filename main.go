package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
)

var httpUrl = []string{
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

func main() {

	for _, url := range httpUrl {

		jsonToken, err := http.Get(url)
		if err != nil {
			return
		}

		body, err := io.ReadAll(jsonToken.Body)
		if err != nil {

			return
		}

		var tokenData TokenData
		err = json.Unmarshal(body, &tokenData)
		if err != nil {

			return
		}
		for address, token := range tokenData.Tokens {
			token.CmcID = ""
			tokenData.Tokens[address] = token
		}

		updatedBody, err := json.MarshalIndent(tokenData, "", "  ")
		if err != nil {
			return
		}

		err = os.WriteFile("tokens_output.json", updatedBody, 0644)
		if err != nil {
			return
		}
	}

}
