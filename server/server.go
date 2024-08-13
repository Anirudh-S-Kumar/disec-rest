package server

import (
	"log"
	"net/http"

	"github.com/Anirudh-S-Kumar/disec/cmd/server/initializer"
	db_server "github.com/Anirudh-S-Kumar/disec/common/database/server"
	"github.com/Anirudh-S-Kumar/disec/server/user_mgmt"
	"github.com/gin-gonic/gin"
)

func NewRouter(testMode bool) *gin.Engine {
	r := gin.Default()

	var db *db_server.ServerDB
	var err error

	if testMode {
		db, err = db_server.CreateTempDB("/home/afterchange/goland_projects/disec-rest/data/schema/server.sql")
	} else {
		db, err = db_server.CreateDB("data/db/server.db", "data/schema/server.sql", true)
	}
	if err != nil {
		log.Fatalf("cannot create database: %v", err)
	}

	r.Use(initializer.DBMiddleWare(db))
	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "Hello world!")
	})
	r.POST("/register", user_mgmt.Register)

	return r
}
