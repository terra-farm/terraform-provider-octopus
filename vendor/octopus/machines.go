package octopus

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Machine represents a machine that Octopus can target for deployment.
type Machine struct {
	ID                string            `json:"Id"`
	Name              string            `json:"Name"`
	Thumbprint        string            `json:"Thumbprint"`
	URI               string            `json:"Uri"`
	IsDisabled        bool              `json:"IsDisabled"`
	EnvironmentIDs    []string          `json:"EnvironmentIds"`
	Roles             []string          `json:"Roles"`
	Status            string            `json:"Status"`
	HasLatestCalamari bool              `json:"HasLatestCalamari"`
	Endpoint          Endpoint          `json:"Endpoint"`
	Links             map[string]string `json:"Links"`
}

// Endpoint represents an Octopus deployment end-point.
type Endpoint struct {
	ID                     *string                `json:"Id,omitempty"`
	CommunicationsStyle    string                 `json:"CommunicationsStyle"`
	URI                    string                 `json:"Uri"`
	Thumbprint             string                 `json:"Thumbprint"`
	TentacleVersionDetails TentacleVersionDetails `json:"TentacleVersionDetails"`
	LastModifiedOn         *string                `json:"LastModifiedOn,omitempty"`
	LastModifiedBy         *string                `json:"LastModifiedBy,omitempty"`
	Links                  map[string]string      `json:"Links"`
}

// TentacleVersionDetails represents version information for an Octopus tentacle.
type TentacleVersionDetails struct {
	Version          string `json:"Version"`
	UpgradeSuggested bool   `json:"UpgradeSuggested"`
	UpgradeRequired  bool   `json:"UpgradeRequired"`
	UpgradeLocked    bool   `json:"UpgradeLocked"`
}

// GetMachine retrieves an Octopus machine by Id or slug.
func (client *Client) GetMachine(idOrSlug string) (machine *Machine, err error) {
	var (
		request       *http.Request
		statusCode    int
		responseBody  []byte
		errorResponse *APIErrorResponse
	)

	requestURI := fmt.Sprintf("machines/%s", idOrSlug)
	request, err = client.newRequest(requestURI, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	responseBody, statusCode, err = client.executeRequest(request)
	if err != nil {
		err = fmt.Errorf("Error invoking request to read variable set '%s': %s", idOrSlug, err.Error())

		return
	}

	if statusCode == http.StatusNotFound {
		// Environment not found.
		return nil, nil
	}

	if statusCode != http.StatusOK {
		errorResponse, err = readAPIErrorResponseAsJSON(responseBody, statusCode)
		if err != nil {
			return nil, err
		}

		return nil, errorResponse.ToError("Request to retrieve machine '%s' failed with status code %d.", idOrSlug, statusCode)
	}

	machine = &Machine{}
	err = json.Unmarshal(responseBody, machine)
	if err != nil {
		err = fmt.Errorf("Invalid response detected when retrieving machine '%s': %s", idOrSlug, err.Error())
	}

	return
}
