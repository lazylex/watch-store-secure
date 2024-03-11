package colorlog

// Базируется на пакете https://pkg.go.dev/github.com/lmittmann/tint
// Опция NoColor не работает и уровень сообщения не выводится после внесённых мной изменений

import (
	"context"
	"encoding"
	"fmt"
	"io"
	"log/slog"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"time"
	"unicode"
)

// ANSI modes
const (
	ansiReset                      = "\033[0m"
	ansiTextBlack                  = "\033[30m"
	ansiBackgroundBlack            = "\033[40m"
	ansiTextDarkRed                = "\033[31m"
	ansiBackgroundDarkRed          = "\033[41m"
	ansiTextDarkGreen              = "\033[32m"
	ansiBackgroundDarkGreen        = "\033[42m"
	ansiTextDarkYellowOrange       = "\033[33m"
	ansiBackgroundDarkYellowOrange = "\033[43m"
	ansiTextDarkBlue               = "\033[34m"
	ansiBackgroundDarkBlue         = "\033[44m"
	ansiTextDarkPurple             = "\033[35m"
	ansiBackgroundDarkPurple       = "\033[45m"
	ansiTextDarkCyan               = "\033[36m"
	ansiBackgroundDarkCyan         = "\033[46m"
	ansiTextLightGray              = "\033[37m"
	ansiBackgroundLightGray        = "\033[47m"
	ansiTextDarkGray               = "\033[90m"
	ansiBackgroundDarkGray         = "\033[100m"
	ansiTextRed                    = "\033[91m"
	ansiBackgroundRed              = "\033[101m"
	ansiTextGreen                  = "\033[92m"
	ansiBackgroundGreen            = "\033[102m"
	ansiTextOrange                 = "\033[93m"
	ansiBackgroundOrange           = "\033[103m"
	ansiTextBlue                   = "\033[94m"
	ansiBackgroundBlue             = "\033[104m"
	ansiTextPurple                 = "\033[95m"
	ansiBackgroundPurple           = "\033[105m"
	ansiTextCyan                   = "\033[96m"
	ansiBackgroundCyan             = "\033[106m"
	ansiTextWhite                  = "\033[97m"
	ansiBackgroundWhite            = "\033[107m"
)

var (
	currentTextColor = ansiTextWhite
	currentBackColor = ansiBackgroundBlack
)

// getTextColor возвращает текущий цвет для текста
func getTextColor(inverted bool) string {
	if inverted {
		return ansiTextBlack
	}
	return ansiTextWhite
}

// getBackColor возвращает текущий цвет для фона
func getBackColor(inverted bool) string {
	if inverted {
		return ansiBackgroundDarkGray
	}
	return ansiBackgroundBlack
}

const errKey = "err"

var (
	defaultLevel      = slog.LevelInfo
	defaultTimeFormat = time.StampMilli
)

// Options for a slog.Handler that writes tinted logs. A zero Options consists
// entirely of default values.
//
// Options can be used as a drop-in replacement for [slog.HandlerOptions].
type Options struct {
	// Enable source code location (Default: false)
	AddSource bool

	// Minimum level to log (Default: slog.LevelInfo)
	Level slog.Leveler

	// ReplaceAttr is called to rewrite each non-group attribute before it is logged.
	// See https://pkg.go.dev/log/slog#HandlerOptions for details.
	ReplaceAttr func(groups []string, attr slog.Attr) slog.Attr

	// Time format (Default: time.StampMilli)
	TimeFormat string

	// Disable color (Default: false)
	NoColor bool
}

// NewHandler creates a [slog.Handler] that writes tinted logs to Writer w,
// using the default options. If opts is nil, the default options are used.
func NewHandler(w io.Writer, opts *Options) slog.Handler {
	h := &handler{
		w:          w,
		level:      defaultLevel,
		timeFormat: defaultTimeFormat,
	}
	if opts == nil {
		return h
	}

	h.addSource = opts.AddSource
	if opts.Level != nil {
		h.level = opts.Level
	}
	h.replaceAttr = opts.ReplaceAttr
	if opts.TimeFormat != "" {
		h.timeFormat = opts.TimeFormat
	}
	h.noColor = opts.NoColor
	return h
}

// handler implements a [slog.Handler].
type handler struct {
	attrsPrefix string
	groupPrefix string
	groups      []string

	mu sync.Mutex
	w  io.Writer

	addSource   bool
	level       slog.Leveler
	replaceAttr func([]string, slog.Attr) slog.Attr
	timeFormat  string
	noColor     bool
}

func (h *handler) clone() *handler {
	return &handler{
		attrsPrefix: h.attrsPrefix,
		groupPrefix: h.groupPrefix,
		groups:      h.groups,
		w:           h.w,
		addSource:   h.addSource,
		level:       h.level,
		replaceAttr: h.replaceAttr,
		timeFormat:  h.timeFormat,
		noColor:     h.noColor,
	}
}

func (h *handler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level.Level()
}

func (h *handler) Handle(_ context.Context, r slog.Record) error {
	// get a buffer from the sync pool
	buf := newBuffer()
	defer buf.Free()

	rep := h.replaceAttr

	// write time

	if !r.Time.IsZero() {
		val := r.Time.Round(0) // strip monotonic to match Attr behavior
		if rep == nil {
			h.appendTime(buf, r.Time)
			buf.WriteByte(' ')
		} else if a := rep(nil /* groups */, slog.Time(slog.TimeKey, val)); a.Key != "" {
			if a.Value.Kind() == slog.KindTime {
				h.appendTime(buf, a.Value.Time())
			} else {
				h.appendValue(buf, a.Value, false)
			}
			buf.WriteByte(' ')
		}
	}

	switch r.Level {
	case slog.LevelError:
		buf.WriteString(ansiTextWhite)
		buf.WriteString(ansiBackgroundDarkRed)
	case slog.LevelInfo:
		buf.WriteString(ansiTextOrange)
		buf.WriteString(ansiBackgroundDarkGreen)
	case slog.LevelDebug:
		buf.WriteString(ansiTextWhite)
		buf.WriteString(ansiBackgroundCyan)
	case slog.LevelWarn:
		buf.WriteString(ansiTextBlack)
		buf.WriteString(ansiBackgroundOrange)
	}

	// write level
	if rep == nil {
		//h.appendLevel(buf, r.Level)
		buf.WriteByte(' ')
	} else if a := rep(nil /* groups */, slog.Any(slog.LevelKey, r.Level)); a.Key != "" {
		h.appendValue(buf, a.Value, false)
		buf.WriteByte(' ')
	}

	// write source
	if h.addSource {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()
		if f.File != "" {
			src := &slog.Source{
				Function: f.Function,
				File:     f.File,
				Line:     f.Line,
			}

			if rep == nil {
				h.appendSource(buf, src)
				buf.WriteByte(' ')
			} else if a := rep(nil /* groups */, slog.Any(slog.SourceKey, src)); a.Key != "" {
				h.appendValue(buf, a.Value, false)
				buf.WriteByte(' ')
			}
		}
	}

	// write message
	if rep == nil {
		buf.WriteString(fmt.Sprintf("💬 %s ", r.Message))

	} else if a := rep(nil /* groups */, slog.String(slog.MessageKey, r.Message)); a.Key != "" {
		h.appendValue(buf, a.Value, false)
		buf.WriteByte(' ')
	}

	// write handler attributes
	if len(h.attrsPrefix) > 0 {
		buf.WriteString(h.attrsPrefix)
	}

	// write attributes
	r.Attrs(func(attr slog.Attr) bool {
		h.appendAttr(buf, attr, h.groupPrefix, h.groups)
		return true
	})

	if len(buf.Text) == 0 {
		return nil
	}
	buf.WriteString(ansiReset + " ")
	(buf.Text)[len(buf.Text)-1] = '\n' // replace last space with newline

	h.mu.Lock()
	defer h.mu.Unlock()

	_, err := h.w.Write(buf.Text)
	return err
}

func (h *handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}
	h2 := h.clone()

	buf := newBuffer()
	defer buf.Free()

	// write attributes to buffer
	for _, attr := range attrs {
		h.appendAttr(buf, attr, h.groupPrefix, h.groups)
	}
	h2.attrsPrefix = h.attrsPrefix + string(buf.Text)
	return h2
}

func (h *handler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	h2 := h.clone()
	h2.groupPrefix += name + "."
	h2.groups = append(h2.groups, name)
	return h2
}

func (h *handler) appendTime(buf *buffer, t time.Time) {
	buf.Text = t.AppendFormat(buf.Text, h.timeFormat)
}

func (h *handler) appendLevel(buf *buffer, level slog.Level) {
	switch {
	case level < slog.LevelInfo:
		buf.WriteString("DBG")
		appendLevelDelta(buf, level-slog.LevelDebug)
	case level < slog.LevelWarn:
		buf.WriteString("INF")
		appendLevelDelta(buf, level-slog.LevelInfo)
	case level < slog.LevelError:
		buf.WriteString("WRN")
		appendLevelDelta(buf, level-slog.LevelWarn)
	default:
		buf.WriteString("ERR")
		appendLevelDelta(buf, level-slog.LevelError)
	}
}

func appendLevelDelta(buf *buffer, delta slog.Level) {
	if delta == 0 {
		return
	} else if delta > 0 {
		buf.WriteByte('+')
	}
	buf.Text = strconv.AppendInt(buf.Text, int64(delta), 10)
}

func (h *handler) appendSource(buf *buffer, src *slog.Source) {
	dir, file := filepath.Split(src.File)

	buf.WriteStringIf(!h.noColor, ansiBackgroundDarkGray)
	buf.WriteString(filepath.Join(filepath.Base(dir), file))
	buf.WriteByte(':')
	buf.WriteString(strconv.Itoa(src.Line))

}

func (h *handler) appendAttr(buf *buffer, attr slog.Attr, groupsPrefix string, groups []string) {
	attr.Value = attr.Value.Resolve()
	if rep := h.replaceAttr; rep != nil && attr.Value.Kind() != slog.KindGroup {
		attr = rep(groups, attr)
		attr.Value = attr.Value.Resolve()
	}

	if attr.Equal(slog.Attr{}) {
		return
	}

	if attr.Value.Kind() == slog.KindGroup {
		if attr.Key != "" {
			groupsPrefix += attr.Key + "."
			groups = append(groups, attr.Key)
		}
		for _, groupAttr := range attr.Value.Group() {
			h.appendAttr(buf, groupAttr, groupsPrefix, groups)
		}
	} else if err, ok := attr.Value.Any().(tintError); ok {
		// append tintError
		h.appendTintError(buf, err, groupsPrefix)
		buf.WriteByte(' ')
	} else {
		buf.WriteString(getBackColor(buf.inverted))
		h.appendKey(buf, attr.Key, groupsPrefix)
		h.appendValue(buf, attr.Value, true)
		buf.WriteByte(' ')
	}
}

func (h *handler) appendKey(buf *buffer, key, groups string) {
	if key == "op" {
		appendString(buf, " 👀 ", false)
	} else {
		buf.WriteString(ansiTextOrange)
		appendString(buf, groups+key, true)
		buf.WriteByte('=')
	}
}

func (h *handler) appendValue(buf *buffer, v slog.Value, quote bool) {
	buf.WriteString(getTextColor(buf.Inverse()))
	switch v.Kind() {
	case slog.KindString:
		appendString(buf, v.String(), quote)
	case slog.KindInt64:
		buf.Text = strconv.AppendInt(buf.Text, v.Int64(), 10)
	case slog.KindUint64:
		buf.Text = strconv.AppendUint(buf.Text, v.Uint64(), 10)
	case slog.KindFloat64:
		buf.Text = strconv.AppendFloat(buf.Text, v.Float64(), 'g', -1, 64)
	case slog.KindBool:
		buf.Text = strconv.AppendBool(buf.Text, v.Bool())
	case slog.KindDuration:
		appendString(buf, v.Duration().String(), quote)
	case slog.KindTime:
		appendString(buf, v.Time().String(), quote)
	case slog.KindAny:
		switch cv := v.Any().(type) {
		case slog.Level:
			h.appendLevel(buf, cv)
		case encoding.TextMarshaler:
			data, err := cv.MarshalText()
			if err != nil {
				break
			}
			appendString(buf, string(data), quote)
		case *slog.Source:
			h.appendSource(buf, cv)
		default:
			appendString(buf, fmt.Sprintf("%+v", v.Any()), quote)
		}
	}
}

func (h *handler) appendTintError(buf *buffer, err error, groupsPrefix string) {
	appendString(buf, groupsPrefix+errKey, true)
	buf.WriteByte('=')
	appendString(buf, err.Error(), true)
}

func appendString(buf *buffer, s string, quote bool) {
	if quote && needsQuoting(s) {
		buf.Text = strconv.AppendQuote(buf.Text, s)
	} else {
		buf.WriteString(s)
	}
}

func needsQuoting(s string) bool {
	if len(s) == 0 {
		return true
	}
	for _, r := range s {
		if unicode.IsSpace(r) || r == '"' || r == '=' || !unicode.IsPrint(r) {
			return true
		}
	}
	return false
}

type tintError struct{ error }

// Err returns a tinted (colorized) [slog.Attr] that will be written in red color
// by the [tint.Handler]. When used with any other [slog.Handler], it behaves as
//
//	slog.Any("err", err)
func Err(err error) slog.Attr {
	if err != nil {
		err = tintError{err}
	}
	return slog.Any(errKey, err)
}
