// Package log is a singleton package for logging either remotely or
// locally. If provided a host, it will connect to a service utilizing the
// syslogsvc protocol for tinyCI, and will alloow it to send log transmissions
// there.
package log

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	transport "github.com/erikh/go-transport"
	_struct "github.com/golang/protobuf/ptypes/struct"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/sirupsen/logrus"
	logsvc "github.com/tinyci/ci-agents/api/logsvc/processors"
	"github.com/tinyci/ci-agents/errors"
	"github.com/tinyci/ci-agents/grpc/services/log"
	"github.com/tinyci/ci-agents/model"
	"google.golang.org/grpc"
)

// RemoteClient is the swagger-based syslogsvc client.
var RemoteClient log.LogClient

// FieldMap is just a type alias for map[string]string to keep me from
// breaking my fingers.
type FieldMap map[string]string

// Fields is just shorthand for some protobuf stuff.
type Fields _struct.Struct

func toValue(str string) *_struct.Value {
	return &_struct.Value{Kind: &_struct.Value_StringValue{StringValue: str}}
}

// NewFields creates a compatible *Fields
func NewFields() *Fields {
	return &Fields{Fields: map[string]*_struct.Value{}}
}

// ToFields converts a FieldMap to a *Fields
func (f FieldMap) ToFields() *Fields {
	fields := map[string]*_struct.Value{}

	for k, v := range f {
		fields[k] = toValue(v)
	}

	return &Fields{Fields: fields}
}

// ToLogrus just casts stuff to make logrus happy
func (f *Fields) ToLogrus() map[string]interface{} {
	m := map[string]interface{}{}
	for key, value := range f.Fields {
		s, ok := value.Kind.(*_struct.Value_StringValue)
		if ok {
			m[key] = s.StringValue
		} else {
			m[key] = fmt.Sprintf("%v", value.Kind)
		}
	}

	return m
}

// Set sets an item in the fields.
func (f *Fields) Set(key, value string) {
}

// ConfigureRemote configures the remote endpoint with a provided URL.
func ConfigureRemote(addr string, cert *transport.Cert) error {
	client, err := transport.GRPCDial(cert, addr)
	if err != nil {
		return err
	}

	RemoteClient = log.NewLogClient(client)
	return nil
}

// SubLogger is a handle to cached parameters for logging.
type SubLogger struct {
	Service string
	Fields  *Fields
}

// New creates a new SubLogger which can be primed with cached values for each log entry.
func New() *SubLogger {
	return &SubLogger{Service: "n/a", Fields: &Fields{Fields: map[string]*_struct.Value{}}}
}

// NewWithData returns a SubLogger already primed with cached data.
func NewWithData(svc string, params FieldMap) *SubLogger {
	if params == nil {
		params = FieldMap{}
	}

	p := params.ToFields()

	p.Fields["service"] = toValue(svc)
	return &SubLogger{svc, p}
}

// WithService is a SubLogger version of package-level WithService. They call the same code.
func (sub *SubLogger) WithService(svc string) *SubLogger {
	sub2 := *sub
	sub2.Service = svc
	sub2.Fields = NewFields()

	if sub.Fields != nil {
		params := sub.Fields

		for k, v := range params.Fields {
			sub2.Fields.Fields[k] = toValue(v.String())
		}
	}

	sub2.Fields.Fields["service"] = toValue(svc)
	return &sub2
}

// WithFields is a SubLogger version of package-level WithFields. They call the same code.
func (sub *SubLogger) WithFields(params FieldMap) *SubLogger {
	sub2 := *sub

	sub2.Fields = NewFields()

	for k, v := range sub.Fields.Fields {
		sub2.Fields.Fields[k] = v
	}

	for k, v := range params.ToFields().Fields {
		sub2.Fields.Fields[k] = v
	}

	return &sub2
}

// WithRequest is a wrapper for WithFields() that handles *http.Request data.
func (sub *SubLogger) WithRequest(req *http.Request) *SubLogger {
	raddr := req.Header.Get("X-Forwarded-For")
	if raddr == "" {
		raddr = strings.Split(req.RemoteAddr, ":")[0]
	} else {
		raddr = strings.TrimSpace(strings.SplitN(raddr, ",", 2)[0])
	}

	fm := FieldMap{
		"remote_addr":    raddr,
		"request_method": req.Method,
		"request_url":    req.URL.String(),
	}

	return sub.WithFields(fm)
}

// WithUser includes user information
func (sub *SubLogger) WithUser(user *model.User) *SubLogger {
	fm := FieldMap{
		"username": user.Username,
		"user_id":  fmt.Sprintf("%v", user.ID),
	}

	return sub.WithFields(fm)
}

func (sub *SubLogger) makeMsg(level, msg string, values []interface{}) *log.LogMessage {
	if values != nil {
		msg = fmt.Sprintf(msg, values...)
	}

	now := time.Now()
	ts := timestamp.Timestamp{}
	ts.Seconds = now.Unix()
	ts.Nanos = int32(now.Nanosecond())

	return &log.LogMessage{
		At:      &ts,
		Fields:  (*_struct.Struct)(sub.Fields),
		Message: msg,
		Level:   level,
		Service: sub.Service,
	}
}

// Logf logs a thing with formats!
func (sub *SubLogger) Logf(level string, msg string, values []interface{}, localLog func(string, ...interface{})) *errors.Error {
	if RemoteClient != nil {
		_, err := RemoteClient.Put(context.Background(), sub.makeMsg(level, msg, values), grpc.WaitForReady(true))
		return errors.New(err)
	}

	localLog(msg, values...)
	return nil
}

// Log logs a thing
func (sub *SubLogger) Log(level string, msg interface{}, localLog func(...interface{})) *errors.Error {
	if RemoteClient != nil {
		_, err := RemoteClient.Put(context.Background(), sub.makeMsg(level, fmt.Sprintf("%v", msg), nil), grpc.WaitForReady(true))
		return errors.New(err)
	}

	switch msg := msg.(type) {
	case *errors.Error:
		if msg.Log {
			localLog(msg)
		}
	default:
		localLog(msg)
	}

	return nil
}

// Info prints an info message
func (sub *SubLogger) Info(msg interface{}) error {
	sub.Log(logsvc.LevelInfo, msg, logrus.WithFields(sub.Fields.ToLogrus()).Info)
	return nil
}

// Infof is the format-capable version of Info
func (sub *SubLogger) Infof(msg string, values ...interface{}) error {
	sub.Logf(logsvc.LevelInfo, msg, values, logrus.WithFields(sub.Fields.ToLogrus()).Infof)
	return nil
}

// Error prints an error message
func (sub *SubLogger) Error(msg interface{}) error {
	sub.Log(logsvc.LevelError, msg, logrus.WithFields(sub.Fields.ToLogrus()).Error)
	return nil
}

// Errorf is the format-capable version of Error
func (sub *SubLogger) Errorf(msg string, values ...interface{}) error {
	sub.Logf(logsvc.LevelError, msg, values, logrus.WithFields(sub.Fields.ToLogrus()).Errorf)
	return nil
}

// Debug prints a debug message
func (sub *SubLogger) Debug(msg interface{}) error {
	sub.Log(logsvc.LevelDebug, msg, logrus.WithFields(sub.Fields.ToLogrus()).Debug)
	return nil
}

// Debugf is the format-capable version of Debug
func (sub *SubLogger) Debugf(msg string, values ...interface{}) error {
	sub.Logf(logsvc.LevelDebug, msg, values, logrus.WithFields(sub.Fields.ToLogrus()).Debugf)
	return nil
}
