package appoptics

import "fmt"

// Metric represents a Librato Metric.
type Metric struct {
	Name        string           `json:"name"`
	Description string           `json:"description,omitempty"`
	Type        string           `json:"type"`
	Period      int              `json:"period,omitempty"`
	DisplayName string           `json:"display_name,omitempty"`
	Composite   string           `json:"composite,omitempty"`
	Attributes  MetricAttributes `json:"attributes,omitempty"`
}

// MetricAttributes are named attributes as key:value pairs.
type MetricAttributes struct {
	Color             *string     `json:"color,omitempty"`
	DisplayMax        interface{} `json:"display_max,omitempty"`
	DisplayMin        interface{} `json:"display_min,omitempty"`
	DisplayUnitsLong  string      `json:"display_units_long,omitempty"`
	DisplayUnitsShort string      `json:"display_units_short,omitempty"`
	DisplayStacked    bool        `json:"display_stacked,omitempty"`
	CreatedByUA       string      `json:"created_by_ua,omitempty"`
	GapDetection      bool        `json:"gap_detection,omitempty"`
	Aggregate         bool        `json:"aggregate,omitempty"`
}

type MetricsResponse struct {
	Query   QueryInfo `json:"query,omitempty"`
	Metrics []*Metric `json:"metrics,omitempty"`
}

type MetricsService struct {
	client *Client
}

type MetricsCommunicator interface {
	List() (*MetricsResponse, error)
	Retrieve(string) (*Metric, error)
	Create(*Metric) (*Metric, error)
	Update(*Metric) (*Metric, error)
	Delete(int) error
}

func NewMetricsService(c *Client) *MetricsService {
	return &MetricsService{c}
}

func (ms *MetricsService) List() (*MetricsResponse, error) {
	req, err := ms.client.NewRequest("GET", "metrics", nil)
	if err != nil {
		return nil, err
	}

	metricsResponse := &MetricsResponse{}

	_, err = ms.client.Do(req, &metricsResponse)

	if err != nil {
		return nil, err
	}

	return metricsResponse, nil
}

func (ms *MetricsService) Retrieve(name string) (*Metric, error) {
	metric := &Metric{}
	path := fmt.Sprintf("metrics/%s", name)
	req, err := ms.client.NewRequest("GET", path, nil)

	if err != nil {
		return nil, err
	}

	_, err = ms.client.Do(req, metric)
	if err != nil {
		return nil, err
	}

	return metric, nil
}

func (ms *MetricsService) Create(m *Metric) (*Metric, error) {
	path := fmt.Sprintf("metrics/%s", m.Name)
	req, err := ms.client.NewRequest("PUT", path, m)
	if err != nil {
		return nil, err
	}

	createdMetric := &Metric{}

	_, err = ms.client.Do(req, createdMetric)
	if err != nil {
		return nil, err
	}

	return createdMetric, nil
}

func (ms *MetricsService) Update(m *Metric) (*MetricsResponse, error) {
	return nil, nil
}

func (ms *MetricsService) Delete(name string) error {
	path := fmt.Sprintf("metrics/%s", name)
	req, err := ms.client.NewRequest("DELETE", path, nil)
	if err != nil {
		return err
	}

	_, err = ms.client.Do(req, nil)

	return err
}
