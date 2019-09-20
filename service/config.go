package service

import (
	"github.com/sirupsen/logrus"
)

type Config struct {
	Name    string `envconfig:"-"`
	Version string `envconfig:"-"`

	RPCAddr         string `default:"0.0.0.0:5002" envconfig:"RPC_ADDR"`
	PrivateRPCAddr  string `default:"0.0.0.0:5102" envconfig:"PRIVATE_RPC_ADDR"`
	AccountsRPCAddr string `default:"0.0.0.0:5001" envconfig:"ACCOUNTS_RPC_ADDR"`
	EmitterRPCAddr  string `default:"0.0.0.0:5003" envconfig:"EMITTER_RPC_ADDR"`

	DBURI string `default:"root:root@/videocoin?charset=utf8&parseTime=True&loc=Local" envconfig:"DBURI"`

	AuthTokenSecret string `default:"" envconfig:"AUTH_TOKEN_SECRET"`

	BaseInputURL  string `default:"" envconfig:"BASE_INPUT_URL"`
	BaseOutputURL string `default:"" envconfig:"BASE_OUTPUT_URL"`

	Logger *logrus.Entry `envconfig:"-"`
}
