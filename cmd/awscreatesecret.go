package cmd

import (
	logging "Con_Utils/Logging"
	"Con_Utils/apps"
	"fmt"
	"github.com/spf13/cobra"
)

// Command argument parser for aws secret creation template

var TagawsCreateSecretCmd = &cobra.Command{
	Use:   "awscrtsecret",
	Short: "This is to create a secret in aws secret manager",
	Long:  "\nGiven input a string, it will create all the necessary aws resources needed to store it in the aws secret manager",
	//Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		logging.Application = "awscrtsecret"
		fmt.Println("the secret key is " + apps.TagSecretName + "\n\n")
		tagExtraEnvVar := map[string]string{"TAG_SECRET_NAME": apps.TagSecretName, "TAG_SECRET": apps.TagSecret}
		apps.TagawsEnvRun(tagExtraEnvVar)
	},
}

func init() {
	RootCmd.AddCommand(TagawsCreateSecretCmd)
	TagawsCreateSecretCmd.Flags().StringVarP(&apps.Region, "region", "r", "", "region name (required)")
	TagawsCreateSecretCmd.MarkFlagRequired("region")
	TagawsCreateSecretCmd.Flags().StringVarP(&apps.AwsProfile, "aws_profile", "p", "", "aws profile name (required)")
	TagawsCreateSecretCmd.MarkFlagRequired("secret")
	TagawsCreateSecretCmd.Flags().StringVarP(&apps.TagProject, "tag_project", "t", "", "tag project name (required)")
	TagawsCreateSecretCmd.MarkFlagRequired("tag_project")
	TagawsCreateSecretCmd.Flags().StringVarP(&apps.TagSecretName, "tag_secret_name", "m", "", "tag secret name (required)")
	TagawsCreateSecretCmd.MarkFlagRequired("tag_secret_name")
	TagawsCreateSecretCmd.Flags().StringVarP(&apps.TagSecret, "tag_secret", "s", "", "tag secret (required)")
	TagawsCreateSecretCmd.MarkFlagRequired("tag_secret")
}
