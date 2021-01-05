package main

//简单易用的基于go的多线程批量ip源代码泄露目录扫描工具Example: scan-j.exe url.txt 20 dict.txt
//支持根据域名针对性生成字典爆破

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"sync"
	"color"
)

func banner(){
	fmt.Println(`
                                             __ 
                                            |__|
   /  ___// ___\\__  \  /    \   ______     |  |
   \___ \\  \___ / __ \|   |  \ /_____/     |  |
  /____  >\___  >____  /___|  /         /\__|  |
	   \/     \/     \/     \/          \______|`)
	   fmt.Println()
	   fmt.Println()
}

var result []string
var wg sync.WaitGroup

func main() {
	banner()
	if len(os.Args) != 4 {
		fmt.Printf("Example: %s url.txt 20 dict.txt\r\n", os.Args[0])
		os.Exit(1)
	}
	runtime.GOMAXPROCS(runtime.NumCPU())
	threa := os.Args[2]
	thread, _ := strconv.Atoi(threa)
	url_file := os.Args[1]
	dict_file := os.Args[3]

	urls, e := readfile(url_file)
	dicts, e := readfile(dict_file)
	if e != nil {
		e.Error()
	}
	urls_arr := strings.Split(string(urls), "\r\n")
	dict_arr := strings.Split(string(dicts), "\r\n")

	ip_num := len(urls_arr) 

	for i := 0; i < ip_num; i++ {
		ip := urls_arr[i]
		new_dict_arr := []string{}
		u, err := url.Parse(ip)
		if err != nil {
			panic(err)
		}
		var host string
		host = u.Host //获取ip的host
		if strings.Contains(u.Host, ":"){
			host = strings.Split(u.Host, ":")[0]
		}
		host_arr := strings.Split(host, ".")
		end := []string{".rar", ".zip", ".7z", ".tar", ".gz", ".tar.gz", ".bz2", ".tar.bz2", ".sql", ".bak", ".dat", ".txt", ".log", ".mdb"}
		for _, value := range end{
			new_dict_arr = append(new_dict_arr, "/" + host + value)
			new_dict_arr = append(new_dict_arr, "/" + host_arr[len(host_arr)-2]+"."+host_arr[len(host_arr)-1] + value)
			new_dict_arr = append(new_dict_arr, "/" + host_arr[len(host_arr)-2] + value)
		} 

		for _, value := range dict_arr{
			new_dict_arr = append(new_dict_arr,value)
		}
		run_ip(ip, new_dict_arr,thread)
	}
}

func run_ip(ip string, dir []string , thread int){
	lens := len(dir)
	new_lens := lens
	var task int
	if lens < thread {
		thread = lens
		task = 1
	}else{
		for i := 1; i< thread ; i++{
			if new_lens % thread != 0{
				new_lens += 1
			}else{
				break
			}
		}

		task = new_lens / thread
	}
	fmt.Printf("[+]当前ip: %s 设置了%d 个线程, 共 %d 个目录，每线程%d 个目录进行爆破\r\n",ip,thread,lens,task)
	wg.Add(thread)
	for i := 0; i < thread; i++ {
		go run(ip, dir, i, task)
	}
	wg.Wait()
	

}

func run(urls string, dir []string, tnum int, task int) {
	//fmt.Printf("[+]当前第 %d 个线程，ip为 %s \r\n",tnum,urls)
	for i := tnum*task + 1; i < (tnum*task)+task+1; i++ {
		if i > len(dir){
			break
		}
		dir[i-1] = strings.TrimSpace(dir[i-1])
		url := urls + dir[i-1]
		code, err := scandir(url)
		if err != nil {
			continue
		}
		if code == 403 || code == 404 {
			fmt.Printf("                                                                              \r")
			fmt.Printf("Checking: %s ...\r",dir[i-1])
		} else {
			color.Yellow.Printf("Found: %s [%d]!!!\r\n",dir[i-1],code)
			result = append(result,dir[i-1])
		}

	}
	wg.Done()
}

func readfile(dir string) ([]byte, error) {
	f, err := os.Open(dir)
	if err != nil {
		err.Error()
	}
	return ioutil.ReadAll(f)
}
func scandir(url string) (int, error) {
	resp, err := http.Get(url)
	var status int
	if err != nil {
		status = 404
	} else {
		status = resp.StatusCode
	}

	return status, err
}
