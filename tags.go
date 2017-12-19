package appoptics

import (
	"bytes"
	"fmt"
	"strings"
)

func MetricWithTags(name string, tags map[string]interface{}) string {
	if tags == nil {
		return name
	}

	b := bytes.NewBufferString(name)

	for k, v := range tags {
		b.WriteString("::")
		b.WriteString(k)
		b.WriteString("::")
		b.WriteString(fmt.Sprint(v))
	}

	return b.String()
}

func parseMeasurementKey(key string) (string, map[string]string) {
	var (
		nameParts  = strings.Split(key, "::")
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
