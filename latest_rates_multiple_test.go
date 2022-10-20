package exchange

import (
	"context"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_LatestRatesMultiple_OnBadCurrencyCode_ShouldReturnError(t *testing.T) {
	ex := New("USD")
	symbols := []string{"UNKNOWN"}

	rates, err := ex.LatestRatesMultiple(symbols)
	assert.IsType(t, map[string]float64{}, rates)
	assert.Len(t, rates, 0)
	assert.ErrorIs(t, err, ErrInvalidCode)
}

func Test_LatestRatesMultiple_ShouldReturnData(t *testing.T) {
	ex := New("USD")
	symbols := []string{"EUR", "JPY", "USD"}

	rates, err := ex.LatestRatesMultiple(symbols)
	assert.IsType(t, map[string]float64{}, rates)
	assert.Len(t, rates, len(symbols))
	assert.NoError(t, err)
}

func Test_LatestRatesMultiple_WithContext_ShouldReturnData(t *testing.T) {
	ctx, cancelFunction := context.WithCancel(context.Background())
	defer cancelFunction()

	ex := New("USD")
	ex.SetContext(ctx)
	symbols := []string{"EUR", "JPY", "USD"}

	rates, err := ex.LatestRatesMultiple(symbols)
	assert.IsType(t, map[string]float64{}, rates)
	assert.Len(t, rates, len(symbols))
	assert.NoError(t, err)
}

func Test_LatestRatesMultiple_WithCanceledContext_ShouldReturnError(t *testing.T) {
	ctx, cancelFunction := context.WithCancel(context.Background())
	cancelFunction()

	ex := New("USD")
	ex.SetContext(ctx)
	symbols := []string{"EUR", "JPY", "USD"}

	rates, err := ex.LatestRatesMultiple(symbols)
	assert.IsType(t, map[string]float64{}, rates)
	assert.Len(t, rates, 0)
	expectedErr := &url.Error{}
	assert.ErrorAs(t, err, &expectedErr)
}
