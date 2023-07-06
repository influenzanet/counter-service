package types
import(
	"time"
)
// Counter metric
type Counter struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
	Type  string      `json:"type"`
	Time  time.Time   `json:"time"`
}

const COUNTER_TYPE_COUNT = "count"
const COUNTER_TYPE_MAP = "map"
