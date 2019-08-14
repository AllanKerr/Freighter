package log

import (
	"bufio"
	"encoding/json"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
)

type unixRemoteLogger struct {
	logger *logrus.Entry
	logC   *os.File
	logP   *os.File
}

func NewRemoteLogger(name string) (RemoteLogger, error) {

	loggerC := logrus.New()
	loggerC.SetLevel(childLogLevel)
	loggerC.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	stdC := loggerC.WithField("_name", name)

	fds, err := unix.Socketpair(unix.AF_LOCAL, unix.SOCK_STREAM, 0)
	if err != nil {
		return nil, err
	}
	logC := os.NewFile(uintptr(fds[1]), "log-c")
	logP := os.NewFile(uintptr(fds[0]), "log-p")
	unix.CloseOnExec(fds[1])
	unix.CloseOnExec(fds[0])

	r := &unixRemoteLogger{
		logger: stdC,
		logC:   logC,
		logP:   logP,
	}
	return r, nil
}

func (r *unixRemoteLogger) Listen() {
	go func() {
		scanner := bufio.NewScanner(r.logP)
		for scanner.Scan() {
			r.replayLogMessage(scanner.Text())
		}
	}()
}

func (r *unixRemoteLogger) Child() *os.File {
	return r.logC
}

func (r *unixRemoteLogger) replayLogMessage(text string) {

	var (
		result map[string]interface{}
		level  logrus.Level
		msg    string
	)
	json.Unmarshal([]byte(text), &result)

	fieldLogger := r.logger
	for k := range result {
		if k == "msg" {
			msg = result["msg"].(string)
		} else if k == "level" {
			levelText := result["level"].(string)
			lvl, err := logrus.ParseLevel(levelText)
			if err != nil {
				r.logger.WithField("text", text).WithError(err).Error("Invalid log level received from child")
				return
			}
			level = lvl
		} else if k == "time" {
			timestampText := result["time"].(string)
			timestamp, err := time.Parse(time.RFC3339, timestampText)
			if err != nil {
				r.logger.WithField("text", text).WithError(err).Error("Invalid time received from child")
				return
			}
			fieldLogger = fieldLogger.WithTime(timestamp)
		} else {
			fieldLogger = fieldLogger.WithField(k, result[k])
		}
	}
	fieldLogger.Log(level, msg)
}
