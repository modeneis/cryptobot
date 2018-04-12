package markets_test

import (
	"testing"

	"github.com/modeneis/cryptobot/src/markets"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/thrasher-/gocryptotrader/config"
)

func init() {
	SetupConfig = config.GetConfig()
	//bittrexSetupConfig.EncryptConfig = -1
	SetupConfig.LoadConfig("../config.dat")
}
func TestGoGeANXtMarkets(t *testing.T) {

	Convey("Given I go get all available markets for ANX ", t, func() {

		//TODO: assert response
		markets.GoGetMarkets("ANX", SetupConfig)
	})
}
