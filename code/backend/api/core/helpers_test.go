package core_test

import (
	"digital-innovation/gostrategy/api/core"
	"digital-innovation/gostrategy/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHelpers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("SendError", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		core.SendError(c, "Test error", http.StatusBadRequest)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("SendJSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		core.SendJSON(c, gin.H{"test": "ok"}, http.StatusOK)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("EnsureAuthenticated", func(t *testing.T) {
		t.Run("Unauthorized", func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			user := core.EnsureAuthenticated(c)
			assert.Nil(t, user)
			assert.Equal(t, http.StatusUnauthorized, w.Code)
		})

		t.Run("Authorized", func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			userObj := &models.User{ID: 1, Username: "test"}
			c.Set("user", userObj)
			user := core.EnsureAuthenticated(c)
			assert.Equal(t, userObj, user)
		})
	})

	t.Run("ParseID", func(t *testing.T) {
		t.Run("FromParam", func(t *testing.T) {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Params = []gin.Param{{Key: "id", Value: "123"}}
			id, err := core.ParseID(c, "id")
			assert.NoError(t, err)
			assert.Equal(t, 123, id)
		})

		t.Run("FromQuery", func(t *testing.T) {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Request, _ = http.NewRequest("GET", "/?id=456", nil)
			id, err := core.ParseID(c, "id")
			assert.NoError(t, err)
			assert.Equal(t, 456, id)
		})
	})
}
