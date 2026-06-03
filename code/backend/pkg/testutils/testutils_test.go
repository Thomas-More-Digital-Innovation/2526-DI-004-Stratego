package testutils_test

import (
	"digital-innovation/gostrategy/pkg/testutils"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetupHandlerTest(t *testing.T) {
	r, h, server := testutils.SetupHandlerTest(t)
	require.NotNil(t, r)
	require.NotNil(t, h)
	require.NotNil(t, server)
}

func TestPerformRequest(t *testing.T) {
	r := gin.New()
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	w := testutils.PerformRequest(r, "GET", "/ping", nil, nil)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}

func TestPerformRequestWithBody(t *testing.T) {
	r := gin.New()
	type reqStruct struct {
		Message string `json:"message"`
	}
	r.POST("/echo", func(c *gin.Context) {
		var body reqStruct
		if err := c.ShouldBindJSON(&body); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		c.JSON(http.StatusOK, body)
	})

	body := reqStruct{Message: "hello"}
	w := testutils.PerformRequest(r, "POST", "/echo", body, map[string]string{"X-Test-Header": "yes"})
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "hello")
}

func TestSetupDBTest(t *testing.T) {
	db := testutils.SetupDBTest(t)
	require.NotNil(t, db)
}

func TestSetupTestGame(t *testing.T) {
	g, p1, p2 := testutils.SetupTestGame()
	require.NotNil(t, g)
	require.NotNil(t, p1)
	require.NotNil(t, p2)

	assert.Equal(t, testutils.Player1Name, p1.GetName())
	assert.Equal(t, testutils.Player2Name, p2.GetName())
	assert.Equal(t, 0, p1.GetID())
	assert.Equal(t, 1, p2.GetID())
}
