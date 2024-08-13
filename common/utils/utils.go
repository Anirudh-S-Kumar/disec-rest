package utils

import (
	"fmt"
	"io"
	"net/http"

	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
)

func ProcessRequest(c *gin.Context, req any) error {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return fmt.Errorf("error reading request body: %v", err)
	}
	if err := sonic.Unmarshal(body, req); err != nil {
		return fmt.Errorf("error in unmarshaling: %v", err)
	}
	return nil
}

func ProcessResponse(c *gin.Context, resp any) error {
	respData, err := sonic.Marshal(resp)
	if err != nil {
		return fmt.Errorf("cannot marshal response")
	}
	c.Data(http.StatusOK, "application/json", respData)
	return nil
}
