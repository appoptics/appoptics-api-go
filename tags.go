package appoptics

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
)

// Tag represents an AppOptics Tag, used in Measurements, Charts, etc
type Tag struct {
	Name    string   `json:"name,omitempty"`
	Values  []string `json:"values,omitempty"`
	Grouped bool     `json:"grouped,omitempty"`
	Dynamic bool     `json:"dynamic,omitempty"`
}

// MetricTagSeparator is used by MetricWithTags as a separator when serializing
// the metric name and tags as a key for aggregation with measurements matching
// the same metric name and tags.
//
// Users can build these strings internally if desired, using a format like below:
// - "metric_name"
// - "metric_name\x00tag_1_key\x00tag_1_value"
// - "metric_name\x00tag_1_key\x00tag_1_value\x00tag_2_key\x00tag_2_value" ...
const MetricTagSeparator = "\x00"

func MetricWithTags(name string, tags map[string]interface{}) string {
	if tags == nil {
		return name
	}
	b := bytes.NewBufferString(name)
	keys := make([]string, 0, len(tags))
	for k := range tags {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := tags[k]
		b.WriteString(MetricTagSeparator)
		b.WriteString(k)
		b.WriteString(MetricTagSeparator)
		b.WriteString(fmt.Sprint(v))
	}

	return b.String()
}

func parseMeasurementKey(key string) (string, map[string]string) {
	var (
		nameParts  = strings.Split(key, MetricTagSeparator)
		metricName = nameParts[0]
	)

	if len(nameParts) < 3 {
		return metricName, nil
	}

	tags := make(map[string]string)
	for n := 1; n < len(nameParts); n += 2 {
		tags[nameParts[n]] = nameParts[n+1]
	}

	return metricName, tags
}
