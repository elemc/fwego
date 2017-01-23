package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"mime"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

const (
	// BufferSizeConst is a base buffer size
	BufferSizeConst = 4096
)

var (
	rootPath     string
	bufferSize   uint64
	listenString string
	showHidden   bool

	filesHandler fasthttp.RequestHandler
	logLevel     string
)

func parentPath(p string) string {
	parentDir, _ := filepath.Split(p)
	return parentDir
}

func getPathWithoutRoot(path string) string {
	if !strings.HasPrefix(path, rootPath) {
		return path
	}
	return strings.TrimPrefix(path, rootPath)
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
	if evLogLevel := os.Getenv("FWEGO_LOG_LEVEL"); evLogLevel != "" {
		logLevel = evLogLevel
	}
}

func init() {
	var ip = flag.String("address", "127.0.0.1", "Listen IP address")
	var port = flag.Uint("port", 4000, "Listen IP port")
	var bs = flag.Uint64("block-size", uint64(BufferSizeConst), "Block size for download files")
	var rp = flag.String("root-path", "", "Root path for browse")
	var sh = flag.Bool("show-hidden", false, "Show hidden files")
	flag.StringVar(&logLevel, "loglevel", "warn", "Set log level")

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

	if loglevel, err := log.ParseLevel(logLevel); err != nil {
		log.Warnf("unable to parse log level \"%s\" to be used \"warn\"", logLevel)
	} else {
		log.SetLevel(loglevel)
	}
}

func main() {
	filesHandler = fasthttp.FSHandler(rootPath, 0)
	if err := fasthttp.ListenAndServe(listenString, requestHandler); err != nil {
		log.Fatalf("error in listen and server: %s", err)
		return
	}
}

func requestHandler(ctx *fasthttp.RequestCtx) {
	if !ctx.IsGet() {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	var (
		pathStat os.FileInfo
		err      error
	)
	urlPath := string(ctx.Path())
	if urlPath == "/favicon.ico" {
		return
	}
	path := filepath.Join(rootPath, urlPath)
	log.Debugf("GET [%s]", path)

	if pathStat, err = os.Stat(path); err != nil {
		log.Errorf("unable to get stats for path [%s]: %s", path, err)
		return
	}

	if pathStat.IsDir() {
		ctx.SetContentType("contentType")
		ctx.WriteString(htmlHeader)
		ctx.WriteString(tableHeader)
		readDir(ctx, path)
		ctx.WriteString(tableFooter)
		ctx.WriteString(htmlFooter)
		ctx.SetStatusCode(fasthttp.StatusOK)

	} else {
		filesHandler(ctx)
	}
}

func readDir(ctx *fasthttp.RequestCtx, path string) {
	var (
		dir      []os.FileInfo
		err      error
		fullPath string
		stat     os.FileInfo
	)
	// fullPath = filepath.Join(rootPath, path)
	if dir, err = ioutil.ReadDir(path); err != nil {
		ctx.WriteString("error reading path:")
		ctx.WriteString(err.Error())
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
	if fullPath != "/" {
		ctx.WriteString(getTableStringUp(path))
	}

	var dirs, files []string

	for _, subDir := range dir {
		name := subDir.Name()
		if name[0] == '.' {
			if !showHidden {
				continue
			}
		}
		npath := filepath.Join(path, name)
		if stat, err = os.Stat(npath); err != nil {
			log.Errorf("unable to get stats for [%s]: %s", getPathWithoutRoot(npath), err)
			continue
		}
		if stat.IsDir() {
			dirs = append(dirs, npath)
		} else {
			files = append(files, npath)
		}

	}
	sort.Strings(dirs)
	sort.Strings(files)
	for _, dnpath := range dirs {
		ctx.WriteString(getTableString(dnpath))
	}
	for _, fnpath := range files {
		ctx.WriteString(getTableString(fnpath))
	}
}
