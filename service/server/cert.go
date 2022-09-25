package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ApplyCertReq struct {
	UID           string `json:"UID"`
	NID           string `json:"NID"`
	CertMsg       string `json:"CertMsg"`
	KeyGenSeed    int64  `json:"KeyGenSeed"`
	EncryptSeed   int64  `json:"EncryptSeed"`
	OperationUser string `json:"OperationUser"`
	OperationNFT  string `json:"OperationNFT"`
}

func (server *MyServer) ApplyCert(c *gin.Context) {
	var req ApplyCertReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}

	result, err := server.contract.SubmitTransaction("ApplyCert", req.UID, req.NID, req.CertMsg,
		strconv.FormatInt(req.KeyGenSeed, 10), strconv.FormatInt(req.EncryptSeed, 10), req.OperationNFT, req.OperationUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	var res Cert
	err = json.Unmarshal(result, &res)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": res,
		"tx_id":  res.TxIDCreated,
	})
}

func (server *MyServer) GetCert(c *gin.Context) {
	if c.Query("CID") == "" {
		c.JSON(http.StatusInternalServerError, ErrorResponse(fmt.Errorf("CID is null")))
		return
	}

	result, err := server.contract.EvaluateTransaction("GetCert", c.Query("CID"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	var res Cert
	err = json.Unmarshal(result, &res)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, res)
}

func (server *MyServer) VerifyCert(c *gin.Context) {
	if c.Query("CID") == "" || c.Query("CertSign") == "" || c.Query("CertMsg") == "" {
		c.JSON(http.StatusInternalServerError, ErrorResponse(fmt.Errorf("CID is null")))
		return
	}

	result, err := server.contract.EvaluateTransaction("VerifyCert", c.Query("CID"), c.Query("CertSign"), c.Query("CertMsg"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": string(result),
	})
}
