package main

import (
    "log"
    "os"
    "os/exec"
    "net/http"
    "fmt"
    "io/ioutil"
    "regexp"
	"encoding/json"
    "flag"
)

type StageData struct {
	Error   int    `json:"error"`
	Message string `json:"message"`
	HTML    string `json:"html"`
	Data    struct {
		ID       string `json:"id"`
		Stage    string `json:"stage"`
		Callback string `json:"callback"`
	} `json:"data"`
}
//http://list.youku.com/show/episode?id=306170&stage=reload_1&callback=html

func main() {

    var showid = flag.Int("i",0,"showid for show in youku")
    var dir = flag.String("d","./","dir to download")
    
    flag.Parse()

    if flag.NArg() != 0 {
        flag.Usage()
        os.Exit(1)
    }

    fmt.Printf("ShowID: %d\n",*showid)
    fmt.Printf("Dir: %s\n",*dir)
    _ = os.Mkdir(*dir, os.ModePerm)
    for i:=1; i<100; i+=10 {
        url := fmt.Sprintf("http://list.youku.com/show/episode?id=%d&stage=reload_%d&callback=html",*showid,i)
        
        resp, err := http.Get(url)
        fmt.Printf("%s\n",url)

        if err != nil {
            fmt.Fprintf(os.Stderr, "fetch: %v\n",err)
            os.Exit(1)
        }

        b, err := ioutil.ReadAll(resp.Body)

        resp.Body.Close()

        if err != nil {

            fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n",url,err)
            os.Exit(1)
        }

        s := string(b[:])

        re:= regexp.MustCompile("^window.html && html\\((.*)\\);")
        n := re.ReplaceAllString(s,"$1")

		var sd StageData

		json.Unmarshal([]byte(n),&sd)

        if sd.Error == 1 {
            break
        }

        re = regexp.MustCompile("v.youku.com/v_show/id_\\w+==.html")
        urls := re.FindAllString(sd.HTML,-1)

        if len(urls)<=0 {
            break
        }

        for _, u := range urls {
            url = fmt.Sprintf("http://%s",u)
            fmt.Println(url)

            //cmd := exec.Command("you-get", "-i", url, "-o",*dir)
            cmd := exec.Command("you-get", "-o",*dir,url)
            out, err := cmd.Output()
            if err != nil {
                log.Fatal(err)
            }
            fmt.Printf("%s",out)
        }

    }
}
