package log

import (
	"bytes"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	nested "github.com/antonfisher/nested-logrus-formatter"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

// IPAddress : c.ClientIP() set from controller before entering service
// Session   : ksuid.New().String() or gonanoid.Generate(x, 12) set from first caller
// ActorID   : from client
// ActorType : MOB/BOF/MSQ/SYS/SCH --> mobile / backoffice / message queuing / system / scheduller

// Data is
type Data struct {
	IPAddress string `` // the ip address of caller
	Session   string `` // id that generated in controller and passed to service to service for flow tracking purpose
	ActorID   string `` // could be userID from Apps or backoffice
	ActorType string `` // MOB (Mobile Apps) / BOF (Backoffice) / MSQ (message queuing) / SCH (scheduller) / SYS (system)
}

// ILogger is
type ILogger interface {
	Debug(data interface{}, description string, args ...interface{})
	Info(data interface{}, description string, args ...interface{})
	Warn(data interface{}, description string, args ...interface{})
	Error(data interface{}, description string, args ...interface{})
	Fatal(data interface{}, description string, args ...interface{})
	Panic(data interface{}, description string, args ...interface{})
}

// LogrusImpl is
type logger struct {
	appName    string
	appVersion string

	filePath string
	level    logrus.Level
	maxAge   time.Duration

	theLogger *logrus.Logger
}

var defaultLogger logger
var defaultLoggerOnce sync.Once

func createLogger() {
	// formatter := logrus.JSONFormatter{}

	formatter := nested.Formatter{
		HideKeys:        true,
		TimestampFormat: "2006-01-02 15:04:05",
		FieldsOrder:     []string{"func"},
	}

	defaultLogger.theLogger = logrus.New()
	defaultLogger.theLogger.SetLevel(defaultLogger.level)
	defaultLogger.theLogger.SetFormatter(&formatter)

	filename := defaultLogger.appName + ".%Y%m%d.log"
	if len(defaultLogger.appVersion) > 0 {
		filename = defaultLogger.appName + "-" + defaultLogger.appVersion + ".%Y%m%d.log"
	}

	writer, _ := rotatelogs.New(
		filepath.Join(defaultLogger.filePath, filename),
		// rotatelogs.WithLinkName(path),
		rotatelogs.WithMaxAge(defaultLogger.maxAge),
		rotatelogs.WithRotationTime(time.Duration(24)*time.Hour),
	)

	defaultLogger.theLogger.AddHook(lfshook.NewHook(
		lfshook.WriterMap{
			logrus.InfoLevel:  writer,
			logrus.WarnLevel:  writer,
			logrus.ErrorLevel: writer,
			logrus.DebugLevel: writer,
		},
		defaultLogger.theLogger.Formatter,
	))
}

// WithFile is command to state the log will printing to files
// the rolling log file will put in logs/ directory
//
// filename is just a name of log file without any extension
//
// maxAge is age (in days) of the logs file before it gets purged from the file system
func Init() {
	defaultLogger.appName = "STRIPE-SUBSCRIPTION"
	defaultLogger.appVersion = "BACKEND"
	defaultLogger.filePath = "logs/"
	defaultLogger.level = logrus.DebugLevel
	defaultLogger.maxAge = time.Duration(24*7) * time.Hour
}

// GetLog is
func GetLog() ILogger {
	defaultLoggerOnce.Do(createLogger)
	return &defaultLogger
}

func (l *logger) getLogEntry(extraInfo interface{}) *logrus.Entry {
	pc, _, _, _ := runtime.Caller(2)
	funcName := runtime.FuncForPC(pc).Name()

	var buffer bytes.Buffer

	buffer.WriteString("fn:")

	x := strings.LastIndex(funcName, "/")
	buffer.WriteString(funcName[x+1:])

	if extraInfo == nil {
		return l.theLogger.WithField("info", buffer.String())
	}

	data, ok := extraInfo.(Data)
	if !ok {
		return l.theLogger.WithField("info", buffer.String())
	}

	if data.IPAddress != "" {
		buffer.WriteString("|ip:")
		buffer.WriteString(data.IPAddress)
	}

	if data.Session != "" {
		buffer.WriteString("|ss:")
		buffer.WriteString(data.Session)
	}

	if data.ActorID != "" {
		buffer.WriteString("|id:")
		buffer.WriteString(data.ActorID)
	}

	if data.ActorType != "" {
		buffer.WriteString("|tp:")
		buffer.WriteString(data.ActorType)
	}

	return l.theLogger.WithField("info", buffer.String())
}

// Debug is
func (l *logger) Debug(data interface{}, description string, args ...interface{}) {
	l.getLogEntry(data).Debugf(description+"\n", args...)
}

// Info is
func (l *logger) Info(data interface{}, description string, args ...interface{}) {
	l.getLogEntry(data).Infof(description+"\n", args...)
}

// Warn is
func (l *logger) Warn(data interface{}, description string, args ...interface{}) {
	l.getLogEntry(data).Warnf(description+"\n", args...)
}

// Error is
func (l *logger) Error(data interface{}, description string, args ...interface{}) {
	l.getLogEntry(data).Errorf(description+"\n", args...)
}

// Fatal is
func (l *logger) Fatal(data interface{}, description string, args ...interface{}) {
	l.getLogEntry(data).Fatalf(description+"\n", args...)
}

// Panic is
func (l *logger) Panic(data interface{}, description string, args ...interface{}) {
	l.getLogEntry(data).Panicf(description+"\n", args...)
}
