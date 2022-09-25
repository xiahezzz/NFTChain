package model

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func (s *SmartContract) updateCert(ctx contractapi.TransactionContextInterface, cid string, cert *Cert) error {
	certJSON, err := json.Marshal(cert)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(cid, certJSON)
}

func (s *SmartContract) GetCert(ctx contractapi.TransactionContextInterface, cid string) (*Cert, error) {
	certJSON, err := ctx.GetStub().GetState(cid)
	if err != nil {
		return nil, err
	}

	var cert Cert
	err = json.Unmarshal(certJSON, &cert)
	if err != nil {
		return nil, err
	}
	return &cert, nil
}

func (s *SmartContract) ApplyCert(ctx contractapi.TransactionContextInterface, uid string, nid string, certMsg string, keyGenSeed int64, encryptSeed int64, operationNFT string, operationUser string) (*Cert, error) {
	if !s.checkUserKeyHash(ctx, uid, keyGenSeed) {
		return nil, fmt.Errorf("keyGenSeed check false")
	}
	time, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return nil, err
	}

	s.CertNum++
	id := strings.Join([]string{"C", strconv.Itoa(s.CertNum)}, "")
	cert := NewCert(id, nid, time.String())
	cert.TxIDCreated = ctx.GetStub().GetTxID()
	certHash, err := GetHash(certMsg)
	if err != nil {
		return nil, err
	}
	cert.CertHash = certHash

	certSign, err := SM2Sign(certMsg, 444, 333)
	if err != nil {
		return nil, err
	}
	cert.CertSign = certSign

	err = s.updateCert(ctx, id, &cert)
	if err != nil {
		return nil, err
	}

	nft, err := s.GetNFTByNID(ctx, nid)
	if err != nil {
		return nil, err
	}
	nft.CertID = id
	err = s.registerNFTOp(nft, operationNFT)
	if err != nil {
		return nil, err
	}
	err = s.updateNFT(ctx, nid, nft)
	if err != nil {
		return nil, err
	}

	user, err := s.GetUser(ctx, uid)
	if err != nil {
		return nil, err
	}
	err = s.registerUserOp(user, operationUser, keyGenSeed, encryptSeed)
	if err != nil {
		return nil, err
	}
	err = s.updateUser(ctx, uid, user)
	if err != nil {
		return nil, err
	}

	return &cert, nil
}

func (s *SmartContract) VerifyCert(ctx contractapi.TransactionContextInterface, cid string, certSign string, certMsg string) (string, error) {
	cert, err := s.GetCert(ctx, cid)
	if err != nil {
		return "", err
	}

	if !CheckHash(certMsg, cert.CertHash) {
		return "", fmt.Errorf("check hash false")
	}

	ok, err := SM2Verify(certMsg, cert.CertSign, 444)
	if err != nil {
		return "", err
	}

	if !ok {
		return "", fmt.Errorf("VerifyCert False")
	}
	return "OK", nil
}
