package types

import (
	"time"
)

// Metric holds data of one counter metric
type Metric struct {
	Name  string      `json:"name"`  // Name of the metric
	Value interface{} `json:"value"` // Value of the metric
	Type  string      `json:"type"`  // Type of the metric
	Time  time.Time   `json:"time"`  // Update time
}

const METRIC_TYPE_COUNT = "count"
const METRIC_TYPE_MAP = "map"

type StatFilter struct {
	From  int64
	Until int64
}
