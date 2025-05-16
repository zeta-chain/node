package main

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/zeta-chain/node/app"
	"github.com/zeta-chain/node/pkg/parsers"
)

func CollectObserverInfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collect-observer-info [folder]",
		Short: "collect observer info into the genesis from a folder , default path is ~/.zetacored/os_info/ \n",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			defaultHome := app.DefaultNodeHome
			defaultFile := filepath.Join(defaultHome, "os_info")
			if len(args) == 0 {
				args = append(args, defaultFile)
			}
			directory := args[0]
			files, err := os.ReadDir(directory)
			if err != nil {
				return err
			}
			var observerInfoList []parsers.ObserverInfoReader
			err = os.Chdir(directory)
			if err != nil {
				return err
			}
			for _, file := range files {
				var observerInfo parsers.ObserverInfoReader
				info, err := file.Info()
				if err != nil {
					return err
				}
				f, err := os.ReadFile(info.Name())
				if err != nil {
					return err
				}
				err = json.Unmarshal(f, &observerInfo)
				if err != nil {
					return err
				}
				observerInfoList = append(observerInfoList, observerInfo)
			}
			file, err := json.MarshalIndent(observerInfoList, "", " ")
			if err != nil {
				return err
			}
			err = os.WriteFile("observer_info.json", file, 0600)
			if err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}
