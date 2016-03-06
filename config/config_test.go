package config

import (
	"testing"
	"github.com/spf13/viper"
)

func TestReadConfig(t *testing.T){
	readConfig(".")

	isSet := viper.IsSet(GRPC_SERVER_PORT)
	if !isSet {
		t.Fatal("no rpc_port variable found")
	}
	port := viper.GetInt(GRPC_SERVER_PORT)
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

	isSet = viper.IsSet(SERF_RPC)
	if !isSet {
		t.Fatal("Serf rpc port is not set")
	}

	port = viper.GetInt(SERF_RPC)
	if port != 7373 {
		t.Fatalf("Serf port is not 7373: %d\n", port)
	}
}
