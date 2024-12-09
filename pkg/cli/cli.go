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
	rootCmd.AddCommand(dryCommand)
	dryCommand.Flags().BoolVar(&prodFlag, "prod", false, "Use 'prod' settings instead of 'dev'. Default is 'dev'")
	dryCommand.Flags().StringVarP(&fileFlag, "file", "f", "", "Specify the TOML configuration file. Default is quadlets.toml")

	rootCmd.AddCommand(createCommand)
	createCommand.Flags().BoolVar(&prodFlag, "prod", false, "Use 'prod' settings instead of 'dev'. Default is 'dev'")
	createCommand.Flags().BoolVar(&startFlag, "start", false, "Start the pod immediately after creation")
	createCommand.Flags().StringVarP(&fileFlag, "file", "f", "", "Specify the TOML configuration file. Default is quadlets.toml")

	rootCmd.AddCommand(destroyCommand)
	destroyCommand.Flags().BoolVar(&prodFlag, "prod", false, "Use 'prod' settings instead of 'dev'. Default is 'dev'")
	destroyCommand.Flags().StringVarP(&fileFlag, "file", "f", "", "Specify the TOML configuration file. Default is quadlets.toml")

	rootCmd.AddCommand(putCommand)
	putCommand.Flags().BoolVar(&prodFlag, "prod", false, "Use 'prod' settings instead of 'dev'. Default is 'dev'")
	putCommand.Flags().StringVarP(&fileFlag, "file", "f", "", "Specify the TOML configuration file. Default is quadlets.toml")
	putCommand.Flags().StringVarP(&locationFlag, "location", "l", "", "Specify the directory to save the quadlets. Default is ./podcraft")
}

var rootCmd = &cobra.Command{
	Use:   "podcraft",
	Short: "podcraft - Run podman quadlets locally with ease",
	Long: `.--. -.-.

podcraft is a CLI utility designed to enable local execution of podman quadlets, 
addressing the lack of support for local quadlet execution. This tool lets you configure, simulate, and run pods locally.`,
}

var dryCommand = &cobra.Command{
	Use:   "dry",
	Short: "Simulate pod commands",
	Long: `The 'dry' command simulates and generates the podman commands required to run the configured quadlets locally.
This is particularly useful for development workflows where you need to inspect or debug the generated commands 
before actual execution.`,
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

var createCommand = &cobra.Command{
	Use:   "create",
	Short: "Create and optionally start the pod",
	Long: `The 'create' command generates and executes the necessary podman commands
to create a pod. It can optionally start the pod if the --start flag is provided.`,
	Run: func(cmd *cobra.Command, args []string) {
		if fileFlag == "" {
			fileFlag = "./quadlets.toml"
		}

		podname, err := execs.CreatePodman(fileFlag, prodFlag)
		if err != nil {
			log.Fatalf("%v", err)
		}

		if startFlag {
			fmt.Println("\nPod created successfully, attempting to start...")
			err = execs.TryStartPod(*podname)
			if err != nil {
				log.Fatalf("%v", err)
			}
			fmt.Println(color.ColorGreen + "Started!" + color.ColorReset)
		} else {
			fmt.Printf("\nPod created successfully. Start it with 'podman pod start %s'.\n", *podname)
			fmt.Printf("View logs with 'podman pod logs -f -c <container name> %s'.\n", *podname)
		}
	},
}

var destroyCommand = &cobra.Command{
	Use:   "destroy",
	Short: "Remove the pod and clean up resources",
	Long: `The 'destroy' command forcefully removes the specified pod and prunes any associated networks.
Note: This does not remove volumes. Use 'podman volume rm <volume name>' to remove volumes.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if fileFlag == "" {
			fileFlag = "./quadlets.toml"
		}

		err := execs.PodmanRmf(fileFlag, prodFlag)
		if err != nil {
			return err
		}

		fmt.Println(color.ColorGreen + "\nPod destroyed successfully!" + color.ColorReset)
		fmt.Println("\nVolumes were not removed. Use 'podman volume ls' to list and 'podman volume rm <volume name>' to remove them.")

		return nil
	},
}

var putCommand = &cobra.Command{
	Use:   "put",
	Short: "Prepare and save quadlets for production",
	Long: `The 'put' command processes your custom configuration file and generates standard quadlet unit files
based on the defined quadlets. These files are then saved to the specified directory, making them ready
for deployment in a production environment. If no location is provided, the default directory './podcraft' is used.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if fileFlag == "" {
			fileFlag = "./quadlets.toml"
		}

		err := execs.PutQuadlets(fileFlag, prodFlag, locationFlag)
		if err != nil {
			return err
		}

		fmt.Println(color.ColorGreen + "\nQuadlets prepared and saved successfully!" + color.ColorReset)
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
