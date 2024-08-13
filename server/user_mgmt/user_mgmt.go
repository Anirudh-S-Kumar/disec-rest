package user_mgmt

import (
	"net/http"

	db_server "github.com/Anirudh-S-Kumar/disec/common/database/server"
	"github.com/Anirudh-S-Kumar/disec/common/utils"
	"github.com/Anirudh-S-Kumar/disec/types"
	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	db := c.MustGet("server_db").(*db_server.ServerDB)
	var req types.RegisterReq
	if err := utils.ProcessRequest(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Username == "" || req.Email == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username/email/password is empty"})
	}

	if err := db.UserExists(req.Username, req.Email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.AddUser(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	resp := &types.EmptyResp{}

	if err := utils.ProcessResponse(c, resp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}
