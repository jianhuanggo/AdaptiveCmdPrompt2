package process

import (
	"Con_Utils/util"
	"bytes"
	"context"
	"errors"
	"fmt"
	tagerror "github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Artificial Timeout between each command
var pgTimeoutBetweenCmd = 5

// TagCommandContext takes in a command and associated arguments and perform a system call accordingly
// It returns error from os else null

func TagCommandContext(command string, args ...string) error {

	log.Println("command: " + command)
	log.Println("args: " + strings.Join(args, " "))
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err := exec.CommandContext(ctx, command, args...).Run(); err != nil {
		return tagerror.Wrap(err, "In TagCommandContext, Error!!!")
	}

	return nil
}

// TagCommandShell takes in a series of commands as input and their associated arguments and perform a system call accordingly
// This function also supports setting such setting working directory and environment variables
// It returns error from os else null

func TagCommandShell(cmdArgs string, cmdEnv []string) error {

	cmd := exec.Command("sh", "-c", cmdArgs)
	for _, _ev := range cmdEnv {
		//cmd.Env = append(cmd.Env, "MY_VAR=some_value")
		cmd.Env = append(cmd.Env, _ev)
	}

	log.Println("executing: ", cmdArgs)
	log.Println("environment variable(s): ", strings.Join(cmdEnv, " "))

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		return tagerror.Wrap(err, "In TagCommandShell, Error!!!")

	}

	if strings.Contains(cmdArgs, "aws ") {
		if len(stderr.String()) > 0 {
			return tagerror.Wrap(errors.New("aws errors"+stderr.String()), "In TagCommandShell, Error!!!")
		}
	}

	log.Printf("translated phrase: %q\n", stdout.String())

	time.Sleep(time.Duration(pgTimeoutBetweenCmd) * time.Second)

	return nil
}

// TagCommandShell is public facing wrapper around tagCommandShell
// This function supports setting such setting working directory and environment variables
// This function support advance features such as conditional statements, retries, random sequence generators, etc...
// It returns error from os else null

func TagCommandShell2(cmdArgs map[string]string, cmdEnv []string) (string, error) {

	tagIgnoreFlg := 0
	retryCmd := ""
	exportEnvVarName := ""
	tagWorkingDir := ""

	log.Info("\n")
	log.Info("\n")
	log.Info(cmdArgs)
	log.Info("view all attributes")

	if _tagIf, ok := cmdArgs["IGNORE_FLG"]; ok {
		val, e := strconv.Atoi(_tagIf)
		if e != nil {
			log.Fatal("Ignore Flag is corrupted...")
		} else {
			tagIgnoreFlg = val
		}
	}

	if _tagRT, ok := cmdArgs["RETRY"]; ok {
		retryCmd = _tagRT
	}

	if _tagExp, ok := cmdArgs["EXPORT"]; ok {
		exportEnvVarName = _tagExp
	}

	if _tagWD, ok := cmdArgs["WORKDIR"]; ok {
		tagWorkingDir = _tagWD
	}

	if _tagIf, ok := cmdArgs["IF"]; ok {
		if _tagIf != "" {
			_, err := tagCommandShell2(_tagIf, cmdEnv, 0, "", tagWorkingDir)
			if err != nil {
				log.Info("IF criterion evaluated to false, skipping the command...")
				return "", nil
			}
		}
	}

	if _tagIfNot, ok := cmdArgs["IFNOT"]; ok {

		if _tagIfNot != "" {
			_, err := tagCommandShell2(_tagIfNot, cmdEnv, 0, "", tagWorkingDir)

			if err == nil {

				log.Info("IFNOT criterion evaluated to false, skipping the command...")
				return "", nil
			}
		}
	}

	if _tagRSTR, ok := cmdArgs["RSTR"]; ok {
		if exportEnvVarName != "" && _tagRSTR != "" {
			tagInput := strings.Split(_tagRSTR, ";")
			log.Info(tagInput)
			val, err := strconv.Atoi(tagInput[1])
			if err != nil {
				log.Fatal("2nd Argument for random string generator needs to be an integer")
			}
			return exportEnvVarName + "=" + util.GenRandomString(tagInput[0], val), nil
		}
	}

	if tagCmd, ok := cmdArgs["COMMAND"]; ok {

		exportEnvPair, err := tagCommandShell2(tagCmd, cmdEnv, tagIgnoreFlg, exportEnvVarName, tagWorkingDir)
		if err != nil {
			if retryCmd != "" {
				log.Info(err)
				log.Info("\nRetrying ... running " + retryCmd)
				_, err = tagCommandShell2(retryCmd, cmdEnv, 1, "", tagWorkingDir)

				exportEnvPair, err = tagCommandShell2(tagCmd, cmdEnv, tagIgnoreFlg, exportEnvVarName, tagWorkingDir)
				if err != nil {
					log.Info("Retrying failed ... ")
					log.Fatal(err)
				}

			} else {
				log.Fatal(err)
			}
		}
		log.Info("this export env pair :" + exportEnvPair)
		return exportEnvPair, nil

	}

	return "", tagerror.Wrap(errors.New("can not parse command"), "Tag_ErrCode_101")
}

// tagCommandShell2 is internal facing function that takes in a series of commands as input and their associated arguments and perform a system call accordingly
// This function also supports setting such setting working directory and environment variables
// It returns error from os else null

func tagCommandShell2(cmdArgs string, cmdEnv []string, ignoreFlg int, export string, workingDir string) (string, error) {

	a, b := TagCmdTransform(cmdArgs, cmdEnv)
	cmdArgs = a.(string)

	cmd := exec.Command("sh", "-c", cmdArgs)

	for _, _ev := range b.([]string) {
		cmd.Env = append(cmd.Env, _ev)
	}

	if workingDir != "" {
		cmd.Dir = workingDir
	}

	log.Println("executing: ", cmdArgs+"\n")
	log.Println("environment variable(s): ", strings.Join(cmdEnv, " "))

	var stdout, stderr bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdout)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderr)

	err := cmd.Run()

	if ignoreFlg == 0 {
		if err != nil {
			//fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
			return "", tagerror.Wrap(err, fmt.Sprint(err)+": "+stderr.String())
			//log.Fatal(fmt.Sprint(err) + ": " + stderr.String())
		}
		fmt.Println(cmdArgs)
		if strings.Contains(cmdArgs, "aws ") {
			if len(stderr.String()) > 0 {
				return "", tagerror.Wrap(errors.New(stderr.String()), "Tag_ErrCode_101")
				//log.Fatal(stderr.String())
			}
		}
	}
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	log.Printf("\n\n%s\nerr:\n%s\n", outStr, errStr)

	if export != "" {
		return export + "=" + stdout.String(), nil

	}

	time.Sleep(time.Duration(pgTimeoutBetweenCmd) * time.Second)
	return "", nil
}

// This internal function extracts the command from command strings
// It returns the command

func tagParseCmd(cmd string) string {
	return util.FirstWords(strings.TrimLeft(cmd, "\r\n\t "), 1)
}

// This is expandable public facing wrapper currently hosts sedCmd
// It returns environment variable name and environment variable value

func TagCmdTransform(s string, a []string) (interface{}, interface{}) {
	m := map[string]interface{}{
		"sed": sedCmd,
	}

	if val, ok := m[tagParseCmd(s)]; ok {
		// pure function piercing
		return val.(func(string, []string) (string, []string))(s, a)
	}
	return s, a
}

// This internal function to perform character escaping on environment variables

func sedCmd(tagxCmd string, tagEnv []string) (string, []string) {

	retEnv := []string{}
	for _, v := range tagEnv {
		retEnv = append(retEnv, strings.Replace(v, "/", "\\/", -1))
	}
	return tagxCmd, retEnv

}
