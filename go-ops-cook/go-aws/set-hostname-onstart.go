package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/codeskyblue/go-sh"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func main() {

	var count int
	var final string
	hn := Get_Local_HostName_to_Judge()
	region := Get_Aws_Region_Name()
	ss := strings.Contains(strings.TrimSpace(hn), "tiejin.cn")
	ss1 := strings.Contains(strings.TrimSpace(hn), "ip-10-")
	if len(region) > 2 && ss != true && (strings.TrimSpace(hn) == "localhost.localdomain" || ss1 == true) {
		final = Final_HostName_Ensure(region)
		real_get_tag_name := strings.Contains(strings.TrimSpace(final), "tiejin.cn")
		if real_get_tag_name != true {
			for count = 0; count <= 6; count++ {
				fmt.Println("didn't get tag the count is:", count)
				time.Sleep(time.Second * 2)
				final = Final_HostName_Ensure(region)
				real_get_tag_name = strings.Contains(strings.TrimSpace(final), "tiejin.cn")
				if real_get_tag_name == true || count == 6 {
					fmt.Println("Will break out of it")
					break
				}
			}
		}
		err := Set_Local_HostName(final)
		Check_Error(err)
		fmt.Println(final)
	} else {
		fmt.Println("args is error")
	}
}

func Get_Aws_Region_Name() (region string) {

	out, _ := sh.Command("curl", "-s", "http://169.254.169.254/latest/meta-data/placement/availability-zone").Command("sed", "s#.$##g").Output()
	return string(out)
	//resp, err := http.Get("curl -s http://169.254.169.254/latest/meta-data/placement/availability-zone |sed 's/.$//g'")
}

func Set_Local_HostName(hostname string) (err error) {

	s := sh.NewSession().SetEnv("hostname1", hostname)
	nis := s.Env["hostname1"]
	sh.Command("echo", nis).Run()
	err = sh.Command("hostnamectl", "set-hostname", nis).Run()
	err = sh.Command("sed", "-i", "s#HOSTNAME=.*#HOSTNAME="+nis+"#g", "/etc/sysconfig/network").Run()
	Check_Error(err)
	return err
}

func Get_Local_HostName_to_Judge() (hsname string) {
	opt_total, err := sh.Command("hostname").Output()
	Check_Error(err)
	return string(opt_total)
}

func Final_HostName_Ensure(region string) (hostname string) {

	var tgs string
	resp, err := http.Get("http://169.254.169.254/latest/meta-data/instance-id")
	Check_Error(err)
	defer resp.Body.Close()
	data, _ := ioutil.ReadAll(resp.Body)
	tgs = GetTagFromId(region, strings.TrimSpace(string(data)))
	if strings.TrimSpace(tgs) == "" {
		hname := Default_HostName()
		hostname = strings.TrimSpace(hname)
	} else {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		a := r.Intn(10000)
		//      b := time.Now().Format("20060102150405")
		rdm := strconv.FormatInt(int64(a), 10)
		if strings.Contains(tgs, "java") || strings.Contains(tgs, "app") {
			hostname = strings.TrimSpace(tgs) + "-" + rdm + ".tiejin.cn"
		} else {
			hostname = strings.TrimSpace(tgs) + ".tiejin.cn"
		}
	}
	return hostname
}

func Default_HostName() (hstname string) {

	resp, err := http.Get("http://169.254.169.254/latest/meta-data/local-hostname")
	Check_Error(err)
	defer resp.Body.Close()
	data, _ := ioutil.ReadAll(resp.Body)
	return string(data)
}

func Check_Error(err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
}

func AwsConfig(region string) *aws.Config {
	sess := session.Must(session.NewSession())
	ec2m := ec2metadata.New(sess,
		&aws.Config{
			HTTPClient: &http.Client{
				Timeout: 20 * time.Second},
		})
	cr := credentials.NewCredentials(&ec2rolecreds.EC2RoleProvider{
		Client: ec2m,
	})
	return &aws.Config{
		Region:      aws.String(region),
		Credentials: cr,
	}
}
func GetTagFromId(region string, insid string) (tag string) {
	sess := session.Must(session.NewSession(AwsConfig(region)))
	svc := ec2.New(sess, AwsConfig(region))
	//svc := ec2.New(sess, &aws.Config{Region: aws.String(region)})
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("instance-id"),
				Values: []*string{
					aws.String(insid),
				},
			},
		},
	}
	resp, err := svc.DescribeInstances(params)
	if err != nil {
		fmt.Println("didn't get that tags")
		tag = ""
		//panic(err)
	}
	//fmt.Println("rtag is:", resp)
	for idx, _ := range resp.Reservations {
		for _, inst := range resp.Reservations[idx].Instances {
			for _, res := range inst.Tags {
				if *res.Key == "Name" {
					tag = *res.Value
				}
			}
		}
	}
	return tag

}
