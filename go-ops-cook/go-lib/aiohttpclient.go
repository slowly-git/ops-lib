package main  
  
import (  
    "flag"  
    "fmt"  
    "io/ioutil"  
    "net/http"  
    "sync/atomic"  
    "time"  
)  
  
func main() {  
    connection := flag.Int("c", 200, "-c N")  
    timeout := flag.Int("o", 5, "-o N")  
    timeover := flag.Int("t", 5, "-t N")  
    printresult := flag.Bool("p", false, "-p false")  
    method := flag.String("m", "GET", "-m GET")  
    url := flag.String("u", "http://www.baidu.com", "-u http://www.baidu.com")  
    flag.Parse()  
    var Count int32  
    defer func() {  
        if !*printresult {  
            fmt.Println("成功响应：", Count)  
        }  
    }()  
    T := time.Tick(time.Duration(*timeover) * time.Second)  
    var result chan string = make(chan string, 10)  
    t := time.Duration(*timeout) * time.Second  
    Client := http.Client{Timeout: t}  
    for i := 0; i < *connection; i++ {  
        go func() {  
            start := time.Now()
            req, _ := http.NewRequest(*method, *url, nil)  
            resp, _ := Client.Do(req)  
            defer resp.Body.Close()  
            if resp.StatusCode == 200 {  
                b, _ := ioutil.ReadAll(resp.Body)  
                result <- string(b)  
                atomic.AddInt32(&Count, int32(1))  
            	end := time.Now()
	        delta := end.Sub(start)
	        fmt.Println("休眠时间:", delta)
            }  
        }()  
    }  
    for {  
        select {  
        case x := <-result:  
            if *printresult {  
                fmt.Print(x)  
            }  
        case <-T:  
            return  
        }  
    }  
}  
