package model

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/tjfoc/gmsm/sm2"
)

type SmartContract struct {
	contractapi.Contract

	NFTNum  int
	UserNum int
	CertNum int
	Sk      *sm2.PrivateKey
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

type Cert struct {
	ID          string `json:"ID"`
	WithNID     string `json:"WithNID"`
	CertSign    string `json:"CertSign"`
	CertHash    string `json:"CertHash"`
	CreatedAt   string `json:"CreatedAt"`
	TxIDCreated string `json:"TxIDCreated"`
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
	Status       int      `json:"status"`
}

func NewUser(id string, money float64, time string, TxID string) User {
	return User{
		ID:             id,
		Money:          money,
		NFTCreatedBy:   make([]string, 0),
		NFTOwnedBy:     make([]string, 0),
		CreatedAt:      time,
		OpHistory:      make([]string, 0),
		OpHash:         make([]string, 0),
		TxIDCreated:    TxID,
		KeyGenSeedHash: "",
	}
}

func (s *SmartContract) NewNFT(id string, creator string, price float64, time string, TxID string, uri string, operations string) NFT {
	priceH := make([]string, 1)
	hPrice, err := SM2Encrypt(strconv.FormatFloat(price, 'f', -1, 32), 222, 111)
	if err != nil {
		log.Println(err)
		return NFT{}
	}
	priceH[0] = hPrice

	nft := NFT{
		ID:           id,
		Creator:      creator,
		Owner:        creator,
		Price:        price,
		PriceHistory: priceH,
		CreatedAt:    time,
		TxIDCreated:  TxID,
		Uri:          uri,
		OpHistory:    make([]string, 1),
		OpHash:       make([]string, 1),
		CertID:       "",
		Status:       0,
	}

	cText, err := SM2Encrypt(operations, 222, 111)
	if err != nil {
		fmt.Println(err)
		return NFT{}
	}
	nft.OpHistory[0] = cText

	hText, err := GetHash(operations)
	if err != nil {
		return NFT{}
	}
	nft.OpHash[0] = hText
	return nft
}

func NewCert(id string, nid string, time string) Cert {
	return Cert{
		ID:        id,
		WithNID:   nid,
		CreatedAt: time,
		CertSign:  "",
		CertHash:  "",
	}
}

func NewTxInfo(CT6 string, DC string) TxInfo {
	return TxInfo{
		CT6: CT6,
		DC:  DC,
	}
}
