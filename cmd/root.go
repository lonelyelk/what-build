package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/lonelyelk/what-build/github"
	"github.com/lonelyelk/what-build/what"

	"github.com/manifoldco/promptui"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const awsRegionConfigName = "aws_region"
const awsSSMParameterNameConfigName = "aws_ssm_configuration"

func promptAndWriteConfig(promptStr string, key string) {
	validate := func(input string) error {
		if input == "" {
			return errors.New("parameter empty")
		}
		return nil
	}
	prompt := promptui.Prompt{
		Label:    promptStr,
		Validate: validate,
	}
	result, err := prompt.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	viper.Set(key, result)
	if err = viper.WriteConfig(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initConfig() {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println("Can't find HOME dir:", err)
		os.Exit(1)
	}
	configPath := filepath.Join(home, ".what-build.yaml")
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("No configuration found:", err)
		promptAndWriteConfig("Enter AWS region", awsRegionConfigName)
		promptAndWriteConfig("Enter path for SSM configuration", awsSSMParameterNameConfigName)
	}
	if viper.GetString(awsRegionConfigName) == "" {
		promptAndWriteConfig("Enter AWS region", awsRegionConfigName)
	}
	if viper.GetString(awsSSMParameterNameConfigName) == "" {
		promptAndWriteConfig("Enter path for SSM configuration", awsSSMParameterNameConfigName)
	}
	err = github.Auth()
	if err != nil {
		fmt.Println("Can't request github token:", err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "what-build",
	Short: "What Build is a CircleCI API wrapper to look for build by Build Parameters",
	Long: `
	what-build

A tool to search for project builds on Circle CI with given Build Parameters

  Created in free time by Sergey Kruk <sergey.kruk@gmail.com>`,
	Run: func(cmd *cobra.Command, args []string) {
		version, err := cmd.Flags().GetBool("version")
		if err == nil && version {
			what.PrintVersion()
			return
		}
		cmd.Usage()
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.Flags().BoolP("version", "v", false, "output version")
}

// Execute runs everytime to trigger cobra init
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
