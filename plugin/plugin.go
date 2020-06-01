package plugin

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/aws/aws-sdk-go/aws"
	_ "github.com/aws/aws-sdk-go/aws/session"
	_ "github.com/aws/aws-sdk-go/service/s3"
)

type Settings struct {
	Mode     string  `json:"mode"`
	Region   string  `json:"region"`
	Parallel bool    `json:"parallel"`
	Batch    []Batch `json:"batch"`
}

type Batch struct {
	StackName string   `json:"stackname"`
	Template  string   `json:"template"`
	Params    []Params `json:"params"`
}

type Params struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func Run() {
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if strings.HasPrefix(pair[0], "PLUGIN_") {
			fmt.Println(pair[0])
			fmt.Println(pair[1])
			var message = []byte(pair[1])
			var batch []Batch

			err := json.Unmarshal(message, &batch)
			if err != nil {
				log.Println("Error unmarshal")
			}
			log.Println(batch)
		}
	}

	// sess := session.Must(session.NewSession(&aws.Config{
	// 	Region: aws.String("eu-west-1"),
	// }))
	// svc := s3.New(sess)
	// input := &s3.ListBucketsInput{}
	// result, _ := svc.ListBuckets(input)
	// fmt.Println(result)
}
