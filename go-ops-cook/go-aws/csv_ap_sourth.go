package main

import (
	"archive/zip"
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/configservice"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/codeskyblue/go-sh"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {

	d, _ := time.ParseDuration("-24h")
	//Get_billing_From_S3
	fname := "421921736743-aws-billing-detailed-line-items-with-resources-and-tags-" + time.Now().Add(d*30).Format("2006-01") + ".csv.zip"
	fmt.Println(fname)
	ret := Get_Aws_S3_File("us-west-2", "AKIAIZ4U7FSPO4CQY5WQ", "s/hPxbEpZ2Wa5H8WIn4ypKcaApjv0Rzg3e9YtbMR", "camera360-billing", fname)
	Check_Error(ret)

	//DeCompress_files
	err := DeCompress_File("bills.tgz", "camera360_02.csv")
	Check_Error(err)

	var elbname, unicelbname, hostname, unichostname, volumeid, unicvolumeid []string
	file := Open_Files("./camera360_02.csv")
	elb := Create_Files("./bill_elbs")

	cfile := Create_Files("./bill_public")
	c1file := Create_Files("./bill_common")
	c2file := Create_Files("./bill_volumes")
	cmr_file := Create_Files("./bill_common_mr_instance")
	kinesis_file := Create_Files("./bill_common_kinesis")
	rds_file := Create_Files("./bill_common_rds")
	route53_file := Create_Files("./bill_common_route53")
	datatranfer_instance := Create_Files("./bill_datatransfer_instance")
	s3_file := Create_Files("./bill_common_s3")
	defer datatranfer_instance.Close()
	defer file.Close()
	defer cfile.Close()
	defer c1file.Close()
	defer c2file.Close()
	defer elb.Close()
	defer cmr_file.Close()
	defer kinesis_file.Close()
	defer route53_file.Close()
	defer s3_file.Close()

	reader := csv.NewReader(file)
	for {
		var proces_value []string
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return
		}
		for _, v := range record {
			proces_value = append(proces_value, strings.TrimSpace(v))
		}
		wstring := strings.Join(proces_value, "&")
		//wstring := strings.Join(record, "&")

		result0 := strings.Contains(record[5], "Compute")
		result_mr := strings.Contains(record[5], "MapReduce")
		result_kinesis := strings.Contains(record[5], "Kinesis")
		result_rds := strings.Contains(record[5], "Relational")
		result_route53 := strings.Contains(record[5], "Route")

		result_notifi := strings.Contains(record[5], "Notification")
		result_queue := strings.Contains(record[5], "Queue")
		result_lambda := strings.Contains(record[5], "Lambda")
		result_support := strings.Contains(record[5], "Support")
		result_direct := strings.Contains(record[5], "Direct")
		result_cloudtrail := strings.Contains(record[5], "Trail")
		result_Watch := strings.Contains(record[5], "Watch")

		result_elastic_cache := strings.Contains(record[5], "Amazon ElastiCache")
		result_aws_config := strings.Contains(record[5], "AWS Config")

		result_s3 := strings.Contains(record[5], "Storage")
		result1 := strings.Contains(wstring, "InvoiceTotal")
		result2 := strings.Contains(wstring, "StatementTotal")
		result3 := strings.Contains(wstring, "Rounding")
		result4 := strings.Contains(record[9], "HeavyUsage") //RI Instances Machines
		result5 := strings.Contains(record[9], "BoxUsage")   // On Demand Instances Machines
		result16 := strings.Contains(record[9], "SpotUsage") // Spot Instances Machines
		result6 := strings.Contains(record[9], "Volume")     // All The DIsk
		result7 := strings.Contains(record[9], "DataProcessing")
		result8 := strings.Contains(record[9], "DataTransfer")
		result9 := strings.Contains(record[10], "LoadBalancing")
		result10 := strings.Contains(record[9], "AWS")
		result11 := strings.Contains(record[10], "RunInstances")
		result12 := strings.HasPrefix(record[10], "PublicIP")
		result13 := strings.Contains(record[9], "MetricMonitorUsage")
		result14 := strings.Contains(record[9], "LoadBalancerUsage")
		result15 := strings.Contains(record[9], "EBSOptimized")
		num, _ := strconv.ParseFloat(record[20], 64)

		//Other Service calculer into bill_public
		if num > 0 && (result_elastic_cache || result_aws_config || result_notifi || result_queue || result_lambda || result_support || result_direct || result_cloudtrail || result_Watch) {
			_, err := cfile.WriteString(wstring + "\n")
			if err != nil {
				fmt.Println("Write file error", err)
				return
			}

		}
		//Public Resource Usage
		if num > 0 && !result2 && !result1 && !result3 && record[21] == "" && !result4 && result0 { //generate ResourceId is null and UnBlendedCost is not null
			_, err := cfile.WriteString(wstring + "\n")
			if err != nil {
				fmt.Println("Write file error", err)
				return
			}
		}
		// AWS Compute EBSOptimized EC2 Usage
		if num > 0 && !result2 && !result1 && !result3 && record[21] != "" && record[23] != "" && result0 && result15 {
			_, err := c1file.WriteString(wstring + "\n")
			if err != nil {
				fmt.Println("Write file error", err)
				return
			}
			hostname = append(hostname, strings.TrimSpace(record[23]))
		}
		// AWS Compute MetricMonitorUsage For EC2 with no_tags and no_resource_ID place into bill_public
		if num > 0 && !result2 && !result1 && !result3 && record[21] == "" && record[23] == "" && result0 && result13 {
			_, err := cfile.WriteString(wstring + "\n")
			if err != nil {
				fmt.Println("Write file error", err)
				return
			}
		}
		// AWS Compute MetricMonitorUsage For EC2
		if num > 0 && !result2 && !result1 && !result3 && record[21] != "" && record[23] != "" && result0 && result13 {
			_, err := c1file.WriteString(wstring + "\n")
			if err != nil {
				fmt.Println("Write file error", err)
				return
			}
			hostname = append(hostname, strings.TrimSpace(record[23]))
		}
		// Get All AWS In-Bytes or Out-Bytes
		if num > 0 && !result2 && !result1 && !result3 && record[21] != "" && record[23] != "" && result10 && (result11 || result12) && result0 {
			_, err := datatranfer_instance.WriteString(wstring + "\n")
			if err != nil {
				fmt.Println("Write file error", err)
				return
			}
			hostname = append(hostname, strings.TrimSpace(record[23]))
		}
		// Get All DataTransfer(out_instance)
		if num > 0 && !result2 && !result1 && !result3 && record[21] != "" && record[23] != "" && result8 && !result9 && result0 {
			_, err := datatranfer_instance.WriteString(wstring + "\n")
			if err != nil {
				fmt.Println("Write file error", err)
				return
			}
			hostname = append(hostname, strings.TrimSpace(record[23]))
		}
		//Get AWS Compute LoadBalancerUsage each elb use with no_tags place into bill_public
		if num > 0 && !result2 && !result1 && !result3 && record[21] != "" && record[23] == "" && result0 && result14 {
			_, err := cfile.WriteString(wstring + "\n")
			if err != nil {
				fmt.Println("Write file error", err)
				return
			}
		}
		//Get AWS Compute LoadBalancerUsage each elb use
		if num > 0 && !result2 && !result1 && !result3 && record[21] != "" && record[23] != "" && result0 && result14 {
			_, err := elb.WriteString(wstring + "\n")
			if err != nil {
				fmt.Println("Write file error", err)
				return
			}
			elbname = append(elbname, strings.TrimSpace(record[23]))
		}
		//Get AWS elb In-Bytes or Out-Bytes no_tags place into bill_public
		if num > 0 && !result2 && !result1 && !result3 && record[21] != "" && record[23] == "" && result10 && result0 {
			_, err := cfile.WriteString(wstring + "\n")
			if err != nil {
				fmt.Println("Write file error", err)
				return
			}
		}
		//Get AWS elb In-Bytes or Out-Bytes
		if num > 0 && !result2 && !result1 && !result3 && record[21] != "" && record[23] != "" && result9 && result10 && result0 {
			_, err := elb.WriteString(wstring + "\n")
			if err != nil {
				fmt.Println("Write file error", err)
				return
			}
			elbname = append(elbname, strings.TrimSpace(record[23]))
		}
		//Get All DataTransfer(out_elb)
		if num > 0 && !result2 && !result1 && !result3 && record[21] != "" && record[23] != "" && result8 && result9 && result0 {
			_, err := elb.WriteString(wstring + "\n")
			if err != nil {
				fmt.Println("Write file error", err)
				return
			}
			elbname = append(elbname, strings.TrimSpace(record[23]))
		}
		//Get All ELB_DataProcessing(in and out) items
		if num > 0 && !result2 && !result1 && !result3 && record[21] != "" && record[23] != "" && result7 && result0 {
			_, err := elb.WriteString(wstring + "\n")
			if err != nil {
				fmt.Println("Write file error", err)
				return
			}
			elbname = append(elbname, strings.TrimSpace(record[23]))
		}
		//Get All ELB_DataTransfer(in and out) no_tags_ items place it in bill_public
		if num > 0 && !result2 && !result1 && !result3 && record[21] != "" && record[23] == "" && result8 && result0 {
			_, err := cfile.WriteString(wstring + "\n")
			if err != nil {
				fmt.Println("Write file error", err)
				return
			}
		}
		//Get All ELB_DataProcessing(in and out) no_tags_ items place it in bill_public
		if num > 0 && !result2 && !result1 && !result3 && record[21] != "" && record[23] == "" && result7 && result0 {
			_, err := cfile.WriteString(wstring + "\n")
			if err != nil {
				fmt.Println("Write file error", err)
				return
			}
		}
		//Get AWS MapReduce
		if num > 0 && result_mr {
			_, err := cmr_file.WriteString(wstring + "\n")
			if err != nil {
				fmt.Println("Write file error", err)
				return
			}
		}
		//Get aws Kinesis
		if num > 0 && result_kinesis {
			_, err := kinesis_file.WriteString(wstring + "\n")
			if err != nil {
				fmt.Println("Write file error", err)
				return
			}
		}
		//Get aws rds
		if num > 0 && result_rds {
			_, err := rds_file.WriteString(wstring + "\n")
			if err != nil {
				fmt.Println("Write file error", err)
				return
			}
		}
		//Get aws Route53
		if num > 0 && result_route53 {
			_, err := route53_file.WriteString(wstring + "\n")
			if err != nil {
				fmt.Println("Write file error", err)
				return
			}
		}
		//Get aws S3
		if num > 0 && result_s3 {
			_, err := s3_file.WriteString(wstring + "\n")
			if err != nil {
				fmt.Println("Write file error", err)
				return
			}
		}
		//Get RI and DI Instances and Tag is null_ like_EMR_INSTANCES
		if num > 0 && !result2 && !result1 && !result3 && record[21] != "" && result4 && record[23] == "" && result0 || num > 0 && !result2 && !result1 && !result3 && record[21] != "" && result5 && record[23] == "" && result0 {
			_, err := cmr_file.WriteString(wstring + "\n")
			if err != nil {
				fmt.Println("Write file error", err)
				return
			}
		}
		//Get RI and DI Instances and Tag not null
		if num > 0 && !result2 && !result1 && !result3 && record[21] != "" && result4 && record[23] != "" && result0 || num > 0 && !result2 && !result1 && !result3 && record[21] != "" && result5 && result16 && record[23] != "" && result0 || num > 0 && !result2 && !result1 && !result3 && record[21] != "" && result16 && record[23] != "" && result0 {
			_, err := c1file.WriteString(wstring + "\n")
			if err != nil {
				fmt.Println("Write file error", err)
				return
			}

			hostname = append(hostname, strings.TrimSpace(record[23]))
		}
		// Get All Volumes
		if num > 0 && !result2 && !result1 && !result3 && record[21] != "" && result6 && result0 {
			_, err := c2file.WriteString(wstring + "\n")
			if err != nil {
				fmt.Println("Write file error", err)
				return
			}
			volumeid = append(volumeid, strings.TrimSpace(record[21]))
		}
	}
	unichostname = Rm_duplicate(&hostname)
	unicvolumeid = Rm_duplicate(&volumeid)
	unicelbname = Rm_duplicate(&elbname)

	fmt.Println(unichostname)
	fmt.Println(unicvolumeid)
	fmt.Println(unicelbname)

	//fmt.Println(len(unichostname))
	//fmt.Println(len(unicvolumeid))
	//fmt.Println(len(unicelbname))

	Get_Price_Ec2_With_Tags("./bill_common", unichostname)
	Get_Price_Volume("./bill_volumes", unicvolumeid)
	Get_Price_Volume_With_Instance_Tags("./bill_per_volume", unichostname)
	Get_Price_DataTransfer_With_Tags("./bill_datatransfer_instance", unichostname)
	Combine_Ec2_With_Volume_Price("./bill_combine_ec2_volume", unichostname)
	Get_Price_Elb_With_Tags("./bill_elbs", unicelbname)

	//Process_Result
	Make_Per_Project_Output()
	Make_Sum_OutPut_Into_Csv_File("total.csv")
	Calcute_Aws_Per_Services_Costs("services.csv")
}

//Get The Volume connect with Instance's price
func Combine_Ec2_With_Volume_Price(outfile string, unichostname []string) {
	combine_fd := Create_Files(outfile)
	for _, hstname := range unichostname {
		fdec2 := Open_Files("./bill_pertags_instance")
		fdec3 := Open_Files("./bill_pertags_volume")
		fdec4 := Open_Files("./bill_pertags_datatransfer")
		buf := bufio.NewReader(fdec2)
		buf1 := bufio.NewReader(fdec3)
		buf2 := bufio.NewReader(fdec4)
		for {
			line, err := buf.ReadString('\n')
			if err == io.EOF {
				break
			} else if err != nil {
				fmt.Println("Error:", err)
			}
			lst := strings.Split(line, "&")
			if hstname == lst[0] {
				for {
					line1, err := buf1.ReadString('\n')
					if err == io.EOF {
						break
					} else if err != nil {
						fmt.Println("Error:", err)
					}
					lst1 := strings.Split(line1, "&")
					if hstname == lst1[0] {
						for {
							line2, err := buf2.ReadString('\n')
							if err == io.EOF {
								break
							} else if err != nil {
								fmt.Println("Error:", err)
							}
							lst2 := strings.Split(line2, "&")
							if hstname == lst2[0] {
								lst_num, _ := strconv.ParseFloat(strings.TrimSpace(lst[1]), 64)
								lst1_num, _ := strconv.ParseFloat(strings.TrimSpace(lst1[1]), 64)
								lst2_num, _ := strconv.ParseFloat(strings.TrimSpace(lst2[1]), 64)
								sums := lst_num + lst1_num + lst2_num
								_, err := combine_fd.WriteString(hstname + "&" + strings.TrimSpace(lst[1]) + "&" + strings.TrimSpace(lst1[1]) + "&" + strings.TrimSpace(lst2[1]) + "&" + strconv.FormatFloat(sums, 'f', -1, 64) + "\n")
								if err != nil {
									fmt.Println(err)
								}

							}
						}
					}
				}
			}

		}

	}

}
func Get_Price_Volume(files string, unicvolumeid []string) {
	// Get All The tags for calate sum
	volumeidfile := Create_Files("./bill_per_volume")
	fmt.Println(unicvolumeid)
	for _, hstname := range unicvolumeid {

		var sum float64

		fd, err := os.Open(files)
		defer fd.Close()
		if err != nil {
			fmt.Println(err)
		}
		buf := bufio.NewReader(fd)
		for {
			line, err := buf.ReadString('\n')
			if err == io.EOF {
				break
			} else if err != nil {
				fmt.Println("Error:", err)
			}
			lst := strings.Split(line, "&")
			//fmt.Println(lst[23])
			if hstname == lst[21] {
				costs, _ := strconv.ParseFloat(lst[20], 64)
				sum += costs
			}

		}
		//fmt.Println(hstname, sum)
		insid := strings.TrimSpace(GetInstanceFromVolumeId(hstname))
		if insid == "null" {
			instag := "null"
			_, err = volumeidfile.WriteString(instag + "&" + insid + "&" + hstname + "&" + strconv.FormatFloat(sum, 'f', -1, 64) + "\n")
			if err != nil {
				fmt.Println(err)
			}
		} else if insid == "null_avalible" {
			instag := "null"
			_, err = volumeidfile.WriteString(instag + "&" + insid + "&" + hstname + "&" + strconv.FormatFloat(sum, 'f', -1, 64) + "\n")
			if err != nil {
				fmt.Println(err)
			}

		} else {
			instag := GetTagsFromInstanceid(insid)
			_, err = volumeidfile.WriteString(strings.TrimSpace(instag) + "&" + insid + "&" + hstname + "&" + strconv.FormatFloat(sum, 'f', -1, 64) + "\n")
			if err != nil {
				fmt.Println(err)
			}

		}
		defer fd.Close()
	}
	defer volumeidfile.Close()
}
func Get_Price_Volume_With_Instance_Tags(files string, unichostname []string) {
	ec2file := Create_Files("./bill_pertags_volume")
	for _, hstname := range unichostname {
		var sum float64

		fd, err := os.Open(files)
		if err != nil {
			fmt.Println(err)
		}
		buf := bufio.NewReader(fd)
		for {
			line, err := buf.ReadString('\n')
			if err == io.EOF {
				break
			} else if err != nil {
				fmt.Println("Error:", err)
			}
			lst := strings.Split(line, "&")
			if hstname == lst[0] {
				costs, err := strconv.ParseFloat(strings.TrimSpace(lst[3]), 64)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				sum += costs
			}

		}
		//fmt.Println(hstname, sum)
		_, err = ec2file.WriteString(hstname + "&" + strconv.FormatFloat(sum, 'f', -1, 64) + "\n")
		if err != nil {
			fmt.Println(err)
		}
	}

}
func Get_Price_Elb_With_Tags(files string, unicelbname []string) {
	// Get All The tags for calate sum
	ec2file := Create_Files("./bill_pertags_elb")
	for _, hstname := range unicelbname {
		var sum float64

		fd, err := os.Open(files)
		if err != nil {
			fmt.Println(err)
		}
		buf := bufio.NewReader(fd)
		for {
			line, err := buf.ReadString('\n')
			if err == io.EOF {
				break
			} else if err != nil {
				fmt.Println("Error:", err)
			}
			lst := strings.Split(line, "&")
			//fmt.Println(lst[23])
			if hstname == lst[23] {
				costs, _ := strconv.ParseFloat(strings.TrimSpace(lst[20]), 64)
				sum += costs
			}

		}
		//fmt.Println(hstname, sum)
		_, err = ec2file.WriteString(hstname + "&" + strconv.FormatFloat(sum, 'f', -1, 64) + "\n")
		if err != nil {
			fmt.Println(err)
		}
	}
}

// Get Price Ec2's DataTransfer traffic
func Get_Price_DataTransfer_With_Tags(files string, unichostname []string) {
	// Get All The tags for calate sum
	ec2file := Create_Files("./bill_pertags_datatransfer")
	for _, hstname := range unichostname {
		var sum float64

		fd, err := os.Open(files)
		if err != nil {
			fmt.Println(err)
		}
		buf := bufio.NewReader(fd)
		for {
			line, err := buf.ReadString('\n')
			if err == io.EOF {
				break
			} else if err != nil {
				fmt.Println("Error:", err)
			}
			lst := strings.Split(line, "&")
			//fmt.Println(lst[23])
			if hstname == lst[23] {
				costs, _ := strconv.ParseFloat(strings.TrimSpace(lst[20]), 64)
				sum += costs
			}

		}
		//fmt.Println(hstname, sum)
		_, err = ec2file.WriteString(strings.TrimSpace(hstname) + "&" + strconv.FormatFloat(sum, 'f', -1, 64) + "\n")
		if err != nil {
			fmt.Println(err)
		}
	}
}

//Get Per Tags Ec2 Instance's Price
func Get_Price_Ec2_With_Tags(files string, unichostname []string) {
	// Get All The tags for calate sum
	//fmt.Println(unichostname)
	ec2file := Create_Files("./bill_pertags_instance")
	for _, hstname := range unichostname {
		var sum float64

		fd, err := os.Open(files)
		if err != nil {
			fmt.Println(err)
		}
		buf := bufio.NewReader(fd)
		for {
			line, err := buf.ReadString('\n')
			if err == io.EOF {
				break
			} else if err != nil {
				fmt.Println("Error:", err)
			}
			lst := strings.Split(line, "&")
			if hstname == lst[23] {
				costs, _ := strconv.ParseFloat(strings.TrimSpace(lst[20]), 64)
				//fmt.Println("debug information is:", lst[23], "=====", hstname, costs)
				sum += costs
			}

		}
		//fmt.Println(hstname, sum)
		_, err = ec2file.WriteString(strings.TrimSpace(hstname) + "&" + strconv.FormatFloat(sum, 'f', -1, 64) + "\n")
		if err != nil {
			fmt.Println(err)
		}
	}
}
func Open_Files(filename string) (fd *os.File) {

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	return file
}
func Create_Files(filename string) (fd *os.File) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Create file error", err)
		return
	}
	return file
}

func Rm_duplicate(list *[]string) []string {
	var x []string = []string{}
	for _, i := range *list {
		if len(x) == 0 {
			x = append(x, i)
		} else {
			for k, v := range x {
				if i == v {
					break
				}
				if k == len(x)-1 {
					x = append(x, i)
				}
			}
		}
	}
	return x
}

//Get Instanceid From Volumeid
func GetInstanceFromVolumeId(volumeid string) (instanceid string) {

	var instid string
	//svc := AwsSession("cn-north-1")
	svc := Aws_Ec2_Session("ap-southeast-1", "AKIAIZ4U7FSPO4CQY5WQ", "s/hPxbEpZ2Wa5H8WIn4ypKcaApjv0Rzg3e9YtbMR")
	//svc := Aws_Ec2_Session("cn-north-1", "AKIAOZRIOKXIE4F2MVMQ", "iIp7dFx7QvoGA2g9Y0bytvVqyAmco8g6EH+BmM9H")
	params := &configservice.GetResourceConfigHistoryInput{
		ResourceId:   aws.String(volumeid),
		ResourceType: aws.String("AWS::EC2::Volume"),
	}

	resp, err := svc.GetResourceConfigHistory(params)

	if err != nil {
		//instid = "null"
		//fmt.Println(err.Error())
		//fmt.Println("will in secondary aws china to get it")
		svc := Aws_Ec2_Session("us-east-1", "AKIAIZ4U7FSPO4CQY5WQ", "s/hPxbEpZ2Wa5H8WIn4ypKcaApjv0Rzg3e9YtbMR")
		params := &configservice.GetResourceConfigHistoryInput{
			ResourceId:   aws.String(volumeid),
			ResourceType: aws.String("AWS::EC2::Volume"),
		}

		resp, err := svc.GetResourceConfigHistory(params)
		for idx, _ := range resp.ConfigurationItems {
			if len(resp.ConfigurationItems[idx].Relationships) == 0 {
				instid = "null_avalible"
			} else {
				for _, inst := range resp.ConfigurationItems[idx].Relationships {
					instid = *inst.ResourceId

				}
			}
		}

		if err != nil {
			//fmt.Println("two aws nodes can't find it")
			instid = "null"
			//fmt.Println(err.Error())
		}

	}

	for idx, _ := range resp.ConfigurationItems {
		if len(resp.ConfigurationItems[idx].Relationships) == 0 {
			instid = "null_avalible"
		} else {
			for _, inst := range resp.ConfigurationItems[idx].Relationships {
				instid = *inst.ResourceId

			}
		}
	}
	return instid

}

//Get Tags From Instanceid
func GetTagsFromInstanceid(instanceid string) (ctags string) {

	//svc := AwsSession("cn-north-1")
	svc := Aws_Ec2_Session("ap-southeast-1", "AKIAIZ4U7FSPO4CQY5WQ", "s/hPxbEpZ2Wa5H8WIn4ypKcaApjv0Rzg3e9YtbMR")
	//svc := Aws_Ec2_Session("cn-north-1", "AKIAOZRIOKXIE4F2MVMQ", "iIp7dFx7QvoGA2g9Y0bytvVqyAmco8g6EH+BmM9H")
	params := &configservice.GetResourceConfigHistoryInput{
		ResourceId:   aws.String(instanceid),
		ResourceType: aws.String("AWS::EC2::Instance"),
	}
	resp, err := svc.GetResourceConfigHistory(params)

	if err != nil {
		//ctags = "null"
		//fmt.Println(err.Error())
		//return
		//fmt.Println("will in secondary aws china to get it tags")
		svc := Aws_Ec2_Session("us-east-1", "AKIAIZ4U7FSPO4CQY5WQ", "s/hPxbEpZ2Wa5H8WIn4ypKcaApjv0Rzg3e9YtbMR")
		params := &configservice.GetResourceConfigHistoryInput{
			ResourceId:   aws.String(instanceid),
			ResourceType: aws.String("AWS::EC2::Instance"),
		}

		resp, err := svc.GetResourceConfigHistory(params)
		//fmt.Println(resp)
		for idx, _ := range resp.ConfigurationItems {
			//fmt.Println("ctags inner is", *resp.ConfigurationItems[idx].ConfigurationItemStatus)
			if *resp.ConfigurationItems[idx].ConfigurationItemStatus == "ResourceDiscovered" || *resp.ConfigurationItems[idx].ConfigurationItemStatus == "OK" {
				//ctags = *resp.ConfigurationItems[idx].Tags["Name"]
				//fmt.Println(resp.ConfigurationItems[idx].Tags)
				_, ok := resp.ConfigurationItems[idx].Tags["Name"]
				if ok {
					ctags = *resp.ConfigurationItems[idx].Tags["Name"]
				} else {
					//fmt.Println("is mapreduce's instance_fuck")
					ctags = "null_mapreduce"
				}
			}

		}
		if err != nil {
			//fmt.Println("two aws nodes can't find it tags")
			ctags = "null"
			//fmt.Println(err.Error())
		}
	}
	//fmt.Println(resp)
	for idx, _ := range resp.ConfigurationItems {
		//fmt.Println("ctags is", *resp.ConfigurationItems[idx].ConfigurationItemStatus)
		if *resp.ConfigurationItems[idx].ConfigurationItemStatus == "ResourceDiscovered" || *resp.ConfigurationItems[idx].ConfigurationItemStatus == "OK" {
			//ctags = *resp.ConfigurationItems[idx].Tags["Name"]
			//fmt.Println(resp.ConfigurationItems[idx].Tags)
			_, ok := resp.ConfigurationItems[idx].Tags["Name"]
			if ok {
				ctags = *resp.ConfigurationItems[idx].Tags["Name"]
				//fmt.Println("ctags 1 is:", ctags)
			} else {
				//fmt.Println("is mapreduce's instance_fuck")
				ctags = "null_mapreduce"
				//fmt.Println("ctags 2 is:", ctags)
			}
		}
	}

	return ctags

}

func Aws_Ec2_Session(regions string, aws_access_key string, aws_secret_access_key string) (svc *configservice.ConfigService) {

	token := ""
	creds := credentials.NewStaticCredentials(aws_access_key, aws_secret_access_key, token)
	_, err := creds.Get()
	Check_Error(err)
	//fmt.Println("Value of credentials:", v)
	sess, err := session.NewSession()
	Check_Error(err)

	svc = configservice.New(sess, &aws.Config{Region: aws.String(regions), Credentials: creds})
	return svc
}

/*
func AwsSession(regions string) (svc *configservice.ConfigService) {
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("failed to create session,", err)
		return
	}
	svc = configservice.New(sess, &aws.Config{Region: aws.String(regions)})
	//svc := configservice.New(sess, &aws.Config{Region: aws.String("cn-north-1")})
	return svc
}
*/
func Check_Error(err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
}

func Get_Aws_S3_File(region string, key string, secert string, bucket string, filename string) (err error) {
	svc := Aws_S3_Session(region, key, secert) //primary

	result, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		//Bucket: aws.String("camera360-billing"),
		Key: aws.String(filename),
		//Key:    aws.String("421921736743-aws-billing-detailed-line-items-with-resources-and-tags-2017-03.csv.zip"),
	})
	Check_Error(err)

	file, err := os.Create("bills.tgz")
	Check_Error(err)

	if _, err := io.Copy(file, result.Body); err != nil {
		fmt.Println("Failed to copy object to file", err)
	}
	result.Body.Close()
	file.Close()
	return err
}
func Aws_S3_Session(regions string, aws_access_key string, aws_secret_access_key string) (svc *s3.S3) {

	token := ""
	creds := credentials.NewStaticCredentials(aws_access_key, aws_secret_access_key, token)
	_, err := creds.Get()
	Check_Error(err)
	//fmt.Println("Value of credentials:", v)
	sess, err := session.NewSession()
	Check_Error(err)

	//svc = ec2.New(sess, &aws.Config{Region: aws.String(regions), Credentials: creds})
	svc = s3.New(sess, &aws.Config{Region: aws.String(regions), Credentials: creds})
	return svc
}
func DeCompress_File(input string, output string) (err error) {

	r, err := zip.OpenReader(input)
	if err != nil {
		log.Fatal(err)
	}

	w, err := os.Create(output)
	if err != nil {
		fmt.Println(err)
	}
	defer w.Close()
	defer r.Close()

	for _, f := range r.File {
		fmt.Printf("File name %s:\n", f.Name)
		rc, err := f.Open()
		if err != nil {
			log.Fatal(err)
		}

		_, err = io.CopyN(w, rc, int64(f.UncompressedSize64))

		if err != nil {
			log.Fatal(err)
		}
		rc.Close()
	}
	return err
}
func Calcute_Aws_Per_Services_Costs(filename string) {
	mvalue := map[string]string{
		"kinesis":   "./bill_common_kinesis",
		"mapreduce": "./bill_common_mr_instance",
		"rds":       "./bill_common_rds",
		"route53":   "./bill_common_route53",
		"s3":        "./bill_common_s3",
		"bpublic":   "./bill_public",
	}

	//ec2_mvalue := map[string]map[string]string{
	//	"public": {"./bill_pertags_public_ec2": "bill_pertags_public_elb"},
	//}
	//sum_total, _ := Calcute_Per_Project_Costs(ec2_mvalue)
	//fmt.Println("sum_total:", sum_total)
	//Get_mapreduce_disk_costs
	bill_volumes, err := sh.Command("cat", "./bill_volumes").Command("awk", "-F", "&", "BEGIN{sum=0}{sum+=$21}END{print sum}").Output()
	Check_Error(err)
	bill_volumes_math, err := strconv.ParseFloat(strings.TrimSpace(string(bill_volumes)), 64)
	Check_Error(err)
	bill_pertags_volumes, err := sh.Command("cat", "./bill_pertags_volume").Command("awk", "-F", "&", "BEGIN{sum=0}{sum+=$2}END{print sum}").Output()
	Check_Error(err)
	bill_pertags_volumes_math, err := strconv.ParseFloat(strings.TrimSpace(string(bill_pertags_volumes)), 64)
	Check_Error(err)
	mapreduce_volumes := bill_volumes_math - bill_pertags_volumes_math

	ofd := Create_Files(filename)
	defer ofd.Close()
	w := csv.NewWriter(ofd)
	w.Write([]string{"Public_Service", "All_Costs"})
	//w.Write([]string{"ec2_total", strconv.FormatFloat(sum_total, 'f', -1, 64)})

	for svs, fname := range mvalue {
		costs, err := sh.Command("cat", fname).Command("awk", "-F", "&", "BEGIN{sum=0}{sum+=$21}END{print sum}").Output()
		Check_Error(err)
		costs_math, err := strconv.ParseFloat(strings.TrimSpace(string(costs)), 64)
		Check_Error(err)
		if svs == "mapreduce" {
			opt := costs_math + mapreduce_volumes
			w.Write([]string{svs, strconv.FormatFloat(opt, 'f', -1, 64)})
			//fmt.Println(svs, costs_math+mapreduce_volumes)
			w.Flush()
		} else {
			w.Write([]string{svs, strconv.FormatFloat(costs_math, 'f', -1, 64)})
			fmt.Println(svs, costs_math)
			w.Flush()
		}
	}
}

func Make_Sum_OutPut_Into_Csv_File(filename string) {
	mvalue := map[string]map[string]string{
		"ads":       {"./bill_pertags_ads_ec2": "./bill_pertags_ads_elb"},
		"platform":  {"./bill_pertags_platform_ec2": "bill_pertags_platform_elb"},
		"phototask": {"./bill_pertags_phototask_ec2": "bill_pertags_phototask_elb"},
		"public":    {"./bill_pertags_public_ec2": "bill_pertags_public_elb"},
	}
	sum_total, all_element := Calcute_Per_Project_Costs(mvalue)

	fmt.Println("sum_total is:", sum_total)
	fmt.Println("all_element is:", all_element)

	ofd := Create_Files(filename)
	defer ofd.Close()
	w := csv.NewWriter(ofd)
	w.Write([]string{"ProjectName", "Compute_Costs", "Storage_Costs", "Network_Costs", "Elb_Costs", "All_Costs", "Percent"})
	for _, v := range all_element {
		st := strings.Split(v, "&")
		fmt.Println("st is:", st)
		tp, err := strconv.ParseFloat(strings.TrimSpace(st[5]), 64)
		Check_Error(err)
		//fmt.Println(st[0], tp/sum_total)
		rs := strconv.FormatFloat(tp/sum_total, 'f', -1, 64)
		st = append(st, rs)
		fmt.Println(st)
		w.Write([]string{st[0], st[1], st[2], st[3], st[4], st[5], st[6]})
		w.Flush()
	}
}

func Calcute_Per_Project_Costs(mvalue map[string]map[string]string) (sum float64, all_element []string) {
	var sum_total []string
	for pjt, vues := range mvalue {
		for ec2, elb := range vues {
			ec2_compute, err := sh.Command("cat", ec2).Command("awk", "-F", "&", "BEGIN{sum=0}{sum+=$2}END{print sum}").Output()
			ec2_ebs, err := sh.Command("cat", ec2).Command("awk", "-F", "&", "BEGIN{sum=0}{sum+=$3}END{print sum}").Output()
			ec2_traffic, err := sh.Command("cat", ec2).Command("awk", "-F", "&", "BEGIN{sum=0}{sum+=$4}END{print sum}").Output()
			elb_total, err := sh.Command("cat", elb).Command("awk", "-F", "&", "BEGIN{sum=0}{sum+=$2}END{print sum}").Output()

			ec2s := string(ec2_compute)
			ec2_ebss := string(ec2_ebs)
			ec2_traffics := string(ec2_traffic)
			elbs := string(elb_total)

			ec2f, err := strconv.ParseFloat(strings.TrimSpace(ec2s), 64)
			Check_Error(err)
			ec2_ebssf, err := strconv.ParseFloat(strings.TrimSpace(ec2_ebss), 64)
			Check_Error(err)
			ec2_trafficsf, err := strconv.ParseFloat(strings.TrimSpace(ec2_traffics), 64)
			Check_Error(err)
			elb2f, err := strconv.ParseFloat(strings.TrimSpace(elbs), 64)
			Check_Error(err)

			total := strconv.FormatFloat(ec2f+elb2f+ec2_ebssf+ec2_trafficsf, 'f', -1, 64)
			sum_total = append(sum_total, total)
			ret, err := sh.Command("echo", strings.TrimSpace(pjt)+"&"+strings.TrimSpace(ec2s)+"&"+strings.TrimSpace(ec2_ebss)+"&"+strings.TrimSpace(ec2_traffics)+"&"+strings.TrimSpace(elbs)+"&"+strings.TrimSpace(total)).Output()
			//ret, err := sh.Command("echo", strings.TrimSpace(pjt)+"&"+strings.TrimSpace(ec2s)+"&"+strings.TrimSpace(elbs)+"&"+strings.TrimSpace(total)).Output()
			Check_Error(err)
			//fmt.Println(strings.TrimSpace(string(ret)))
			tp := strings.TrimSpace(string(ret))
			all_element = append(all_element, tp)
		}
	}

	for _, value := range sum_total {
		vs, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
		Check_Error(err)
		sum += vs
	}

	//fmt.Println(sum_total)
	//fmt.Println(all_element)
	return sum, all_element
}

func Make_Per_Project_Output() {

	//For_Split_Combine_Files_for_per_project_calcute
	err := Split_Combine_Input_File("./bill_combine_ec2_volume", "./bill_pertags_ads_ec2", "./bill_pertags_platform_ec2", "./bill_pertags_phototask_ec2", "bill_pertags_public_ec2")
	Check_Error(err)
	err = Split_Elb_Input_File("./bill_pertags_elb", "./bill_pertags_ads_elb", "./bill_pertags_platform_elb", "./bill_pertags_phototask_elb", "bill_pertags_public_elb")
	Check_Error(err)

	//For_ads_Elb_output
	err = Soft_Pertags_Input_File("./bill_pertags_ads_elb", "./bill_pertags_ads_elb_sorted")
	Check_Error(err)
	err = Generate_Pertags_Csv_Output("./bill_pertags_ads_elb_sorted", "bill_pertags_ads_elb.csv")
	Check_Error(err)
	//For_platform_Elb_output
	err = Soft_Pertags_Input_File("./bill_pertags_platform_elb", "./bill_pertags_platform_elb_sorted")
	Check_Error(err)
	err = Generate_Pertags_Csv_Output("./bill_pertags_platform_elb_sorted", "bill_pertags_platform_elb.csv")
	Check_Error(err)
	//For_platform_Elb_output
	err = Soft_Pertags_Input_File("./bill_pertags_phototask_elb", "./bill_pertags_phototask_elb_sorted")
	Check_Error(err)
	err = Generate_Pertags_Csv_Output("./bill_pertags_phototask_elb_sorted", "bill_pertags_phototask_elb.csv")
	Check_Error(err)
	//For_platform_Elb_output
	err = Soft_Pertags_Input_File("./bill_pertags_public_elb", "./bill_pertags_public_elb_sorted")
	Check_Error(err)
	err = Generate_Pertags_Csv_Output("./bill_pertags_public_elb_sorted", "bill_pertags_public_elb.csv")
	Check_Error(err)

	//For_ads_Combine_output
	err = Soft_Combine_Input_File("./bill_pertags_ads_ec2", "./bill_pertags_ads_ec2_sorted")
	Check_Error(err)
	err = Generate_Combine_Csv_Output("./bill_pertags_ads_ec2_sorted", "./bill_pertags_ads_ec2.csv")
	Check_Error(err)
	//For_platform_Combine_output
	err = Soft_Combine_Input_File("./bill_pertags_platform_ec2", "./bill_pertags_platform_ec2_sorted")
	Check_Error(err)
	err = Generate_Combine_Csv_Output("./bill_pertags_platform_ec2_sorted", "./bill_pertags_platform_ec2.csv")
	Check_Error(err)
	//For_public_Combine_output
	err = Soft_Combine_Input_File("./bill_pertags_public_ec2", "./bill_pertags_public_ec2_sorted")
	Check_Error(err)
	err = Generate_Combine_Csv_Output("./bill_pertags_public_ec2_sorted", "./bill_pertags_public_ec2.csv")
	Check_Error(err)
	//For_phototask_Combine_output
	err = Soft_Combine_Input_File("./bill_pertags_phototask_ec2", "./bill_pertags_phototask_ec2_sorted")
	Check_Error(err)
	err = Generate_Combine_Csv_Output("./bill_pertags_phototask_ec2_sorted", "./bill_pertags_phototask_ec2.csv")
	Check_Error(err)
}

// to_modify_use_self
func Split_Combine_Input_File(filename string, ads, platform, phototask, public string) (err error) {
	err = sh.Command("grep", "-E", "ads|facebook|smallappsapi", filename).Command("grep", "-v", "adv").WriteStdout(ads)
	err = sh.Command("grep", "-E", "hotpot|ms_|phototask", filename).WriteStdout(phototask)
	err = sh.Command("grep", "-E", "mixweb|platform|adv", filename).WriteStdout(platform)
	err = sh.Command("grep", "-Ev", "ads|adv|facebook|smallappsapi|phototask|hotpot|ms_|platform|mixweb", filename).WriteStdout(public)
	return err
}
func Split_Elb_Input_File(filename string, ads, platform, phototask, public string) (err error) {
	err = sh.Command("grep", "-E", "ads", filename).Command("grep", "-v", "adv").WriteStdout(ads)
	err = sh.Command("grep", "-E", "phototask|searchapi", filename).WriteStdout(phototask)
	err = sh.Command("grep", "-E", "platform|adv", filename).WriteStdout(platform)
	err = sh.Command("grep", "-Ev", "ads|adv|phototask|platform|searchapi", filename).WriteStdout(public)
	return err
}
func Generate_Combine_Csv_Output(filename string, outfilename string) (err error) {

	fd := Open_Files(filename)
	ofd := Create_Files(outfilename)
	defer fd.Close()
	defer ofd.Close()
	buf := bufio.NewReader(fd)
	w := csv.NewWriter(ofd)
	w.Write([]string{"Name", "Compute_Costs", "Storage_Costs", "Network_Costs", "All_Costs"})
	for {
		line, err := buf.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error", err)
		}
		lst := strings.Split(strings.TrimSpace(line), "&")
		//fmt.Println(lst)
		w.Write([]string{lst[0], lst[1], lst[2], lst[3], lst[4]})
		w.Flush()

	}
	return err
}
func Generate_Pertags_Csv_Output(filename string, outfilename string) (err error) {

	fd := Open_Files(filename)
	ofd := Create_Files(outfilename)
	defer fd.Close()
	defer ofd.Close()
	buf := bufio.NewReader(fd)
	w := csv.NewWriter(ofd)
	w.Write([]string{"Name", "Costs"})
	for {
		line, err := buf.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error", err)
		}
		lst := strings.Split(strings.TrimSpace(line), "&")
		//fmt.Println(lst)
		w.Write([]string{lst[0], lst[1]})
		w.Flush()

	}
	return err
}

func Soft_Pertags_Input_File(filename string, outfilename string) (err error) {
	err = sh.Command("sort", "-t", "&", "-k2nr", filename).WriteStdout(outfilename)
	return err
}
func Soft_Combine_Input_File(filename string, outfilename string) (err error) {
	err = sh.Command("sort", "-t", "&", "-k5nr", filename).WriteStdout(outfilename)
	return err
}
