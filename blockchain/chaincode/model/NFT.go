package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	math_rand "math/rand"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/tjfoc/gmsm/sm2"
)

func IsExistenceAny(ctx contractapi.TransactionContextInterface, id string) ([]byte, bool) {
	JSON, err := ctx.GetStub().GetState(id)
	if err != nil || JSON == nil {
		log.Println("err:", err, "nftJSON:", JSON)
		return nil, false
	}
	return JSON, true
}

func (s *SmartContract) MintNFT(ctx contractapi.TransactionContextInterface, creator string, price float64, uri string, keyGenSeed int64, encryptSeed int64, operationNFT string, operationUser string) (*NFT, error) {
	time, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if !s.checkUserKeyHash(ctx, creator, keyGenSeed) {
		return nil, fmt.Errorf("keyGenSeed check false")
	}
	s.NFTNum++
	id := strings.Join([]string{"N", strconv.Itoa(s.NFTNum)}, "")
	nft := s.NewNFT(id, creator, price, time.String(), ctx.GetStub().GetTxID(), uri, operationNFT)
	log.Println(nft)
	nftJson, err := json.Marshal(nft)
	if err != nil {
		log.Println(err)
		s.NFTNum--
		return nil, err
	}

	err = ctx.GetStub().PutState(id, nftJson)
	if err != nil {
		s.NFTNum--
		log.Println(err)
		return nil, err
	}

	user, err := s.GetUser(ctx, creator)
	if err != nil {
		s.NFTNum--
		log.Println(err)
		return nil, err
	}
	s.updateUserNFT(user, id, id)
	err = s.registerUserOp(user, operationUser, keyGenSeed, encryptSeed)
	if err != nil {
		s.NFTNum--
		log.Println(err)
		return nil, err
	}
	s.updateUser(ctx, creator, user)
	if err != nil {
		s.NFTNum--
		log.Println(err)
		return nil, err
	}
	return &nft, nil
}

func (s *SmartContract) updateNFT(ctx contractapi.TransactionContextInterface, nid string, nft *NFT) error {
	nftJSON, err := json.Marshal(nft)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(nid, nftJSON)
}

func (s *SmartContract) GetNFTByNID(ctx contractapi.TransactionContextInterface, nid string) (*NFT, error) {
	nftJSON, exist := IsExistenceAny(ctx, nid)
	if !exist {
		return nil, fmt.Errorf(nid, "does not existence")
	}

	var nft NFT
	err := json.Unmarshal(nftJSON, &nft)
	if err != nil {
		return nil, err
	}
	return &nft, nil
}

func (s *SmartContract) GetNFTByUIDCreate(ctx contractapi.TransactionContextInterface, uid string) ([]*NFT, error) {
	userJSON, exist := IsExistenceAny(ctx, uid)
	if !exist {
		return nil, fmt.Errorf(uid, "does not existence")
	}

	var user User
	err := json.Unmarshal(userJSON, &user)
	if err != nil {
		return nil, err
	}

	res := make([]*NFT, len(user.NFTCreatedBy))
	for index, nid := range user.NFTCreatedBy {
		nft, err := s.GetNFTByNID(ctx, nid)
		if err != nil {
			return nil, err
		}
		res[index] = nft
	}
	return res, nil
}

func (s *SmartContract) GetNFTByUIDOwn(ctx contractapi.TransactionContextInterface, uid string) ([]*NFT, error) {
	userJSON, exist := IsExistenceAny(ctx, uid)
	if !exist {
		return nil, fmt.Errorf(uid, "does not existence")
	}

	var user User
	err := json.Unmarshal(userJSON, &user)
	if err != nil {
		return nil, err
	}

	res := make([]*NFT, len(user.NFTOwnedBy))
	for index, nid := range user.NFTOwnedBy {
		nft, err := s.GetNFTByNID(ctx, nid)
		if err != nil {
			return nil, err
		}
		res[index] = nft
	}
	return res, nil
}

type ChangeNFTPriceRes struct {
	Nft  *NFT   `json:"NFT"`
	TxID string `json:"TxID"`
}

func (s *SmartContract) ChangeNFTPrice(ctx contractapi.TransactionContextInterface, nid string, newPrice float64, keyGenSeed int64, encryptSeed int64, operationNFT string, operationUser string) (*ChangeNFTPriceRes, error) {
	nft, err := s.GetNFTByNID(ctx, nid)
	if err != nil {
		return nil, err
	}

	if !s.checkUserKeyHash(ctx, nft.Owner, keyGenSeed) {
		return nil, fmt.Errorf("keyGenSeed check false")
	}

	nft.Price = newPrice
	hNewPrice, err := SM2Encrypt(strconv.FormatFloat(newPrice, 'f', -1, 32), 222, 111)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	nft.PriceHistory = append(nft.PriceHistory, hNewPrice)

	err = s.registerNFTOp(nft, operationNFT)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	nftJSON, err := json.Marshal(nft)
	if err != nil {
		return nil, err
	}

	err = ctx.GetStub().PutState(nid, nftJSON)
	if err != nil {
		return nil, err
	}

	user, err := s.GetUser(ctx, nft.Owner)
	if err != nil {
		return nil, err
	}
	s.registerUserOp(user, operationUser, keyGenSeed, encryptSeed)
	if err != nil {
		return nil, err
	}
	err = s.updateUser(ctx, nft.Owner, user)
	if err != nil {
		return nil, err
	}
	return &ChangeNFTPriceRes{
		Nft:  nft,
		TxID: ctx.GetStub().GetTxID(),
	}, nil
}

func (s *SmartContract) registerNFTOp(nft *NFT, operations string) error {
	cText, err := SM2Encrypt(operations, 222, 111)
	if err != nil {
		fmt.Println(err)
		return err
	}
	nft.OpHistory = append(nft.OpHistory, cText)

	hText, err := GetHash(operations)
	if err != nil {
		return err
	}
	nft.OpHash = append(nft.OpHash, hText)
	return nil
}

type DecryptNFTHistoryRes struct {
	Res []string `json:"Res"`
}

func (s *SmartContract) DecryptNFTHistory(ctx contractapi.TransactionContextInterface, nid string) (*DecryptNFTHistoryRes, error) {
	math_rand.Seed(222)
	keyStr := RandAllString(40)
	keyReader := bytes.NewReader([]byte(keyStr + keyStr))
	priv, err := sm2.GenerateKey(keyReader) // 生成密钥对
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	nft, err := s.GetNFTByNID(ctx, nid)
	if err != nil {
		return nil, err
	}

	msg := make([]string, len(nft.OpHistory))
	for index, c := range nft.OpHistory {
		m := SM2Decrypt(priv, c)
		msg[index] = m
		if !CheckHash(m, nft.OpHash[index]) {
			return nil, fmt.Errorf("DecryptNFTHistory:CheckHash false:%d", index)
		}
	}
	return &DecryptNFTHistoryRes{msg}, nil
}

type ChangeNFTStatusRes struct {
	Nft  *NFT   `json:"NFT"`
	TxID string `json:"TxID"`
}

func (s *SmartContract) ChangeNFTStatus(ctx contractapi.TransactionContextInterface, nid string, status int, keyGenSeed int64, encryptSeed int64, operationNFT string, operationUser string) (*ChangeNFTStatusRes, error) {
	nft, err := s.GetNFTByNID(ctx, nid)
	if err != nil {
		return nil, err
	}

	nft.Status = status

	err = s.registerNFTOp(nft, operationNFT)
	if err != nil {
		return nil, err
	}
	err = s.updateNFT(ctx, nid, nft)
	if err != nil {
		return nil, err
	}

	user, err := s.GetUser(ctx, nft.Owner)
	if err != nil {
		return nil, err
	}
	err = s.registerUserOp(user, operationUser, keyGenSeed, encryptSeed)
	if err != nil {
		return nil, err
	}
	err = s.updateUser(ctx, nft.Owner, user)
	if err != nil {
		return nil, err
	}
	return &ChangeNFTStatusRes{
		Nft:  nft,
		TxID: ctx.GetStub().GetTxID(),
	}, nil
}

type ChangeNFTOwnerRes struct {
	Nft  *NFT   `json:"NFT"`
	TxID string `json:"TxID"`
}

func (s *SmartContract) ChangeNFTOwnerRes(ctx contractapi.TransactionContextInterface, nid string, buyUID string, TID string, keyGenSeedBuy int64, keyGenSeedSold int64, encryptSeedBuy int64, encryptSeedSold int64, operationNFT string, operationBuyUser string, operationSoldUser string, CT6 string, DC string) (*ChangeNFTOwnerRes, error) {
	nft, err := s.GetNFTByNID(ctx, nid)
	if err != nil {
		return nil, err
	}
	if nft.Status != 1 {
		return nil, fmt.Errorf("nft status false")
	}
	soldUID := nft.Owner

	userSold, err := s.GetUser(ctx, soldUID)
	if err != nil {
		return nil, err
	}
	userBuy, err := s.GetUser(ctx, buyUID)
	if err != nil {
		return nil, err
	}
	if userBuy.Money < nft.Price {
		return nil, fmt.Errorf("money not enough")
	}

	if !CheckHash(strconv.FormatInt(keyGenSeedBuy, 10), userBuy.KeyGenSeedHash) {
		return nil, fmt.Errorf("check buy kenGenSeed false")
	}
	if !CheckHash(strconv.FormatInt(keyGenSeedSold, 10), userSold.KeyGenSeedHash) {
		return nil, fmt.Errorf("check sold kenGenSeed false")
	}

	userSold.NFTOwnedBy = DeleteSlice(userSold.NFTOwnedBy, nid)
	userBuy.NFTOwnedBy = append(userBuy.NFTOwnedBy, nid)

	userBuy.Money -= nft.Price
	userSold.Money += nft.Price
	nft.Owner = buyUID

	s.SetCT6AndDC(ctx, TID, CT6, DC)

	err = s.registerNFTOp(nft, operationNFT)
	if err != nil {
		return nil, err
	}
	err = s.updateNFT(ctx, nid, nft)
	if err != nil {
		return nil, err
	}

	err = s.registerUserOp(userBuy, operationBuyUser, keyGenSeedBuy, encryptSeedBuy)
	if err != nil {
		return nil, err
	}
	err = s.updateUser(ctx, buyUID, userBuy)
	if err != nil {
		return nil, err
	}

	err = s.registerUserOp(userSold, operationSoldUser, keyGenSeedSold, encryptSeedSold)
	if err != nil {
		return nil, err
	}
	err = s.updateUser(ctx, soldUID, userSold)
	if err != nil {
		return nil, err
	}

	return &ChangeNFTOwnerRes{
		Nft:  nft,
		TxID: ctx.GetStub().GetTxID(),
	}, nil
}
