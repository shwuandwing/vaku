package cmd

import (
	"fmt"

	"github.com/lingrino/vaku/vaku"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var folderCopy_V2Cmd = &cobra.Command{
	Use:   "copy_v2 [source folder] [target path]",
	Short: "Copy a vault folder from one location to another",
	Long: `Takes in a source path and target path and copies every path in the source to the target.
Note that this will copy the input path if it is a secret and all paths under the input path that
result from calling 'vaku folder list' on that path. Also note that this will overwrite any existing
keys at the target paths.

Example:
  vaku folder copy_v2 secret/foo secret/bar`,

	Args: cobra.ExactArgs(2),

	Run: func(cmd *cobra.Command, args []string) {

		authCopyClient()

		inputSource := vaku.NewPathInput(args[0])
		inputTarget := vaku.NewPathInput(args[1])

		err := vgc.FolderCopy(inputSource, inputTarget)
		if err != nil {
			fmt.Printf("%s", errors.Wrapf(err, "Failed to copy folder %s to %s", args[0], args[1]))
		} else {
			print(map[string]interface{}{
				args[0]: fmt.Sprintf("Successfully copied folder %s to %s", args[0], args[1]),
			})
		}
	},
}

func init() {
	folderCmd.AddCommand(folderCopy_V2Cmd)
	folderCopy_V2Cmd.Flags().StringVarP(&sourceAddress, "source-address", "", "", "The address for the source folder")
	folderCopy_V2Cmd.Flags().StringVarP(&sourceNamespace, "source-namespace", "", "", "The namespace for the source folder")
	folderCopy_V2Cmd.Flags().StringVarP(&sourceToken, "source-token", "", "", "The token for the source folder")
	folderCopy_V2Cmd.Flags().StringVarP(&targetAddress, "target-address", "", "", "The address for the target folder")
	folderCopy_V2Cmd.Flags().StringVarP(&targetNamespace, "target-namespace", "", "", "The namespace for the target folder")
	folderCopy_V2Cmd.Flags().StringVarP(&targetAddress, "target-token", "", "", "The token for the target folder")
}
