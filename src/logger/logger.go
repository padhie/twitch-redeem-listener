package logger

import (
	"io"
	"log"
	"log/syslog"
	"os"
	"path/filepath"

	"twitch-redeem-trigger/src/config"
)

type Logger struct {
	debugLogger    *log.Logger
	infoLogger     *log.Logger
	noticeLogger   *log.Logger
	warningLogger  *log.Logger
	errorLogger    *log.Logger
	criticalLogger *log.Logger
	alertLogger    *log.Logger
	fatalLogger    *log.Logger
}

func Build(cfgLogging config.Logging) *Logger {
	l := &Logger{}

	writer := getWriter(cfgLogging)
	l.debugLogger = log.New(writer, "[DEBUG] ", log.Ldate|log.Ltime|log.Lshortfile)
	l.infoLogger = log.New(writer, "[INFO] ", log.Ldate|log.Ltime)
	l.noticeLogger = log.New(writer, "[NOTICE] ", log.Ldate|log.Ltime|log.Lshortfile)
	l.warningLogger = log.New(writer, "[WARNING] ", log.Ldate|log.Ltime|log.Lshortfile)
	l.errorLogger = log.New(writer, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
	l.criticalLogger = log.New(writer, "[CRITICAL] ", log.Ldate|log.Ltime|log.Lshortfile)
	l.alertLogger = log.New(writer, "[ALERT] ", log.Ldate|log.Ltime|log.Lshortfile)
	l.fatalLogger = log.New(writer, "[FATAL] ", log.Ldate|log.Ltime|log.Lshortfile)

	return l
}

func getWriter(cfgLogging config.Logging) io.Writer {
	if cfgLogging.UseSyslog {
		syslogWriter, err := syslog.New(syslog.LOG_NOTICE, "twitch-redeem")
		if err != nil {
			log.Fatalf("Failed to connect to syslog: %v", err)
		}
		return syslogWriter
	}

	if cfgLogging.LogFile != "" {
		// Verzeichnis erstellen, falls es nicht existiert
		dir := filepath.Dir(cfgLogging.LogFile)
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Printf("Failed to create log directory: %v", err)
		}

		fileWriter, err := os.OpenFile(cfgLogging.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Failed to open log file: %v", err)
		}
		return io.MultiWriter(os.Stdout, fileWriter)
	}

	return os.Stdout
}

func (l *Logger) Debug(format string, v ...interface{}) {
	l.debugLogger.Printf(format, v...)
}

func (l *Logger) Info(format string, v ...interface{}) {
	l.infoLogger.Printf(format, v...)
}

func (l *Logger) Notice(format string, v ...interface{}) {
	l.noticeLogger.Printf(format, v...)
}

func (l *Logger) Warning(format string, v ...interface{}) {
	l.warningLogger.Printf(format, v...)
}

func (l *Logger) Error(format string, v ...interface{}) {
	l.errorLogger.Printf(format, v...)
}

func (l *Logger) Critical(format string, v ...interface{}) {
	l.criticalLogger.Printf(format, v...)
}

func (l *Logger) Alert(format string, v ...interface{}) {
	l.alertLogger.Printf(format, v...)
}

func (l *Logger) Fatal(format string, v ...interface{}) {
	l.fatalLogger.Printf(format, v...)
	os.Exit(1)
}
