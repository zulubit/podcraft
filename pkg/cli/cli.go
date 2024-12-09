package cli

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/zulubit/podcraft/pkg/color"
	"github.com/zulubit/podcraft/pkg/execs"
)

var (
	fileFlag     string
	prodFlag     bool
	startFlag    bool
	locationFlag string
)

func init() {
	rootCmd.AddCommand(dryCoomand)
	dryCoomand.Flags().BoolVar(&prodFlag, "prod", false, "Takes 'prod' replaceables insted of 'dev'")
	dryCoomand.Flags().StringVarP(&fileFlag, "file", "f", "", "Set what TOML file to read. Default is quadlets.toml")

	rootCmd.AddCommand(createCoomand)
	createCoomand.Flags().BoolVar(&prodFlag, "prod", false, "Takes 'prod' replaceables insted of 'dev'")
	createCoomand.Flags().BoolVar(&startFlag, "start", false, "Start also tries to start the pod")
	createCoomand.Flags().StringVarP(&fileFlag, "file", "f", "", "Set what TOML file to read. Default is quadlets.toml")

	rootCmd.AddCommand(destroyCommand)
	destroyCommand.Flags().BoolVar(&prodFlag, "prod", false, "Takes 'prod' replaceables insted of 'dev'")
	destroyCommand.Flags().StringVarP(&fileFlag, "file", "f", "", "Set what TOML file to read. Default is quadlets.toml")

	rootCmd.AddCommand(putCommand)
	putCommand.Flags().BoolVar(&prodFlag, "prod", false, "Takes 'prod' replaceables insted of 'dev'")
	putCommand.Flags().StringVarP(&fileFlag, "file", "f", "", "Set what TOML file to read. Default is quadlets.toml")
	putCommand.Flags().StringVarP(&locationFlag, "location", "l", "", "location specifies the directory to put the quadlets into. Default is ./podcraft")
}

var rootCmd = &cobra.Command{
	Use:   "podcraft",
	Short: "Somewhat like podman compose but for quadlets",
}

var dryCoomand = &cobra.Command{
	Use:   "dry",
	Short: "Print the comands to run the pod locally to the terminal",
	RunE: func(cmd *cobra.Command, args []string) error {
		if fileFlag == "" {
			fileFlag = "./quadlets.toml"
		}

		err := execs.PrintPodman(fileFlag, prodFlag)
		if err != nil {
			return err
		}

		return nil
	},
}

var createCoomand = &cobra.Command{
	Use:   "create",
	Short: "Create command generates commands and tries to run them",
	Run: func(cmd *cobra.Command, args []string) {
		if fileFlag == "" {
			fileFlag = "./quadlets.toml"
		}

		podname, err := execs.CreatePodman(fileFlag, prodFlag)
		if err != nil {
			log.Fatalf("%v", err)
		}

		if startFlag {
			fmt.Println("\nPod created successfully, now trying to start...\n ")
			err = execs.TryStartPod(*podname)
			if err != nil {
				log.Fatalf("%v", err)
			}
			fmt.Println(color.ColorGreen + "Started!" + color.ColorReset)
		} else {

			fmt.Printf("\nCreated successfully, you can start your pod with 'podman pod start %s'\nYou can view your logs by running 'podman pod logs -f -c <container name> %s'\n", *podname, *podname)
		}

	},
}

var destroyCommand = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy forcefully removes the pod and prunes the networks",
	RunE: func(cmd *cobra.Command, args []string) error {

		if fileFlag == "" {
			fileFlag = "./quadlets.toml"
		}

		err := execs.PodmanRmf(fileFlag, prodFlag)
		if err != nil {
			return err
		}
		fmt.Println(color.ColorGreen + "\nSuccess!" + color.ColorReset)

		fmt.Println("\nThis action DID NOT remove your volumes, to remove them find one by running 'podman volume ls' and run 'podman volume rm <volume name>'")

		return nil
	},
}

var putCommand = &cobra.Command{
	Use:   "put [location]",
	Short: "Put puts the quadlets in the desired location",
	RunE: func(cmd *cobra.Command, args []string) error {
		if fileFlag == "" {
			fileFlag = "./quadlets.toml"
		}

		err := execs.PutQuadlets(fileFlag, prodFlag, locationFlag)
		if err != nil {
			return err
		}
		fmt.Println(color.ColorGreen + "\nSuccess!" + color.ColorReset)

		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
