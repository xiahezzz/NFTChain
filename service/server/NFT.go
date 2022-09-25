package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type MintNFTReq struct {
	Creator       string  `json:"Creator"`
	Price         float64 `json:"Price"`
	Uri           string  `json:"Uri"`
	KeyGenSeed    int64   `json:"KeyGenSeed"`
	EncryptSeed   int64   `json:"EncryptSeed"`
	OperationUser string  `json:"OperationUser"`
	OperationNFT  string  `json:"OperationNFT"`
}

func (server *MyServer) MintNFT(c *gin.Context) {
	var req MintNFTReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}

	res, err := server.contract.SubmitTransaction("MintNFT", req.Creator, strconv.FormatFloat(req.Price, 'f', -1, 32), req.Uri,
		strconv.FormatInt(req.KeyGenSeed, 10), strconv.FormatInt(req.EncryptSeed, 10), req.OperationNFT, req.OperationUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	var nft NFT
	err = json.Unmarshal(res, &nft)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"result": nft,
		"tx_id":  nft.TxIDCreated,
	})
}

func (server *MyServer) GetNFTByNID(c *gin.Context) {
	if c.Query("NID") == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse(fmt.Errorf("NID is null")))
		return
	}

	res, err := server.contract.EvaluateTransaction("GetNFTByNID", c.Query("NID"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	var result NFT
	err = json.Unmarshal(res, &result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}
	//txInfo, err := GetTransactionByID(server.qscc, result.TxIDCreated)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"result": result,
	})
}

func (server *MyServer) GetNFTByUIDCreate(c *gin.Context) {
	if c.Query("UID") == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse(fmt.Errorf("UID is null")))
		return
	}

	res, err := server.contract.EvaluateTransaction("GetNFTByUIDCreate", c.Query("UID"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	result := make([]*NFT, len(strings.Split(string(res), "{")))
	err = json.Unmarshal(res, &result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}
	//txInfo, err := GetTransactionByID(server.qscc, result.TxIDCreated)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"result": result,
	})
}

func (server *MyServer) GetNFTByUIDOwn(c *gin.Context) {
	if c.Query("UID") == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse(fmt.Errorf("UID is null")))
		return
	}

	res, err := server.contract.EvaluateTransaction("GetNFTByUIDOwn", c.Query("UID"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	result := make([]*NFT, len(strings.Split(string(res), "{")))
	err = json.Unmarshal(res, &result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}
	//txInfo, err := GetTransactionByID(server.qscc, result.TxIDCreated)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"result": result,
	})
}

type ChangeNFTPriceReq struct {
	NID           string  `json:"NID"`
	Price         float64 `json:"Price"`
	KeyGenSeed    int64   `json:"KeyGenSeed"`
	EncryptSeed   int64   `json:"EncryptSeed"`
	OperationUser string  `json:"OperationUser"`
	OperationNFT  string  `json:"OperationNFT"`
}

type ReadNFTRes struct {
	Nft  *NFT   `json:"NFT"`
	TxID string `json:"TxiD"`
}

func (server *MyServer) ChangeNFTPrice(c *gin.Context) {
	var req ChangeNFTPriceReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}

	result, err := server.contract.SubmitTransaction("ChangeNFTPrice", req.NID, strconv.FormatFloat(req.Price, 'f', -1, 32), strconv.FormatInt(req.KeyGenSeed, 10), strconv.FormatInt(req.EncryptSeed, 10), req.OperationNFT, req.OperationUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	var res ReadNFTRes
	err = json.Unmarshal(result, &res)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"result": res.Nft,
		"tx_id":  res.TxID,
	})
}

type RegisterNFTOpReq struct {
	NID         string `json:"NID"`
	Operations  string `json:"Operations"`
	KeyGenSeed  int64  `json:"KeyGenSeed"`
	EncryptSeed int64  `json:"EncryptSeed"`
}

func (server *MyServer) RegisterNFTOp(c *gin.Context) {
	var req RegisterNFTOpReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}

	fmt.Printf("req.Operations: %v\n", req.Operations)
	result, err := server.contract.SubmitTransaction("RegisterNFTOp", req.NID, req.Operations, strconv.FormatInt(req.KeyGenSeed, 10), strconv.FormatInt(req.EncryptSeed, 10))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	var res ReadNFTRes
	err = json.Unmarshal(result, &res)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"result": res.Nft,
		"tx_id":  res.TxID,
	})
}

type DecryptNFTHistoryRes struct {
	Res []string `json:"Res"`
}

func (server *MyServer) DecryptNFTHistory(c *gin.Context) {
	if c.Query("NID") == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse(fmt.Errorf("NID is null")))
		return
	}

	result, err := server.contract.EvaluateTransaction("DecryptNFTHistory", c.Query("NID"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	res := &DecryptNFTHistoryRes{
		Res: make([]string, len(strings.Split(string(result), ","))),
	}
	err = json.Unmarshal(result, res)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK,
		gin.H{
			"result": res.Res,
		},
	)
}

type ChangeNFTStatusReq struct {
	NID           string `json:"NID"`
	Status        int    `json:"Status"`
	KeyGenSeed    int64  `json:"KeyGenSeed"`
	EncryptSeed   int64  `json:"EncryptSeed"`
	OperationUser string `json:"OperationUser"`
	OperationNFT  string `json:"OperationNFT"`
}

func (server *MyServer) ChangeNFTStatus(c *gin.Context) {
	var req ChangeNFTStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}

	result, err := server.contract.SubmitTransaction("ChangeNFTStatus", req.NID,
		strconv.Itoa(req.Status), strconv.FormatInt(req.KeyGenSeed, 10), strconv.FormatInt(req.EncryptSeed, 10), req.OperationNFT, req.OperationUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	var res ReadNFTRes
	err = json.Unmarshal(result, &res)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": res.Nft,
		"tx_id":  res.TxID,
	})
}

type ChangeNFTOwnerReq struct {
	NID               string `json:"NID"`
	BuyUID            string `json:"BuyUID"`
	TID               string `json:"TID"`
	KeyGenSeedBuy     int64  `json:"KeyGenSeedBuy"`
	KeyGenSeedSold    int64  `json:"KeyGenSeedSold"`
	EncryptSeedBuy    int64  `json:"EncryptSeedBuy"`
	EncryptSeedSold   int64  `json:"EncryptSeedSold"`
	OperationNFT      string `json:"OperationNFT"`
	OperationBuyUser  string `json:"OperationBuyUser"`
	OperationSoldUser string `json:"OperationSoldUser"`
	CT6               string `json:"CT6"`
	DC                string `json:"DC"`
}

func (server *MyServer) ChangeNFTOwner(c *gin.Context) {
	var req ChangeNFTOwnerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}

	result, err := server.contract.SubmitTransaction("ChangeNFTOwnerRes", req.NID, req.BuyUID, req.TID,
		strconv.FormatInt(req.KeyGenSeedBuy, 10), strconv.FormatInt(req.KeyGenSeedSold, 10),
		strconv.FormatInt(req.EncryptSeedBuy, 10), strconv.FormatInt(req.EncryptSeedSold, 10),
		req.OperationNFT, req.OperationBuyUser, req.OperationSoldUser, req.CT6, req.DC)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	var res ReadNFTRes
	err = json.Unmarshal(result, &res)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": res.Nft,
		"tx_id":  res.TxID,
	})
}
