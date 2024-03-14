package cmd

import (
	"envme/lib/forms"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
)

var rootCmd = &cobra.Command{
	Use: "envme",
}

func Execute(version string) error {
	rootCmd.Version = version
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(createCmd, exposeCmd, listCmd)
	createCmd.AddCommand(createServiceCmd, createDevCmd)
	listCmd.AddCommand(listServicesCmd)

	// Add flags to all commands
	rootCmd.PersistentFlags().BoolP("interactive", "i", false, "Use interactive mode")
	_ = viper.BindPFlag("interactive", rootCmd.Flags().Lookup("interactive"))

	// Add flags to the `envme create` command
	createCmd.PersistentFlags().StringArrayP("env", "e", []string{}, "Add environment variables for service")
	_ = viper.BindPFlag("env", createCmd.Flags().Lookup("env"))
	createCmd.PersistentFlags().String("env-file", "", "Read in a file of environment variables")
	_ = viper.BindPFlag("env-file", createCmd.Flags().Lookup("env-file"))
	createCmd.PersistentFlags().StringArrayP("expose", "p", []string{}, "Expose a service to the internet")
	_ = viper.BindPFlag("expose", createCmd.Flags().Lookup("expose"))

	// Add flags to the `envme list` command
	listCmd.PersistentFlags().Bool("no-interactive", false, "List services without interactive mode")
	_ = viper.BindPFlag("no-interactive", listCmd.Flags().Lookup("no-interactive"))
}

// createCmd handles the `envme create` command
var createCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"generate", "g"},
	Short:   "Create a new service or development environment",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

// listCmd handles the `envme list` command
var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls", "ps"},
	Short:   "List services",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

// exposeCmd handles the `envme expose` command
// TODO: Handle interactive mode when no args are provided and the interactive flag is set to true
var exposeCmd = &cobra.Command{
	Use:     "expose <service-name> <port> <hostname>",
	Aliases: []string{"publish", "p"},
	Short:   "Expose a service to the internet",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 3 && !cmd.Flags().Changed("interactive") {
			return fmt.Errorf("\n  Please specify <service-name>, <port> and <hostname> or using interactive mode\n")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Exposing service %s on port %s with hostname %s\n", args[0], args[1], args[2])
		return nil
	},
}

// createServiceCmd handles the `envme create service` command
var createServiceCmd = &cobra.Command{
	Use:     "service <service-name> <image-name>",
	Aliases: []string{"srv", "s"},
	Short:   "Create a new service",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 && !cmd.Flags().Changed("interactive") {
			return fmt.Errorf("\n  Please specify <service-name> and <image-name> or using interactive mode\n")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var name, image string
		if len(args) < 2 && cmd.Flags().Changed("interactive") {
			_, err := tea.NewProgram(forms.NewServiceForm()).Run()
			if err != nil {
				return err
			}
			name = viper.GetString("container_name")
			image = viper.GetString("image")
		} else {
			name = args[0]
			image = args[1]
		}
		fmt.Printf("Creating service %s with image %s\n", name, image)
		return nil
	},
}

// createDevCmd handles the `envme create development` command
var createDevCmd = &cobra.Command{
	Use:     "development <environment-name> <directory>",
	Aliases: []string{"dev", "d"},
	Short:   "Create a new development environment",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 && !cmd.Flags().Changed("interactive") {
			return fmt.Errorf("\n  Please specify <env-name> and <directory> or using interactive mode\n")
		}
		if !cmd.Flags().Changed("interactive") && (args[1] == "." || strings.HasPrefix(args[1], "./")) {
			dir, err := os.Getwd()
			if err != nil {
				return err
			}
			args[1] = strings.Replace(args[1], ".", dir, 1)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var name, dir string
		if len(args) < 2 && cmd.Flags().Changed("interactive") {
			_, err := tea.NewProgram(forms.NewDevelopmentForm()).Run()
			if err != nil {
				return err
			}
			name = viper.GetString("container_name")
			dir = viper.GetString("directory")
		} else {
			name = args[0]
			dir = args[1]
		}
		fmt.Printf("Creating development environment %s build from %s\n", name, dir)
		return nil
	},
}

// listServicesCmd handles the `envme list services` command
// TODO: Handle interactive mode when no args are provided and the interactive flag is set to true
var listServicesCmd = &cobra.Command{
	Use:     "service",
	Aliases: []string{"srv", "s"},
	Short:   "List services",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Listing services")
		return nil
	},
}