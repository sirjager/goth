package tokens

import (
	"os"
	"testing"

	"github.com/sirjager/goth/utils"
)

var tokenBuilder TokenBuilder

func TestMain(t *testing.M) {
	randomSecretKey := utils.RandomString(32)
	builder, err := NewPasetoBuilder(randomSecretKey)
	if err != nil {
		panic(err)
	}
	tokenBuilder = builder
	os.Exit(t.Run())
}
