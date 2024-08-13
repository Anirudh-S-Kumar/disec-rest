package initializer

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"

	db_server "github.com/Anirudh-S-Kumar/disec/common/database/server"
	"github.com/gin-gonic/gin"
)

func SetupPrivateKey() *ecdsa.PrivateKey {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get wd: %+v", err)
	}
	keyData, err := os.ReadFile(cwd + "/certs/signing/ecdsa_private_key.pem")
	if err != nil {
		log.Fatalf("failed to load private key")
	}
	block, _ := pem.Decode(keyData)
	if block == nil || block.Type != "EC PRIVATE KEY" {
		log.Fatalf("block is not private key")
	}

	key, _ := x509.ParseECPrivateKey(block.Bytes)
	return key
}

func CertPaths() (string, string) {
	cwd, _ := os.Getwd()

	certPath := cwd + "/certs/tls/cert.pem"
	keyPath := cwd + "/certs/tls/key.pem"

	return certPath, keyPath
}

func DBMiddleWare(db *db_server.ServerDB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set("server_db", db)
		ctx.Next()
	}
}
