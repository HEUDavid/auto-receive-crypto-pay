package parser

type RawContract struct {
	RawValue string `json:"rawValue"`
	Address  string `json:"address"`
	Decimals int    `json:"decimals"`
}

type Log struct {
	Address          string   `json:"address"`
	Topics           []string `json:"topics"`
	Data             string   `json:"data"`
	BlockNumber      string   `json:"blockNumber"`
	TransactionHash  string   `json:"transactionHash"`
	TransactionIndex string   `json:"transactionIndex"`
	BlockHash        string   `json:"blockHash"`
	LogIndex         string   `json:"logIndex"`
	Removed          bool     `json:"removed"`
}

type Activity struct {
	BlockNum        string      `json:"blockNum"`
	Hash            string      `json:"hash"`
	FromAddress     string      `json:"fromAddress"`
	ToAddress       string      `json:"toAddress"`
	Value           float64     `json:"value"`
	Erc721TokenId   *string     `json:"erc721TokenId"`
	Erc1155Metadata *string     `json:"erc1155Metadata"`
	Asset           string      `json:"asset"`
	Category        string      `json:"category"`
	RawContract     RawContract `json:"rawContract"`
	Log             Log         `json:"log"`
}

type Event struct {
	Network  string     `json:"network"`
	Activity []Activity `json:"activity"`
}

type WebhookData struct {
	WebhookId string `json:"webhookId"`
	Id        string `json:"id"`
	CreatedAt string `json:"createdAt"`
	Type      string `json:"type"`
	Event     Event  `json:"event"`
}
