package main

import (
	"archive/tar"
	"compress/gzip"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/codeskyblue/go-sh"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {

	var err error
	var dt = "nidaye"
	var user *string
	user = &dt
	action := flag.String("action", "", "action have two select, one is: update, other is deploy")
	s3_bucket := flag.String("bucket", "en-ap-deploy", "default S3 bucket is: en-ap-deploy,cn-north-1 S3 bucket is: deploy")
	s3_region := flag.String("region", "ap-southeast-1", "s3 region: ap-southeast-1,cn-north-1,us-east-1 ?")

	public_lib_filename := flag.String("source_lib_filename", "cn-bj-aws-lib.camera360.com.last", "php public libs,cn-north-1 or ap-southeast-1 is the same")
	public_lib_desdir := flag.String("dest_lib_dir", "/home/worker/data/www/lib", "the destation for lib deploy")

	project_filename := flag.String("source_project_filename", "xxxx.last", "the program on s3_bucket")
	project_desdir := flag.String("project_dest_domain", "", "the program deploy to this dir")
	runtime_dir := flag.String("runtime_dir", "/home/worker/data/www/runtime/xxx", "the app_log_dir,need to ask devloper")
	logstash_conf := flag.String("logstash_conf", "logstash_ap_conf.tar.gz", "default logstash conf is ap-southeast-1,other will be create need")
	//user := flag.String("user", "", "user for get username")
	password := flag.String("password", "", "password for auth username")
	start_nginx := flag.String("start_nginx", "false", "if start nginx,two choice true or false,default is false")
	start_fpm := flag.String("start_fpm", "false", "if start php-fpm,two choice true or false,default is false")
	auth_file := flag.String("auth_file", "", "authfile for auth,this param is must")
	get_public_lib := flag.String("get_public_lib", "false", "if install public lib,default is false,if need set 'true'")

	flag.Parse()

	fmt.Println("Usage: ./update_php_fpm_code -h")
	fmt.Println()
	fmt.Println("Usage For PHP-FPM Update: ./update_php_fpm_code -get_public_lib true -auth_file en-ap-auth.txt  -password xx -action update -bucket en-ap-deploy -region ap-southeast-1  -source_project_filename exp.360in.com.last -project_dest_domain exp.360in.com -start_nginx true -start_fpm true")
	fmt.Println("Usage For PHP-FPM Deploy: ./update_php_fpm_code -get_public_lib false -auth_file en-ap-auth.txt   -password xx -action deploy -bucket en-ap-deploy -region ap-southeast-1 -source_project_filename exp.360in.com.last -project_dest_domain exp.360in.com -runtime_dir ads -start_nginx true -start_fpm true -logstash_conf logstash_ap_conf.tar.gz")

	fmt.Println()
	fmt.Println("Usage For Msf-Swoole Update: ./update_php_fpm_code -get_public_lib true -auth_file en-ap-auth.txt  -password xx -action update -bucket en-ap-deploy -region ap-southeast-1  -source_project_filename exp.360in.com.last -project_dest_domain exp.360in.com")
	fmt.Println("Usage For Msf-Swoole Deploy: ./update_php_fpm_code -get_public_lib false -auth_file en-ap-auth.txt   -password xx -action deploy -bucket en-ap-deploy -region ap-southeast-1 -source_project_filename exp.360in.com.last -project_dest_domain exp.360in.com -runtime_dir ads -logstash_conf logstash_ap_conf.tar.gz")
	fmt.Println("Usage For Rollback: ./update_php_fpm_code -action rollback -get_public_lib true -project_dest_domain effectapi.camera360.com")
	fmt.Println()

	if *action == "rollback" {

		fmt.Println("Will Print The Input Params!!!")
		fmt.Println("action:", *action)
		fmt.Println("get_public_lib:", *get_public_lib)
		fmt.Println("project_desdir:", *project_desdir)
		fmt.Println()
		// Roll Back For Public Libs
		if *get_public_lib == "true" {
			if _, err := os.Stat("/home/worker/data/www/lib_rollbk"); err == nil {
				err = Copy_File_Dir_To_Dest("/home/worker/data/www/lib_rollbk", "/home/worker/data/www/lib")
				Check_Error(err)
			}
		}
		// Roll Back For Project Files
		if *project_desdir != "" {
			if _, err := os.Stat("/home/worker/data/www/" + *project_desdir + "_rollbk"); err == nil {
				err = Copy_File_Dir_To_Dest("/home/worker/data/www/"+*project_desdir+"_rollbk", "/home/worker/data/www/"+*project_desdir)
				Check_Error(err)
			}

		}
	}

	if *action == "update" {
		if *password == "" || *auth_file == "" {
			fmt.Println("Must indicate password,auth_file to get authorize")
			os.Exit(100)
		}
		fmt.Println("Will Print The Input Params!!!")
		fmt.Println("action:", *action)
		fmt.Println("s3_bucket:", *s3_bucket)
		fmt.Println("s3_region:", *s3_region)
		fmt.Println("public_lib_filename:", *public_lib_filename)
		fmt.Println("public_lib_desdir:", *public_lib_desdir)
		fmt.Println("project_filename:", *project_filename)
		fmt.Println("project_desdir:", *project_desdir)
		fmt.Println("runtime_dir:", *runtime_dir)
		fmt.Println("logstash_conf:", *logstash_conf)
		//fmt.Println("user:", *user)
		fmt.Println("password:", *password)
		fmt.Println("start_nginx:", *start_nginx)
		fmt.Println("start_fpm:", *start_fpm)
		fmt.Println("auth_file:", *auth_file)
		fmt.Println("get_public_lib:", *get_public_lib)
		fmt.Println()

		//Get_Php-Fpm Public lib files
		if *get_public_lib == "true" {
			fmt.Println("Start get public lib from S3")
			err = Get_Aws_S3_File(*s3_region, *s3_bucket, *public_lib_filename, *public_lib_filename, *user, *password, *auth_file)
			Check_Error(err)
			fmt.Println("End get public lib from S3")
		}
		//DeCompress PHP-FPM Public Lib Files
		if *get_public_lib == "true" {
			fmt.Println("Start DeCompress Public lib")
			err = DeCompress_TarFile(*public_lib_filename, "/dev/shm/")
			Check_Error(err)
			//Delete_Public .git
			Delete_Git_Resp("cn-bj-aws-lib.camera360.com")
			fmt.Println("End DeCompress Public lib")
		}
		//Copy and fugai Public lib files to destination dir
		if *get_public_lib == "true" {
			fmt.Println("Start_Copy Public Libs to Destination dir")
			//Use_For_Rollback
			if _, err := os.Stat("/home/worker/data/www/lib"); err == nil {
				err = Copy_File_Dir_To_Dest("/home/worker/data/www/lib", "/home/worker/data/www/lib_rollbk")
				Check_Error(err)
			}
			err = Copy_File_Dir_To_Dest("/dev/shm/cn-bj-aws-lib.camera360.com", "/home/worker/data/www/lib")
			Check_Error(err)
			fmt.Println("End_Copy Public Libs to Destination dir")
		}

		//Get Php-Fprm Ptroject Codes
		if *project_desdir != "" && *project_filename != "" {
			fmt.Println("Start get Project Codes from S3")
			err = Get_Aws_S3_File(*s3_region, *s3_bucket, *project_filename, *project_filename, *user, *password, *auth_file)
			Check_Error(err)
			fmt.Println("End get Project Codes from S3")
		}

		//DeCompress PHP-FPM Project Files
		if *project_desdir != "" && *project_filename != "" {
			fmt.Println("Start DeCompress Project Codes ")
			err = DeCompress_TarFile(*project_filename, "/dev/shm/")
			//Delete Project .git
			Delete_Git_Resp(*project_desdir)
			Check_Error(err)
			fmt.Println("End DeCompress Project Codes ")
		}

		//Copy and fugai Project files to destination dir
		if *project_desdir != "" && *project_filename != "" {
			fmt.Println("Start_Copy Project Codes to Destination dir")
			//Use_For_Rollback_project
			if _, err := os.Stat("/home/worker/data/www/" + *project_desdir); err == nil {
				err = Copy_File_Dir_To_Dest("/home/worker/data/www/"+*project_desdir, "/home/worker/data/www/"+*project_desdir+"_rollbk")
				Check_Error(err)
			}
			err = Copy_File_Dir_To_Dest("/dev/shm/"+*project_desdir, "/home/worker/data/www/"+*project_desdir)
			Check_Error(err)
			fmt.Println("End_Copy Project Codes to Destination dir")
		}

		//Delete Source tar or dir after deploy
		Delete_Deployed_Code(*project_filename, *project_desdir, *public_lib_filename, "cn-bj-aws-lib.camera360.com")

		//for avaliable to start nginx again
		if *start_nginx == "true" {
			err = Start_Nginx()
			Check_Error(err)
		}
		//for avaliable to start php-fpm
		if *start_fpm == "true" {
			err = Start_Php_Fpm()
			Check_Error(err)
		}
	}

	if *action == "deploy" {
		if *logstash_conf == "" || *runtime_dir == "" || *password == "" || *auth_file == "" {
			fmt.Println("Must indicate logstash_conf,runtime_dir,user,password,auth_file")
			os.Exit(100)
		}
		fmt.Println()
		fmt.Println("Will Print The Input Params!!!")
		fmt.Println("action:", *action)
		fmt.Println("s3_bucket:", *s3_bucket)
		fmt.Println("s3_region:", *s3_region)
		fmt.Println("public_lib_filename:", *public_lib_filename)
		fmt.Println("public_lib_desdir:", *public_lib_desdir)
		fmt.Println("project_filename:", *project_filename)
		fmt.Println("project_desdir:", *project_desdir)
		fmt.Println("runtime_dir:", *runtime_dir)
		fmt.Println("logstash_conf:", *logstash_conf)
		//fmt.Println("user:", *user)
		fmt.Println("password:", *password)
		fmt.Println("start_nginx:", *start_nginx)
		fmt.Println("start_fpm:", *start_fpm)
		fmt.Println("auth_file:", *auth_file)
		fmt.Println("get_public_lib:", *get_public_lib)
		fmt.Println()

		//if have old code then delete it
		if *action == "deploy" {
			fmt.Println("Start Del Old codes")
			Delete_Already_Deployed_Code(*project_desdir, "lib")
			fmt.Println("End Del Old codes")
		}

		//Get_Php-Fpm Public lib files
		if *get_public_lib == "true" {
			fmt.Println("Start get public lib from S3")
			err = Get_Aws_S3_File(*s3_region, *s3_bucket, *public_lib_filename, *public_lib_filename, *user, *password, *auth_file)
			Check_Error(err)
			fmt.Println("End get public lib from S3")
		}

		//DeCompress PHP-FPM Public Lib Files
		if *get_public_lib == "true" {
			fmt.Println("Start DeCompress Public lib")
			err = DeCompress_TarFile(*public_lib_filename, "/dev/shm/")
			Check_Error(err)
			//Delete_Public .git
			Delete_Git_Resp("cn-bj-aws-lib.camera360.com")
			fmt.Println("End DeCompress Public lib")
		}

		//Copy and fugai Public lib files to destination dir
		if *get_public_lib == "true" {
			fmt.Println("Start_Copy Public Libs to Destination dir")
			err = Copy_File_Dir_To_Dest("/dev/shm/cn-bj-aws-lib.camera360.com", "/home/worker/data/www/lib")
			Check_Error(err)
			fmt.Println("End_Copy Public Libs to Destination dir")
		}

		//Get Php-Fprm Ptroject Codes
		if *project_desdir != "" && *project_filename != "" {
			fmt.Println("Start get Project Codes from S3")
			err = Get_Aws_S3_File(*s3_region, *s3_bucket, *project_filename, *project_filename, *user, *password, *auth_file)
			Check_Error(err)
			fmt.Println("End get Project Codes from S3")
		}

		//DeCompress PHP-FPM Project Files
		if *project_desdir != "" && *project_filename != "" {
			fmt.Println("Start DeCompress Project Codes ")
			err = DeCompress_TarFile(*project_filename, "/dev/shm/")
			Check_Error(err)
			//Delete Project .git
			Delete_Git_Resp(*project_desdir)
			fmt.Println("End DeCompress Project Codes ")
		}

		//Copy and fugai Project files to destination dir
		if *project_desdir != "" && *project_filename != "" {
			fmt.Println("Start_Copy Project Codes to Destination dir")
			err = Copy_File_Dir_To_Dest("/dev/shm/"+*project_desdir, "/home/worker/data/www/"+*project_desdir)
			Check_Error(err)
			fmt.Println("End_Copy Project Codes to Destination dir")
			fmt.Println("Start Create runtime_dir")
			//make runtime_dir for app store data
			err = os.MkdirAll("/home/worker/data/www/runtime/"+*runtime_dir, 0755)
			Check_Error(err)
			fmt.Println("End Create runtime_dir")
		}

		// Get logstash_conf to deploy
		if *logstash_conf != "" {
			fmt.Println("Start Get Logstash config from S3")
			err = Get_Aws_S3_File(*s3_region, *s3_bucket, *logstash_conf, *logstash_conf, *user, *password, *auth_file)
			Check_Error(err)
			fmt.Println("End Get Logstash config from S3")
		}

		if *logstash_conf != "" {
			//Deploy Logstash Conf to Dest
			fmt.Println("Start Deploy logstash_conf to Dest dir")
			err = sh.Command("tar", "xfP", "/dev/shm/"+*logstash_conf).Run()
			Check_Error(err)
			fmt.Println("End Deploy logstash_conf to Dest dir")
		}

		// Delete code when deploy is over
		Delete_Deployed_Code(*project_filename, *project_desdir, *public_lib_filename, "cn-bj-aws-lib.camera360.com")

		//for avaliable to start nginx again
		if *start_nginx == "true" {
			err = Start_Nginx()
			Check_Error(err)
		}
		//for avaliable to start php-fpm
		if *start_fpm == "true" {
			err = Start_Php_Fpm()
			Check_Error(err)
		}
	}
}

func Copy_File_Dir_To_Dest(source_dir, dest_dir string) (err error) {

	// check if the source dir exist
	src, err := os.Stat(source_dir)
	if err != nil {
		panic(err)
	}

	if !src.IsDir() {
		fmt.Println("Source is not a directory")
		os.Exit(1)
	}

	// create the destination directory
	fmt.Println("Destination :" + dest_dir)

	err = CopyDir(source_dir, dest_dir)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Directory copied")
	}
	return err

}

func DeCompress_TarFile(input, dst string) (err error) {
	r, err := os.Open("/dev/shm/" + input)
	Check_Error(err)
	gzr, err := gzip.NewReader(r)
	defer gzr.Close()
	if err != nil {
		return err
	}

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()

		switch {

		// if no more files are found return
		case err == io.EOF:
			return nil

		// return any other error
		case err != nil:
			return err

		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		// the target location where the dir/file should be created
		target := filepath.Join(dst, header.Name)

		// the following switch could also be done using fi.Mode(), not sure if there
		// a benefit of using one vs. the other.
		// fi := header.FileInfo()

		// check the file type
		switch header.Typeflag {

		// if its a dir and it doesn't exist create it
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}

		// if it's a file create it
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			defer f.Close()

			// copy over contents
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}
		}
	}
}

func Get_Aws_S3_File(region string, bucket string, filename string, savename string, user string, password string, authfile string) (err error) {
	svc := Aws_S3_Session(region, user, password, authfile) //primary

	result, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		//Bucket: aws.String("camera360-billing"),
		Key: aws.String(filename),
		//Key:    aws.String("421921736743-aws-billing-detailed-line-items-with-resources-and-tags-2017-03.csv.zip"),
	})
	Check_Error(err)

	file, err := os.Create("/dev/shm/" + savename)
	//file, err := os.Create("/dev/shm/cn-bj-aws-lib.camera360.com.last")
	Check_Error(err)

	if _, err := io.Copy(file, result.Body); err != nil {
		fmt.Println("Failed to copy object to file", err)
	}
	result.Body.Close()
	file.Close()
	return err
}

func Aws_S3_Session(regions string, user string, password string, authfile string) (svc *s3.S3) {

	aws_access_key, aws_secret_access_key, token := Get_S3_Access_Key_From_Server(user, password, authfile)
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
	svc = s3.New(sess, &aws.Config{Region: aws.String(regions), Credentials: creds})
	return svc
}
func Get_S3_Access_Key_From_Server(user string, password string, authfile string) (aws_access_key, aws_secret_access_key, token string) {

	var aws_access_k, aws_secret_access_k []byte
	token = ""
	aws_access_k, err := sh.Command("curl", "-s", "-u", user+":"+password, "https://seckey.360in.com/secure_key/"+authfile).Command("grep", "-w", "aws_access_key").Command("awk", "-F", "=", "{print $2}").Output()
	aws_secret_access_k, err = sh.Command("curl", "-s", "-u", user+":"+password, "https://seckey.360in.com/secure_key/"+authfile).Command("grep", "-w", "aws_secret_access_key").Command("awk", "-F", "=", "{print $2}").Output()
	Check_Error(err)
	return string(aws_access_k), string(aws_secret_access_k), token
}
func Check_Error(err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
}
func CopyFile(source string, dest string) (err error) {
	sourcefile, err := os.Open(source)
	if err != nil {
		return err
	}

	defer sourcefile.Close()

	destfile, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer destfile.Close()

	_, err = io.Copy(destfile, sourcefile)
	if err == nil {
		sourceinfo, err := os.Stat(source)
		if err != nil {
			err = os.Chmod(dest, sourceinfo.Mode())
		}

	}

	return
}

func CopyDir(source string, dest string) (err error) {

	// get properties of source dir
	sourceinfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	// create dest dir

	err = os.MkdirAll(dest, sourceinfo.Mode())
	if err != nil {
		return err
	}

	directory, _ := os.Open(source)

	objects, err := directory.Readdir(-1)

	for _, obj := range objects {

		sourcefilepointer := source + "/" + obj.Name()

		destinationfilepointer := dest + "/" + obj.Name()

		if obj.IsDir() {
			// create sub-directories - recursively
			err = CopyDir(sourcefilepointer, destinationfilepointer)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			// perform copy
			err = CopyFile(sourcefilepointer, destinationfilepointer)
			if err != nil {
				fmt.Println(err)
			}
		}

	}
	return
}

func Judge_If_Php_Fpm_Started() (status string) {
	ot, err := sh.Command("ps", "-ef").Command("grep", "-i", "php-fpm").Command("grep", "-v", "grep").Command("wc", "-l").Output()
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

func Judge_If_Nginx_Started() (status string) {
	ot, err := sh.Command("ps", "-ef").Command("grep", "-i", "nginx").Command("grep", "-v", "grep").Command("wc", "-l").Output()
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

func Start_Nginx() (err error) {
	os.Chdir("/home/worker/nginx/sbin")
	err = sh.Command("./admin-nginx.sh", "start").Run()
	return err
}

func Start_Php_Fpm() (err error) {
	os.Chdir("/home/worker/php/sbin")
	err = sh.Command("./php-fpm.sh", "start").Run()
	return err

}

func Create_Files(filename string) (fd *os.File) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Create file error", err)
		return
	}
	return file
}

func Delete_Already_Deployed_Code(file_dir string, file_lib_dir string) {
	os.Chdir("/home/worker/data/www")
	if _, err := os.Stat(file_dir); err == nil {
		err := os.RemoveAll(file_dir)
		Check_Error(err)
	}
	if _, err := os.Stat(file_lib_dir); err == nil {
		err = os.RemoveAll(file_lib_dir)
		Check_Error(err)
	}
}

func Delete_Git_Resp(file_dir string) {
	os.Chdir("/dev/shm")
	if _, err := os.Stat(file_dir); err == nil {
		err = os.RemoveAll(file_dir + "/.git")
		Check_Error(err)
	}

}
func Delete_Deployed_Code(filelast string, file_dir string, file_liblast string, file_lib_dir string) {
	os.Chdir("/dev/shm")
	if _, err := os.Stat(filelast); err == nil {
		err = os.Remove(filelast)
		Check_Error(err)
	}
	if _, err := os.Stat(file_dir); err == nil {
		err = os.RemoveAll(file_dir)
		Check_Error(err)
	}
	if _, err := os.Stat(file_liblast); err == nil {
		err = os.Remove(file_liblast)
		Check_Error(err)
	}
	if _, err := os.Stat(file_lib_dir); err == nil {
		err = os.RemoveAll(file_lib_dir)
		Check_Error(err)
	}

}
