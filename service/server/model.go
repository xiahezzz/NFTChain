package server

type Asset struct {
	ID             string `json:"ID"`
	Color          string `json:"color"`
	Size           int    `json:"size"`
	Owner          string `json:"owner"`
	AppraisedValue int    `json:"appraisedValue"`
}

type CreateUserReq struct {
	Money float64 `json:"Money"`
}

type ReadUserReq struct {
	ID string `json:"ID"`
}

type User struct {
	ID             string   `json:"ID"`
	NFTCreatedBy   []string `json:"NFTCreatedBy"`
	NFTOwnedBy     []string `json:"NFTOwnedBy"`
	Money          float64  `json:"Money"`
	CreatedAt      string   `json:"CreatedAt"`
	OpHistory      []string `json:"OpHistory"`
	OpHash         []string `json:"OpHash"`
	TxIDCreated    string   `json:"TxIDCreated"`
	KeyGenSeedHash string   `json:"KeyGenSeedHash"`
}

type TxInfo struct {
	CT6 string `json:"CT6"`
	DC  string `json:"DC"`
}

type NFT struct {
	ID           string   `json:"ID"`
	Creator      string   `json:"Creator"`
	Owner        string   `json:"Owner"`
	Price        float64  `json:"Price"`
	PriceHistory []string `json:"PriceHistory"`
	CreatedAt    string   `json:"CreatedAt"`
	TxIDCreated  string   `json:"TxIDCreated"`
	Uri          string   `json:"Uri"`
	OpHistory    []string `json:"OpHistory"`
	OpHash       []string `json:"OpHash"`
	CertID       string   `json:"CertID"`
	Status       int      `json:"Status"`
}

type Cert struct {
	ID          string `json:"ID"`
	WithNID     string `json:"WithNID"`
	CertSign    string `json:"CertSign"`
	CertHash    string `json:"CertHash"`
	CreatedAt   string `json:"CreatedAt"`
	TxIDCreated string `json:"TxIDCreated"`
}
