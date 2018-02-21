package wrapped

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type wrapper map[string]logrus.Fields

const defaultTimestampFormat = time.RFC3339

const (
	fieldKeyMsg   = "msg"
	fieldKeyLevel = "level"
	fieldKeyTime  = "time"
	defaultKind   = "log"
)

type WrappedJSONFormatter struct {
	// TimestampFormat sets the format used for marshaling timestamps.
	TimestampFormat string

	// DisableTimestamp allows disabling automatic timestamps in output
	DisableTimestamp bool

	// FieldMap allows users to customize the names of keys for various fields.
	// As an example:
	// formatter := &JSONFormatter{
	//   	FieldMap: FieldMap{
	// 		 FieldKeyTime: "@timestamp",
	// 		 FieldKeyLevel: "@level",
	// 		 FieldKeyLevel: "@message",
	//    },
	// }
	FieldMap FieldMap
}

type fieldKey string

// FieldMap allows customization of the key names for default fields.
type FieldMap map[fieldKey]string

func (f *WrappedJSONFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	wrp := wrapper{}
	var wrapKey string
	var preserveKind bool
	kind, ok := entry.Data["kind"]
	if ok {
		wrapKey, ok = kind.(string)
		if !ok {
			wrapKey = defaultKind
		} else {
			preserveKind = true
		}
		delete(entry.Data, "kind")
	} else {
		wrapKey = defaultKind
	}
	data := make(logrus.Fields, len(entry.Data)+3)
	for k, v := range entry.Data {
		switch v := v.(type) {
		case error:
			// Otherwise errors are ignored by `encoding/json`
			// https://github.com/sirupsen/logrus/issues/137
			data[k] = v.Error()
		default:
			data[k] = v
		}
	}
	prefixFieldClashes(data)

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = defaultTimestampFormat
	}

	if !f.DisableTimestamp {
		data[f.FieldMap.resolve(fieldKeyTime)] = entry.Time.Format(timestampFormat)
	}
	data[f.FieldMap.resolve(fieldKeyMsg)] = entry.Message
	data[f.FieldMap.resolve(fieldKeyLevel)] = entry.Level.String()
	wrp[wrapKey] = data
	serialized, err := json.Marshal(wrp)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal fields to JSON, %v", err)
	}
	if preserveKind {
		entry.Data["kind"] = wrapKey
	}
	return append(serialized, '\n'), nil
}

func prefixFieldClashes(data logrus.Fields) {
	if t, ok := data["time"]; ok {
		data["fields.time"] = t
	}

	if m, ok := data["msg"]; ok {
		data["fields.msg"] = m
	}

	if l, ok := data["level"]; ok {
		data["fields.level"] = l
	}
}

func (f FieldMap) resolve(key fieldKey) string {
	if k, ok := f[key]; ok {
		return k
	}

	return string(key)
}
