/*
Copyright 2018 Spotify AB. All rights reserved.

The contents of this file are licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spotify/protoman/cli/protoman"
	"github.com/spotify/protoman/cli/validator"
)

var rootCmd = &cobra.Command{
	Use:   "protoman",
	Short: "Protoman is a fantastic way of managing your protos",
}

func init() {
	rootCmd.AddCommand(versionCmd, validateCmd, publishCmd, genCmd, getCmd)
	publishCmd.PersistentFlags().StringP("server", "s", "", "Protoman server address")
	rootCmd.PersistentFlags().StringP("proto-dir", "p", "", "Root directory where protos will be stored")
	getCmd.PersistentFlags().StringP("server", "s", "", "Protoman server address")
}

func exitOnErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of protoman",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Protoman version 0.0.1")
	},
}

var validateCmd = &cobra.Command{
	Use:   "validate [path]",
	Short: "Validate proto defintion file(s)",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Path required")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(args[0]); err == nil {
			exitOnErr(validator.Validate(args[0]))
		}
	},
}

var publishCmd = &cobra.Command{
	Use:   "publish [protos]",
	Short: "Publish proto defintion file(s) to a protoman registry.",
	Long: `
	Publish proto defintion file(s) to a protoman registry.
	Providing no arguments will upload local packages defined in .protoman`,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := verifyFlags(cmd); err != nil {
			return err
		}
		if len(args) == 0 {
			return nil
		}
		for _, proto := range args {
			if !strings.HasSuffix(proto, ".proto") {
				return errors.New("must be .proto files")
			}
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		client, err := protoman.NewRegistryClient(cmd.Flag("server").Value.String())
		if err != nil {
			exitOnErr(err)
		}
		exitOnErr(protoman.Publish(
			args, cmd.Flag("proto-dir").Value.String(),
			client))
	},
}

var genCmd = &cobra.Command{
	Use:   "generate [package name] [service name] [root path]",
	Short: "generate protocol stanza",
	Long: `
	Generate protocol stanza
	Example:
		package name: spotify.protoman.v1
		service name: registry
	The following input example will create
	spotify/protoman/v1/registry.proto in your current directory alongside
	a .protoman file that tracks your dependencies.
	`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 3 {
			return errors.New("Missing parameters")
		}
		if strings.HasPrefix(args[2], "/") {
			return errors.New("Root path must be relative to project")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		exitOnErr(protoman.Generate(args[0], args[1], args[2]))
	},
}

func verifyFlags(cmd *cobra.Command) error {
	if cmd.Flag("server").Value.String() == "" {
		return errors.New("--server must be set to the protoman registry")
	}
	if cmd.Flag("proto-dir").Value.String() == "" {
		return errors.New("--proto-dir must be specified")
	}
	if strings.HasPrefix(cmd.Flag("proto-dir").Value.String(), "/") {
		return errors.New("Proto dir must be relative to project")
	}
	return nil
}

var getCmd = &cobra.Command{
	Use:   "get [package names]",
	Short: "Get package",
	Args: func(cmd *cobra.Command, args []string) error {
		return verifyFlags(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		client, err := protoman.NewRegistryClient(cmd.Flag("server").Value.String())
		if err != nil {
			exitOnErr(err)
		}
		exitOnErr(protoman.Get(args, cmd.Flag("proto-dir").Value.String(), client))
	},
}

// Execute CLI
func Execute() {
	exitOnErr(rootCmd.Execute())
}
