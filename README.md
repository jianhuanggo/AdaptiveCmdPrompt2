# AdaptiveCmdPrompt 

AdaptiveCmdPrompt is a gluing tool and a companion to terraform, not only it can be an effective wrapper, replacing the function of shell wrapper, because it is written in a modern program language and agnostic to cloud vendors, we can extend it to serve any function of future, measured or far away.  Such flexibility is welcomed since we don't need to fit all infrastructure as code exclusively in Terraform, SDKs, CDKs or CLI in one implementation, rather, we can pick and choose the best of all worlds and then use AdaptiveCmdPrompt to integrate into one piece.  

In addition, it currently provides prebuild tagging to facilitate some activities. 

<TAG_IF> conditional logic to only execute subsequent statement if previous statement returns True  
<TAG_IFNOT> conditional logic to only execute subsequent statement if previous statement returns False
<TAG_RETRY> retry the statement, it doesn't have to be the original statement
<TAG_EXPORT> Inject additional variable to environment variable bank
<TAG_WORKDIR> Setting working directory
<TAG_RSTR> Random string generate, can be used to generate password on the fly
<TAG_COMMAND> regular command line statement
  
  
  
Example 1:
  
  - "<TAG_IFNOT>terraform -version<TAG_TONFI><TAG_CMD>git clone https://github.com/tfutils/tfenv.git /home/ec2-user/.tfenv && sudo ln -s /home/ec2-user/.tfenv/bin/* /usr/local/bin && /usr/local/bin/tfenv install latest && /usr/local/bin/tfenv install 0.11.6 && tfenv use latest<TAG_DMC>"
  
if condition statement "terraform -version" resolves to False which means terraform is not installed, then will run subsequent commands to install tfenv and both 0.11.6 and latest versions of terraform
  
  











