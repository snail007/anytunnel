package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

const (
	DebugLevel uint8 = 1
	InfoLevel  uint8 = 1 << 1
	WarnLevel  uint8 = 1 << 2
	ErrorLevel uint8 = 1 << 3
	FatalLevel uint8 = 1 << 4
	AllLevels  uint8 = byte(DebugLevel | InfoLevel | WarnLevel | ErrorLevel | FatalLevel)
)

var (
	flush   = &sync.WaitGroup{}
	exitchn = make(chan bool, 1)
)

type Fields map[string]string
type Writer interface {
	Write(e Entry)
	Init() error
}

type Logger struct {
	writersMap   *map[int][]*WrappedWriter
	mu           sync.Mutex
	modeSafe     bool
	wait         sync.WaitGroup
	fields       Fields
	fieldsString string
}

type WrappedWriter struct {
	lock   *sync.Mutex
	chn    chan Entry
	writer Writer
	wait   sync.WaitGroup
}
type Entry struct {
	Content               string
	Timestamp             int64
	Milliseconds          int64
	TimestampMilliseconds int64
	Level                 uint8
	LevelString           string
	Fields                map[string]string
	FieldsString          string
}
type MiniLogger interface {
	Debug(v ...interface{}) MiniLogger
	Info(v ...interface{}) MiniLogger
	Warn(v ...interface{}) MiniLogger
	Error(v ...interface{}) MiniLogger
	Fatal(v ...interface{})
	Debugf(format string, v ...interface{}) MiniLogger
	Infof(format string, v ...interface{}) MiniLogger
	Warnf(format string, v ...interface{}) MiniLogger
	Errorf(format string, v ...interface{}) MiniLogger
	Fatalf(format string, v ...interface{})
	Debugln(v ...interface{}) MiniLogger
	Infoln(v ...interface{}) MiniLogger
	Warnln(v ...interface{}) MiniLogger
	Errorln(v ...interface{}) MiniLogger
	Fatalln(v ...interface{})
	AddWriter(w Writer, levels uint8) MiniLogger
	Safe() MiniLogger
	Unsafe() MiniLogger
	With(fields Fields) MiniLogger
}

func getLevelString(level uint8) string {
	switch level {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case FatalLevel:
		return "FATAL"
	}
	return "UNKOWN"
}

// New returns a new logger
//modeSafe when false : must call logger.Flush() in main defer
//if not do that, may be lost message,but this mode has a highest performance
//when true : each message can be processed , but this mode may be has a little lower performance
//beacuse of logger must wait for  all writers process done with each messsage.
//note: if you do logging before call os.Exit(), you had better to set safeMode to true before call os.Exit().
//fields same as With()'s args,more about fields to see With().
func New(modeSafe bool, fields Fields) MiniLogger {
	if fields == nil {
		fields = Fields{}
	}
	return &Logger{
		fields:       fields,
		fieldsString: "",
		mu:           sync.Mutex{},
		modeSafe:     modeSafe,
		wait:         sync.WaitGroup{},
		writersMap: &map[int][]*WrappedWriter{
			int(DebugLevel): []*WrappedWriter{},
			int(InfoLevel):  []*WrappedWriter{},
			int(WarnLevel):  []*WrappedWriter{},
			int(ErrorLevel): []*WrappedWriter{},
			int(FatalLevel): []*WrappedWriter{},
		},
	}
}

//Flush wait for process all left entry
func Flush() {
	exitchn <- true
	flush.Wait()
}

//With return a new MiniLogger,which it's modeSafe inited same as caller MiniLogger,
//and revecived levels  same as caller MiniLogger,
//and set it's fields.
//you can use it as a standard MiniLogger.
func (l *Logger) With(fields Fields) MiniLogger {
	for k, v := range l.fields {
		fields[k] = v
	}
	s, _ := json.Marshal(fields)
	return &Logger{
		mu:           sync.Mutex{},
		modeSafe:     l.modeSafe,
		wait:         sync.WaitGroup{},
		writersMap:   l.writersMap,
		fields:       fields,
		fieldsString: string(s),
	}
}

//Unsafe : if you call Unsafe() , then you must call logger.Flush() in main defer
//if not do that, may be lost message,but this mode has a highest performance
//note: if you do logging before call os.Exit(), you had better to call Safe() before call os.Exit().
//you call call func Safe() or Unsafe() in any where any time to switch safe mode.

func (l *Logger) Unsafe() MiniLogger {
	l.modeSafe = false
	return l
}

//Safe : if you call Safe(), each message can be processed immediately,
//but this mode may be has  lower performance,beacuse of logger must wait for
//all writers process done with each messsage.
//note: if you do logging before call os.Exit(), you had better to call Safe() before call os.Exit().
//you call func Safe() or Unsafe() in any where any time to switch safe mode.
func (l *Logger) Safe() MiniLogger {
	l.modeSafe = true
	return l
}
func (l *Logger) AddWriter(writer Writer, levels byte) MiniLogger {
	w := &WrappedWriter{
		lock:   &sync.Mutex{},
		writer: writer,
		chn:    make(chan Entry, 1024),
	}
	m := *l.writersMap
	if DebugLevel&levels == DebugLevel {
		m[int(DebugLevel)] = append(m[int(DebugLevel)], w)
	}
	if InfoLevel&levels == InfoLevel {
		m[int(InfoLevel)] = append(m[int(InfoLevel)], w)
	}
	if WarnLevel&levels == WarnLevel {
		m[int(WarnLevel)] = append(m[int(WarnLevel)], w)
	}
	if ErrorLevel&levels == ErrorLevel {
		m[int(ErrorLevel)] = append(m[int(ErrorLevel)], w)
	}
	if FatalLevel&levels == FatalLevel {
		m[int(FatalLevel)] = append(m[int(FatalLevel)], w)
	}
	flush.Add(1)
	go func() {
		defer func() {
			_ = recover()
			flush.Done()
		}()
		if err := w.writer.Init(); err != nil {
			fmt.Printf("init writer fail,%s", err)
			return
		}
		for {
			select {
			case entry, ok := <-w.chn:
				if ok {
					w.writer.Write(entry)
					if l.modeSafe {
						l.wait.Done()
					}
					if entry.Level == FatalLevel {
						os.Exit(0)
					}
				} else {
					return
				}
			case exit := <-exitchn:
				if exit {
					return
				}
			}
		}
	}()
	return l
}
func (l *Logger) callWriter(level byte, t, foramt string, v ...interface{}) {
	if level == FatalLevel {
		l.Safe()
	}
	c := ""
	if t == "f" && len(v) == 0 {
		v = append(v, foramt)
		foramt = "%s"
	}
	switch t {
	case "ln":
		c = fmt.Sprintln(v...)
	case "f":
		c = fmt.Sprintf(foramt, v...)
	default:
		c = fmt.Sprint(v...)
	}
	m := *l.writersMap
	for _, w := range m[int(level)] {
		now := time.Now().UnixNano()
		nowUnix := time.Now().Unix()
		mili := (now / 1000000) - nowUnix*1000
		w.chn <- Entry{
			Timestamp:             nowUnix,
			TimestampMilliseconds: now / 10000000,
			Milliseconds:          mili,
			Content:               c,
			Level:                 level,
			LevelString:           getLevelString(level),
			Fields:                l.fields,
			FieldsString:          l.fieldsString,
		}
		if l.modeSafe {
			l.wait.Add(1)
		}
	}
	if l.modeSafe {
		l.wait.Wait()
	}
}

func (l *Logger) Fatal(v ...interface{}) {
	l.callWriter(FatalLevel, "", "", v...)
}

func (l *Logger) Error(v ...interface{}) MiniLogger {
	l.callWriter(ErrorLevel, "", "", v...)
	return l
}

func (l *Logger) Warn(v ...interface{}) MiniLogger {
	l.callWriter(WarnLevel, "", "", v...)
	return l
}

func (l *Logger) Info(v ...interface{}) MiniLogger {
	l.callWriter(InfoLevel, "", "", v...)
	return l
}

func (l *Logger) Debug(v ...interface{}) MiniLogger {
	l.callWriter(DebugLevel, "", "", v...)
	return l
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.callWriter(FatalLevel, "f", format, v...)
}

func (l *Logger) Errorf(format string, v ...interface{}) MiniLogger {
	l.callWriter(ErrorLevel, "f", format, v...)
	return l
}

func (l *Logger) Warnf(format string, v ...interface{}) MiniLogger {
	l.callWriter(WarnLevel, "f", format, v...)
	return l
}

func (l *Logger) Infof(format string, v ...interface{}) MiniLogger {
	l.callWriter(InfoLevel, "f", format, v...)
	return l
}

func (l *Logger) Debugf(format string, v ...interface{}) MiniLogger {
	l.callWriter(DebugLevel, "f", format, v...)
	return l
}

func (l *Logger) Fatalln(v ...interface{}) {
	l.callWriter(FatalLevel, "ln", "", v...)
}

func (l *Logger) Errorln(v ...interface{}) MiniLogger {
	l.callWriter(ErrorLevel, "ln", "", v...)
	return l
}

func (l *Logger) Warnln(v ...interface{}) MiniLogger {
	l.callWriter(WarnLevel, "ln", "", v...)
	return l
}

func (l *Logger) Infoln(v ...interface{}) MiniLogger {
	l.callWriter(InfoLevel, "ln", "", v...)
	return l
}

func (l *Logger) Debugln(v ...interface{}) MiniLogger {
	l.callWriter(DebugLevel, "ln", "", v...)
	return l
}
