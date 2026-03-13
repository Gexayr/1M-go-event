package rules

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&RiskRule{})
	return db
}

func TestRuleHandlers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	database := setupTestDB()
	repo := NewRepository(database)
	service := NewService(repo)
	handler := NewHandler(service)

	r := gin.Default()
	RegisterRoutes(r, handler)

	t.Run("Create Rule", func(t *testing.T) {
		reqBody := CreateRuleRequest{
			Name:              "Large Withdrawal",
			EventType:         "withdrawal",
			ConditionField:    "amount",
			ConditionOperator: ">",
			ConditionValue:    "5000",
			Score:             70,
			Enabled:           true,
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/api/rules", bytes.NewBuffer(body))
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusCreated, resp.Code)

		var rule RuleResponse
		json.Unmarshal(resp.Body.Bytes(), &rule)
		assert.Equal(t, "Large Withdrawal", rule.Name)
		assert.Equal(t, uint(1), rule.ID)
	})

	t.Run("Get All Rules", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/rules", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		var rules []RuleResponse
		json.Unmarshal(resp.Body.Bytes(), &rules)
		assert.Len(t, rules, 1)
	})

	t.Run("Toggle Rule", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/rules/1/toggle", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		var rule RuleResponse
		json.Unmarshal(resp.Body.Bytes(), &rule)
		assert.False(t, rule.Enabled)
	})

	t.Run("Delete Rule", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/api/rules/1", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNoContent, resp.Code)

		// Verify it's gone
		var count int64
		database.Model(&RiskRule{}).Count(&count)
		assert.Equal(t, int64(0), count)
	})
}
