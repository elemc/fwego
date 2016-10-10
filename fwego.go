package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"os/user"
	"path"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
)

const (
	// BufferSizeConst is a base buffer size
	BufferSizeConst = 4096
	tableHeader     = `    <table>
    <thead>
        <tr>
            <th align="center" width="5%%">Type</th>
            <th align="left" width="50%%">Name</th>
            <th align="right" width="20%%">Size</th>
            <th width="25%%">Date</th>
        </tr>
    </thead>
    <tbody>
`
	tableFooter = `    </tbody>
</table>
`
	htmlHeader = `<!DOCTYPE html>
<html>
<head>
<style type="text/css">
    body { background-color:#fff; color:#333; font-family:verdana, arial, helvetica, sans-serif; font-size:13px; line-height:18px }
    p,ol,ul,td { font-family: verdana, arial, helvetica, sans-serif;font-size:13px; line-height:18px}
    a { color:#000 }
    a:visited { color:#666 }
    a:hover{ color:#fff; background-color:#000 }
    tr.dir { font-weight: bold }
    td.icon { font-size: 20px; }
    a.icon { font-size: 27px; text-decoration: none; }
</style>
<meta charset="UTF-8" />
<title>Files</title>
</head>
<body>
<h1><a href="/">Files</a></h1>

`
	htmlFooter = `</body>
</html>
`
)

var rootPath string
var bufferSize uint64
var listenString string
var showHidden bool

func getTableStringUp(target string) string {
	result := "<tr>" +
		"<td align=\"center\">" +
		"<a class=\"icon\" href=\"%s\">%s</a></td>" +
		"<td colspan=\"4\"></td>" +
		"</tr>"
	result = fmt.Sprintf(result, parentPath(target), "")
	return result
}

func getTableString(target string) string {
	var tClass, tType, tName,
		tSize, tDate, result string

	readPath := path.Join(rootPath, target)
	stat, err := os.Stat(readPath)

	if err != nil {
		log.Printf("Error in getTableString -> os.Stat: %s", err)
		return ""
	}

	if stat.IsDir() {
		tType = ""
		tClass = "dir"
	} else {
		tType = ""
		tSize = humanize.Bytes(uint64(stat.Size()))
		tClass = "file"
	}

	tDate = humanize.Time(stat.ModTime())
	tName = stat.Name()

	result = "<tr class=\"%s\">" +
		"<td class=\"icon\" align=\"center\">%s</td>" +
		"<td><a href=\"%s\">%s</a></td>" +
		"<td align=\"right\">%s</td>" +
		"<td align=\"center\">%s</td></tr>"
	result = fmt.Sprintf(result, tClass, tType, target, tName, tSize, tDate)
	return result
}

func parentPath(p string) string {
	cp := []rune(p)
	cpLen := len(cp)
	if cp[cpLen-1] == '/' {
		cp = cp[:cpLen-1]
	}

	parentDir := path.Dir(string(cp))
	return parentDir
}

func readFile(w http.ResponseWriter, r *http.Request, pathPart string) {

	readPath := path.Join(rootPath, pathPart)

	stat, err := os.Stat(readPath)
	if err != nil {
		log.Printf("Error in readFile -> os.Stat: %s", err)
		return
	}

	mimeType := mime.TypeByExtension(path.Ext(readPath))
	if mimeType == "" {
		mimeType = "application/download"
	}

	octLength := fmt.Sprintf("%d", stat.Size())
	downloadType := []string{mimeType}
	//"application/download"}
	length := []string{octLength}
	w.Header()["Content-Type"] = downloadType
	w.Header()["Content-Length"] = length

	log.Printf("Start download file %s from %s\n", readPath, r.RemoteAddr)

	fileForRead, err := os.Open(readPath)
	if err != nil {
		log.Printf("Error in readFile -> os.Open: %s", err)
		return
	}

	defer func() {
		if err := fileForRead.Close(); err != nil {
			log.Printf("Error in readFile -> fileForRead.Close: %s", err)
			return
		}
	}()

	buf := make([]byte, BufferSizeConst)
	for {
		n, err := fileForRead.Read(buf)
		if err != nil && err != io.EOF {
			log.Printf("Error in readFile -> fileForRead.Read: %s", err)
			return
		}
		if n == 0 {
			break
		}

		if _, err := w.Write(buf[:n]); err != nil {
			log.Printf("Error in writing buffer to web response: %s", err)
			return
		}
	}
	log.Printf("Finished download file %s from %s\n", readPath, r.RemoteAddr)
}

func readDir(w http.ResponseWriter, pathPart string) {
	readPath := path.Join(rootPath, pathPart)
	dir, derr := ioutil.ReadDir(readPath)
	if derr != nil {
		fmt.Fprintf(w, "Oops reading %s", readPath)
		return
	}
	if pathPart != "/" {
		fmt.Fprintf(w, getTableStringUp(pathPart))
	}

	var dirs, files []string

	for _, subDir := range dir {
		name := subDir.Name()
		if name[0] == '.' {
			if !showHidden {
				continue
			}
		}
		npath := path.Join(pathPart, name)

		rNpath := path.Join(rootPath, npath)
		stat, err := os.Stat(rNpath)
		if err != nil {
			continue
		}
		if stat.IsDir() {
			dirs = append(dirs, npath)
		} else {
			files = append(files, npath)
		}

	}
	for _, dnpath := range dirs {
		fmt.Fprintf(w, getTableString(dnpath))
	}
	for _, fnpath := range files {
		fmt.Fprintf(w, getTableString(fnpath))
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	htmlType := []string{"text/html"}
	w.Header()["Content-Type"] = htmlType

	pathPart := r.URL.Path
	readPath := path.Join(rootPath, pathPart)

	rpStat, err := os.Stat(readPath)
	if err != nil {
		fmt.Fprintf(w, "No directory or file")
		return
	}
	if rpStat.IsDir() {
		fmt.Fprintf(w, htmlHeader)
		fmt.Fprintf(w, tableHeader)
		readDir(w, pathPart)
		fmt.Fprintf(w, tableFooter)
		fmt.Fprintf(w, htmlFooter)
	} else {
		readFile(w, r, pathPart)
	}
}

func getRealPath(rp string) string {
	if rp[0] != '~' {
		return rp
	}

	usr, err := user.Current()
	if err != nil {
		log.Printf("Error in getRealPath -> user.Current: %s", err)
		return ""
	}

	homeDir := usr.HomeDir
	newRP := strings.Replace(rp, "~", homeDir, 1)
	return newRP
}

func parseEnvVars() {
	if evAddress := os.Getenv("FWEGO_LISTEN"); evAddress != "" {
		listenString = evAddress
	}
	if evBlockSize := os.Getenv("FWEGO_BLOCK_SIZE"); evBlockSize != "" {
		tBufferSize, err := strconv.ParseUint(evBlockSize, 10, 64)
		if err == nil {
			bufferSize = tBufferSize
		}
	}
	if evShowHidden := os.Getenv("FWEGO_SHOW_HIDDEN"); evShowHidden != "" {
		if evShowHidden == "true" || evShowHidden == "on" || evShowHidden == "1" {
			showHidden = true
		} else if evShowHidden == "false" || evShowHidden == "off" || evShowHidden == "0" {
			showHidden = false
		}
	}
	if evRootPath := os.Getenv("FWEGO_ROOT_PATH"); evRootPath != "" {
		rootPath = getRealPath(evRootPath)
	}
}

func init() {
	var ip = flag.String("address", "127.0.0.1", "Listen IP address")
	var port = flag.Uint("port", 4000, "Listen IP port")
	var bs = flag.Uint64("block-size", uint64(BufferSizeConst), "Block size for download files")
	var rp = flag.String("root-path", "", "Root path for browse")
	var sh = flag.Bool("show-hidden", false, "Show hidden files")

	flag.Parse()

	if *rp == "" && os.Getenv("FWEGO_ROOT_PATH") == "" {
		fmt.Printf("!!! Root path (-root-path) doesn't set.\n!!! Or environment variable FWEGO_ROOT_PATH is not set.\n!!! Please set it.\n")
		flag.PrintDefaults()
		os.Exit(1)
	} else if *rp != "" {
		rootPath = getRealPath(*rp)
	}

	listenString = fmt.Sprintf("%s:%d", *ip, *port)
	bufferSize = *bs
	showHidden = *sh

	// MIME-types
	mime.AddExtensionType(".sh", "text/plain; charset=UTF-8")
	mime.AddExtensionType(".repo", "text/plain; charset=UTF-8")

	// Change variables if environment variables is set
	parseEnvVars()
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(listenString, nil)
}
