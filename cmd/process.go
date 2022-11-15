package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
)

func bbbContentProcess(folder, link string) error {
	var bc bbbContent

	u, err := url.Parse(link)
	if err != nil {
		return err
	}

	bc.folder = folder
	bc.meetingId = strings.Split(strings.Split(u.RawQuery, "meetingId=")[1], "&")[0]
	bc.jsonName = "presentation_text.json"
	bc.rawUrl = fmt.Sprintf("%s://%s/presentation/%s", u.Scheme, u.Host, bc.meetingId)
	bc.downloadLinks = []downloadLink{{name: "webcams", ext: "mp4", link: bc.rawUrl + "/video/webcams.mp4"}, {name: "deskshare", ext: "mp4", link: bc.rawUrl + "/deskshare/deskshare.mp4"}, {name: "presentation_text", ext: "json", link: bc.rawUrl + "/" + bc.jsonName}}

	return bc.start()
}

func (bc *bbbContent) start() error {
	fmt.Println("Downloads started...")
	var wg sync.WaitGroup
	for _, dl := range bc.downloadLinks {
		wg.Add(1)
		go dl.downloadFile(bc.folder, &wg)
	}
	wg.Wait()

	ss, err := bc.slaytFiles()
	if err == nil && len(ss) > 0 {
		for _, s := range ss {
			wg.Add(1)
			go s.downloadFile(bc.folder+"/slayt", &wg)
		}
		wg.Wait()
	}

	return nil
}

func (d downloadLink) downloadFile(folderpath string, wg *sync.WaitGroup) {
	fmt.Println(d.link)
	defer wg.Done()
	r, err := http.Get(d.link)
	if err != nil {
		return
	}

	defer r.Body.Close()

	out, err := os.Create(folderpath + "/" + d.name + "." + d.ext)
	if err != nil {
		return
	}
	defer out.Close()

	_, err = io.Copy(out, r.Body)
	if err != nil {
		return
	}
}

func (bc *bbbContent) slaytFiles() ([]downloadLink, error) {
	jsonPath := bc.folder + "/" + bc.jsonName
	plan, err := os.ReadFile(jsonPath)
	if err != nil {
		return nil, err
	}

	d := make(map[string]interface{})
	err = json.Unmarshal(plan, &d)
	if err != nil {
		return nil, err
	}

	keys := []string{}

	for k := range d {
		keys = append(keys, k)
	}

	var dls []downloadLink

	for _, k := range keys {
		for kk := range d[k].(map[string]interface{}) {
			dls = append(dls, downloadLink{name: k + "_" + kk, ext: "png", link: bc.rawUrl + "/presentation/" + k + "/" + kk + ".png"})
		}
	}

	if len(dls) > 0 {
		os.Mkdir(bc.folder+"/slayt", 0755)
	}

	return dls, nil
}
