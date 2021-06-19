package log

import "testing"

func TestLog(t *testing.T) {
	Trace("Trace")
	Tracef("%s", "Trace format")
	Debug("Debug")
	Debugf("%s", "Debug format")
	Info("Info")
	Infof("%s", "Info format")
	Warning("Warning")
	Warningf("%s", "Warning format")
	Error("Error")
	Errorf("%s", "Error format")
}
