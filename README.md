# Cadence Samples
These are some samples in go to demonstrate how to write cadence activities and workflows.

## Prerequisite
Run Cadence Server

See instructions for running the Cadence Server: https://github.com/uber/cadence/blob/master/README.md


## Steps to run the samples
Add the repo in your $GOPATH

Run dep-ensure 
```$xslt
dep ensure
```

cd to the project directory

Start the workers and server to trigger workflows: 
```$xslt
go run cmd/workflows/*.go
```

## The Greeting workflow
The greeting workflow will run following activities: 
1. Get the Name to greet (hardcoded to punchh)
2. Get the Greeting to use (hardcoded to hello)
3. Get the greeting message

To start this workflow call following GET api: 
```$xslt
http://localhost:10000/start-greet-workflow
```

The api will return you the RunID of the job it started.

## The Batch Operations workflow
The batch operations workflow takes following inputs: 
1. An array of ints. 
2. Sequence of operations to run on these ints.
3. Number of times the sequence should be executed. 

The API will start 1 workflow per int in the array. 
The workflow will execute the sequence of operations on the int.

This workflow sample is to demonstrate ability of cadence to be able to start multiple instances of the same workflow.
And to be able to dynamically provide the steps of the workflow.


To start this workflow call following POST api: 
```$xslt
http://localhost:10000/start-batch-workflow
```

You will need to add the asked params in request body. Following is an example: 
```$xslt
{
	"Inputs": [1, 2], 
	"Operations": [
		{
			"Type": "add",
			"Value": 10
		}, 
		{
			"Type": "multiply",
			"Value": 2
		}, 
		{
			"Type": "sleep", 
			"Value": 10000
		}, 
		{
			"Type": "subtract", 
			"Value": 4
		}
	],
	"Iterations": 2
}
```

Following type of operations are available : 
* add
* multiply
* subtract
* sleep

The can be mentioned as many number of times and in any order. 

## Checking the workflow execution and result
Cadence comes with a UI to monitor/access the cadence server.

This can be accessed at : 
```$xslt
http://localhost:8088
```

It will ask you to mention the domain, for the given examples the domain is: 
```$xslt
punchh-samples
```

The Cadence UI lets you:
* See the open/completed/canceled/failed workflows.
* For each job you can check what is going on with details of each event and the value returned by each activity.
* A timeline of how the activities went and whats happening right now.
* For running workflows it lets you check the stacktrace 

You can check the network calls to see the calls UI is making to get this info from cadence server. 
This can help us in making our dashboard.

## Things To DO Next
Add more examples demonstrating other capabilities of cadence that can be useful for us.

Also as I am new to go, so would love to better the code structure and practices in this project. 
## References
1. Cadence Go Examples : https://github.com/samarabbas/cadence-samples (I am using the helper files from this project)
2. Cadence Java Examples : https://github.com/uber/cadence-java-samples
3. Cadence Client : https://github.com/uber-go/cadence-client
4. Cadence Architecture : https://www.youtube.com/watch?v=5M5eiNBUf4Q
5. Cadence Use case : https://www.youtube.com/watch?v=-LRghQzfF8k
6. Cadence, how to write a workflow : https://www.youtube.com/watch?v=Nbz6XUBKdbM


##### Happy Coding!!!
