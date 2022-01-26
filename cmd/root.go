package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"git.cryptic.systems/volker.raschek/civ/pkg/config"
	"git.cryptic.systems/volker.raschek/civ/pkg/docker"
	"git.cryptic.systems/volker.raschek/civ/pkg/usecases"
	"git.cryptic.systems/volker.raschek/dockerutils"
	"github.com/spf13/cobra"
)

// Execute a
func Execute(version string) error {
	rootCmd := &cobra.Command{
		Use:     "civ",
		Short:   "go container label checker",
		Version: version,
		Args:    cobra.MinimumNArgs(1),
		RunE:    runE,
	}

	err := rootCmd.Execute()
	if err != nil {
		return err
	}

	return nil
}

func runE(cmd *cobra.Command, args []string) error {
	readConfigFile := args[0]
	var writeConfigFile string

	if len(args) == 2 {
		writeConfigFile = args[1]
	} else {
		writeConfigFile = fmt.Sprintf("config_result%s", filepath.Ext(readConfigFile))
	}

	if _, err := os.Stat(readConfigFile); os.IsNotExist(err) {
		return err
	}

	fileReader := config.NewFileReader(readConfigFile)
	cnf, err := fileReader.ReadFile()
	if err != nil {
		return err
	}

	dockerClient, err := dockerutils.New()
	if err != nil {
		return err
	}

	dockerRuntime, err := docker.NewRuntime(dockerClient)
	if err != nil {
		return err
	}

	labelVerifier, err := usecases.NewLabelVerifier(cnf, dockerRuntime)
	if err != nil {
		return err
	}

	err = labelVerifier.Run(cmd.Context())
	if err != nil {
		return err
	}

	fileWriter := config.NewFileWriter(writeConfigFile)
	err = fileWriter.WriteFile(cnf)
	if err != nil {
		return err
	}

	return nil
}
