package main

import (
	"context"
	"io/ioutil"

	"github.com/pelletier/go-toml"
	"github.com/rekby/zapcontext"
	"go.uber.org/zap"
)

type configType struct {
	IssueTimeout           int    `default:"300" comment:"Seconds for issue every certificate. Cancel issue and return error if timeout."`
	AutoIssueForSubdomains string `default:"www" comment:"Comma separated for subdomains for try get common used subdomains in one certificate."`
	HttpsListeners         string `default:"[::]:443" comment:"Comma-separated bindings for listen https.\nSupported formats:\n1.2.3.4:443,0.0.0.0:443,[::]:443,[2001:db8::a123]:443"`
	StorageDir             string `default:"storage" comment:"Path to dir, which will store state and certificates"`
	AcmeServer             string `default:"https://acme-v01.api.letsencrypt.org/directory" comment:"Directory url of acme server.\nTest server: https://acme-staging-v02.api.letsencrypt.org/directory"`
}

var (
	_config *configType
)

func getConfig(ctx context.Context) *configType {
	if _config == nil {
		logger := zc.LNop(ctx).With(zap.String("config_file", *configFileP))
		logger.Info("Read config")
		config, err := readConfig(ctx, *configFileP)
		if err == nil {
			_config = &config
		} else {
			logger.Fatal("Error while read config.")
		}
		applyFlags(ctx, _config)
	}
	return _config
}

// Apply command line flags to config
func applyFlags(ctx context.Context, config *configType) {

}

func defaultConfig(ctx context.Context) []byte {
	config, _ := readConfig(ctx, "")
	configBytes, _ := toml.Marshal(&config)
	return configBytes
}

func readConfig(ctx context.Context, file string) (configType, error) {
	logger := zc.LNop(ctx).With(zap.String("config_file", file))
	var fileBytes []byte
	var err error
	if file == "" {
		logger.Info("Use default config.")
	} else {
		fileBytes, err = ioutil.ReadFile(file)
	}
	if err != nil {
		logger.Error("Can't read config", zap.Error(err))
		return configType{}, err
	}

	var res configType
	err = toml.Unmarshal(fileBytes, &res)
	if err != nil {
		logger.Error("Can't unmarshal config.", zap.Error(err))
		return configType{}, err
	}

	readedConfig, err := toml.Marshal(res)
	if err == nil {
		logger.Info("Read config.", zap.ByteString("config_content", readedConfig))
	} else {
		logger.Error("Can't marshal config", zap.Error(err))
	}
	return res, nil
}
