package main

import (
	"fmt"
	"go.uber.org/cadence/activity"
	"go.uber.org/cadence/workflow"
	"time"
)

/**
	The following is a greeting sample, this submits a workflow in cadence which does following :
	1. Get the name to greet
	2. Get the greeting to greet with
	3. Get the final greeting message
 */

/**
	Registration of activities and workflow
 */
func init(){
	activity.Register(getNameActivity)
	activity.Register(getGreetingActivity)
	activity.Register(setGreetingActivity)
	workflow.Register(greetXWorkflow)
}

/**
	Activity to return the name to greet
 */
func getNameActivity() (string, error){
	return "Punchh", nil
}

/**
	Activity to return the greeting
 */
func getGreetingActivity() (string, error){
	return "Hello", nil
}

/**
	Activity to return the greeting message
 */
func setGreetingActivity(greeting string, name string) (string, error){
	result := fmt.Sprintf("%s %s !\n", greeting, name)
	return result, nil
}


/**
	Workflow to perform:
	1. Get the name
	2. Get the greeting
	3. return the greeting message
 */
func greetXWorkflow(ctx workflow.Context) error{
	logger := workflow.GetLogger(ctx)

	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute * 10,
		StartToCloseTimeout: time.Minute * 10,
		HeartbeatTimeout: time.Minute,
	}

	ctx = workflow.WithActivityOptions(ctx, ao)

	var nameResult string
	err := workflow.ExecuteActivity(ctx, getNameActivity).Get(ctx, &nameResult)
	if err != nil{
		logger.Error("Get Name Failed.")
		return err
	}

	var greetResult string
	err = workflow.ExecuteActivity(ctx, getGreetingActivity).Get(ctx, &greetResult)
	if err != nil{
		logger.Error("Get Greeting Failed.")
		return err
	}

	var result string
	err = workflow.ExecuteActivity(ctx, setGreetingActivity, greetResult, nameResult).Get(ctx, &result)
	if err != nil{
		logger.Error("Set Greeting message Failed.")
		return err
	}

	logger.Info("Greeting Message: " + result )


	return nil

}