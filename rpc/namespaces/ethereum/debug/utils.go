// Copyright 2021 Evmos Foundation
// This file is part of Evmos' Ethermint library.
//
// The Ethermint library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The Ethermint library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the Ethermint library. If not, see https://github.com/zeta-chain/ethermint/blob/main/LICENSE
package debug

import (
	"os"
	"runtime/pprof"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/server"

	zetaos "github.com/zeta-chain/node/pkg/os"
)

// isCPUProfileConfigurationActivated checks if cpuprofile was configured via flag
func isCPUProfileConfigurationActivated(ctx *server.Context) bool {
	// TODO: use same constants as server/start.go
	// constant declared in start.go cannot be imported (cyclical dependency)
	return ctx.Viper.GetString("cpu-profile") != ""
}

// writeProfile writes the data to a file
func writeProfile(name, file string, log log.Logger) error {
	p := pprof.Lookup(name)
	log.Info("Writing profile records", "count", p.Count(), "type", name, "dump", file)
	fp, err := zetaos.ExpandHomeDir(file)
	if err != nil {
		return err
	}
	// #nosec G304 variable value is controlled
	f, err := os.Create(fp)
	if err != nil {
		return err
	}

	if err := p.WriteTo(f, 0); err != nil {
		if err := f.Close(); err != nil {
			return err
		}
		return err
	}

	return f.Close()
}
