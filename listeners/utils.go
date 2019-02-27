package listeners

import "github.com/sirupsen/logrus"

func logError(err error, msg string) {
	logrus.WithFields(logrus.Fields{
		"error": err,
	}).Error(msg)
}
