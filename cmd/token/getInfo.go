package dexApi

import (
	"encoding/json"
	"github/Mogza/Goofy-Bot/cmd/utils"
	"io"
	"net/http"
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
	var pairs Pairs

	// Call Dex Screener API
	response, err := http.Get("https://api.dexscreener.com/latest/dex/tokens/" + addr)
	utils.CheckError(err, "Error while calling the API")

	// Get Response
	responseData, err := io.ReadAll(response.Body)
	utils.CheckError(err, "Error while reading the response")

	// Parse Response
	err = json.Unmarshal(responseData, &pairs)
	utils.CheckError(err, "Error while parsing json")

	// Return if no token found
	if len(pairs.Pairs) == 0 {
		return "# :x: Token not found", ""
	}

	// Retrieve information on the first pair
	tokenAll := pairs.Pairs[0]
	finalResult := getString(pairs)

	return finalResult, tokenAll.Info.ImageURL
}

func getString(pairs Pairs) string {
	tokenAll := pairs.Pairs[0]

	finalResult := "### :coin: " + tokenAll.BaseToken.Name + " | $" + tokenAll.BaseToken.Symbol
	finalResult += "\n### :globe_with_meridians: " + utils.MakeFirstUpper(tokenAll.ChainId) + " @ " + utils.MakeFirstUpper(tokenAll.DexId)
	finalResult += "\n### :moneybag: USD: $" + tokenAll.PriceUSD
	finalResult += "\n### :gem: FDV: $" + utils.FormatValue(tokenAll.FDV)
	finalResult += "\n### :sweat_drops: Liq: $" + utils.FormatValue(tokenAll.Liquidity.USD)
	finalResult += "\n### :bar_chart: Vol: $" + utils.FormatValue(tokenAll.Volume.H24)

	finalResult += "\n### :link: Links: "
	for _, sites := range tokenAll.Info.Websites {
		finalResult += "[" + sites.Label + "]" + "(" + sites.URL + ") | "
	}

	finalResult += "\n###  :black_bird: Socials: "
	for _, socials := range tokenAll.Info.Socials {
		finalResult += "[" + utils.MakeFirstUpper(socials.Type) + "]" + "(" + utils.MakeFirstUpper(socials.URL) + ") | "
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
