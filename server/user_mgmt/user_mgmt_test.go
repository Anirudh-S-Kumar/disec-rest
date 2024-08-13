package user_mgmt_test

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Anirudh-S-Kumar/disec/server"
	"github.com/Anirudh-S-Kumar/disec/types"
	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestDuplicateUser(t *testing.T) {
	assert := assert.New(t)

	gin.SetMode(gin.TestMode)
	router := server.NewRouter(true)

	req := types.RegisterReq{
		Username: "test",
		Email:    "test@example.com",
		Password: "123",
	}
	reqBody, err := sonic.Marshal(req)
	assert.Nil(err)

	log.Println(req)

	httpreq := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(string(reqBody)))

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, httpreq)
	assert.Equal(http.StatusOK, rr.Code)

	testcases := []struct {
		name       string
		req        types.RegisterReq
		statuscode int
		err        string
	}{
		{
			name: "Duplicate username and email",
			req: types.RegisterReq{
				Username: "test",
				Email:    "test@example.com",
				Password: "123",
			},
			statuscode: http.StatusBadRequest,
			err:        "username already exists",
		},
		{
			name: "Duplicate username only",
			req: types.RegisterReq{
				Username: "test",
				Email:    "test1@example.com",
				Password: "123",
			},
			statuscode: http.StatusBadRequest,
			err:        "username already exists",
		},
		{
			name: "Duplicate email only",
			req: types.RegisterReq{
				Username: "test1",
				Email:    "test@example.com",
				Password: "123",
			},
			statuscode: http.StatusBadRequest,
			err:        "email already exists",
		},
		{
			name: "Different user registering",
			req: types.RegisterReq{
				Username: "test2",
				Email:    "test2@example.com",
				Password: "123",
			},
			statuscode: http.StatusOK,
			err:        "",
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, err := json.Marshal(tt.req)
			assert.Nil(err)

			rr = httptest.NewRecorder()

			httpreq, err = http.NewRequest(http.MethodPost, "/register", strings.NewReader(string(reqBody)))
			assert.Nil(err)
			router.ServeHTTP(rr, httpreq)
			assert.Equal(tt.statuscode, rr.Code)

			respBody, err := io.ReadAll(rr.Body)
			assert.Nil(err)

			var resp struct {
				Error string `json:"error"`
			}
			err = json.Unmarshal(respBody, &resp)
			assert.Nil(err)

			assert.Equal(tt.err, resp.Error)
		})
	}

}
