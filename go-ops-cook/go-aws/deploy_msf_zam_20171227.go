package main

import (
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/codeskyblue/go-sh"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	var dt = "nidaye"
	var user *string
	user = &dt

	password := flag.String("password", "", "The Password To Get Auth File")
	authfile := flag.String("authfile", "", "The File Have Aws Authorized Information")
	region := flag.String("region", "ap-southeast-1", "Aws region: ap-southeast-1,cn-north-1,us-east-1 ?")
	elbname := flag.String("elbname", "", "The Elb Will Be Operatei,If Have Multiple Elb Split With ','")
	elbv2name := flag.String("elbv2name", "", "The Elb V2 Target Group, If Have Multiple Elb V2 Target Group Split With ','")
	elbv2port := flag.Int64("elbv2port", 80, "The Elb V2 Target Group Register Port , Default is 80")
	msfname := flag.String("msfname", "", "The Msf Name In supervisor")
	action := flag.String("action", "", "The Action For msf service stop,start,restart")
	flag.Parse()

	fmt.Println()
	fmt.Println("password is:", *password)
	fmt.Println("authfile is:", *authfile)
	fmt.Println("region is:", *region)
	fmt.Println("elbname is:", *elbname)
	fmt.Println("elbv2name is:", *elbv2name)
	fmt.Println("elbv2port is:", *elbv2port)
	fmt.Println("msfname is:", *msfname)
	fmt.Println("action is:", *action)
	fmt.Println()

	sid := GetInstanceId()

	//Remove instance From Elb And Elbv2
	if *elbname != "" && *elbv2name == "" {
		fmt.Println("Start DeReginsterInstance_Form_Classic_Elbs.......")
		DeRegisterInstance_From_Elbs(*region, *elbname, sid, *user, *password, *authfile)
		fmt.Println("End DeReginsterInstance_Form_Classic_Elbs.......")
	} else if *elbname == "" && *elbv2name != "" {
		fmt.Println("Start DeReginsterInstance_From_Elbv2.......")
		DeRegisterInstance_From_Elbv2s(*region, *elbv2name, sid, elbv2port, *user, *password, *authfile)
		fmt.Println("End DeReginsterInstance_From_Elbv2.......")
	} else {
		fmt.Println("This Project Have Elb And Elbv2 Will Operate seq")
		fmt.Println("Start DeReginsterInstance_Form_Classic_Elbs.......")
		DeRegisterInstance_From_Elbs(*region, *elbname, sid, *user, *password, *authfile)
		fmt.Println("End DeReginsterInstance_Form_Classic_Elbs.......")
		fmt.Println()
		fmt.Println("Start DeReginsterInstance_From_Elbv2.......")
		DeRegisterInstance_From_Elbv2s(*region, *elbv2name, sid, elbv2port, *user, *password, *authfile)
		fmt.Println("End DeReginsterInstance_From_Elbv2.......")
	}

	//Operate The Msf Service
	RestartMsfService(*action, *msfname)

	// Check if msf Service is started
	for {
		st := JudgeIfMsfStarted()
		time.Sleep(5 * time.Second)
		if st == "succ" {
			fmt.Println("msf is started")
			break
		}
	}

	//Add Instance To Elb And Elbv2
	if *elbname != "" && *elbv2name == "" {
		fmt.Println("Start ReginsterInstance_To_Classic_Elbs.......")
		RegisterInstance_To_Elbs(*region, *elbname, sid, *user, *password, *authfile)
		fmt.Println("End ReginsterInstance_to_Classic_Elbs.......")
	} else if *elbname == "" && *elbv2name != "" {
		fmt.Println("Start ReginsterInstance_To_Elbv2.......")
		RegisterInstance_To_Elbv2s(*region, *elbv2name, sid, elbv2port, *user, *password, *authfile)
		fmt.Println("End ReginsterInstance_To_Elbv2.......")
	} else {
		fmt.Println("This Project Have Elb And Elbv2 Will Operate seq")
		fmt.Println("Start ReginsterInstance_To_Classic_Elbs.......")
		RegisterInstance_To_Elbs(*region, *elbname, sid, *user, *password, *authfile)
		fmt.Println("End ReginsterInstance_to_Classic_Elbs.......")
		fmt.Println()
		fmt.Println("Start ReginsterInstance_To_Elbv2.......")
		RegisterInstance_To_Elbv2s(*region, *elbv2name, sid, elbv2port, *user, *password, *authfile)
		fmt.Println("End ReginsterInstance_To_Elbv2.......")
	}

}

func RegisterInstance_To_Elbv2s(region, elbname, sid string, port *int64, user string, password string, authfile string) {
	if strings.Contains(elbname, ",") {
		lst := strings.FieldsFunc(elbname, func(c rune) bool {
			if c == ',' {
				return true
			}
			return false
		})
		fmt.Println("register elbv2 lst is:", lst)
		for _, elname := range lst {
			for {
				add_res_elbv2 := RegisterInstanceToElbv2(region, elname, strings.TrimSpace(sid), port, user, password, authfile)
				ret_elbv2 := DescribeElbv2InstanceHealth(region, elname, strings.TrimSpace(sid), port, user, password, authfile)
				time.Sleep(5 * time.Second)
				fmt.Println("add_res_elbv2:", add_res_elbv2)
				fmt.Println("ret_elbv2;", ret_elbv2)
				if strings.TrimSpace(add_res_elbv2) == "succ" && strings.TrimSpace(ret_elbv2) == "healthy" {
					fmt.Println("instance is add to elbv2s")
					break
				}
			}
		}
	} else {
		for {
			add_res_elbv2 := RegisterInstanceToElbv2(region, elbname, strings.TrimSpace(sid), port, user, password, authfile)
			ret_elbv2 := DescribeElbv2InstanceHealth(region, elbname, strings.TrimSpace(sid), port, user, password, authfile)
			time.Sleep(5 * time.Second)
			fmt.Println("add_res_elbv2:", add_res_elbv2)
			fmt.Println("ret_elbv2;", ret_elbv2)
			if strings.TrimSpace(add_res_elbv2) == "succ" && strings.TrimSpace(ret_elbv2) == "healthy" {
				fmt.Println("instance is add to elbv2s")
				break
			}

		}
	}
}
func RegisterInstance_To_Elbs(region string, elbname string, sid string, user string, password string, authfile string) {
	if strings.Contains(elbname, ",") {
		lst := strings.FieldsFunc(elbname, func(c rune) bool {
			if c == ',' {
				return true
			}
			return false
		})
		fmt.Println("register elb lst is:", lst)
		for _, elname := range lst {
			for {
				add_res := RegisterInstanceToElb(region, elname, strings.TrimSpace(sid), user, password, authfile)
				ret := DescribeElbInstanceHealth(region, elname, strings.TrimSpace(sid), user, password, authfile)
				time.Sleep(5 * time.Second)
				fmt.Println("add_res is:", add_res)
				fmt.Println("ret is:", ret)
				if strings.TrimSpace(add_res) == "succ" && strings.TrimSpace(ret) == "InService" {
					fmt.Println("instance is add to elbs")
					break
				}
			}
		}
	} else {
		for {
			add_res := RegisterInstanceToElb(region, elbname, strings.TrimSpace(sid), user, password, authfile)
			ret := DescribeElbInstanceHealth(region, elbname, strings.TrimSpace(sid), user, password, authfile)
			time.Sleep(5 * time.Second)
			fmt.Println("add_res is:", add_res)
			fmt.Println("ret is:", ret)
			if strings.TrimSpace(add_res) == "succ" && strings.TrimSpace(ret) == "InService" {
				fmt.Println("instance is add to elbs")
				break
			}

		}
	}
}
func DeRegisterInstance_From_Elbv2s(region string, elbname string, sid string, port *int64, user string, password string, authfile string) {
	if strings.Contains(elbname, ",") {
		lst := strings.FieldsFunc(elbname, func(c rune) bool {
			if c == ',' {
				return true
			}
			return false
		})
		fmt.Println("elbv2 lst is:", lst)
		for _, elname := range lst {
			for {
				del_res_elbv2 := DeRegisterInstanceToElbv2(region, elname, strings.TrimSpace(sid), user, password, authfile)
				ret_elbv2 := DescribeElbv2InstanceHealth(region, elname, strings.TrimSpace(sid), port, user, password, authfile)
				time.Sleep(5 * time.Second)
				fmt.Println("del_res_elbv2 is:", del_res_elbv2)
				fmt.Println("ret_elbv2 is:", ret_elbv2)
				if strings.TrimSpace(del_res_elbv2) == "succ" && strings.TrimSpace(ret_elbv2) == "unused" {
					fmt.Println("instance is detach from elbv2")
					break
				}
			}

		}

	} else {
		for {
			del_res_elbv2 := DeRegisterInstanceToElbv2(region, elbname, strings.TrimSpace(sid), user, password, authfile)
			ret_elbv2 := DescribeElbv2InstanceHealth(region, elbname, strings.TrimSpace(sid), port, user, password, authfile)
			time.Sleep(5 * time.Second)
			fmt.Println("del_res_elbv2 is:", del_res_elbv2)
			fmt.Println("ret_elbv2 is:", ret_elbv2)
			if strings.TrimSpace(del_res_elbv2) == "succ" && strings.TrimSpace(ret_elbv2) == "unused" {
				fmt.Println("instance is detach from elbv2")
				break
			}
		}

	}
}

func DeRegisterInstance_From_Elbs(region string, elbname string, sid string, user string, password string, authfile string) {
	if strings.Contains(elbname, ",") {
		lst := strings.FieldsFunc(elbname, func(c rune) bool {
			if c == ',' {
				return true
			}
			return false
		})
		fmt.Println("lst is:", lst)
		for _, elname := range lst {
			for {
				del_res := DeRegisterInstanceToElb(region, elname, strings.TrimSpace(sid), user, password, authfile) //ap-southeast-1
				ret := DescribeElbInstanceHealth(region, elname, strings.TrimSpace(sid), user, password, authfile)   //ap-southeast-1
				time.Sleep(5 * time.Second)
				fmt.Println("del_res is:", del_res)
				fmt.Println("ret is:", ret)
				if strings.TrimSpace(del_res) == "succ" && strings.TrimSpace(ret) == "OutOfService" {
					fmt.Println("instance is detach from elb")
					break
				}
			}

		}

	} else {
		for {
			del_res := DeRegisterInstanceToElb(region, elbname, strings.TrimSpace(sid), user, password, authfile) //ap-southeast-1
			ret := DescribeElbInstanceHealth(region, elbname, strings.TrimSpace(sid), user, password, authfile)   //ap-southeast-1
			time.Sleep(5 * time.Second)
			fmt.Println("del_res is:", del_res)
			fmt.Println("ret is:", ret)
			if strings.TrimSpace(del_res) == "succ" && strings.TrimSpace(ret) == "OutOfService" {
				fmt.Println("instance is detach from elb")
				break
			}
		}

	}

}

func RestartMsfService(action string, msfname string) {
	os.Chdir("/home/worker/supervisor")
	err := sh.Command("supervisorctl", action, msfname).Run()
	Check_Error(err)
}

func JudgeIfMsfStarted() (status string) {
	ot, err := sh.Command("ps", "-ef").Command("grep", "-i", "msf").Command("grep", "-v", "grep").Command("wc", "-l").Output()
	Check_Error(err)
	//fmt.Println(string(ot))
	num, err := strconv.ParseFloat(strings.TrimSpace((string(ot))), 64)
	Check_Error(err)
	if num > 1 {
		status = "succ"
	} else {
		status = "fail"
	}
	return status
}

func DescribeElbInstanceHealth(region string, elbname string, insid string, user string, password string, authfile string) (state string) {

	svc := Aws_Elb_Session(region, user, password, authfile) //primary
	params := &elb.DescribeInstanceHealthInput{
		LoadBalancerName: aws.String(elbname),
		Instances: []*elb.Instance{
			{
				InstanceId: aws.String(insid),
			},
		},
	}
	resp, err := svc.DescribeInstanceHealth(params)
	Check_Error(err)
	//fmt.Println(*resp.InstanceStates[0].State)
	return *resp.InstanceStates[0].State
}
func DescribeElbv2InstanceHealth(region string, elbname string, insid string, inport *int64, user string, password string, authfile string) (state string) {

	svc := Aws_Elbv2_Session(region, user, password, authfile)
	params := &elbv2.DescribeTargetHealthInput{

		TargetGroupArn: aws.String(elbname),
		Targets: []*elbv2.TargetDescription{
			{
				Id:   aws.String(insid),
				Port: aws.Int64(*inport),
			},
		},
	}
	resp, err := svc.DescribeTargetHealth(params)
	Check_Error(err)
	return *resp.TargetHealthDescriptions[0].TargetHealth.State

}

func RegisterInstanceToElbv2(region string, elbname string, insid string, inport *int64, user string, password string, authfile string) (ret string) {

	svc := Aws_Elbv2_Session(region, user, password, authfile)
	params := &elbv2.RegisterTargetsInput{
		TargetGroupArn: aws.String(elbname),
		Targets: []*elbv2.TargetDescription{
			{
				Id:   aws.String(insid),
				Port: aws.Int64(*inport),
			},
		},
	}

	resp, err := svc.RegisterTargets(params)

	fmt.Println(resp.String())
	if err != nil {
		fmt.Println(err.Error())
		ret = "fail"
	} else {
		ret = "succ"
	}

	return ret
}

func RegisterInstanceToElb(region string, elbname string, insid string, user string, password string, authfile string) (ret string) {

	var inelbs []string
	var count int64
	svc := Aws_Elb_Session(region, user, password, authfile) //primary
	params := &elb.RegisterInstancesWithLoadBalancerInput{
		Instances: []*elb.Instance{
			{
				InstanceId: aws.String(insid),
			},
		},
		LoadBalancerName: aws.String(elbname),
	}
	resp, err := svc.RegisterInstancesWithLoadBalancer(params)
	Check_Error(err)
	fmt.Println(resp)
	for idx, _ := range resp.Instances {
		inelbs = append(inelbs, *resp.Instances[idx].InstanceId)
	}
	//fmt.Println("add fck", inelbs)
	for _, value := range inelbs {
		if insid == value {
			count++
		}
	}

	if count >= 1 {
		ret = "succ"
	} else {
		ret = "fail"
	}
	return ret

}

func DeRegisterInstanceToElbv2(region string, elbname string, insid string, user string, password string, authfile string) (ret string) {

	svc := Aws_Elbv2_Session(region, user, password, authfile)
	params := &elbv2.DeregisterTargetsInput{
		TargetGroupArn: aws.String(elbname),
		Targets: []*elbv2.TargetDescription{
			{
				Id: aws.String(insid),
			},
		},
	}

	resp, err := svc.DeregisterTargets(params)
	if err != nil {
		fmt.Println(err.Error())
		ret = "fail"
	} else {

		ret = "succ"

	}
	fmt.Println(resp.GoString())
	return ret
}

func DeRegisterInstanceToElb(region string, elbname string, insid string, user string, password string, authfile string) (ret string) {

	var inelbs []string
	var count int64
	svc := Aws_Elb_Session(region, user, password, authfile) //primary
	params := &elb.DeregisterInstancesFromLoadBalancerInput{
		Instances: []*elb.Instance{
			{
				InstanceId: aws.String(insid),
			},
		},
		LoadBalancerName: aws.String(elbname),
	}
	resp, err := svc.DeregisterInstancesFromLoadBalancer(params)
	Check_Error(err)
	fmt.Println(resp)
	for idx, _ := range resp.Instances {
		inelbs = append(inelbs, *resp.Instances[idx].InstanceId)
	}
	//fmt.Println(*resp.Instances[0].InstanceId)
	fmt.Println("fck", inelbs)
	for _, value := range inelbs {
		if insid == value {
			count++
		}
	}
	if count >= 1 {
		ret = "fail"
	} else {
		ret = "succ"
	}
	return ret

}

func Aws_Elbv2_Session(regions string, user string, password string, authfile string) (svc *elbv2.ELBV2) {
	aws_access_key, aws_secret_access_key, token := Get_Access_Key_From_Server(user, password, authfile)
	creds := credentials.NewStaticCredentials(strings.TrimSpace(string(aws_access_key)), strings.TrimSpace(string(aws_secret_access_key)), strings.TrimSpace(string(token)))
	if aws_access_key == "" || aws_secret_access_key == "" {
		fmt.Println("Get Aws auth_file Fail,May Be Password is error")
		os.Exit(100)
	}
	_, err := creds.Get()
	Check_Error(err)
	//fmt.Println("Value of credentials:", v)
	sess, err := session.NewSession()
	Check_Error(err)
	//svc = ec2.New(sess, &aws.Config{Region: aws.String(regions), Credentials: creds})
	svc = elbv2.New(sess, &aws.Config{Region: aws.String(regions), Credentials: creds})
	return svc
}
func Aws_Elb_Session(regions string, user string, password string, authfile string) (svc *elb.ELB) {
	aws_access_key, aws_secret_access_key, token := Get_Access_Key_From_Server(user, password, authfile)
	creds := credentials.NewStaticCredentials(strings.TrimSpace(string(aws_access_key)), strings.TrimSpace(string(aws_secret_access_key)), strings.TrimSpace(string(token)))
	if aws_access_key == "" || aws_secret_access_key == "" {
		fmt.Println("Get Aws auth_file Fail,May Be Password is error")
		os.Exit(100)
	}
	_, err := creds.Get()
	Check_Error(err)
	//fmt.Println("Value of credentials:", v)
	sess, err := session.NewSession()
	Check_Error(err)
	//svc = ec2.New(sess, &aws.Config{Region: aws.String(regions), Credentials: creds})
	svc = elb.New(sess, &aws.Config{Region: aws.String(regions), Credentials: creds})
	return svc
}

func Get_Access_Key_From_Server(user string, password string, authfile string) (aws_access_key, aws_secret_access_key, token string) {

	var aws_access_k, aws_secret_access_k []byte
	token = ""
	aws_access_k, err := sh.Command("curl", "-s", "-u", user+":"+password, "https://seckey.360in.com/secure_key/"+authfile).Command("grep", "-w", "aws_access_key").Command("awk", "-F", "=", "{print $2}").Output()
	aws_secret_access_k, err = sh.Command("curl", "-s", "-u", user+":"+password, "https://seckey.360in.com/secure_key/"+authfile).Command("grep", "-w", "aws_secret_access_key").Command("awk", "-F", "=", "{print $2}").Output()
	Check_Error(err)
	return string(aws_access_k), string(aws_secret_access_k), token
}
func GetInstanceId() (insid string) {

	resp, err := http.Get("http://169.254.169.254/latest/meta-data/instance-id")
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
