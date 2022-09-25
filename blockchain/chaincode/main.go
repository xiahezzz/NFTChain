package main

import (
	"bytes"
	"log"
	math_rand "math/rand"
	"nft/model"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/tjfoc/gmsm/sm2"
)

// SmartContract provides functions for managing an Asset

/*type TXInfo struct {
	CT_6 string `json:"ct_6"`
	DC   string `json:"dc"`
}

type Cert struct {
	ID        string               `json:"ID"`
	CreatedAt *timestamp.Timestamp `json:"created_at"`
}

type NFT struct {
	ID        string               `json:"ID"`
	Creator   string               `json:"creator"`
	Owner     string               `json:"owner"`
	Price     string               `json:"price"`
	TxInfo    TXInfo               `json:"tx_info"`
	CreatedAt *timestamp.Timestamp `json:"created_at"`
}

func TimestampProto(t time.Time) *timestamp.Timestamp {
	ts := &timestamp.Timestamp{
		Seconds: t.Unix(),
		Nanos:   int32(t.Nanosecond()),
	}
	return ts
}

func TimestampNow() *timestamp.Timestamp {
	ts := TimestampProto(time.Now())
	return ts
}

func NewUser(id string, money int) *User {
	return &User{
		ID:           id,
		Money:        money,
		NFTCreatedBy: "",
		NFTOwnedBy:   "",
		CreatedAt:    TimestampNow(),
		Operations:   "",
	}
}

type Asset struct {
	AppraisedValue int    `json:"AppraisedValue"`
	Color          string `json:"Color"`
	ID             string `json:"ID"`
	Owner          string `json:"Owner"`
	Size           int    `json:"Size"`
}

// InitLedger adds a base set of assets to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	assets := []Asset{
		{ID: "asset1", Color: "blue", Size: 5, Owner: "Tomoko", AppraisedValue: 300},
		{ID: "asset2", Color: "red", Size: 5, Owner: "Brad", AppraisedValue: 400},
		{ID: "asset3", Color: "green", Size: 10, Owner: "Jin Soo", AppraisedValue: 500},
		{ID: "asset4", Color: "yellow", Size: 10, Owner: "Max", AppraisedValue: 600},
		{ID: "asset5", Color: "black", Size: 15, Owner: "Adriana", AppraisedValue: 700},
		{ID: "asset6", Color: "white", Size: 15, Owner: "Michel", AppraisedValue: 800},
	}

	for _, asset := range assets {
		assetJSON, err := json.Marshal(asset)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(asset.ID, assetJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Asset, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*Asset
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset Asset
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}

/*
// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract

	NFTNum  int
	UserNum int
}

type TXInfo struct {
	CT_6 string `json:"ct_6"`
	DC   string `json:"dc"`
}

type Cert struct {
	ID int `json:"ID"`
}

type NFT struct {
	ID      int    `json:"ID"`
	Creator string `json:"creator"`
	Owner   string `json:"owner"`
	Price   string `json:"price"`
	TxInfo  TXInfo `json:"tx_info"`
}

type User struct {
	ID           int   `json:"ID"`
	NFTCreatedBy []int `json:"nft_created_by"`
	NFTOwnedBy   []int `json:"nft_own_by"`

}

// Asset describes basic details of what makes up a simple asset
type Asset struct {
	ID             string `json:"ID"`
	Color          string `json:"color"`
	Size           int    `json:"size"`
	Owner          string `json:"owner"`
	AppraisedValue int    `json:"appraisedValue"`
}

// InitLedger adds a base set of assets to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {

	return nil
}

// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, id string, color string, size int, owner string, appraisedValue int) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", id)
	}

	asset := Asset{
		ID:             id,
		Color:          color,
		Size:           size,
		Owner:          owner,
		AppraisedValue: appraisedValue,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJSON)
}

// ReadAsset returns the asset stored in the world state with given id.
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, id string) (*Asset, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	var asset Asset
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

// UpdateAsset updates an existing asset in the world state with provided parameters.
func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, id string, color string, size int, owner string, appraisedValue int) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", id)
	}

	// overwriting original asset with new asset
	asset := Asset{
		ID:             id,
		Color:          color,
		Size:           size,
		Owner:          owner,
		AppraisedValue: appraisedValue,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJSON)
}

// DeleteAsset deletes an given asset from the world state.
func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
}

// AssetExists returns true when asset with given ID exists in world state
func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}

// TransferAsset updates the owner field of asset with given id in world state.
func (s *SmartContract) TransferAsset(ctx contractapi.TransactionContextInterface, id string, newOwner string) error {
	asset, err := s.ReadAsset(ctx, id)
	if err != nil {
		return err
	}

	asset.Owner = newOwner
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJSON)
}*/

/*
func (s *SmartContract) CreateUser(ctx contractapi.TransactionContextInterface, money int) error {
	user := NewUser(strconv.Itoa(1), money)
	userJson, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(strconv.Itoa(1), userJson)

}

func (s *SmartContract) GetUser(ctx contractapi.TransactionContextInterface, id string) (*User, error) {
	userJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, err
	}

	var user User
	err = json.Unmarshal(userJSON, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}*/

func InitContract() *model.SmartContract {
	math_rand.Seed(20)
	str := model.RandAllString(40)
	sReader := bytes.NewReader([]byte(str + str))
	contract := &model.SmartContract{}
	contract.Sk, _ = sm2.GenerateKey(sReader)
	return contract
}

func main() {
	assetChaincode, err := contractapi.NewChaincode(InitContract())
	if err != nil {
		log.Panicf("Error creating asset-transfer-basic chaincode: %v", err)
	}

	if err := assetChaincode.Start(); err != nil {
		log.Panicf("Error starting asset-transfer-basic chaincode: %v", err)
	}
}
