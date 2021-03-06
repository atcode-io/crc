package cmd

import (
	"fmt"
	"os"

	"github.com/code-ready/crc/pkg/crc/constants"
	"github.com/code-ready/crc/pkg/crc/exit"
	"github.com/code-ready/crc/pkg/crc/input"
	"github.com/code-ready/crc/pkg/crc/machine"
	"github.com/code-ready/crc/pkg/crc/output"
	"github.com/spf13/cobra"
)

func init() {
	deleteCmd.Flags().BoolVarP(&clearCache, "clear-cache", "", false,
		fmt.Sprintf("Clear the OpenShift cluster cache at: %s", constants.MachineCacheDir))
	rootCmd.AddCommand(deleteCmd)
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete the OpenShift cluster",
	Long:  "Delete the OpenShift cluster",
	Run: func(cmd *cobra.Command, args []string) {
		if err := runDelete(args); err != nil {
			exit.WithMessage(1, err.Error())
		}
	},
}

var clearCache bool

func runDelete(arguments []string) error {
	deleteConfig := machine.DeleteConfig{
		Name: constants.DefaultName,
	}

	if clearCache {
		deleteCache()
	}

	client := machine.NewClient()
	if err := checkIfMachineMissing(client, deleteConfig.Name); err != nil {
		return err
	}

	yes := input.PromptUserForYesOrNo("Do you want to delete the OpenShift cluster", globalForce)
	if yes {
		_, err := client.Delete(deleteConfig)
		if err != nil {
			return err
		}
		output.Outln("Deleted the OpenShift cluster")
	}
	return nil
}

func deleteCache() {
	yes := input.PromptUserForYesOrNo("Do you want to delete the OpenShift cluster cache", globalForce)
	if yes {
		os.RemoveAll(constants.MachineCacheDir)
	}
}
