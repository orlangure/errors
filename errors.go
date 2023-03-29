package errors

import (
	goerrors "errors"
	"fmt"
	"runtime"
	"strings"
)

type builder struct {
	fields map[string]any
}

func WithField(key string, value any) (b *builder) { return b.WithField(key, value) }
func (b *builder) WithField(key string, value any) *builder {
	if b == nil {
		b = &builder{}
	}

	if b.fields == nil {
		b.fields = make(map[string]any)
	}

	b.fields[key] = value
	if _, ok := b.fields["codePos"]; ok {
		return b
	}

	pc, file, line, ok := runtime.Caller(1)
	if ok {
		fn := runtime.FuncForPC(pc)
		if strings.HasSuffix(fn.Name(), ".WithField") {
			pc, file, line, ok = runtime.Caller(2)
		}

		if ok {
			file = fmt.Sprintf("%s:%d", shortenFilePath(file), line)
			b.fields["codePos"] = fmt.Sprintf("%s/%s:%d", packageName(pc), file, line)
		}
	}

	return b
}
func packageName(pc uintptr) string {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "unknown"
	}

	name := fn.Name()
	index := strings.LastIndex(name, ".")
	if index == -1 {
		return "unknown"
	}

	return name[:index]
}
func shortenFilePath(file string) string {
	if idx := strings.LastIndex(file, "/"); idx >= 0 {
		return file[idx+1:]
	}
	return file
}

func (b *builder) Wrap(err error, message string) error {
	return &MyError{
		cause:      err,
		message:    message,
		properties: b.fields,
	}
}

func (b *builder) New(message string) error {
	return &MyError{
		message:    message,
		properties: b.fields,
	}
}

type MyError struct {
	message    string
	properties map[string]any
	cause      error
}

func (e *MyError) Error() string {
	msg := e.message
	if e.cause != nil {
		msg += fmt.Sprintf("; %v", e.cause)
	}
	return msg
}

func (e *MyError) Unwrap() error {
	return e.cause
}

func Trace(e error) (fields map[string]interface{}, stackTrace []string) {
	fields = make(map[string]interface{})
	stackTrace = make([]string, 0)

	if e == nil {
		return
	}

	err := e
	for err != nil {
		var targetErr *MyError
		if goerrors.As(err, &targetErr) {
			for k, v := range targetErr.properties {
				if k != "codePos" {
					fields[k] = v
				}
			}
			if codePos, ok := targetErr.properties["codePos"]; ok {
				stackTrace = append(stackTrace, fmt.Sprintf("%v", codePos))
			}
		}
		err = goerrors.Unwrap(err)
	}

	return fields, stackTrace
}
