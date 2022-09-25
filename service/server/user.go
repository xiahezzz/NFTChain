package server

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func (server *MyServer) CreateUser(c *gin.Context) {
	var req CreateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}
	keyGenSeed := rand.Int63()
	res, err := server.contract.SubmitTransaction("CreateUser", strconv.FormatFloat(req.Money, 'f', -1, 32), strconv.FormatInt(keyGenSeed, 10))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	var result User
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
		"result":     result,
		"keyGenSeed": keyGenSeed,
		"tx_id":      result.TxIDCreated,
	})
}

func (server *MyServer) GetUser(c *gin.Context) {
	if c.Query("UID") == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse(fmt.Errorf("NO UID")))
		return
	}

	res, err := server.contract.EvaluateTransaction("GetUser", c.Query("UID"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	var result User
	err = json.Unmarshal(res, &result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, result)

}

type DecryptUserHistoryReq struct {
	ID         string `json:"ID"`
	KeyGenSeed int64  `json:"KeyGenSeed"`
}

type DecryptUserHistoryRes struct {
	Res []string `json:"Res"`
}

func (server *MyServer) DecryptUserHistory(c *gin.Context) {
	var req DecryptUserHistoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}

	result, err := server.contract.EvaluateTransaction("DecryptUserHistory", req.ID, strconv.FormatInt(req.KeyGenSeed, 10))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	res := &DecryptUserHistoryRes{
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
