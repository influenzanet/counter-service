package types

import (
	"encoding/json"
	"time"
)

// Metric holds data of one counter metric
type Metric struct {
	Name  string      `json:"name"`  // Name of the metric
	Value interface{} `json:"value"` // Value of the metric
	Type  string      `json:"type"`  // Type of the metric
	Time  Timestamp   `json:"time"`  // Update time
}

const METRIC_TYPE_COUNT = "count"
const METRIC_TYPE_MAP = "map"

type StatFilter struct {
	From  int64
	Until int64
}

// Timestamp is a time serialized as timestamp
type Timestamp struct {
	time.Time
}

func (d *Timestamp) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Unix())
}
