package cmd

import (
	logging "Con_Utils/Logging"
	apps "Con_Utils/apps"
	"github.com/spf13/cobra"
)

// Command argument parser for aws secret creation template

var TagawsenvsetupCmd = &cobra.Command{
	Use:   "awsenv",
	Short: "This is to setup Tag AWS envrionment",
	Long:  "\nThis is to create all the aws resources needed to setup tag environment. See details at...",
	//Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		logging.Application = "awsenv"
		//fmt.Println(Region)
		apps.TagawsEnvRun(nil)
	},
}

func init() {
	RootCmd.AddCommand(TagawsenvsetupCmd)
	TagawsenvsetupCmd.Flags().StringVarP(&apps.AwsProfile, "aws_profile", "p", "", "aws profile name (required)")
	TagawsenvsetupCmd.MarkFlagRequired("aws_profile")
	TagawsenvsetupCmd.Flags().StringVarP(&apps.TagProject, "tag_project", "t", "", "tag project name (required)")
	TagawsenvsetupCmd.MarkFlagRequired("tag_project")
	TagawsenvsetupCmd.Flags().BoolVarP(&apps.TagNewStartFlag, "new_start", "n", false, "new start flag (option)")
}
