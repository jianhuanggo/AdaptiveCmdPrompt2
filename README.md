# AdaptiveCmdPrompt 

AdaptiveCmdPrompt is a gluing tool and a companion to terraform, not only it can be an effective wrapper, replacing the function of shell wrapper, because it is written in a modern program language and agnostic to cloud vendors, we can extend it to serve any function of future, measured or far away.  Such flexibility is welcomed since we don't need to fit all infrastructure as code exclusively in Terraform, SDKs, CDKs or CLI in one implementation, rather, we can pick and choose the best of all worlds and then use AdaptiveCmdPrompt to integrate into one piece.  

In addition, it currently provides advance tagging following 

<IF> conditional logic to only execute subsequent statement if previous statement returns True
  
  
  
<IFNOT> conditional logic to only execute subsequent statement if previous statement returns False
<RETRY> retry the statement, it doesn't have to be the original statement
<EXPORT> Inject additional variable to environment variable bank
<WORKDIR> Setting working directory
<RSTR> Random string generate, can be used to generate password on the fly
<COMMAND> regular command line statement








