package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/manifoldco/promptui"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func promptAndSetConfig(promptStr string, key string) error {
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
		return err
	}
	viper.Set(key, result)
	return nil
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
		if err = promptAndSetConfig("Enter AWS region", "aws_region"); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if err = promptAndSetConfig("Enter path for SSM configuration", "aws_ssm_configuration"); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if err = viper.WriteConfigAs(configPath); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	if viper.GetString("aws_region") == "" {
		if err = promptAndSetConfig("Enter AWS region", "aws_region"); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if err = viper.WriteConfigAs(configPath); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	if viper.GetString("aws_ssm_configuration") == "" {
		if err = promptAndSetConfig("Enter path for SSM configuration", "aws_ssm_configuration"); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if err = viper.WriteConfigAs(configPath); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

var rootCmd = &cobra.Command{
	Use:   "what-build",
	Short: "What Build is a CircleCI API wrapper to look for build by Build Parameters",
	Long: `
	what-build

A tool to search for project builds on Circle CI with given Build Parameters

  Created in free time by Sergey Kruk <sergey.kruk@gmail.com>`,
}

func init() {
	cobra.OnInitialize(initConfig)
}

// Execute runs everytime to trigger cobra init
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
