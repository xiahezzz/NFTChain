package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (server *MyServer) GetTxINfo(c *gin.Context) {
	if c.Query("TID") == "" {
		c.JSON(http.StatusInternalServerError, ErrorResponse(fmt.Errorf("TID is null")))
		return
	}

	result, err := server.contract.EvaluateTransaction("GetCT6AndDC", c.Query("TID"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	var res TxInfo
	err = json.Unmarshal(result, &res)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, res)
}
