package entities

type NFTTokenData struct {
	TokenID     string  `json:"token_id"`
	BlockNumber string  `json:"block_number"`
	Content     Content `json:"content"`
}

type Content struct {
	Data Data `json:"data"`
}
type Data struct {
	Name                 string    `json:"name"`
	Symbol               string    `json:"symbol"`
	URI                  string    `json:"uri"`
	SellerFeeBasisPoints int       `json:"sellerFeeBasisPoints"`
	Creators             []Creator `json:"creators"`
}

type Creator struct {
	Address string `json:"address"`
	Share   int    `json:"share"`
}
