package service

import (
	"time"
)

type Config struct {
	Name    string `envconfig:"-"`
	Version string `envconfig:"-"`

	RPCAddr         string `default:"0.0.0.0:5002" envconfig:"RPC_ADDR"`
	PrivateRPCAddr  string `default:"0.0.0.0:5102" envconfig:"PRIVATE_RPC_ADDR"`
	UsersRPCAddr    string `default:"0.0.0.0:5000" envconfig:"USERS_RPC_ADDR"`
	AccountsRPCAddr string `default:"0.0.0.0:5001" envconfig:"ACCOUNTS_RPC_ADDR"`
	EmitterRPCAddr  string `default:"0.0.0.0:5003" envconfig:"EMITTER_RPC_ADDR"`
	BillingRPCAddr  string `default:"0.0.0.0:5120" envconfig:"BILLING_RPC_ADDR"`

	DBURI    string `default:"root:root@/videocoin?charset=utf8&parseTime=True&loc=Local" envconfig:"DBURI"`
	MQURI    string `default:"amqp://guest:guest@127.0.0.1:5672" envconfig:"MQURI"`
	RedisURI string `default:"redis://:@127.0.0.1:6379/1" envconfig:"REDISURI"`

	AuthTokenSecret string `required:"true" envconfig:"AUTH_TOKEN_SECRET"`

	BaseInputURL  string `required:"true" envconfig:"BASE_INPUT_URL"`
	BaseOutputURL string `required:"true" envconfig:"BASE_OUTPUT_URL"`
	RTMPURL       string `required:"true" envconfig:"RTMP_URL"`

	MaxLiveStreamTime time.Duration `required:"true" envconfig:"MAX_LIVESTREAM_TIME" default:"43200s"`
}
