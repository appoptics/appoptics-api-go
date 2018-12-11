package appoptics

import (
	"fmt"
	"net/http"
	"time"
)

type AnnotationStream struct {
	Name        string                       `json:"name"`
	DisplayName string                       `json:"display_name,omitempty"`
	Events      map[string][]AnnotationEvent `json:"events,omitempty"`
}

// AnnotationEvent is the main data structure for the Annotations API
type AnnotationEvent struct {
	Title       string           `json:"title"`
	Source      string           `json:"source,omitempty"`
	Description string           `json:"description,omitempty"`
	Links       []AnnotationLink `json:"links,omitempty"`
	StartTime   int64            `json:"start_time,omitempty"`
	EndTime     int64            `json:"end_time,omitempty"`
}

// AnnotationLink represents the Link metadata for on the AnnotationEvent
type AnnotationLink struct {
	Rel   string `json:"rel"`
	Href  string `json:"href"`
	Label string `json:"label,omitempty"`
}

type RetrieveAnnotationsRequest struct {
	Name      string
	StartTime time.Time
	EndTime   time.Time
	Sources   []string
}

type ListAnnotationsResponse struct {
	AnnotationStreams []*AnnotationStream `json:"annotations"`
	Query             QueryInfo           `json:"query"`
}

type AnnotationsCommunicator interface {
	List(string) (*ListAnnotationsResponse, error)
	Retrieve(*RetrieveAnnotationsRequest) (*AnnotationStream, error)
	RetrieveEvent(string, int) (*AnnotationEvent, error)
	Create(*AnnotationEvent, string) (*AnnotationEvent, error)
	Update(string, int, *AnnotationLink) (*AnnotationEvent, error)
}

type AnnotationsService struct {
	client *Client
}

func NewAnnotationsService(c *Client) *AnnotationsService {
	return &AnnotationsService{c}
}

// List retrieves all AnnotationEvents for the provided stream name
func (as *AnnotationsService) List(streamName string) (*ListAnnotationsResponse, error) {
	path := fmt.Sprintf("annotations?name=%s", streamName)
	var annotations *ListAnnotationsResponse
	req, err := as.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	_, err = as.client.Do(req, &annotations)

	if err != nil {
		return nil, err
	}

	return annotations, nil
}

// Retrieve fetches all AnnotationEvents matching the provided sources
func (as *AnnotationsService) Retrieve(retReq *RetrieveAnnotationsRequest) (*AnnotationStream, error) {
	stream := &AnnotationStream{}
	path := fmt.Sprintf("annotations/%s", retReq.Name)
	req, err := as.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = retReq.queryString(req)

	_, err = as.client.Do(req, stream)
	if err != nil {
		return nil, err
	}

	return stream, nil
}

func (retReq *RetrieveAnnotationsRequest) queryString(req *http.Request) string {
	var (
		stringStartTime string
		stringEndTime   string
	)
	q := req.URL.Query()

	if !retReq.StartTime.IsZero() {
		stringStartTime = fmt.Sprintf("%s", retReq.StartTime.Unix())
		q.Add("start_time", stringStartTime)
	}

	if !retReq.EndTime.IsZero() {
		stringEndTime = fmt.Sprintf("%s", retReq.EndTime.Unix())
		q.Add("end_time", stringEndTime)
	}

	if len(retReq.Sources) > 0 {
		for _, source := range retReq.Sources {
			q.Add("sources[]", source)
		}
	}
	return q.Encode()
}

// RetrieveEvent returns a single event identified by an integer ID from a given stream
func (as *AnnotationsService) RetrieveEvent(streamName string, id int) (*AnnotationEvent, error) {
	event := &AnnotationEvent{}
	path := fmt.Sprintf("annotations/%s/%d", streamName, id)
	req, err := as.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	_, err = as.client.Do(req, event)
	if err != nil {
		return nil, err
	}
	return event, nil
}

// Create makes an AnnotationEvent on the stream with the given name
func (as *AnnotationsService) Create(event *AnnotationEvent, streamName string) (*AnnotationEvent, error) {
	path := fmt.Sprintf("annotations/%s", streamName)
	req, err := as.client.NewRequest("POST", path, event)
	if err != nil {
		return nil, err
	}

	createdEvent := &AnnotationEvent{}

	_, err = as.client.Do(req, createdEvent)

	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Update adds a link to an annotation Event
func (as *AnnotationsService) Update(streamName string, id int, link *AnnotationLink) (*AnnotationEvent, error) {
	return nil, nil
}
