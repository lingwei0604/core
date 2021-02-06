package core

import (
	stdlog "log"

	"github.com/DoNewsCode/std/pkg/config"
	"github.com/DoNewsCode/std/pkg/contract"
	"github.com/DoNewsCode/std/pkg/event"
	"github.com/DoNewsCode/std/pkg/logging"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/rawbytes"
	"go.uber.org/dig"
)

const defaultConfig = `
name: skeleton
version: 0.1.0
Env: local
http:
  addr: :8080
grpc:
  addr: :9090
log:
  level: debug
  format: logfmt
redis:
  addrs:
    - 127.0.0.1:6379
  database: 0
gorm:
  database: mysql
  dsn: root@tcp(127.0.0.1:3306)/skeleton?charset=utf8mb4&parseTime=True&loc=Local
`

func ProvideConfig(configStack []config.ProviderSet, configWatcher contract.ConfigWatcher) contract.ConfigAccessor {
	var (
		stack []config.Option
		err   error
		cfg   contract.ConfigAccessor
	)

	for _, layer := range configStack {
		stack = append(stack, config.WithProviderLayer(layer.Provider, layer.Parser))
	}
	stack = append(stack, config.WithProviderLayer(rawbytes.Provider([]byte(defaultConfig)), yaml.Parser()))
	if configWatcher != nil {
		stack = append(stack, config.WithWatcher(configWatcher))
	}

	cfg, err = config.NewConfig(stack...)
	if err != nil {
		stdlog.Fatal(err)
	}
	return cfg
}

func ProvideEnv(conf contract.ConfigAccessor) contract.Env {
	var env config.Env
	err := conf.Unmarshal("Env", &env)
	if err != nil {
		return config.NewEnv("local")
	}
	return env
}

func ProvideAppName(conf contract.ConfigAccessor) contract.AppName {
	var appName config.AppName
	err := conf.Unmarshal("name", &appName)
	if err != nil {
		return config.AppName("default")
	}
	return appName
}

func ProvideLogger(conf contract.ConfigAccessor, appName contract.AppName, env contract.Env) log.Logger {
	var (
		lvl    string
		format string
		err    error
	)
	err = conf.Unmarshal("level", &lvl)
	if err != nil {
		lvl = "debug"
	}
	err = conf.Unmarshal("format", &format)
	if err != nil {
		format = "logfmt"
	}
	logger := logging.NewLogger(format)
	logger = level.NewInjector(logger, level.DebugValue())
	return level.NewFilter(logger, logging.LevelFilter(lvl))
}

func ProvideDi(conf contract.ConfigAccessor) *dig.Container {
	return dig.New()
}

func ProvideEventDispatcher(conf contract.ConfigAccessor) contract.Dispatcher {
	return &event.SyncDispatcher{}
}
