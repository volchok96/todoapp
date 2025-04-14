package domain_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/volchok96/todoapp/internal/domain"
)

func TestDateOnly_MarshalUnmarshal(t *testing.T) {
	dateStr := `"2024-04-01"`
	var d domain.DateOnly

	err := json.Unmarshal([]byte(dateStr), &d)
	assert.NoError(t, err)
	assert.Equal(t, "2024-04-01", d.Format("2006-01-02"))

	bytes, err := json.Marshal(d)
	assert.NoError(t, err)
	assert.JSONEq(t, dateStr, string(bytes))
}

func TestDateOnly_UnmarshalInvalid(t *testing.T) {
	var d domain.DateOnly
	err := json.Unmarshal([]byte(`"invalid-date"`), &d)
	assert.Error(t, err)
}
