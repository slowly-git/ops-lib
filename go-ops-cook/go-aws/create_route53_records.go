package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/codeskyblue/go-sh"
	"net"
	"os"
	"strings"
)

func main() {

	hstname, ip := Get_Hostname_And_Ipaddress()
	svc := Aws_Route53_Session("ap-southeast-1")
	params := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action: aws.String("UPSERT"),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name: aws.String(hstname),
						Type: aws.String("A"),
						ResourceRecords: []*route53.ResourceRecord{
							{
								Value: aws.String(ip),
							},
						},
						TTL: aws.Int64(30),
					},
				},
			},
		},
		HostedZoneId: aws.String("Z12BMDK8OYPP9Q"),
	}

	resp, err := svc.ChangeResourceRecordSets(params)
	fmt.Println(resp)
	Check_Error(err)
}

func Get_Hostname_And_Ipaddress() (hostname string, ret string) {

	hostname, _ = os.Hostname()
	ntn, _ := net.Interfaces()
	for _, v := range ntn {
		adr, _ := v.Addrs()
		for _, vv := range adr {
			ip, _, _ := net.ParseCIDR(vv.String())
			if strings.Contains(v.Name, "eth") && (strings.HasPrefix(ip.String(), "10.100") || strings.HasPrefix(ip.String(), "10.200") || strings.HasPrefix(ip.String(), "10.90") || strings.HasPrefix(ip.String(), "10.80") || strings.HasPrefix(ip.String(), "172.") || strings.HasPrefix(ip.String(), "192.")) {
				ret = ip.String()
			}
		}
	}
	return hostname, ret

}

func Aws_Route53_Session(regions string) (svc *route53.Route53) {

	//fmt.Println(string(boutput))
	aws_access_key, err := sh.Command("curl", "-s", "http://seckey.360in.com/secure_key/route53.txt").Command("grep", "-w", "aws_access_key").Command("awk", "-F", "=", "{print $2}").Output()
	aws_secret_access_key, err := sh.Command("curl", "-s", "http://seckey.360in.com/secure_key/route53.txt").Command("grep", "-w", "aws_secret_access_key").Command("awk", "-F", "=", "{print $2}").Output()
	token := ""
	//fmt.Println("aws_access_key", string(aws_access_key))
	//fmt.Println("aws_secret_access_key", string(aws_secret_access_key))

	creds := credentials.NewStaticCredentials(strings.TrimSpace(string(aws_access_key)), strings.TrimSpace(string(aws_secret_access_key)), strings.TrimSpace(string(token)))
	_, err = creds.Get()
	Check_Error(err)
	//fmt.Println("Value of credentials:", v)
	sess, err := session.NewSession()
	Check_Error(err)

	svc = route53.New(sess, &aws.Config{Region: aws.String(regions), Credentials: creds})
	return svc
}

func Check_Error(err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
}
