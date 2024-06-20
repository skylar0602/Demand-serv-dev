package strategy

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCrosschain(t *testing.T) {
	ok, err := ableToCrossChain("mumbai", "fuji", "USDC")
	assert.Nil(t, err)
	assert.Truef(t, ok, "")
}
