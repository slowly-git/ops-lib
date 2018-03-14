package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	//	"io/ioutil"
	//	"net/http"
	"strings"
)

func main() {
	//svc := Aws_Ec2_Session("cn-north-1", "AKIAOZRIOKXIE4F2MVMQ", "iIp7dFx7QvoGA2g9Y0bytvVqyAmco8g6EH+BmM9H") //primary

	var all_tags []string
	ltags_cn_primary := GetTagFromRegion("cn-north-1", "AKIAOZRIOKXIE4F2MVMQ", "iIp7dFx7QvoGA2g9Y0bytvVqyAmco8g6EH+BmM9H")    //primary_cn_north
	ltags_cn_secondary := GetTagFromRegion("cn-north-1", "AKIAOEUYH4TNOBO33Y7A", "OKTcVwMp37B5YcDS72niVaPUvxRXzAoHJiM7YJtz")  //primary_cn_secondary
	ltags_ap_sourth := GetTagFromRegion("ap-southeast-1", "AKIAIZ4U7FSPO4CQY5WQ", "s/hPxbEpZ2Wa5H8WIn4ypKcaApjv0Rzg3e9YtbMR") //ap-southeast-1

	two_tags := append(ltags_cn_primary, ltags_ap_sourth...)
	all_tags = append(two_tags, ltags_cn_secondary...)
	for _, info := range all_tags {
		fmt.Println(info)
	}
}
func GetTagFromRegion(region string, key string, secert string) (rtags []string) {
	//var rtags []string

	svc := Aws_Ec2_Session(region, key, secert) //primary
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("instance-state-name"),
				Values: []*string{
					aws.String("running"),
				},
			},
		},
	}
	resp, err := svc.DescribeInstances(params)
	if err != nil {
		panic(err)
	}
	//      fmt.Println("rtag is:", resp)
	for idx, _ := range resp.Reservations {
		for _, inst := range resp.Reservations[idx].Instances {
			//fmt.Println(inst.Tags)
			for key, _ := range inst.Tags {
				//fmt.Println(*inst.Tags[key].Key)
				if *inst.Tags[key].Key == "Name" {
					//fmt.Println(*inst.Tags[key].Value)
					tags := *inst.Tags[key].Value
					rtags = append(rtags, strings.TrimSpace(tags)+"\t"+*inst.PrivateIpAddress)

				}
			}
		}
	}
	return rtags

}
func Aws_Ec2_Session(regions string, aws_access_key string, aws_secret_access_key string) (svc *ec2.EC2) {

	token := ""
	creds := credentials.NewStaticCredentials(aws_access_key, aws_secret_access_key, token)
	_, err := creds.Get()
	Check_Error(err)
	//fmt.Println("Value of credentials:", v)
	sess, err := session.NewSession()
	Check_Error(err)

	svc = ec2.New(sess, &aws.Config{Region: aws.String(regions), Credentials: creds})
	return svc
}
func Check_Error(err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
}
