package main

import (

	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	mspID        = "Org1MSP"
	cryptoPath   = "/home/srivani/workspace/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com"
	certPath     = cryptoPath + "/users/User1@org1.example.com/msp/signcerts/cert.pem"
	keyPath      = cryptoPath + "/users/User1@org1.example.com/msp/keystore/"
	tlsCertPath  = cryptoPath + "/peers/peer0.org1.example.com/tls/ca.crt"
	peerEndpoint = "localhost:7051"
	gatewayPeer  = "peer0.org1.example.com"
)

var now = time.Now()
var assetId = fmt.Sprintf("asset%d", now.Unix()*1e3+int64(now.Nanosecond())/1e6)

func main() {
	r := gin.Default()
    r.POST("/assets", CreateAssetHandler)
	r.GET("/assets", GetAllAssetsHandler)
	r.GET("/assets/:assetID", ReadAssetHandler)

	r.Run(":8080")
}


func CreateAssetHandler(c *gin.Context) {
	clientConnection := newGrpcConnection()
	defer clientConnection.Close()

	id := newIdentity()
	sign := newSign()

	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error connecting to Hyperledger Fabric: %s", err)})
		return
	}
	defer gw.Close()

	chaincodeName := "basic"
	if ccname := os.Getenv("CHAINCODE_NAME"); ccname != "" {
		chaincodeName = ccname
	}

	channelName := "mychannel"
	if cname := os.Getenv("CHANNEL_NAME"); cname != "" {
		channelName = cname
	}

	network := gw.GetNetwork(channelName)
	contract := network.GetContract(chaincodeName)

	// Get asset data from the request body
	var assetData map[string]interface{}
	if err := c.ShouldBindJSON(&assetData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Call createAsset function
	result, err := createAsset(contract, assetData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error creating asset: %s", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": result})
}

// createAsset submits a transaction to create a new asset.
func createAsset(contract *client.Contract, data map[string]interface{}) (string, error) {
	args := make([]string, 0)
	for _, v := range data {
		args = append(args, fmt.Sprint(v))
	}

	// Call the contract's SubmitTransaction method with the transaction name and string arguments
	_, err := contract.SubmitTransaction("CreateAsset", "asset11", "tom", "hyderabad", "9087654321", "CY", "delivered", "cargomovement", "servicereference", "servicename", "servicecode", "servicereference", "exportvoyagenumber", "voyagereference", "currency", "true", "true", "declarationreference", "true", "quotationreference", "channelreference", "terms", "true")
	if err != nil {
		return "", fmt.Errorf("failed to submit transaction: %w", err)
	}

	return "Transaction committed successfully", nil
}

// GetAllAssetsHandler handles the retrieval of all assets.
func GetAllAssetsHandler(c *gin.Context) {
	clientConnection := newGrpcConnection()
	defer clientConnection.Close()

	id := newIdentity()
	sign := newSign()

	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error connecting to Hyperledger Fabric: %s", err)})
		return
	}
	defer gw.Close()

	chaincodeName := "basic"
	if ccname := os.Getenv("CHAINCODE_NAME"); ccname != "" {
		chaincodeName = ccname
	}

	channelName := "mychannel"
	if cname := os.Getenv("CHANNEL_NAME"); cname != "" {
		channelName = cname
	}

	network := gw.GetNetwork(channelName)
	contract := network.GetContract(chaincodeName)

	// Call getAllAssets function
	result, err := getAllAssets(contract)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error getting all assets: %s", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": result})
}

// ReadAssetHandler handles the retrieval of asset attributes by assetID.
func ReadAssetHandler(c *gin.Context) {
	clientConnection := newGrpcConnection()
	defer clientConnection.Close()

	id := newIdentity()
	sign := newSign()

	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error connecting to Hyperledger Fabric: %s", err)})
		return
	}
	defer gw.Close()

	chaincodeName := "basic"
	if ccname := os.Getenv("CHAINCODE_NAME"); ccname != "" {
		chaincodeName = ccname
	}

	channelName := "mychannel"
	if cname := os.Getenv("CHANNEL_NAME"); cname != "" {
		channelName = cname
	}

	network := gw.GetNetwork(channelName)
	contract := network.GetContract(chaincodeName)

	// Get assetID from URL parameter
	assetID := c.Param("assetID")

	// Call readAssetByID function with the provided assetID
	result, err := readAssetByID(contract, assetID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error reading asset: %s", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": result})
}

// readAssetByID evaluates a transaction to query ledger state by assetID.
func readAssetByID(contract *client.Contract, assetID string) (string, error) {
	result, err := contract.EvaluateTransaction("ReadAsset", assetID)
	if err != nil {
		return "", fmt.Errorf("failed to evaluate transaction: %w", err)
	}

	return string(result), nil
}

// getAllAssets evaluates a transaction to query all assets on the ledger.
func getAllAssets(contract *client.Contract) (string, error) {
	result, err := contract.EvaluateTransaction("GetAllAssets")
	if err != nil {
		return "", fmt.Errorf("failed to evaluate transaction: %w", err)
	}

	return string(result), nil
}

// newGrpcConnection creates a gRPC connection to the Gateway server.
func newGrpcConnection() *grpc.ClientConn {
	certificate, err := loadCertificate(tlsCertPath)
	if err != nil {
		panic(err)
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(certificate)
	transportCredentials := credentials.NewClientTLSFromCert(certPool, gatewayPeer)

	connection, err := grpc.Dial(peerEndpoint, grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		panic(fmt.Errorf("failed to create gRPC connection: %w", err))
	}

	return connection
}

// newIdentity creates a client identity for this Gateway connection using an X.509 certificate.
func newIdentity() *identity.X509Identity {
	certificate, err := loadCertificate(certPath)
	if err != nil {
		panic(err)
	}

	id, err := identity.NewX509Identity(mspID, certificate)
	if err != nil {
		panic(err)
	}

	return id
}

func loadCertificate(filename string) (*x509.Certificate, error) {
	certificatePEM, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate file: %w", err)
	}
	return identity.CertificateFromPEM(certificatePEM)
}

// newSign creates a function that generates a digital signature from a message digest using a private key.
func newSign() identity.Sign {
	files, err := os.ReadDir(keyPath)
	if err != nil {
		panic(fmt.Errorf("failed to read private key directory: %w", err))
	}
	privateKeyPEM, err := os.ReadFile(path.Join(keyPath, files[0].Name()))

	if err != nil {
		panic(fmt.Errorf("failed to read private key file: %w", err))
	}

	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		panic(err)
	}

	sign, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		panic(err)
	}

	return sign
}
