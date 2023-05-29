package bdpan

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

var (
	Log     *logrus.Logger
	logFile *os.File
)

func initLogger() {
	GetLogger()
	logFile, _ = os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModeAppend|os.ModePerm)
}

func createLogger() *logrus.Logger {
	log := logrus.New()
	// log.SetReportCaller(true)
	log.SetFormatter(&LogFormatter{
		TextFormatter: &logrus.TextFormatter{
			ForceColors:     true,
			FullTimestamp:   true,
			TimestampFormat: "01-02 15:04:05",
			// CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			// filename := filepath.Base(frame.File)
			// return "", fmt.Sprintf("[%s:%d]", filename, frame.Line)
			// },
		},
	})
	Log = log
	return log
}

type LogFormatter struct {
	*logrus.TextFormatter
}

func (f *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return f.TextFormatter.Format(entry)
}

func SetLogLevel(level logrus.Level) {
	Log.SetLevel(level)
}

// 设置日志输出文件
func SetOutputFile() {
	Log.SetOutput(logFile)
}

func GetLogger() *logrus.Logger {
	if Log == nil {
		Log = createLogger()
	}
	return Log
}

func IsLoggerDebug() bool {
	return Log.GetLevel() == logrus.DebugLevel
}

func LogInfoString(w io.Writer, s string) {
	Log.SetOutput(w)
	Log.Info(s)
}
