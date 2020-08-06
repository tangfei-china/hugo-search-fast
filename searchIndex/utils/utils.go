package utils

import (
	log "github.com/sirupsen/logrus"
	"os"
	"regexp"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true,})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
}

//正则提取内容
func MatchString(p string, str string) (bool, string) {
	flowRegexp := regexp.MustCompile(p)
	params := flowRegexp.FindStringSubmatch(str)
	if len(params) > 0 {
		return true, params[1]
	}
	return false, ""
}

