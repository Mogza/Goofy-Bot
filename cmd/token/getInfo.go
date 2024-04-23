package dexApi

import (
	"encoding/json"
	"fmt"
	"github/Mogza/Goofy-Bot/cmd/utils"
	"io"
	"net/http"
	"strings"
)

type Pairs struct {
	Pairs []Pair `json:"pairs"`
}

type Pair struct {
	ChainId   string    `json:"chainId"`
	DexId     string    `json:"dexId"`
	URL       string    `json:"url"`
	BaseToken BaseToken `json:"baseToken"`
	PriceUSD  string    `json:"priceUsd"`
	FDV       float64   `json:"fdv"`
	Liquidity Liquidity `json:"liquidity"`
	Volume    Volume    `json:"volume"`
	Info      Info      `json:"info"`
}

type BaseToken struct {
	Address string `json:"address"`
	Name    string `json:"name"`
	Symbol  string `json:"symbol"`
}

type Liquidity struct {
	USD   float64 `json:"usd"`
	Base  float64 `json:"base"`
	Quote float64 `json:"quote"`
}

type Volume struct {
	H24 float64 `json:"h24"`
}

type Info struct {
	ImageURL string    `json:"imageUrl"`
	Websites []Website `json:"websites"`
	Socials  []Social  `json:"socials"`
}

type Website struct {
	Label string `json:"Label"`
	URL   string `json:"url"`
}

type Social struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

func GetToken(addr string) (string, string) {
	response, err := http.Get("https://api.dexscreener.com/latest/dex/tokens/" + addr)
	utils.CheckError(err, "Error while calling the API")

	responseData, err := io.ReadAll(response.Body)
	utils.CheckError(err, "Error while reading the response")

	var pairs Pairs

	err = json.Unmarshal(responseData, &pairs)
	utils.CheckError(err, "Error while parsing json")

	if len(pairs.Pairs) == 0 {
		return "# :x: Token not found", ""
	}

	tokenAll := pairs.Pairs[0]
	finalResult := getString(pairs)

	return string(finalResult), tokenAll.Info.ImageURL
}

func getString(pairs Pairs) string {
	tokenAll := pairs.Pairs[0]

	finalResult := "### :coin: " + tokenAll.BaseToken.Name + " | $" + tokenAll.BaseToken.Symbol
	finalResult += "\n### :globe_with_meridians: " + makeFirstUpper(tokenAll.ChainId) + " @ " + makeFirstUpper(tokenAll.DexId)
	finalResult += "\n### :moneybag: USD: $" + tokenAll.PriceUSD
	finalResult += "\n### :gem: FDV: $" + formatValue(tokenAll.FDV)
	finalResult += "\n### :sweat_drops: Liq: $" + formatValue(tokenAll.Liquidity.USD)
	finalResult += "\n### :bar_chart: Vol: $" + formatValue(tokenAll.Volume.H24)

	finalResult += "\n### :link: Links: "
	for _, sites := range tokenAll.Info.Websites {
		finalResult += "[" + sites.Label + "]" + "(" + sites.URL + ") | "
	}

	finalResult += "\n###  :black_bird: Socials: "
	for _, socials := range tokenAll.Info.Socials {
		finalResult += "[" + makeFirstUpper(socials.Type) + "]" + "(" + makeFirstUpper(socials.URL) + ") | "
	}

	finalResult += "\n### :satellite: [DEX](" + tokenAll.URL + ")"
	if tokenAll.ChainId == "solana" {
		finalResult += " | [SOLSCAN](https://solscan.io/token/" + tokenAll.BaseToken.Address + ")"
		finalResult += " | [BIRDEYE](https://birdeye.so/token/" + tokenAll.BaseToken.Address + "?chain=solana)"
		finalResult += " | [DEXTOOLS](https://www.dextools.io/app/en/solana/pair-explorer/" + tokenAll.BaseToken.Address + ")"
		finalResult += " | [RUGCHECK](https://rugcheck.xyz/tokens/" + tokenAll.BaseToken.Address + ")"
	}

	finalResult += "\nby [Mogza](https://github.com/Mogza)"

	return finalResult
}

func makeFirstUpper(s string) string {
	return strings.ToUpper(s[:1]) + s[1:]
}

func formatValue(value float64) string {
	var suffix string
	var divisor float64

	switch {
	case value >= 1e6:
		suffix = "M"
		divisor = 1e6
	case value >= 1e3:
		suffix = "K"
		divisor = 1e3
	default:
		return fmt.Sprintf("%f", value)
	}

	return fmt.Sprintf("%.1f%s", value/divisor, suffix)
}
