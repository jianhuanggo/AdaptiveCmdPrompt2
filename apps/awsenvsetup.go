package apps

import (
	"Con_Utils/process"
	"Con_Utils/util"
	"encoding/json"
	"fmt"
	"github.com/bigkevmcd/go-configparser"
	tagerror "github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	//"strconv"
	"strings"
)

var Region string
var AwsProfile string
var TagProject string
var TagNewStartFlag = false
var tagRemoveStateF string

// This public facing function starts the workflow for aws environment setup
// 1) parses the template which is in yaml format
// 2) invokes TagawsEnvPrepare function to parse the template, the template can be break down into multiple sections
// 3) checks whether there is running version of template already, if so, valids the integrity of templates and resume from previous completed step
// 4) invokes TagawsEnvExecute function to run the statement in the template one by one, auto retry each statement if it is specified in the template
// 5) Once the statement is completed successful, record the state to support restartability

func TagawsEnvRun(tagEnvV map[string]string) {

	tagCurrDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(tagCurrDir)

	fmt.Println(TagProject)

	err = util.Makedirectory(filepath.Join(tagCurrDir, ".tag"))
	if err != nil {
		util.Check(err)
	}

	if TagNewStartFlag {
		if util.IsFileExist(filepath.Join(tagCurrDir, ".tag", TagProject+".state")) ||
			util.IsFileExist(filepath.Join(tagCurrDir, ".tag", TagProject+".env")) ||
			util.IsFileExist(filepath.Join(tagCurrDir, ".tag", TagProject+".chechsum")) {
			for ok := true; ok; {
				fmt.Print("New start flag is specified.  Remove current state files (Y/n):")
				fmt.Scanln(&tagRemoveStateF)

				if tagRemoveStateF == "Y" {
					util.TagRemoveFile(filepath.Join(tagCurrDir, ".tag", TagProject+".state"))
					util.TagRemoveFile(filepath.Join(tagCurrDir, ".tag", TagProject+".env"))
					util.TagRemoveFile(filepath.Join(tagCurrDir, ".tag", TagProject+".checksum"))
					break
				} else if tagRemoveStateF == "n" {
					log.Println("Exiting...")
					os.Exit(0)
				} else {
					log.Println("Choice is not found, try again")
				}
			}
		} else {
			log.Info("No state files are detected, ignore new start flag....")

		}
	}

	log.Info("Running....")

	err = tag_validate(AwsProfile)
	if err != nil {
		util.Check(err)
	}

	tagYamlConfig, err := util.TagGetYaml3(filepath.Join(tagCurrDir, "project", "awsenvsetup", TagProject+".yml"))
	if err != nil {
		util.Check(err)
	}

	tagYamlConf, tagEnvSlice := TagawsEnvPrepare(tagYamlConfig, tagEnvV)

	// Verify the integrity of project configuration file.  If some steps have already ran which implies the intermediate state has been saved,
	// adding/removing/modifying project configuration file will be dangerous since it can inadvertently net us unforeseeable results.
	// for example, adding a step to prior to the latest successful step saved in the state file will cause an already completed command to run again

	if !util.VerifyConfigFileIntegrity(filepath.Join(tagCurrDir, ".tag", TagProject+".state"),
		filepath.Join(tagCurrDir, ".tag", TagProject+".env"),
		filepath.Join(tagCurrDir, ".tag", TagProject+".checksum"),
		filepath.Join(tagCurrDir, "project", "awsenvsetup", TagProject+".yml")) {
		log.Fatal("Configuration file " + filepath.Join(tagCurrDir, "project", "awsenvsetup", TagProject+".yml") + " Integrity has been promised, either revert back the change or start the run with a clean state using -n flag")
	}

	for i, v := range tagYamlConf {
		for k, val := range v {
			fmt.Println(i, k, val)
		}

	}

	for m, n := range tagEnvSlice {
		fmt.Println(m, n)
	}

	err = TagawsEnvExecute(filepath.Join(tagCurrDir, ".tag", TagProject+".state"),
		filepath.Join(tagCurrDir, ".tag", TagProject+".env"),
		tagYamlConf, tagEnvSlice)
	if err != nil {
		util.Check(err)
	}
}

// This public facing function parses environment variable section of aws environment setup template and
// returns a map to the main program

func TagawsEnvPrepare(TagConfigYamlContent []map[string][]string, TagFlagContent map[string]string) ([]map[string][]string, []string) {

	_envYamlMap := map[string][]string{}

	for i, v := range TagConfigYamlContent {
		if _, ok := v["Envvariable"]; ok {
			_envYamlMap = v
			TagConfigYamlContent = util.RemoveSliceItemYamlMapIndOrder(TagConfigYamlContent, i)
			break
		}
	}

	_tagEnvList := make([]string, 0, len(TagFlagContent))

	for k, v := range TagFlagContent {
		_tagEnvList = append(_tagEnvList, k+"="+v)
	}

	_s := util.ReplaceMerge(_envYamlMap["Envvariable"], _tagEnvList)

	for k, v := range _s {
		log.Println(k, v)
	}

	return TagConfigYamlContent, _s

}

// This public facing function prepares and executes statement
// It currently supports following special tags to enhance the execution of command line statement.  And tags can be extended when needed.
// <IF> conditional logic to only execute subsequent statement if previous statement returns True
// <IFNOT> conditional logic to only execute subsequent statement if previous statement returns False
// <RETRY> retry the statement, it doesn't have to be the original statement
// <EXPORT> Inject additional variable to environment variable bank
// <WORKDIR> Setting working directory
// <RSTR> Random string generate, can be used to generate password on the fly
// <COMMAND> regular command line statement
// returns error object

func TagawsEnvExecute(stateFilepath string, envFilepath string, TagCmdSliceMap []map[string][]string, TagEnvSlice []string) error {

	pipeline := map[string]int{}
	TagInstMap := map[string][]string{"IF": {"<TAG_IF>", "<TAG_FI>"},
		"IFNOT":      {"<TAG_IFNOT>", "<TAG_TONFI>"},
		"IGNORE_FLG": {"<TAG_IG>", "<TAG_GI>"},
		"COMMAND":    {"<TAG_CMD>", "<TAG_DMC>"},
		"RETRY":      {"<TAG_RT>", "<TAG_TR>"},
		"EXPORT":     {"<TAG_EXP>", "<TAG_PXE>"},
		"WORKDIR":    {"<TAG_WD>", "<TAG_DW"},
		"RSTR":       {"<TAG_RSTR>", "<TAG_RTSR>"},
	}

	TagCmdMap := map[string]string{"IF": "", "IFNOT": "", "IGNORE_FLG": "0", "COMMAND": "", "RETRY": "", "EXPORT": "", "TAG_WD": "", "RSTR": ""}
	TagAllEnvVar := TagEnvSlice

	stepTracking := func(section string) int {
		if _, ok := pipeline[section]; ok {
			pipeline[section]++
		} else {
			pipeline[section] = 1
		}
		return pipeline[section]
	}

	// anonymous function to save states
	saveState := func() {
		fmt.Println("\n\npipeline debugging")
		fmt.Println(pipeline)
		fmt.Println("pipeline debugging\n\n")
		j, err := json.Marshal(pipeline)
		util.Check(err)
		util.TagWriteFile(stateFilepath, string(j))

	}

	// anonymous function to save environment variables
	saveEnvVar := func() {
		j, err := json.Marshal(TagAllEnvVar)
		util.Check(err)
		util.TagWriteFile(envFilepath, string(j))
	}

	state, err := tagGetState(stateFilepath)
	if err != nil {
		util.Check(err)
	}

	env, err := tagGetEnv(envFilepath)
	if err != nil {
		util.Check(err)
	}

	TagAllEnvVar = util.ReplaceMerge(TagAllEnvVar, env)

	log.Info(TagCmdSliceMap)
	fmt.Println(reflect.TypeOf(TagCmdSliceMap))

	for _, tag_phases := range TagCmdSliceMap {
		for tag_phase, tag_cmd_list := range tag_phases {
			fmt.Println("\n\n\n\n\n")

			for _, cmd := range tag_cmd_list {

				_step := stepTracking(tag_phase)

				if procState, ok := state[tag_phase]; ok {
					if _step <= procState {
						log.Info("In " + tag_phase + ": Step " + strconv.Itoa(_step) + " is already completed.  Skipping... ")

						continue
					}
				}

				log.Println("In " + tag_phase + ": Step " + strconv.Itoa(_step) + " / " + strconv.Itoa(len(tag_cmd_list)))

				for key, val := range TagInstMap {
					cmdString := util.GetStringInBetween(cmd, val[0], val[1])
					if cmdString != "" {
						TagCmdMap[key] = cmdString
					}
				}

				exportEnvPair, err := process.TagCommandShell2(TagCmdMap, TagAllEnvVar)
				if err != nil {
					log.Fatal(err)
				}

				if exportEnvPair != "" {
					for index, val := range TagAllEnvVar {
						if strings.Split(exportEnvPair, "=")[0] == strings.Split(val, "=")[0] {
							TagAllEnvVar = util.RemoveStrSliceByIndex(TagAllEnvVar, index)
							break
						}
					}
					TagAllEnvVar = append(TagAllEnvVar, strings.TrimRight(exportEnvPair, "\r\n"))
				}

				// clean the command structure
				TagCmdMap = map[string]string{"IF": "", "IFNOT": "", "IGNORE_FLG": "0", "COMMAND": "", "RETRY": "", "EXPORT": "", "TAG_WD": "", "RSTR": ""}

				saveState()
				saveEnvVar()

			}
			fmt.Println("\n\n*****Completed " + tag_phase + " stage*****\n\n")

		}

	}

	return nil

}

// This public facing function gets state information which is stored in the JSON format
// returns a map of phases with corresponding number of the steps with in each phase that is completed

func tag_validate(profile_name string) error {
	dirname, err := os.UserHomeDir()
	if err != nil {
		return tagerror.Wrap(err, "can not retrieve home directory")
	}

	p, err := configparser.NewConfigParserFromFile(filepath.Join(dirname, ".aws", "credentials"))
	if err != nil {
		return tagerror.Wrap(err, "can not find aws credentials at "+filepath.Join(dirname, ".aws", "credentials"+
			" \nPlease make sure aws client is installed and profile is setup appropriately.  See doc at..."))
	}

	_, err = p.Get(profile_name, "aws_access_key_id")
	if err != nil {
		return tagerror.Wrap(err, "can not find specified profile name, here are the valid profile name(s): ["+strings.Join(p.Sections(), " ")+"]\n")
	}

	return nil

}

// This public facing function gets state information which is stored in the JSON format
// returns a map of phases with corresponding number of the steps with in each phase that is completed

func tagGetState(filepath string) (map[string]int, error) {
	state := map[string]int{}

	//ignore if the state file doesn't exist which will be created later
	if !util.IsFileExist(filepath) {
		return nil, nil
	}
	if util.CheckEmptyFile(filepath) == false {
		tagProcState := util.TagReadFile(filepath)
		log.Info("content of file :" + tagProcState)

		err := json.Unmarshal([]byte(tagProcState), &state)
		if err != nil {
			return nil, tagerror.Wrap(err, "TagErrCode_1000")
		}
	} else {
		log.Info("warning: state file" + filepath + " is empty. ")

	}

	return state, nil
}

// This public facing function gets saved environment variable info which is stored in the JSON format
// returns a map of environment variables with corresponding values

func tagGetEnv(filepath string) ([]string, error) {
	env := []string{}

	//ignore if the state file doesn't exist which will be created later
	if !util.IsFileExist(filepath) {
		return nil, nil
	}
	if util.CheckEmptyFile(filepath) == false {
		tagProcEnv := util.TagReadFile(filepath)
		log.Info("content of file :" + tagProcEnv)

		err := json.Unmarshal([]byte(tagProcEnv), &env)
		if err != nil {
			return nil, tagerror.Wrap(err, "TagErrCode_1000")
		}
	} else {
		log.Info("warning: state file" + filepath + " is empty. ")

	}

	return env, nil
}
