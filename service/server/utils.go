package server

import (
	"fmt"
	"nftservice/chaincode"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

func GetTransactionByID(qscc *gateway.Contract, TxID string) (*Response, error) {
	txn, err := qscc.CreateTransaction("GetTransactionByID", gateway.WithEndorsingPeers("peer0.org1.example.com:7051"))
	if err != nil {
		fmt.Printf("Failed to create transaction: %s\n", err)
		return nil, err
	}

	k, err := txn.Evaluate(chaincode.CHANNEL_NAME, TxID)
	if err != nil {
		fmt.Printf("Failed to evaluate transaction: %s\n", err)
		return nil, err
	}

	res := &Response{}
	proto.Unmarshal(k, res)
	fmt.Println(res.Status, res.Payload)
	if err != nil {
		fmt.Printf("Failed to unmarshal res: %s\n", err)
		return nil, err
	}
	return res, nil
}
