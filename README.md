# AdaptiveCmdPrompt 

AdaptiveCmdPrompt is a gluing tool and a companion to terraform, not only it can be an effective wrapper, replacing the function of shell wrapper, because it is written in a modern program language and agnostic to cloud vendors, we can extend it to serve any function of future, measured or far away.  Such flexibility is welcomed since we don't need to fit all infrastructure as code exclusively in Terraform, SDKs, CDKs or CLI in one implementation, rather, we can pick and choose the best of all worlds and then use AdaptiveCmdPrompt to integrate into one piece.  

In addition, it currently provides prebuild tagging to facilitate some activities. 

- <TAG_IF> conditional logic to only execute subsequent statement if previous statement returns True  
- <TAG_IFNOT> conditional logic to only execute subsequent statement if previous statement returns False
- <TAG_RETRY> retry the statement, it doesn't have to be the original statement

- <TAG_EXP> Inject additional variable to environment variable bank
- <TAG_WORKDIR> Setting working directory
- <TAG_RSTR> Random string generate, can be used to generate password on the fly
- <TAG_COMMAND> regular command line statement
  
  
  
Example 1:
  
  - "<TAG_IFNOT>terraform -version<TAG_TONFI><TAG_CMD>git clone https://github.com/tfutils/tfenv.git /home/ec2-user/.tfenv && sudo ln -s /home/ec2-user/.tfenv/bin/* /usr/local/bin && /usr/local/bin/tfenv install latest && /usr/local/bin/tfenv install 0.11.6 && tfenv use latest<TAG_DMC>"
  
if condition statement "terraform -version" resolves to False which means terraform is not installed, then will run subsequent commands to install tfenv and both 0.11.6 and latest versions of terraform


Example 2:

 - "<TAG_EXP>TAG_PASS_KMS_ID<TAG_PXE><TAG_CMD>aws kms create-key --profile ${TAG_AWS_CLI_PROFILE} --tags TagKey=environment,TagValue=${TAG_ENVIRONMENT} TagKey=tag_app,TagValue=${TAG_APP_NAME} TagKey=usage,TagValue=ec2-passwd --description \"kms key for encrypt ec2 bastion key\"  | jq -r '.KeyMetadata.\"KeyId\"'<TAG_DMC>"

The command between <TAG_CMD> and <TAG_DMC> creates KMS key using AWS CLI and its corresponding key id is stored as a variable TAG_PASS_KMS_ID which can be used in future statements

Example 3:

 - "<TAG_CMD>aws kms create-alias --profile $${TAG_AWS_CLI_PROFILE} --alias-name alias/kms-4_ec2_app_${TAG_APP_NAME}_${TAG_ENVIRONMENT} --target-key-id ${TAG_PASS_KMS_ID}<TAG_DMC><TAG_RT>aws kms delete-alias --profile ${TAG_AWS_CLI_PROFILE} --alias-name alias/kms-4_ec2_app_${TAG_APP_NAME}_${TAG_ENVIRONMENT}<TAG_TR>"

The AWS CLI command between <TAG_CMD> and <TAG_DMC> is trying to create an alias for KMS key and if that is not successful, another AWS CLI command between <TAG_RT> and <TAG_TR> gets invoked to delete any existing alias with the same name and then AdaptiveCmdPrompt will retry the original statement and fails only if retry also fails.


Example 4:

"<TAG_EXP>TAG_RDS_PASSWD<TAG_PXE><TAG_RSTR>tag;12<TAG_RTSR>"

with <TAG_RSTR>tag;12<TAG_RTSR>, it tells AdaptiveCmdPrompt to generate a random string 12 characters long with "tag" as prefix and then store it into the variable TAG_RDS_PASSWD which later can be used in the database creation command as a password.
















