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

func (s *SmartContract) CreateUser(ctx contractapi.TransactionContextInterface, money float64, keyGenSeed int64) (*User, error) {
	time, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	s.UserNum++
	id := strings.Join([]string{"U", strconv.Itoa(s.UserNum)}, "")
	user := NewUser(id, money, time.String(), ctx.GetStub().GetTxID())
	keyHash, err := GetHash(strconv.FormatInt(keyGenSeed, 10))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	user.KeyGenSeedHash = keyHash
	log.Println(user)
	userJson, err := json.Marshal(user)
	if err != nil {
		log.Println(err)
		s.UserNum--
		return nil, err
	}

	err = ctx.GetStub().PutState(id, userJson)
	if err != nil {
		s.UserNum--
		log.Println(err)
		return nil, err
	}
	return &user, nil
}

func (s *SmartContract) GetUser(ctx contractapi.TransactionContextInterface, id string) (*User, error) {
	userJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if userJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	var user User
	err = json.Unmarshal(userJSON, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *SmartContract) checkUserKeyHash(ctx contractapi.TransactionContextInterface, uid string, keyGenSeed int64) bool {
	user, err := s.GetUser(ctx, uid)
	if err != nil {
		log.Fatal("checkUserKeyHash Failed:", err)
		return false
	}
	return CheckHash(strconv.FormatInt(keyGenSeed, 10), user.KeyGenSeedHash)
}

func (s *SmartContract) updateUserNFT(user *User, nftCreated string, nftOwned string) error {
	if nftCreated != "" {
		user.NFTCreatedBy = append(user.NFTCreatedBy, nftCreated)
	}
	if nftOwned != "" {
		user.NFTOwnedBy = append(user.NFTOwnedBy, nftOwned)
	}
	return nil
}

type RegisterUserOpRes struct {
	USER *User  `json:"User"`
	TxID string `json:"TxID"`
}

/*func (s *SmartContract) RegisterUserOp(ctx contractapi.TransactionContextInterface, uid string, operation string, keyGenSeed int64, encryptSeed int64) (*RegisterUserOpRes, error) {
	user, err := s.GetUser(ctx, uid)
	if err != nil {
		return nil, err
	}

	if !s.checkUserKeyHash(ctx, uid, keyGenSeed) {
		return nil, fmt.Errorf("keyGenSeed check false")
	}

	cText, err := SM2Encrypt(operation, keyGenSeed, encryptSeed)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("RegisterNFTOp:SM2Encrypt:Enc Error")
	}
	user.OpHistory = append(user.OpHistory, cText)

	hText, err := GetHash(operation)
	if err != nil {
		return nil, fmt.Errorf("RegisterNFTOp:GetHash:Hash Error")
	}
	user.OpHash = append(user.OpHash, hText)

	userJSON, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	err = ctx.GetStub().PutState(uid, userJSON)
	if err != nil {
		return nil, err
	}

	return &RegisterUserOpRes{
		USER: user,
		TxID: ctx.GetStub().GetTxID(),
	}, nil
}*/

func (s *SmartContract) registerUserOp(user *User, operation string, keyGenSeed int64, encryptSeed int64) error {
	cText, err := SM2Encrypt(operation, keyGenSeed, encryptSeed)
	if err != nil {
		fmt.Println(err)
		return err
	}
	user.OpHistory = append(user.OpHistory, cText)

	hText, err := GetHash(operation)
	if err != nil {
		return err
	}
	user.OpHash = append(user.OpHash, hText)

	return nil
}

func (s *SmartContract) updateUser(ctx contractapi.TransactionContextInterface, uid string, user *User) error {
	userJSON, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(uid, userJSON)
}

type DecryptUserHistoryRes struct {
	Res []string `json:"Res"`
}

func (s *SmartContract) DecryptUserHistory(ctx contractapi.TransactionContextInterface, uid string, keyGenSeed int64) (*DecryptUserHistoryRes, error) {
	user, err := s.GetUser(ctx, uid)
	if err != nil {
		return nil, err
	}
	math_rand.Seed(keyGenSeed)
	keyStr := RandAllString(40)
	keyReader := bytes.NewReader([]byte(keyStr + keyStr))
	priv, err := sm2.GenerateKey(keyReader) // 生成密钥对
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	msg := make([]string, len(user.OpHistory))
	for index, c := range user.OpHistory {
		m := SM2Decrypt(priv, c)
		msg[index] = m
		if !CheckHash(m, user.OpHash[index]) {
			return nil, fmt.Errorf("DecryptNFTHistory:CheckHash false:%d", index)
		}
	}
	return &DecryptUserHistoryRes{msg}, nil
}
