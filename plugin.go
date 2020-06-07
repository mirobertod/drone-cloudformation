package main

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/ecs"
)

type Settings struct {
	Mode     string  `json:"mode"`
	Region   string  `json:"region"`
	Parallel bool    `json:"parallel"`
	Batch    []Batch `json:"batch"`
}

type Batch struct {
	StackName          string   `json:"stackname"`
	ServiceLogicalName string   `json:"servicelogicalname"`
	ClusterStackName   string   `json:"clusterstackname"`
	Template           string   `json:"template"`
	Params             []Params `json:"params"`
}

type Params struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

var sess = session.Must(session.NewSessionWithOptions(session.Options{
	SharedConfigState: session.SharedConfigEnable,
}))
var cfSvc = cloudformation.New(sess, &aws.Config{
	Region: aws.String("eu-west-1"),
})
var ecsSvc = ecs.New(sess, &aws.Config{
	Region: aws.String("eu-west-1"),
})

func run() {

	batch := parseBatch()
	log.Println(batch[0].Params[0].Key)
	log.Println(batch[0].Params[0].Value)

	action := getAction("stackname", "logicalservice", "clusterstackname")
	log.Println(action)
}

func parseBatch() []Batch {
	data := `[{"params":[{"key":"environment","value":"staging"}],"stackname":"my-database-stack","template":"templates/db.yml"},{"params":[{"key":"version","value":"123"},{"key":"environment","value":"staging"}],"stackname":"my-app-stack","template":"templates/app.yml"}]`

	// pair := os.Getenv("PLUGIN_BATCH")
	var message = []byte(data)
	var batch []Batch

	err := json.Unmarshal(message, &batch)
	if err != nil {
		log.Println("Error unmarshal", err)
	}
	return batch
}

func stackExists(stackName string) bool {
	stacksInput := &cloudformation.DescribeStacksInput{StackName: aws.String(stackName)}
	result, _ := cfSvc.DescribeStacks(stacksInput)
	if len(result.Stacks) == 1 {
		return true
	}
	return false
	// fmt.Println(awsutil.StringValue(result.Stacks[0].StackStatus))
}

func getAction(stackName string, serviceLogicalName string, clusterStackName string) string {
	if stackExists(stackName) {
		getTaskDesiredCount(stackName, serviceLogicalName, clusterStackName)
		return "update"
	}
	return "create"
}

func getTaskDesiredCount(stackName string, serviceLogicalName string, clusterStackName string) string {

	stackResourceInput := &cloudformation.DescribeStackResourceInput{StackName: aws.String(clusterStackName), LogicalResourceId: aws.String("ECSCluster")}
	stackResource, _ := cfSvc.DescribeStackResource(stackResourceInput)

	ecsCluster := stackResource.StackResourceDetail.PhysicalResourceId

	stackResourceInput = &cloudformation.DescribeStackResourceInput{StackName: aws.String(stackName), LogicalResourceId: aws.String(serviceLogicalName)}
	stackResource, _ = cfSvc.DescribeStackResource(stackResourceInput)

	var ecsService []*string
	ecsService = append(ecsService, stackResource.StackResourceDetail.PhysicalResourceId)

	describeServicesInput := &ecs.DescribeServicesInput{Services: ecsService, Cluster: ecsCluster}
	services, _ := ecsSvc.DescribeServices(describeServicesInput)

	return awsutil.StringValue(services.Services[0].RunningCount)

}
