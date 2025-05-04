package cmd

import (
	"fmt"
	"github.com/P001water/P1finger/cmd/vars"
	cc "github.com/ivanpirog/coloredcobra"
	"github.com/projectdiscovery/gologger"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var version = "1.0.0" // 定义版本号

func init() {
	RootCmd.CompletionOptions.DisableDefaultCmd = true

	// 添加版本标志
	RootCmd.PersistentFlags().StringP("version", "v", "", "显示版本信息")
	RootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if versionFlag, _ := cmd.Flags().GetString("version"); versionFlag != "" {
			fmt.Printf("P1finger version: %s\n", version)
			return nil
		}
		return nil
	}

	RootCmd.Flags().StringVarP(&vars.Options.Proxy, "proxy", "p", "", "proxy eg: -proxy socks5://127.0.0.1")
	RootCmd.Flags().StringVarP(&vars.Options.Output, "output", "o", "p1finger.json", "output file name: [-o p1finger.xlsx] / [-o p1finger.json]")
	RootCmd.Flags().BoolVar(&vars.Options.Debug, "debug", false, "http debug info, eg:[-debug]")
}

var RootCmd = &cobra.Command{
	Use:   "P1finger",
	Short: "一款红队行动下的重点资产指纹识别工具",
	Long: `
	██████╗  ██╗███████╗██╗███╗   ██╗ ██████╗ ███████╗██████╗ 
	██╔══██╗███║██╔════╝██║████╗  ██║██╔════╝ ██╔════╝██╔══██╗
	██████╔╝╚██║█████╗  ██║██╔██╗ ██║██║  ███╗█████╗  ██████╔╝
	██╔═══╝  ██║██╔══╝  ██║██║╚██╗██║██║   ██║██╔══╝  ██╔══██╗
	██║      ██║██║     ██║██║ ╚████║╚██████╔╝███████╗██║  ██║
	╚═╝      ╚═╝╚═╝     ╚═╝╚═╝  ╚═══╝ ╚═════╝ ╚══════╝╚═╝  ╚═
 
        一款红队行动下的重点资产指纹识别工具, Powered by P001water


`,
	PreRun: func(cmd *cobra.Command, args []string) {

		// load P1finger config from p1fingerConf.yaml file
		filePath := filepath.Join(vars.ExecDir, "p1fingerConf.yaml")
		err := vars.LoadAppConf(filePath, &vars.AppConf)
		if err != nil {
			gologger.Error().Msgf("%v", err)
			return
		}

		err = vars.LoadHttpClient()
		if err != nil {
			gologger.Error().Msgf("%v", err)
			return
		}

	},
}

func Execute() {
	cc.Init(&cc.Config{
		RootCmd:  RootCmd,
		Headings: cc.Red + cc.Bold + cc.Underline,
		Commands: cc.Cyan + cc.Bold,
		Example:  cc.Italic,
		ExecName: cc.Magenta + cc.Bold,
		Flags:    cc.Cyan + cc.Bold,
	})
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
