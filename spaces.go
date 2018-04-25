package appoptics

import (
	"fmt"
)

// ListSpacesResponse represents the returned data payload from Spaces API's List command (/spaces)
type ListSpacesResponse struct {
	Query  map[string]int `json:"query"`
	Spaces []*Space       `json:"spaces"`
}

// RetrieveSpaceResponse represents the returned data payload from Spaces API's Retrieve command (/spaces/:id)
type RetrieveSpaceResponse struct {
	Space
	Charts []map[string]int `json:"charts","omitempty"`
}

// Space represents a single AppOptics Space
type Space struct {
	// ID is the unique identifier of the Space
	ID int `json:"id","omitempty"`
	// Name is the name of the Space
	Name string `json:"name","omitempty"`
}

// SpacesCommunicator defines the interface for the Spaces API
type SpacesCommunicator interface {
	Create(string) (*Space, error)
	List(*RequestParameters) ([]*Space, error)
	Retrieve(int) (*RetrieveSpaceResponse, error)
	Update(int, string) (*Space, error)
	Delete(int) error
}

type SpacesService struct {
	client *Client
}

func NewSpacesService(c *Client) *SpacesService {
	return &SpacesService{c}
}

// Create creates the space with the given name
func (s *SpacesService) Create(name string) (*Space, error) {
	requestedSpace := &Space{Name: name}
	createdSpace := &Space{}
	req, err := s.client.NewRequest("POST", "spaces", requestedSpace)
	if err != nil {
		return nil, err
	}

	_, err = s.client.Do(req, createdSpace)
	if err != nil {
		return nil, err
	}
	return createdSpace, nil
}

// List implements the  Spaces API's List command
func (s *SpacesService) List(rp *RequestParameters) ([]*Space, error) {
	var spaces []*Space
	req, err := s.client.NewRequest("GET", "spaces", nil)

	if err != nil {
		return spaces, err
	}

	if rp != nil {
		rp.AddToRequest(req)
	}

	var spacesResponse ListSpacesResponse
	_, err = s.client.Do(req, &spacesResponse)

	if err != nil {
		return spaces, err
	}

	spaces = spacesResponse.Spaces
	return spaces, nil
}

// Retrieve implements the Spaces API's Retrieve command
func (s *SpacesService) Retrieve(id int) (*RetrieveSpaceResponse, error) {
	retrievedSpace := &RetrieveSpaceResponse{}
	spacePath := fmt.Sprintf("spaces/%d", id)
	req, err := s.client.NewRequest("GET", spacePath, nil)

	if err != nil {
		return nil, err
	}

	_, err = s.client.Do(req, retrievedSpace)

	if err != nil {
		return nil, err
	}

	return retrievedSpace, nil
}

// Update implements the Spaces API's Update command
func (s *SpacesService) Update(id int, name string) (*Space, error) {
	requestedSpace := &Space{Name: name}
	updatedSpace := &Space{}
	spacePath := fmt.Sprintf("spaces/%d", id)
	req, err := s.client.NewRequest("PUT", spacePath, requestedSpace)

	if err != nil {
		return nil, err
	}

	_, err = s.client.Do(req, updatedSpace)

	if err != nil {
		return nil, err
	}

	return updatedSpace, nil
}

// Delete implements the Spaces API's Delete command
func (s *SpacesService) Delete(id int) error {
	spacePath := fmt.Sprintf("spaces/%d", id)
	req, err := s.client.NewRequest("DELETE", spacePath, nil)

	if err != nil {
		return err
	}

	_, err = s.client.Do(req, nil)

	if err != nil {
		return err
	}

	return nil
}
