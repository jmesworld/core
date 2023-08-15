package main

import (
	serverconfig "github.com/cosmos/cosmos-sdk/server/config"
	"github.com/ignite/cli/app/wasmconfig"
)

// JMESAppConfig jmes specify app config
type JMESAppConfig struct {
	serverconfig.Config

	WASMConfig wasmconfig.Config `mapstructure:"wasm"`
}

// initAppConfig helps to override default appConfig template and configs.
// return "", nil if no custom configuration is required for the application.
func initAppConfig() (string, interface{}) {
	// Optionally allow the chain developer to overwrite the SDK's default
	// server config.
	srvCfg := serverconfig.DefaultConfig()

	// The SDK's default minimum gas price is set to "" (empty value) inside
	// app.toml. If left empty by validators, the node will halt on startup.
	// However, the chain developer can set a default app.toml value for their
	// validators here.
	//
	// In summary:
	// - if you leave srvCfg.MinGasPrices = "", all validators MUST tweak their
	//   own app.toml config,
	// - if you set srvCfg.MinGasPrices non-empty, validators CAN tweak their
	//   own app.toml to override, or use this default value.
	//
	// In simapp, we set the min gas prices to 0.
	srvCfg.MinGasPrices = "0ujmes"

	jmesAppConfig := JMESAppConfig{
		Config:     *srvCfg,
		WASMConfig: *wasmconfig.DefaultConfig(),
	}

	jmesAppTemplate := serverconfig.DefaultConfigTemplate + wasmconfig.DefaultConfigTemplate

	return jmesAppTemplate, jmesAppConfig
}
