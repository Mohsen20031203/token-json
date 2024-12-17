package main

import (
	"encoding/json"
	"fmt"
	"io"
	"my-echo-app/utils"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Chain int

const (
	Ethereum Chain = iota
	Avalanche
	Base
	//Blast
	Arbitrum
	Polygon
	Optimism
	BscChain
	//CronosChain
)

func (c Chain) GetOneInchChainId() int {
	if c == Ethereum {
		return 1
	}
	if c == Avalanche {
		return 43114
	}
	if c == Base {
		return 8453
	}
	/*if c == Blast {
		return 81457
	}*/
	if c == Arbitrum {
		return 42161
	}
	if c == Polygon {
		return 137
	}
	if c == Optimism {
		return 10
	}
	if c == BscChain {
		return 56
	}
	/*
		if c == CronosChain {
			return 25
		}
	*/
	return 0
}

func (c Chain) GetCMCName() string {
	if c == Ethereum {
		return "Ethereum"
	}
	if c == Avalanche {
		return "Avalanche"
	}
	if c == Base {
		return "Base"
	}
	/*
		if c == Blast {
			return ?
		}
	*/
	if c == Arbitrum {
		return "Arbitrum"
	}
	if c == Polygon {
		return "Polygon"
	}
	if c == Optimism {
		return "Optimism"
	}
	if c == BscChain {
		return "BSC"
	}
	/*
		if c == CronosChain {
			return ?
		}
	*/
	return ""
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

func fetchAndProcessURL(ch Chain) (utils.TokenData, error) {

	url := fmt.Sprintf("https://api.vultisig.com/1inch/swap/v6.0/%v/tokens", ch.GetOneInchChainId())
	resp, err := http.Get(url)
	if err != nil {
		return utils.TokenData{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return utils.TokenData{}, err
	}

	var tokenData utils.TokenData
	if err := json.Unmarshal(body, &tokenData); err != nil {
		return utils.TokenData{}, err
	}

	return tokenData, nil
}
func main() {
	ch := []Chain{
		Ethereum,
		Avalanche,
		Base,
		//Blast,
		Arbitrum,
		Polygon,
		Optimism,
		BscChain,
		//CronosChain,
	}
	for _, nameChain := range ch {
		coine := nameChain
		var tokenData utils.TokenData
		var err error
		nameFile := fmt.Sprintf("tokens%v.json", coine.GetCMCName())
		file, err := os.OpenFile(nameFile, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			fmt.Println("Error opening/creating file:", err)
			return
		}
		defer file.Close()

		fileInfo, err := file.Stat()
		if err != nil {
			fmt.Println(err)
		}

		if fileInfo.Size() == 0 {
			tokenData, err = fetchAndProcessURL(coine)
			if err != nil {
				return
			}
		} else {
			utils.ReadFile(file, &tokenData)
		}
		var numberWrite int

		for key, value := range tokenData.Tokens {

			if value.CmcID != "" {
				continue
			}

			var cmcid int
			time.Sleep(time.Second * 1)

			for i := 0; i < 5; i++ {

				cmcid, err = fetchTokenPrice(coine, tokenData.Tokens[key].Address)
				if err != nil {

					if strings.Contains(err.Error(), "error for request api https://api.coinmarketcap.com") {
						fmt.Println(err)
						time.Sleep(time.Second * 20)
						if i == 4 {
							return
						}
						continue
					} else {
						fmt.Println(err)
						return

					}

				}
				break
			}
			fmt.Println(cmcid)

			token := tokenData.Tokens[key]
			token.CmcID = strconv.Itoa(cmcid)
			tokenData.Tokens[key] = token

			numberWrite++
			if numberWrite > 20 {
				numberWrite = 0
				utils.WriteFile(file, &tokenData)
			}
		}
		utils.WriteFile(file, &tokenData)
	}

}

func fetchTokenPrice(ch Chain, address string) (int, error) {
	var tokenData utils.TokenPairResponse

	apiUrl := fmt.Sprintf("https://api.coinmarketcap.com/dexer/v3/dexer/search/main-site?keyword=%s&all=false", address)
	resp, err := http.Get(apiUrl)
	if err != nil || resp.StatusCode != 200 {
		return 0, fmt.Errorf("error for request api %s", apiUrl)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	if err := json.Unmarshal(body, &tokenData); err != nil {
		fmt.Print(body)
		return 0, err
	}

	for _, pair := range tokenData.Data.Pairs {
		if pair.PlatformName == ch.GetCMCName() {
			if strings.EqualFold(pair.BaseToken.Address, address) {
				return pair.BaseToken.ID, nil
			}
			if strings.EqualFold(pair.QuoteToken.Address, address) {
				return pair.QuoteToken.ID, nil
			}
		}

	}
	return 0, nil

}
