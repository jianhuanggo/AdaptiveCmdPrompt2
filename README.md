<!--
SPDX-FileCopyrightText: 2019–2023 Pynguin Contributors

SPDX-License-Identifier: CC-BY-4.0
-->

# Pynguin

Pynguin (IPA: ˈpɪŋɡuiːn),
the
PYthoN
General
UnIt
test
geNerator,
is a tool that allows developers to generate unit tests automatically.

Testing software is often considered to be a tedious task.
Thus, automated generation techniques have been proposed and mature tools exist—for
statically typed languages, such as Java.
There is, however, no fully-automated tool available that produces unit tests for
general-purpose programs in a dynamically typed language.
Pynguin is, to the best of our knowledge, the first tool that fills this gap
and allows the automated generation of unit tests for Python programs.

<details>
<summary>Internal Pipeline Status</summary>

[![pipeline status](https://gitlab.infosun.fim.uni-passau.de/se2/pynguin/pynguin/badges/main/pipeline.svg)](https://gitlab.infosun.fim.uni-passau.de/se2/pynguin/pynguin/-/commits/main)
[![coverage report](https://gitlab.infosun.fim.uni-passau.de/se2/pynguin/pynguin/badges/main/coverage.svg)](https://gitlab.infosun.fim.uni-passau.de/se2/pynguin/pynguin/-/commits/main)  

</details>

[![License MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Code style: black](https://img.shields.io/badge/code%20style-black-000000.svg)](https://github.com/ambv/black)
[![PyPI version](https://badge.fury.io/py/pynguin.svg)](https://badge.fury.io/py/pynguin)
[![Supported Python Versions](https://img.shields.io/pypi/pyversions/pynguin.svg)](https://github.com/se2p/pynguin)
[![Documentation Status](https://readthedocs.org/projects/pynguin/badge/?version=latest)](https://pynguin.readthedocs.io/en/latest/?badge=latest)
[![DOI](https://zenodo.org/badge/DOI/10.5281/zenodo.3989840.svg)](https://doi.org/10.5281/zenodo.3989840)
[![REUSE status](https://api.reuse.software/badge/github.com/se2p/pynguin)](https://api.reuse.software/info/github.com/se2p/pynguin)
[![Downloads](https://static.pepy.tech/personalized-badge/pynguin?period=total&units=international_system&left_color=grey&right_color=blue&left_text=Downloads)](https://pepy.tech/project/pynguin)


![Pynguin Logo](https://raw.githubusercontent.com/se2p/pynguin/master/docs/source/_static/pynguin-logo.png "Pynguin Logo")

## Attention

*Please Note:*

**Pynguin executes the module under test!**
As a consequence, depending on what code is in that module,
running Pynguin can cause serious harm to your computer,
for example, wipe your entire hard disk!
We recommend running Pynguin in an isolated environment;
use, for example, a Docker container to minimize the risk of damaging
your system.

**Pynguin is only a research prototype!**
It is not tailored towards production use whatsoever.
However, we would love to see Pynguin in a production-ready stage at some point;
please report your experiences in using Pynguin to us.


## Prerequisites

Before you begin, ensure you have met the following requirements:
- You have installed Python 3.10 (we have not yet tested with Python
  3.11, there might be some problems due to changed internals regarding the byte-code
  instrumentation).

  **Attention:** Pynguin now requires Python 3.10!  Older versions are no longer 
  supported!
- You have a recent Linux/macOS/Windows machine.

Please consider reading the [online documentation](https://pynguin.readthedocs.io)
to start your Pynguin adventure.
 
## Installing Pynguin

Pynguin can be easily installed using the `pip` tool by typing:
```bash
pip install pynguin
```

Make sure that your version of `pip` is that of a supported Python version, as any 
older version is not supported by Pynguin!

## Using Pynguin

Before you continue, please read the [quick start guide](https://pynguin.readthedocs.io/en/latest/user/quickstart.html)

Pynguin is a command-line application.
Once you installed it to a virtual environment, you can invoke the tool by typing
`pynguin` inside this virtual environment.
Pynguin will then print a list of its command-line parameters.

A minimal full command line to invoke Pynguin could be the following,
where we assume that a project `foo` is located in `/tmp/foo`,
we want to store Pynguin's generated tests in `/tmp/testgen`,
and we want to generate tests using a whole-suite approach for the module `foo.bar`
(wrapped for better readability):
```bash
pynguin \
  --project-path /tmp/foo \
  --output-path /tmp/testgen \
  --module-name foo.bar
```


# AdaptiveCmdPrompt 

AdaptiveCmdPrompt is a gluing tool and a companion to terraform, not only it can be an effective wrapper, replacing the function of shell wrapper, because it is written in a modern program language and agnostic to cloud vendors, we can extend it to serve any function of future, measured or far away.  Such flexibility is welcomed since we don't need to fit all infrastructure as code exclusively in Terraform, SDKs, CDKs or CLI in one implementation, rather, we can pick and choose the best of all worlds (including taking the advantages of certain language version of SDKs are better structured than others) and then use AdaptiveCmdPrompt to integrate into one piece.  In addition, AdaptiveCmdPrompt can be extended to call any API free from terraform constrains and can even modify terraform scripts on the fly if needed. 

In addition, it currently provides prebuild tagging to facilitate some activities. 

- <TAG_IF> conditional logic to only execute subsequent statement if previous statement returns True  
- <TAG_IFNOT> conditional logic to only execute subsequent statement if previous statement returns False
- <TAG_RETRY> retry the statement, it doesn't have to be the original statement
- <TAG_EXP> Inject additional variable to environment variable bank
- <TAG_WORKDIR> Setting working directory
- <TAG_RSTR> Random string generate, can be used to generate password on the fly
- <TAG_COMMAND> regular command line statement
  
  
  
Example 1:
```
  - "<TAG_IFNOT>terraform -version<TAG_TONFI><TAG_CMD>git clone https://github.com/tfutils/tfenv.git /home/ec2-user/.tfenv && sudo ln -s /home/ec2-user/.tfenv/bin/* /usr/local/bin && /usr/local/bin/tfenv install latest && /usr/local/bin/tfenv install 0.11.6 && tfenv use latest<TAG_DMC>"

```
If condition statement "terraform -version" resolves to False which means terraform is not installed, then will run subsequent commands to install tfenv and both 0.11.6 and latest versions of terraform


Example 2:

```
 - "<TAG_EXP>TAG_PASS_KMS_ID<TAG_PXE><TAG_CMD>aws kms create-key --profile ${TAG_AWS_CLI_PROFILE} --tags TagKey=environment,TagValue=${TAG_ENVIRONMENT} TagKey=tag_app,TagValue=${TAG_APP_NAME} TagKey=usage,TagValue=ec2-passwd --description \"kms key for encrypt ec2 bastion key\"  | jq -r '.KeyMetadata.\"KeyId\"'<TAG_DMC>"
```

The command between <TAG_CMD> and <TAG_DMC> creates KMS key using AWS CLI and its corresponding key id is stored as a variable TAG_PASS_KMS_ID which can be used in future statements

Example 3:

```
 - "<TAG_CMD>aws kms create-alias --profile ${TAG_AWS_CLI_PROFILE} --alias-name alias/kms-4_ec2_app_${TAG_APP_NAME}_${TAG_ENVIRONMENT} --target-key-id ${TAG_PASS_KMS_ID}<TAG_DMC><TAG_RT>aws kms delete-alias --profile ${TAG_AWS_CLI_PROFILE} --alias-name alias/kms-4_ec2_app_${TAG_APP_NAME}_${TAG_ENVIRONMENT}<TAG_TR>"
```

The AWS CLI command between <TAG_CMD> and <TAG_DMC> is trying to create an alias for KMS key and if that is not successful, another AWS CLI command between <TAG_RT> and <TAG_TR> gets invoked to delete any existing alias with the same name and then AdaptiveCmdPrompt will retry the original statement and fails only if retry also fails.


Example 4:
```
 - "<TAG_EXP>TAG_RDS_PASSWD<TAG_PXE><TAG_RSTR>tag;12<TAG_RTSR>"
```
with <TAG_RSTR>tag;12<TAG_RTSR>, it tells AdaptiveCmdPrompt to generate a random string 12 characters long with "tag" as prefix and then store it into the variable TAG_RDS_PASSWD which later can be used in the database creation command as a password.



In templates folder:

We will find a set of terraform modules to terraform a small AWS environment with compute and store.

We will also find a set of k8s templates for variety of settings.  












