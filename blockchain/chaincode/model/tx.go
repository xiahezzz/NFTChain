package model

import (
	"encoding/json"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func (s *SmartContract) SetCT6AndDC(ctx contractapi.TransactionContextInterface, TID string, CT6 string, DC string) error {
	txInfo := NewTxInfo(CT6, DC)

	txInfoJSON, err := json.Marshal(txInfo)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(TID, txInfoJSON)
}

func (s *SmartContract) GetCT6AndDC(ctx contractapi.TransactionContextInterface, TID string) (*TxInfo, error) {
	txInfoJSON, err := ctx.GetStub().GetState(TID)
	if err != nil {
		return nil, err
	}

	var txInfo TxInfo
	err = json.Unmarshal(txInfoJSON, &txInfo)
	if err != nil {
		return nil, err
	}

	return &txInfo, nil
}
