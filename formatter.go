package main

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type Formatter struct {
	TimestampFormat string
	NoColors        bool
	TrimMessages    bool
}

func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	levelColor := getColorByLevel(entry.Level)

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = time.StampMilli
	}

	// output buffer
	b := &bytes.Buffer{}

	// write time
	b.WriteString(colorBlack)
	b.WriteString("[")
	b.WriteString(entry.Time.Format(timestampFormat))
	b.WriteString("]")
	b.WriteString(colorReset)

	b.WriteString(" ")

	// write level
	if !f.NoColors {
		_, _ = fmt.Fprint(b, levelColor)
	}

	b.WriteString(strings.ToUpper(entry.Level.String())[:4])
	b.WriteString(": ")

	// write caller
	if entry.HasCaller() {
		if !f.NoColors {
			b.WriteString(colorMagenta)
		}

		_, _ = fmt.Fprintf(b, "<%s> ", entry.Caller.Function)
		if !f.NoColors {
			b.WriteString(colorReset)
		}
	}

	// write message
	if f.TrimMessages {
		b.WriteString(strings.TrimSpace(entry.Message))
	} else {
		b.WriteString(entry.Message)
	}

	b.WriteByte('\n')

	return b.Bytes(), nil
}

const (
	colorBlack   = "\u001b[90m"
	colorRed     = "\u001b[91m"
	colorGreen   = "\u001b[92m"
	colorYellow  = "\u001b[93m"
	colorBlue    = "\u001b[94m"
	colorMagenta = "\u001b[95m"
	colorReset   = "\u001b[0m"
)

func getColorByLevel(level logrus.Level) string {
	switch level {
	case logrus.DebugLevel, logrus.TraceLevel:
		return colorBlue
	case logrus.WarnLevel:
		return colorYellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		return colorRed
	default:
		return colorGreen
	}
}
