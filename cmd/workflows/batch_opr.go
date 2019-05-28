package main

import (
	"go.uber.org/cadence/activity"
	"go.uber.org/cadence/workflow"
	"time"
)

/**
	Following sample lets user trigger a workflow which takes following inputs :
	1. An integer
	2. an Array of Operations to perform on the integer, available operations are : ADD, SUBTRACT, MULTIPLY, SLEEP
	3. Number of times this sequence of operation needs to be run.

	The workflow will return the final value after performing the operations.
 */


/**
	Enum to define Operation type
 */
type OperationType string

/**
	All possible Operations
 */
const (
	ADD OperationType = "add"
	SUBTRACT OperationType = "subtract"
	MULTIPLY OperationType = "multiply"
	SLEEP OperationType = "sleep"
)

/**
	Struct to define operation
 */
type Operation struct {
	Type OperationType
	Value int
}

/**
	Registering the activities and workflow
 */
func init(){
	activity.Register(add)
	activity.Register(subtract)
	activity.Register(multiply)

	workflow.Register(batchOperationWorkflow)
}

/**
	Activity to add two numbers
 */
func add(x int, y int) (int, error){
	return x+y, nil
}

/**
	Activity to subtract two numbers
 */
func subtract(x int, y int) (int, error){
	return x-y, nil
}

/**
	Activity to multiply two numbers
 */
func multiply(x int, y int) (int, error){
	return x*y, nil
}

/**
	Workflow to run the given operations on input int
 */
func batchOperationWorkflow(ctx workflow.Context, input int, operations []Operation, iterations int) (int, error){
	logger := workflow.GetLogger(ctx)

	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute * 10,
		StartToCloseTimeout: time.Minute * 10,
		HeartbeatTimeout: time.Minute,
	}

	ctx = workflow.WithActivityOptions(ctx, ao)

	val := input
	var err error
	for i:=0; i< iterations; i++{
		for _, operation := range operations{
			switch operation.Type{
			case ADD:
				err = workflow.ExecuteActivity(ctx, add, val, operation.Value).Get(ctx, &val)
				if err != nil {
					logger.Error("Not able to execute Add")
					return 0, err
				}
			case SUBTRACT:
				err = workflow.ExecuteActivity(ctx, subtract, val, operation.Value).Get(ctx, &val)
				if err != nil {
					logger.Error("Not able to execute Subtract")
					return 0, err
				}
			case MULTIPLY:
				err = workflow.ExecuteActivity(ctx, multiply, val, operation.Value).Get(ctx, &val)
				if err != nil {
					logger.Error("Not able to execute Multiply")
					return 0, err
				}
			case SLEEP:
				err = workflow.Sleep(ctx, time.Millisecond * time.Duration(operation.Value))
				if err != nil {
					logger.Error("Not able to execute Sleep")
					return 0, err
				}
			}
		}
	}

	return val, nil
}