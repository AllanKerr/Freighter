package log

import (
	"github.com/sirupsen/logrus"
)

var (
	parentLogLevel = logrus.GetLevel()
	childLogLevel  = logrus.GetLevel()
)
