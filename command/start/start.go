package start

import (
	_ "github.com/Clay294/forum/apps"
	"github.com/Clay294/forum/config"
	"github.com/Clay294/forum/flog"
	"github.com/Clay294/forum/ioc"
	"github.com/Clay294/forum/protocol"
	"github.com/Clay294/forum/rsakeys"
	_ "github.com/Clay294/forum/rsakeys/impl"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var (
	ConfigFile           string
	LogFile              string
	RSAKeyPairsUUIDsFile string
)

func runStartCmd(cmd *cobra.Command, args []string) {
	err := flog.Init(LogFile)
	if err != nil {
		cobra.CheckErr(err)
	}

	err = config.LoadConfig(ConfigFile)
	if err != nil {
		cobra.CheckErr(err)
	}

	err = rsakeys.InitGlobalRSAKeyPairsUUIds(RSAKeyPairsUUIDsFile)
	if err != nil {
		cobra.CheckErr(err)
	}

	err = ioc.Controllers().Init()
	if err != nil {
		cobra.CheckErr(err)
	}

	engine := gin.Default()

	err = ioc.Handlers().Init(engine)
	if err != nil {
		cobra.CheckErr(err)
	}

	hs := protocol.NewHttpServer(engine)

	err = hs.Start()
	if err != nil {
		cobra.CheckErr(err)
	}
}

var StartCmd = cobra.Command{
	Use:   "start",
	Short: "start forum project api",
	Long:  "start forum project api v2",
	Run:   runStartCmd,
}

func init() {
	StartCmd.Flags().StringVarP(&ConfigFile, "--config-file", "C", config.DefaultConfigFile, "specify the configuration file")
	StartCmd.Flags().StringVarP(&LogFile, "--log-file", "L", flog.DefalutLogFile, "specify the log file")
	StartCmd.Flags().StringVarP(&RSAKeyPairsUUIDsFile, "--rsakeypairs-uuids-file", "U", rsakeys.DefaultRSAKeyPairsUUIDsFile, "specify the rsa key pairs uuids file")
}
