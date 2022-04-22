package main

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type Formatter struct {
	logrus.Formatter

	TimestampFormat string
	PrintColors     bool
	TrimMessages    bool
}

func (f Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	levelColor := getColorByLevel(entry.Level)

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = time.StampMilli
	}

	/* output buffer */
	b := &bytes.Buffer{}

	/* write timestamp */
	if f.PrintColors {
		b.WriteString(colorBlack)
	}

	b.WriteString("[")
	b.WriteString(entry.Time.Format(timestampFormat))
	b.WriteString("]")
	if f.PrintColors {
		b.WriteString(colorReset)
	}

	b.WriteString(" ")

	/* write log level */
	if f.PrintColors {
		_, _ = fmt.Fprint(b, levelColor)
	}

	b.WriteString(strings.ToUpper(entry.Level.String())[:4])
	b.WriteString(" ")
	if f.PrintColors {
		b.WriteString(colorReset)
	}

	/* write log message */
	if f.TrimMessages {
		b.WriteString(strings.TrimSpace(entry.Message))
	} else {
		b.WriteString(entry.Message)
	}

	b.WriteByte('\n')

	return b.Bytes(), nil
}

const (
	colorBlack  = "\u001b[90m"
	colorRed    = "\u001b[91m"
	colorGreen  = "\u001b[92m"
	colorYellow = "\u001b[93m"
	colorBlue   = "\u001b[94m"
	colorReset  = "\u001b[0m"
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
