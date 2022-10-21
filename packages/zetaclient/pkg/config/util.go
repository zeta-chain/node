package config

// GetChainConfig returns a chain configuration
func GetChainConfig(cfg *Configuration, chain string) *ChainConfig {
	switch chain {
	case "GOERLI":
		return &cfg.Goerli
	case "ROPSTEN":
		return &cfg.Ropsten
	case "MUMBAI":
		return &cfg.Mumbai
	case "BAOBAB":
		return &cfg.Baobab
	case "BITCOIN":
		return &cfg.Bitcoin
	}
	return nil
}
