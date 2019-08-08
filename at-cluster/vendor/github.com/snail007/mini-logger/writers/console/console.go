package console

import (
	"fmt"
	"os"
	"strings"
	"time"

	"encoding/json"

	"github.com/fatih/color"
	"github.com/snail007/mini-logger"
)

const (
	T_JSON = iota
	T_TEXT
)

type ConsoleWriterConfig struct {
	Format string
	Type   int
}
type ConsoleWriter struct {
	c ConsoleWriterConfig
}

func New(config ConsoleWriterConfig) *ConsoleWriter {
	return &ConsoleWriter{
		c: config,
	}
}
func NewDefault() *ConsoleWriter {
	return &ConsoleWriter{
		c: ConsoleWriterConfig{
			Format: "[{level}] [{date} {time}.{mili}] {fields} {text}",
			Type:   T_TEXT,
		},
	}
}
func (w *ConsoleWriter) Init() (err error) {
	return
}
func (w *ConsoleWriter) Write(e logger.Entry) {
	fg := color.New(color.FgHiWhite).SprintFunc()
	switch e.Level {
	case logger.InfoLevel:
		fg = color.New(color.FgHiGreen).SprintFunc()
	case logger.WarnLevel:
		fg = color.New(color.FgHiYellow).SprintFunc()
	case logger.ErrorLevel:
		fg = color.New(color.FgHiRed).SprintFunc()
	case logger.FatalLevel:
		fg = color.New(color.FgHiRed, color.FgHiMagenta).SprintFunc()
	}
	date := time.Unix(e.Timestamp, 0).Format("2006/01/02")
	time := time.Unix(e.Timestamp, 0).Format("15:04:05")
	var c string
	if w.c.Type == T_TEXT {
		c = strings.Replace(w.c.Format, "{level}", fmt.Sprintf("%-5s", e.LevelString), -1)
		c = strings.Replace(c, "{date}", date, -1)
		c = strings.Replace(c, "{time}", time, -1)
		c = strings.Replace(c, "{fields}", e.FieldsString, -1)
		c = strings.Replace(c, "{mili}", fmt.Sprintf("%03d", e.Milliseconds), -1)
		c = strings.Replace(c, "{text}", e.Content, -1)
	} else if w.c.Type == T_JSON {
		e.LevelString = strings.TrimRight(e.LevelString, " ")
		e.FieldsString = ""
		v, _ := json.Marshal(e)
		c = string(v)
	} else {
		return
	}
	fmt.Fprintln(os.Stdout, fg(c))
}
