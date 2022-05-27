package Logging

import (
	"Con_Utils/util"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"time"
)

var Application = ""

// This public facing function implements a logger

func Logging(path string, loglevel string) error {
	switch loglevel {
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warning":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "panic":
		log.SetLevel(log.PanicLevel)

	}

	var f *os.File
	var err error
	err = util.Makedirectory(path)
	log.Println(path)
	if err != nil {
		errors.Wrap(err, "could not open directory: "+path)
	}
	Logfilename := path + "/" + Application + " " + strconv.FormatInt(time.Now().UTC().Unix(), 10) + "log"
	log.Printf("Logfile is created\n")
	if f, err = os.Create(Logfilename); err != nil {
		errors.Wrap(err, "could not create log file: "+Logfilename)
	}
	if _, err = os.Stat(Logfilename); os.IsNotExist(err) {
		errors.Wrap(err, "could not see log file: "+Logfilename)
	}
	log.SetOutput(f)
	return nil
}
