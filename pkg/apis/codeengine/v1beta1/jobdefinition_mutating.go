/*******************************************************************************
 * Licensed Materials - Property of IBM
 * IBM Cloud Code Engine, 5900-AB0
 * © Copyright IBM Corp. 2020
 * US Government Users Restricted Rights - Use, duplication or
 * disclosure restricted by GSA ADP Schedule Contract with IBM Corp.
 ******************************************************************************/

package v1beta1

import (
	"fmt"

	uuid "github.com/google/uuid"
	"knative.dev/pkg/apis"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

const (
	// EventingJobrunnerName is the name of the server who receive events
	EventingJobrunnerName = "jobrunner"
	// EventingNamespace is the namespace that contains the events server
	EventingNamespace = "knative-eventing"
)

// MutateJobDefStatus adds Addressable URL status to make jobdefinitions addressable
// Will need to move to controller if a jobdefinition controller is ever created
func (jd *JobDefinition) MutateJobDefStatus() {
	jd.Status = JobDefinitionStatus{}

	host := fmt.Sprintf("%s.%s.svc.cluster.local", EventingJobrunnerName, EventingNamespace)
	uuid := uuid.New().String()
	path := fmt.Sprintf("/%s/%s/%s", jd.Namespace, jd.Name, uuid)
	jd.Status.Address = &duckv1.Addressable{
		URL: &apis.URL{
			Scheme: "http",
			Host:   host,
			Path:   path,
		},
	}
}
