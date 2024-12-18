package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

type NullableField struct {
	Value string
	Set   bool
}

func (n *NullableField) UnmarshalJSON(data []byte) error {
	n.Set = true
	if string(data) == "null" || len(data) == 0 {
		n.Value = ""
	} else {
		n.Value = string(data)
	}
	return nil
}

type TokenPairResponse struct {
	Data struct {
		Total int `json:"total"`
		Pairs []struct {
			PlatformID          NullableField `json:"platformId"`
			PlatformName        string        `json:"platformName"`
			BaseTokenSymbol     string        `json:"baseTokenSymbol"`
			QuoteTokenSymbol    string        `json:"quoteTokenSymbol"`
			Liquidity           string        `json:"liquidity"`
			PairContractAddress string        `json:"pairContractAddress"`
			PlatFormCryptoID    int           `json:"platFormCryptoId"`
			ExchangeID          int           `json:"exchangeId"`
			PoolID              string        `json:"poolId"`
			BaseTokenName       string        `json:"baseTokenName"`
			MarketCap           string        `json:"marketCap"`
			PriceUsd            string        `json:"priceUsd"`
			PriceChange24H      string        `json:"priceChange24h"`
			BaseToken           struct {
				ID           NullableField `json:"id,string"`
				Name         string        `json:"name"`
				Address      string        `json:"address"`
				Symbol       string        `json:"symbol"`
				Slug         string        `json:"slug"`
				CryptoSymbol string        `json:"cryptoSymbol"`
				Decimals     int           `json:"decimals"`
			} `json:"baseToken"`
			QuoteToken struct {
				ID           NullableField `json:"id,string"`
				Name         string        `json:"name"`
				Address      string        `json:"address"`
				Symbol       string        `json:"symbol"`
				Slug         string        `json:"slug"`
				CryptoSymbol string        `json:"cryptoSymbol"`
				Decimals     int           `json:"decimals"`
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

func ReadFile(file *os.File, tokenData *TokenData) {
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&tokenData); err != nil {
		fmt.Println(err)
	}
}

func WriteFile(file *os.File, tokenData *TokenData) {
	file.Seek(0, 0)
	file.Truncate(0)
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(tokenData); err != nil {
		fmt.Println(err)
	}
}
