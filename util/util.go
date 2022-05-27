package util

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/bigkevmcd/go-configparser"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func Makedirectory(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, os.ModePerm)
		log.Info("this directory is created", path)
		return errors.Wrap(err, "could not create the directory")
	}
	return nil
}

func GetConfigFile(config_name string) error {
	dirname, err := os.UserHomeDir()
	if err != nil {
		return errors.Wrap(err, "can not retrieve home directory")
	}

	p, err := configparser.NewConfigParserFromFile(filepath.Join(dirname, "go", "src", "Con_Utils", "project", "awsenvsetup", "awstag.conf"))
	if err != nil {
		return errors.Wrap(err, "can not find aws credentials at "+filepath.Join(dirname, ".aws", "credentials"+
			" \nPlease make sure aws client is installed and profile is setup appropriately.  See doc at..."))
	}

	v, err := p.Get(config_name, "1")

	if err != nil {
		return errors.Wrap(err, "can not find specified profile name, here are the valid profile name(s): ["+strings.Join(p.Sections(), " ")+"]\n")
	}
	fmt.Println(v)

	return nil

}

func Check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func TagWriteFile(filePath string, content string) {

	err := ioutil.WriteFile(filePath, []byte(content), 0644)
	Check(err)
}

func TagReadFile(filePath string) string {
	dat, err := ioutil.ReadFile(filePath)
	Check(err)
	return string(dat[:])
}

func TagJson2File(filePath string, data []byte) {
	file, err := json.MarshalIndent(data, "", " ")
	Check(err)
	err = ioutil.WriteFile(filePath, file, 0644)
	Check(err)
}

func GetStringInBetween(value string, a string, b string) string {
	// Get substring between two strings.
	posFirst := strings.Index(value, a)
	if posFirst == -1 {
		return ""
	}
	posLast := strings.Index(value, b)
	if posLast == -1 {
		return ""
	}
	posFirstAdjusted := posFirst + len(a)
	if posFirstAdjusted >= posLast {
		return ""
	}
	return value[posFirstAdjusted:posLast]
}

func IsFileExist(filepath string) bool {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return false
	}
	return true
}

func CheckEmptyFile(filepath string) bool {
	fi, err := os.Stat(filepath)
	if err != nil {
		Check(err)
	}
	if fi.Size() == 0 {
		return true
	} else {
		return false
	}
}

func RemoveStrSliceByIndex(s []string, i int) []string {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

func RemoveSliceItemYamlMapIndOrder(s []map[string][]string, index int) []map[string][]string {
	return append(s[:index], s[index+1:]...)
}

func FirstWords(value string, count int) string {
	// Loop over all indexes in the string.
	for i := range value {
		// If we encounter a space, reduce the count.
		if value[i] == ' ' {
			count -= 1
			// When no more words required, return a substring.
			if count == 0 {
				return value[0:i]
			}
		}
	}
	// Return the entire string.
	return value
}

func GenRandomString(prefix string, length int) string {
	id, err := uuid.NewRandom()
	if err != nil {
		return (prefix + uuid.New().String())[:length]
	}
	return (prefix + id.String())[:length]
}

func TagEnv2Map(a []string) map[string]string {
	elementMap := make(map[string]string)
	for _, v := range a {
		elementMap[strings.Split(v, "=")[0]] = strings.Split(v, "=")[1]
	}
	return elementMap
}

func TagMap2Env(a map[string]string) []string {
	elementSlice := []string{}
	for key, val := range a {
		elementSlice = append(elementSlice, key+"="+val)
	}
	return elementSlice
}

func ReplaceMerge(a []string, b []string) []string {
	a_ := TagEnv2Map(a)
	b_ := TagEnv2Map(b)
	for key, val := range b_ {
		a_[key] = val
	}
	return TagMap2Env(a_)
}

func GetFileHash(filepath string) string {
	f, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}
	return hex.EncodeToString(h.Sum(nil))
}

func VerifyConfigFileIntegrity(stateFilepath string, envFilepath string, hashFilepath string, ConfigFilepath string) bool {
	if (!IsFileExist(stateFilepath) || CheckEmptyFile(stateFilepath)) && (!IsFileExist(envFilepath) || CheckEmptyFile(envFilepath)) {
		TagWriteFile(hashFilepath, GetFileHash(ConfigFilepath))
		return true
	} else {
		if GetFileHash(ConfigFilepath) == TagReadFile(hashFilepath) {
			return true
		}
	}
	return false
}
