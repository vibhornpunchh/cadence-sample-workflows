package main

import (
	"encoding/json"
	"fmt"
	"github.com/vibhornpunchh/cadence-sample-workflows/cmd/common"
	"go.uber.org/cadence/client"
	"go.uber.org/cadence/worker"
	"log"
	"net/http"
	"time"

	"github.com/pborman/uuid"

	"github.com/gorilla/mux"
)

/**
	Helper class for connecting to cadence-client
 */
var h common.SampleHelper

/**
	Worker Identifier
 */
const ApplicationName = "PunchhCadenceSamples"

/**
	Struct for Batch operation workflow request
	Inputs : Array of input ints, for each int in this array one workflow will be started
	Operations: Array of operations to be performed on the each input. Check the operation struct for more.
	Iterations: Number of iterations of set of operations to be run on each input
 */
type BatchRequest struct{
	Inputs []int
	Operations []Operation
	Iterations int
}

/**
	This needs to be done as part of a bootstrap step when the process starts.
	The workers are supposed to be long running.
 */
func startWorkers(h *common.SampleHelper) {
	// Configure worker options.
	workerOptions := worker.Options{
		MetricsScope: h.Scope,
		Logger:       h.Logger,
	}
	h.StartWorkers(h.Config.DomainName, ApplicationName, workerOptions)
}



/**
	Service method to start the greeting workflow
 */
func startGreetXWorkflow(h *common.SampleHelper) string{
	workflowOptions := client.StartWorkflowOptions{
		ID:                              "punchh_" + uuid.New(),
		TaskList:                        ApplicationName,
		ExecutionStartToCloseTimeout:    time.Minute * 10,
		DecisionTaskStartToCloseTimeout: time.Minute * 10,
	}
	return h.StartWorkflow(workflowOptions, greetXWorkflow)
}

/**
	Service method to start the batch workflow
 */
func startBatchWorkflow(h *common.SampleHelper, batchRequest BatchRequest) []string{
	var runIds []string
	for _, input := range batchRequest.Inputs{
		workflowOptions := client.StartWorkflowOptions{
			ID:                              "punchh_" + uuid.New(),
			TaskList:                        ApplicationName,
			ExecutionStartToCloseTimeout:    time.Minute * 10,
			DecisionTaskStartToCloseTimeout: time.Minute * 10,
		}
		runIds = append(runIds, h.StartWorkflow(workflowOptions, batchOperationWorkflow, input, batchRequest.Operations, batchRequest.Iterations))
	}
	return runIds
}

/**
	Controller for managing trigger for greeting workflow.
 */
func startGreetXWorkflowController(w http.ResponseWriter, r *http.Request){

	_, _ = fmt.Fprint(w, "Started workflow with runId " + startGreetXWorkflow(&h))
}


/**
	Controller for managing the trigger for batch operation work flow.
	It expects a Batch request struct in body
 */
func startBatchOperationWorkflowController(w http.ResponseWriter, r *http.Request){
	batchRequest := BatchRequest{}

	err := json.NewDecoder(r.Body).Decode(&batchRequest)
	if err != nil{
		fmt.Println(err)
	}

	runIDs := startBatchWorkflow(&h,batchRequest)
	_, _ = fmt.Fprint(w, "Started workflow with following runIds \n")
	for _, runID := range runIDs{
		_, _ = fmt.Fprint(w, runID + "\n")
	}
}

/**
	Router for trigger listener server
 */
func handleRequests(){
	myRouter := mux.NewRouter()
	myRouter.HandleFunc("/start-greet-workflow", startGreetXWorkflowController).Methods("GET")
	myRouter.HandleFunc("/start-batch-workflow", startBatchOperationWorkflowController).Methods("POST")
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

/**
	Need to start the workers, these workers will be the ones which will be triggering the workflow instances
	Also starting a server to listen for trigger requests
 */
func main() {
	h.SetupServiceConfig()
	startWorkers(&h)
	handleRequests()
}
