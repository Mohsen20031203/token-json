package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
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

	//if c == Blast {
	//	return 81457
	//}

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
	//if c == CronosChain {
	//	return 25
	//}

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
	//if c == Blast {
	//	return ?
	//}
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

	//if c == CronosChain {
	//	return ?
	//}
	return ""
}

type Coin struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Symbol         string    `json:"symbol"`
	Slug           string    `json:"slug"`
	NumMarketPairs int       `json:"num_market_pairs"`
	DateAdded      time.Time `json:"date_added"`
	Tags           []struct {
		Slug     string `json:"slug"`
		Name     string `json:"name"`
		Category string `json:"category"`
	} `json:"tags"`
	MaxSupply         interface{} `json:"max_supply"`
	CirculatingSupply interface{} `json:"circulating_supply"`
	TotalSupply       float64     `json:"total_supply"`
	Platform          struct {
		ID           int    `json:"id"`
		Name         string `json:"name"`
		Symbol       string `json:"symbol"`
		Slug         string `json:"slug"`
		TokenAddress string `json:"token_address"`
	} `json:"platform"`
	IsActive                      int         `json:"is_active"`
	InfiniteSupply                bool        `json:"infinite_supply"`
	CmcRank                       int         `json:"cmc_rank"`
	IsFiat                        int         `json:"is_fiat"`
	SelfReportedCirculatingSupply interface{} `json:"self_reported_circulating_supply"`
	SelfReportedMarketCap         interface{} `json:"self_reported_market_cap"`
	TvlRatio                      float64     `json:"tvl_ratio"`
	LastUpdated                   time.Time   `json:"last_updated"`
	Quote                         struct {
		USD struct {
			Price                 float64   `json:"price"`
			Volume24H             float64   `json:"volume_24h"`
			VolumeChange24H       float64   `json:"volume_change_24h"`
			PercentChange1H       float64   `json:"percent_change_1h"`
			PercentChange24H      float64   `json:"percent_change_24h"`
			PercentChange7D       float64   `json:"percent_change_7d"`
			PercentChange30D      float64   `json:"percent_change_30d"`
			PercentChange60D      float64   `json:"percent_change_60d"`
			PercentChange90D      float64   `json:"percent_change_90d"`
			MarketCap             float64   `json:"market_cap"`
			MarketCapDominance    float64   `json:"market_cap_dominance"`
			FullyDilutedMarketCap float64   `json:"fully_diluted_market_cap"`
			Tvl                   float64   `json:"tvl"`
			LastUpdated           time.Time `json:"last_updated"`
		} `json:"USD"`
	} `json:"quote"`
}

type CmcResponse struct {
	Status struct {
		Timestamp    time.Time   `json:"timestamp"`
		ErrorCode    int         `json:"error_code"`
		ErrorMessage interface{} `json:"error_message"`
		Elapsed      int         `json:"elapsed"`
		CreditCount  int         `json:"credit_count"`
		Notice       interface{} `json:"notice"`
	} `json:"status"`
	Data map[string]Coin `json:"data"`
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
		var tokenData utils.TokenData
		var err error
		var numberWrite int

		nameFile := fmt.Sprintf("tokens%v.json", nameChain.GetCMCName())
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
			tokenData, err = fetchAndProcessURL(nameChain)
			if err != nil {
				return
			}
		} else {
			utils.ReadFile(file, &tokenData)
		}

		for key, value := range tokenData.Tokens {
			/*
				if value.CmcID != "0" {
					if !checkID(value.Symbol, key, value.CmcID) {
						fmt.Println(key, value.CmcID)
					}
				}
			*/

			if value.CmcID != "" {
				continue
			}

			var cmcid int
			time.Sleep(time.Second * 1)

			for i := 0; i < 5; i++ {
				cmcid, err = fetchTokenPrice(nameChain, tokenData.Tokens[key].Address)
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
		//numberCmcid(tokenData, string(nameFile))

	}
}

func checkID(nameChain, addres string, id string) bool {
	var respostId CmcResponse
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return false
	}
	apiUrl := fmt.Sprintf("https://api.vultisig.com/cmc/v2/cryptocurrency/quotes/latest?id=%d", idInt)
	resp, err := http.Get(apiUrl)
	if err != nil {
		return false
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false
	}

	err = json.Unmarshal(body, &respostId)
	if err != nil {
		return false
	}

	m := respostId.Data[id]

	if m.Symbol == nameChain && m.ID == idInt {
		return true
	}

	return false

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
				if pair.BaseToken.ID.Set {
					pair.BaseToken.ID.Value = strings.Trim(pair.BaseToken.ID.Value, "\"")

					return strconv.Atoi(pair.BaseToken.ID.Value)
				} else {
					numebrName, err := strconv.Atoi(pair.BaseToken.Name)
					if err == nil {
						return numebrName, err
					} else {
						return numebrName, nil
					}
				}
			}
			if strings.EqualFold(pair.QuoteToken.Address, address) {
				if pair.QuoteToken.ID.Set {
					pair.QuoteToken.ID.Value = strings.Trim(pair.QuoteToken.ID.Value, "\"")

					return strconv.Atoi(pair.QuoteToken.ID.Value)

				} else {
					numebrName, err := strconv.Atoi(pair.QuoteToken.Name)
					if err == nil {
						return numebrName, err
					} else {
						return 0, nil
					}
				}
			}
		}

	}
	return 0, nil

}

func numberCmcid(tokenData utils.TokenData, namechain string) {
	var z_cmcid int
	var cmcid int
	for _, value := range tokenData.Tokens {
		if value.CmcID == "0" {
			z_cmcid++
		} else if value.CmcID != "0" {
			cmcid++
		}

	}

	file, err := os.OpenFile("file.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	result := fmt.Sprintf("%s\nallToken: %d\nCorrect cmcids: %d\ncmcids that are zero: %d\n\n",
		namechain, len(tokenData.Tokens), cmcid, z_cmcid)
	_, err = file.WriteString(result)
	if err != nil {
		log.Fatal(err)
	}
}
