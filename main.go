package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/go-resty/resty/v2"
)

type APIRequest struct {
	BaseURL string
	Client  *resty.Client
	APIKey  string
}

type ApiResponse struct {
	Status           string            `json:"status"`
	Code             string            `json:"code"`
	Message          string            `json:"message"`
	CorrelationID    string            `json:"correlationId"`
	DetailedMessage  interface{}       `json:"detailedMessage"`
	ResponseMessages []ResponseMessage `json:"responseMessages"`
	Metadata         interface{}       `json:"metadata"`
}

type ResponseMessage struct {
	Code         string        `json:"code"`
	Level        string        `json:"level"`
	Message      string        `json:"message"`
	Exception    interface{}   `json:"exception"`
	FailureTypes []interface{} `json:"failureTypes"`
}

func main() {
	var projectList []ProjectsContent
	var serviceList []*ServiceClass

	accountArg := flag.String("account", "", "Provide your account ID.")
	apiKeyArg := flag.String("api-key", "", "Provide your API Key.")
	targetString := flag.String("target", "", "Provide your targeted string.")
	replaceString := flag.String("replacement", "", "Provide your replacement string.")

	flag.Parse()

	type Flags struct {
		Account     string
		ApiKey      string
		Target      string
		Replacement string
	}

	f := Flags{
		Account:     *accountArg,
		ApiKey:      *apiKeyArg,
		Target:      *targetString,
		Replacement: *replaceString,
	}

	api := APIRequest{
		BaseURL: "https://app.harness.io",
		Client:  resty.New(),
		APIKey:  f.ApiKey,
	}

	projects, err := api.GetAllProjects(f.Account)
	if err != nil {
		fmt.Print(color.RedString("Unable to get projects - %s", err))
		return
	}
	projectList = append(projectList, projects.Data.Content...)
	for _, project := range projectList {
		p := project.Project
		service, err := api.GetService(f.Account, string(p.OrgIdentifier), p.Identifier)
		if err != nil {
			fmt.Print(color.RedString("Unable to get service - %s", err))
		}
		serviceList = append(serviceList, service...)
	}
	for _, service := range serviceList {
		service.YAML = strings.ReplaceAll(service.YAML, f.Target, f.Target)
		service.UpdateService(&api)
		fmt.Print(color.GreenString("Updated services - %s", service.Name))
	}
}

func (api *APIRequest) GetAllProjects(account string) (Projects, error) {
	resp, err := api.Client.R().
		SetHeader("x-api-key", api.APIKey).
		SetQueryParams(map[string]string{
			"accountIdentifier": account,
			"hasModule":         "true",
			"pageSize":          "500",
		}).
		Get(api.BaseURL + "/ng/api/projects")
	if err != nil {
		return Projects{}, err
	}
	projects := Projects{}
	err = json.Unmarshal(resp.Body(), &projects)
	if err != nil {
		return Projects{}, err
	}

	return projects, nil
}

func (api *APIRequest) GetService(account, org, project string) ([]*ServiceClass, error) {
	params := map[string]string{
		"accountIdentifier": account,
		"orgIdentifier":     org,
		"projectIdentifier": project,
	}

	resp, err := api.Client.R().
		SetHeader("x-api-key", api.APIKey).
		SetHeader("Content-Type", "application/json").
		SetQueryParams(params).
		SetPathParam("org", org).
		SetPathParam("project", project).
		Get(api.BaseURL + "/v1/orgs/{org}/projects/{project}/services")
	if err != nil {
		return []*ServiceClass{}, err
	}

	service := []*Service{}
	err = json.Unmarshal(resp.Body(), &service)
	if err != nil {
		return []*ServiceClass{}, err
	}

	serviceList := []*ServiceClass{}
	for _, s := range service {
		serviceList = append(serviceList, &s.Service)
	}

	return serviceList, nil
}

func (api *APIRequest) UpdateService(service ServiceRequest, account string) error {
	params := map[string]string{
		"accountIdentifier": account,
	}
	resp, err := api.Client.R().
		SetHeader("x-api-key", api.APIKey).
		SetHeader("Content-Type", "application/json").
		SetQueryParams(params).
		SetBody(service).
		Put(api.BaseURL + "/ng/api/servicesV2")
	if err != nil {
		return err
	}

	if resp.StatusCode() != 200 {
		ar := ApiResponse{}
		err = json.Unmarshal(resp.Body(), &ar)
		if err != nil {
			return err
		}
		errMsg := fmt.Sprintf("CorrelationId: %s, ResponseMessages: %+v", ar.CorrelationID, ar.ResponseMessages)
		return fmt.Errorf(errMsg)
	}

	return nil
}
