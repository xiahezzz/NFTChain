package server

import (
	"nftservice/chaincode"

	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

type MyServer struct {
	router *gin.Engine

	wallet   *gateway.Wallet
	gw       *gateway.Gateway
	contract *gateway.Contract
	qscc     *gateway.Contract
}

func NewServer() *MyServer {
	server := &MyServer{}

	server.SetupWallet()
	server.SetupGateWay()
	server.SetupContract()

	server.SetupRouter()

	return server
}

func (server *MyServer) SetupRouter() {
	router := gin.Default()

	router.GET("/user/GetUser", server.GetUser)
	router.POST("/user/CreateUser", server.CreateUser)
	router.GET("/user/DecryptUserHistory", server.DecryptUserHistory)

	router.POST("/nft/MintNFT", server.MintNFT)
	router.GET("/nft/GetNFTByNID", server.GetNFTByNID)
	router.GET("/nft/GetNFTByUIDCreate", server.GetNFTByUIDCreate)
	router.GET("/nft/GetNFTByUIDOwn", server.GetNFTByUIDOwn)
	router.POST("/nft/ChangeNFTPrice", server.ChangeNFTPrice)
	router.POST("/nft/RegisterNFTOp", server.RegisterNFTOp)
	router.GET("/nft/DecryptNFTHistory", server.DecryptNFTHistory)
	router.POST("/nft/ChangeNFTStatus", server.ChangeNFTStatus)
	router.POST("/nft/ChangeNFTOwner", server.ChangeNFTOwner)

	router.POST("/cert/ApplyCert", server.ApplyCert)
	router.GET("/cert/GetCert", server.GetCert)
	router.GET("/cert/VerifyCert", server.VerifyCert)

	router.GET("/tx/GetTxInfo", server.GetTxINfo)

	server.router = router
}

func (server *MyServer) SetupWallet() {
	wallet := chaincode.SetupWallet()

	server.wallet = wallet
}

func (server *MyServer) SetupGateWay() {
	gw := chaincode.SetupGateWay(server.wallet)

	server.gw = gw
}

func (server *MyServer) SetupContract() {
	contract, qscc := chaincode.SetupContract(server.gw)

	server.contract = contract
	server.qscc = qscc
}

func (server *MyServer) Start(addr string) {
	defer server.gw.Close()

	server.router.Run(addr)
}

func ErrorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
