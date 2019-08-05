package commands

import (
	"github.com/abassian/huron/src/huron"
	"github.com/abassian/huron/src/proxy/dummy"
	aproxy "github.com/abassian/huron/src/proxy/socket/app"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//NewRunCmd returns the command that starts a Huron node
func NewRunCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "run",
		Short:   "Run node",
		PreRunE: loadConfig,
		RunE:    runHuron,
	}
	AddRunFlags(cmd)
	return cmd
}

/*******************************************************************************
* RUN
*******************************************************************************/

func runHuron(cmd *cobra.Command, args []string) error {
	if !config.Standalone {
		p, err := aproxy.NewSocketAppProxy(
			config.ClientAddr,
			config.ProxyAddr,
			config.Huron.NodeConfig.HeartbeatTimeout,
			config.Huron.Logger,
		)

		if err != nil {
			config.Huron.Logger.Error("Cannot initialize socket AppProxy:", err)
			return err
		}

		config.Huron.Proxy = p
	} else {
		p := dummy.NewInmemDummyClient(config.Huron.Logger)

		config.Huron.Proxy = p
	}

	engine := huron.NewHuron(&config.Huron)

	if err := engine.Init(); err != nil {
		config.Huron.Logger.Error("Cannot initialize engine:", err)
		return err
	}

	engine.Run()

	return nil
}

/*******************************************************************************
* CONFIG
*******************************************************************************/

//AddRunFlags adds flags to the Run command
func AddRunFlags(cmd *cobra.Command) {

	cmd.Flags().String("datadir", config.Huron.DataDir, "Top-level directory for configuration and data")
	cmd.Flags().String("log", config.Huron.LogLevel, "debug, info, warn, error, fatal, panic")
	cmd.Flags().String("moniker", config.Huron.Moniker, "Optional name")

	// Network
	cmd.Flags().StringP("listen", "l", config.Huron.BindAddr, "Listen IP:Port for huron node")
	cmd.Flags().DurationP("timeout", "t", config.Huron.NodeConfig.TCPTimeout, "TCP Timeout")
	cmd.Flags().DurationP("join-timeout", "j", config.Huron.NodeConfig.JoinTimeout, "Join Timeout")
	cmd.Flags().Int("max-pool", config.Huron.MaxPool, "Connection pool size max")

	// Proxy
	cmd.Flags().Bool("standalone", config.Standalone, "Do not create a proxy")
	cmd.Flags().StringP("proxy-listen", "p", config.ProxyAddr, "Listen IP:Port for huron proxy")
	cmd.Flags().StringP("client-connect", "c", config.ClientAddr, "IP:Port to connect to client")

	// Service
	cmd.Flags().StringP("service-listen", "s", config.Huron.ServiceAddr, "Listen IP:Port for HTTP service")

	// Store
	cmd.Flags().Bool("store", config.Huron.Store, "Use badgerDB instead of in-mem DB")
	cmd.Flags().Bool("bootstrap", config.Huron.NodeConfig.Bootstrap, "Load from database")
	cmd.Flags().Int("cache-size", config.Huron.NodeConfig.CacheSize, "Number of items in LRU caches")

	// Node configuration
	cmd.Flags().Duration("heartbeat", config.Huron.NodeConfig.HeartbeatTimeout, "Time between gossips")
	cmd.Flags().Int("sync-limit", config.Huron.NodeConfig.SyncLimit, "Max number of events for sync")
	cmd.Flags().Bool("fast-sync", config.Huron.NodeConfig.EnableFastSync, "Enable FastSync")
}

func loadConfig(cmd *cobra.Command, args []string) error {

	err := bindFlagsLoadViper(cmd)
	if err != nil {
		return err
	}

	config, err = parseConfig()
	if err != nil {
		return err
	}

	config.Huron.Logger.Level = huron.LogLevel(config.Huron.LogLevel)
	config.Huron.NodeConfig.Logger = config.Huron.Logger

	config.Huron.Logger.WithFields(logrus.Fields{
		"huron.DataDir":               config.Huron.DataDir,
		"huron.BindAddr":              config.Huron.BindAddr,
		"huron.ServiceAddr":           config.Huron.ServiceAddr,
		"huron.MaxPool":               config.Huron.MaxPool,
		"huron.Store":                 config.Huron.Store,
		"huron.LoadPeers":             config.Huron.LoadPeers,
		"huron.LogLevel":              config.Huron.LogLevel,
		"huron.Moniker":               config.Huron.Moniker,
		"huron.Node.HeartbeatTimeout": config.Huron.NodeConfig.HeartbeatTimeout,
		"huron.Node.TCPTimeout":       config.Huron.NodeConfig.TCPTimeout,
		"huron.Node.JoinTimeout":      config.Huron.NodeConfig.JoinTimeout,
		"huron.Node.CacheSize":        config.Huron.NodeConfig.CacheSize,
		"huron.Node.SyncLimit":        config.Huron.NodeConfig.SyncLimit,
		"huron.Node.EnableFastSync":   config.Huron.NodeConfig.EnableFastSync,
		"ProxyAddr":                    config.ProxyAddr,
		"ClientAddr":                   config.ClientAddr,
		"Standalone":                   config.Standalone,
	}).Debug("RUN")

	return nil
}

//Bind all flags and read the config into viper
func bindFlagsLoadViper(cmd *cobra.Command) error {
	// cmd.Flags() includes flags from this command and all persistent flags from the parent
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}

	viper.SetConfigName("huron")              // name of config file (without extension)
	viper.AddConfigPath(config.Huron.DataDir) // search root directory
	// viper.AddConfigPath(filepath.Join(config.Huron.DataDir, "huron")) // search root directory /config

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		config.Huron.Logger.Debugf("Using config file: %s", viper.ConfigFileUsed())
	} else if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		config.Huron.Logger.Debugf("No config file found in: %s", config.Huron.DataDir)
	} else {
		return err
	}

	return nil
}

//Retrieve the default environment configuration.
func parseConfig() (*CLIConfig, error) {
	conf := NewDefaultCLIConfig()
	err := viper.Unmarshal(conf)
	if err != nil {
		return nil, err
	}
	return conf, err
}

func logLevel(l string) logrus.Level {
	switch l {
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "warn":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	case "fatal":
		return logrus.FatalLevel
	case "panic":
		return logrus.PanicLevel
	default:
		return logrus.DebugLevel
	}
}
