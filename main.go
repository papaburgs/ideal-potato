package main

import (
	"fmt"
	"os/exec"

	"github.com/Sirupsen/logrus"
	"github.com/ogier/pflag"
	"github.com/spf13/viper"
)

func main() {
	var (
		flagWriteDefault *bool
		flagPersist      *bool
	)
	flagWriteDefault = pflag.Bool("writedefault", false, "write defaults to file")
	flagPersist = pflag.Bool("persist", false, "do not wipe run directory before restarting")
	pflag.Parse()
	if !*flagPersist {
		if err := clearRunDir(); err != nil {
			return
		}
	}
	if *flagWriteDefault {

		runConfigPath := viper.GetString("appPath") + viper.GetString("runPath") + "config.yaml"
		fmt.Println("writing default config to run directory, " + runConfigPath + "\nThis will not survive a restart")
		viper.SetConfigFile(runConfigPath)

		if err := viper.WriteConfig(); err != nil {
			fmt.Println("could not write config file:\n\t", err)
		}
		fmt.Println("done with config file. Should maybe put some comments in the config file")
	}
	fmt.Println("starting")
	fmt.Println("got this:\n", viper.AllSettings())

}

func init() {
	viper.AddConfigPath("/opt/processwatcher")
	viper.SetConfigName("config")
	// viper.SetEnvPrefix("PSWATCHER")
	// viper.AutomaticEnv() // disable this for now

	viper.SetDefault("appPath", "/opt/processwatcher/")
	viper.SetDefault("runPath", "run/")
	viper.SetDefault("databasePath", "db/")

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Println("no config file found, using defaults where possible")
	}
}

func clearRunDir() error {
	rmCmd := exec.Command("rm", "-rf", viper.GetString("appPath")+viper.GetString("runPath"))
	if err := rmCmd.Run(); err != nil {
		fmt.Println("error removing previous run files: \n\t", err)
		return err
	}
	mkCmd := exec.Command("mkdir", viper.GetString("appPath")+viper.GetString("runPath"))
	if o, err := mkCmd.CombinedOutput(); err != nil {
		logrus.Error(string(o))
		return err
	}
	return nil
}
