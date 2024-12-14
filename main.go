package main

import (
	"encoding/json"
	"fmt"
	"io"
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
	Blast
	Arbitrum
	Polygon
	Optimism
	BscChain
	CronosChain
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
	if c == Blast {
		return 81457
	}
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
	if c == CronosChain {
		return 25
	}
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
	return ""
}

type TokenPairResponse struct {
	Data struct {
		Total int `json:"total"`
		Pairs []struct {
			PlatformID          int    `json:"platformId"`
			PlatformName        string `json:"platformName"`
			BaseTokenSymbol     string `json:"baseTokenSymbol"`
			QuoteTokenSymbol    string `json:"quoteTokenSymbol"`
			Liquidity           string `json:"liquidity"`
			PairContractAddress string `json:"pairContractAddress"`
			PlatFormCryptoID    int    `json:"platFormCryptoId"`
			ExchangeID          int    `json:"exchangeId"`
			PoolID              string `json:"poolId"`
			BaseTokenName       string `json:"baseTokenName"`
			MarketCap           string `json:"marketCap"`
			PriceUsd            string `json:"priceUsd"`
			PriceChange24H      string `json:"priceChange24h"`
			BaseToken           struct {
				ID           int    `json:"id,string"`
				Name         string `json:"name"`
				Address      string `json:"address"`
				Symbol       string `json:"symbol"`
				Slug         string `json:"slug"`
				CryptoSymbol string `json:"cryptoSymbol"`
				Decimals     int    `json:"decimals"`
			} `json:"baseToken"`
			QuoteToken struct {
				ID           int    `json:"id,string"`
				Name         string `json:"name"`
				Address      string `json:"address"`
				Symbol       string `json:"symbol"`
				Slug         string `json:"slug"`
				CryptoSymbol string `json:"cryptoSymbol"`
				Decimals     int    `json:"decimals"`
			} `json:"quoteToken"`
			Volume24H      string `json:"volume24h"`
			VolumeQuote24H string `json:"volumeQuote24h"`
			PlatformIcon   string `json:"platformIcon"`
		} `json:"pairs"`
	} `json:"data"`
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

func fetchAndProcessURL(ch Chain) (TokenData, error) {

	url := fmt.Sprintf("https://api.vultisig.com/1inch/swap/v6.0/%v/tokens", ch.GetOneInchChainId())
	resp, err := http.Get(url)
	if err != nil {
		return TokenData{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return TokenData{}, err
	}

	var tokenData TokenData
	if err := json.Unmarshal(body, &tokenData); err != nil {
		return TokenData{}, err
	}

	return tokenData, nil
}
func main() {
	coine := Ethereum
	var tokenData TokenData
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
		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&tokenData); err != nil {
			fmt.Println(err)
		}
	}

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
					time.Sleep(time.Second * 10)
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

		file.Seek(0, 0)
		file.Truncate(0)
		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(tokenData); err != nil {
			fmt.Println(err)
		}
	}

}

func fetchTokenPrice(ch Chain, address string) (int, error) {
	var tokenData TokenPairResponse

	apiUrl := fmt.Sprintf("https://api.coinmarketcap.com/dexer/v3/dexer/search/main-site?keyword=%s&all=false", address)
	resp, err := http.Get(apiUrl)
	if err != nil || resp.StatusCode >= 400 {
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
