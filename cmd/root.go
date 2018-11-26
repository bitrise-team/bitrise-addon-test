package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/bitrise-team/bitrise-add-on-testing-kit/addonprovisioner"
	"github.com/bitrise-team/bitrise-add-on-testing-kit/addontester"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

const (
	colorRed        = "\x1b[31;1m"
	colorNeutral    = "\x1b[0m"
	numberOfRetries = 2
)

func fail(err error) {
	fmt.Printf("\n%s%s%s\n", colorRed, err, colorNeutral)
	os.Exit(1)
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bitrise-add-on-testing-kit",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file to use (default is ./config.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name "config" (without extension).
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Reading config file:", viper.ConfigFileUsed())
	} else {
		fmt.Printf("Failed to read config file: %s", err.Error())
		os.Exit(1)
	}

	validateConfig()
}

func validateConfig() {
	requiredConfigs := []string{"addon-url", "auth-token", "sso-secret"}

	fmt.Println("\nConfigs:")
	for _, config := range requiredConfigs {
		if !viper.IsSet(config) {
			fmt.Printf("Config %s is required but not set\n", config)
			os.Exit(1)
		}
		fmt.Printf("%s: %s\n", config, viper.Get(config))
	}
}

func addonTesterFromConfig() (*addontester.Tester, error) {
	addonClient, err := addonprovisioner.NewClient(&addonprovisioner.ClientConfig{
		AddonURL:  viper.Get("addon-url").(string),
		AuthToken: viper.Get("auth-token").(string),
		SSOSecret: viper.Get("sso-secret").(string),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return addontester.New(addonClient, log.New(os.Stdout, "", 0))
}
