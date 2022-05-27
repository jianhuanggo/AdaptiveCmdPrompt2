package cmd

import (
	apps "Con_Utils/apps"
	"fmt"
	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var RootCmd = &cobra.Command{
	Use:   "Con_Utils",
	Short: "Tag Utility",
	Long:  "\n Tag Utility Wrapper for day to day tasks",
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initConfig(filepath string) {

	apps.CfgFile = filepath

	if apps.CfgFile != "" {

		log.Printf("CfgFile is... ")
		viper.SetConfigFile(apps.CfgFile)

	} else {
		home, err := homedir.Dir()
		fmt.Printf("Home is #{home}\n")

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err == nil {
		log.Println("using config file: ", viper.ConfigFileUsed())

	}

	err = viper.Unmarshal(&apps.Config)
	if err != nil {
		log.Printf("unable to decode into struct, #{error}")
	}
}
