package config

import (
	"github.com/spf13/viper"
	log "github.com/Sirupsen/logrus"
)

//Global
const HTTP_SERVER_PORT string = "http_server_port"
const GRPC_SERVER_PORT string = "grpc_server_port"
const DEFAULT_TOPIC string = "default_topic"

//Delays
const JOIN_DELAY string = "delays.join"
const SHUTDOWN_DELAY string = "delays.shutdown"
const GRPC_SERVER_START_DELAY string = "delays.grpc_server_start"

//Serf
const SERF_PORT string = "serf.port"
const SERF_RPC string = "serf.rpc"

//Refresh loops
const MEMBER_LIST_REFRESH_SECONDS string = "refresh_loops_delay_seconds.member_list"

//Listeners
const WRITE_TO_FILE_DELAY string = "listeners.file.write_to_file_delay"

//Performance queues
const MESSAGE_SIZE string = "performance_queues.message_size"
const MESSAGE_CLUSTER_SIZE string = "performance_queues.message_cluster_size"

func init(){
	readConfig("config")
}

func readConfig(path string){
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Error("Couldn't read config file: ", err)
	}
}
