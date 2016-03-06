package config

import (
	"testing"
	"github.com/spf13/viper"
)

func TestReadConfig(t *testing.T){
	ReadConfig()

	isSet := viper.IsSet(RPC_PORT)
	if !isSet {
		t.Fatal("no rpc_port variable found")
	}
	port := viper.GetInt(RPC_PORT)
	if port != 5123 {
		t.Fatal("rpc port was not 5123")
	}

	isSet = viper.IsSet(JOIN_DELAY)
	if !isSet {
		t.Fatal("no join delay found")
	}
	joinDelay := viper.GetInt(JOIN_DELAY)
	if joinDelay != 5 {
		t.Fatal("Error getting join delay")
	}
}
