package cli

import (
	"fmt"

	"github.com/0xfe/microstellar"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func (cli *CLI) getAssetCmd() *cobra.Command {
	accountsCmd := &cobra.Command{
		Use:   "asset [set|del|code|issuer|type]",
		Short: "manage stellar assets",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				showError(logrus.Fields{"cmd": "asset"}, "unrecognized command: %s, expecting: set|del|code|issuer|type", args[0])
				return
			}
		},
	}

	accountsCmd.AddCommand(cli.getAssetSetCmd())
	accountsCmd.AddCommand(cli.getAssetCodeCmd())
	accountsCmd.AddCommand(cli.getAssetIssuerCmd())
	accountsCmd.AddCommand(cli.getAssetTypeCmd())
	accountsCmd.AddCommand(cli.getAccountSeedCmd())

	return accountsCmd
}

func (cli *CLI) getAssetSetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set [name] [issuer]",
		Short: "set asset issuer of asset [name]",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			issuer := args[1]
			code := "XLM"
			assetType := "native"

			for _, part := range []string{"issuer", "code", "type"} {
				key := fmt.Sprintf("asset:%s:%s", name, part)
				value := ""

				if part == "issuer" {
					value = issuer
				}

				if part == "code" {
					if cmd.Flag("code").Changed {
						code, _ = cmd.Flags().GetString("code")
					} else {
						code = name
					}
					value = code
				}

				if part == "type" {
					if cmd.Flag("type").Changed {
						assetType, _ = cmd.Flags().GetString("type")
						switch assetType {
						case "credit4", "credit12", "native":
							break
						default:
							showError(logrus.Fields{"cmd": "asset", "subcmd": "set"}, "bad asset type: %s", assetType)
							return
						}
					} else {
						if len(code) > 4 {
							assetType = "credit12"
						} else {
							assetType = "credit4"
						}
					}

					value = assetType
				}

				logrus.WithFields(logrus.Fields{"cmd": "asset", "subcmd": "set"}).Debugf("saving asset %s: %s %s", name, part, value)
				err := cli.SetVar(key, value)

				if err != nil {
					logrus.WithFields(logrus.Fields{"cmd": "asset", "subcmd": "set"}).Debugf("%v", err)
					showError(logrus.Fields{"cmd": "asset", "subcmd": "set"}, "could not save asset: %s", name)
					return
				}
			}
		},
	}

	cmd.Flags().String("code", "XLM", "specify asset code")
	cmd.Flags().String("type", "native", "specify asset type (credit4, credit12, native)")

	return cmd
}

func (cli *CLI) getAssetCodeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "code [name]",
		Short: "get asset code of [name]",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			if asset, err := cli.GetAsset(name); err != nil {
				logrus.WithFields(logrus.Fields{"cmd": "asset", "subcmd": "code"}).Debugf("%v", err)
				showError(logrus.Fields{"cmd": "asset", "subcmd": "code"}, "could not load asset: %s", name)
			} else {
				showSuccess(asset.Code)
			}
		},
	}

	return cmd
}

func (cli *CLI) getAssetIssuerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issuer [name]",
		Short: "get asset issuer of [name]",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			if asset, err := cli.GetAsset(name); err != nil {
				logrus.WithFields(logrus.Fields{"cmd": "asset", "subcmd": "issuer"}).Debugf("%v", err)
				showError(logrus.Fields{"cmd": "asset", "subcmd": "issuer"}, "could not load asset: %s", name)
			} else {
				showSuccess(asset.Issuer)
			}
		},
	}

	return cmd
}

func (cli *CLI) getAssetTypeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "type [name]",
		Short: "get asset type of [name]",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			if asset, err := cli.GetAsset(name); err != nil {
				logrus.WithFields(logrus.Fields{"cmd": "asset", "subcmd": "type"}).Debugf("%v", err)
				showError(logrus.Fields{"cmd": "asset", "subcmd": "type"}, "could not load asset: %s", name)
			} else {
				assetType := "native"
				if asset.Type == microstellar.Credit4Type {
					assetType = "credit4"
				} else if asset.Type == microstellar.Credit12Type {
					assetType = "credit12"
				}
				showSuccess(assetType)
			}
		},
	}

	return cmd
}