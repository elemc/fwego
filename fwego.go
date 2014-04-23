package main

import (
	"flag"
	"fmt"
	"github.com/dustin/go-humanize"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"os/user"
	"path"
	"strconv"
	"strings"
)

const (
	BUFFER_SIZE = 4096
)

var root_path string
var buffer_size uint64
var listen_string string
var show_hidden bool

func get_table_header() string {
	result := `    <table>
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
	return result
}

func get_table_footer() string {
	result := `    </tbody>
</table>
`
	return result
}

func get_html_header() string {
	result := `<!DOCTYPE html>
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
	return result
}

func get_html_footer() string {
	result := `</body>
</html>
`
	return result
}

func get_table_string_up(target string) string {
	result := "<tr>" +
		"<td align=\"center\">" +
		"<a class=\"icon\" href=\"%s\">%s</a></td>" +
		"<td colspan=\"4\"></td>" +
		"</tr>"
	result = fmt.Sprintf(result, parent_path(target), "")
	return result
}

func get_table_string(target string) string {
	var t_class, t_type, t_name,
		t_size, t_date, result string

	read_path := path.Join(root_path, target)
	stat, err := os.Stat(read_path)

	if err != nil {
		panic(err)
	}

	if stat.IsDir() {
		t_type = ""
		t_class = "dir"
	} else {
		t_type = ""
		t_size = humanize.Bytes(uint64(stat.Size()))
		t_class = "file"
	}

	t_date = humanize.Time(stat.ModTime())
	t_name = stat.Name()

	result = "<tr class=\"%s\">" +
		"<td class=\"icon\" align=\"center\">%s</td>" +
		"<td><a href=\"%s\">%s</a></td>" +
		"<td align=\"right\">%s</td>" +
		"<td align=\"center\">%s</td></tr>"
	result = fmt.Sprintf(result, t_class, t_type, target, t_name, t_size, t_date)
	return result
}

func parent_path(p string) string {
	cp := []rune(p)
	cp_len := len(cp)
	if cp[cp_len-1] == '/' {
		cp = cp[:cp_len-1]
	}

	parent_dir := path.Dir(string(cp))
	return parent_dir
}

func read_file(w http.ResponseWriter, r *http.Request, path_part string) {

	read_path := path.Join(root_path, path_part)

	stat, err := os.Stat(read_path)
	if err != nil {
		panic(err)
	}

	mime_type := mime.TypeByExtension(path.Ext(read_path))
	if mime_type == "" {
		mime_type = "application/download"
	}

	oct_length := fmt.Sprintf("%d", stat.Size())
	download_type := []string{mime_type}
	//"application/download"}
	length := []string{oct_length}
	w.Header()["Content-Type"] = download_type
	w.Header()["Content-Length"] = length

	fmt.Printf("Start download file %s from %s\n", read_path, r.RemoteAddr)

	file_for_read, err := os.Open(read_path)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := file_for_read.Close(); err != nil {
			panic(err)
		}
	}()

	buf := make([]byte, BUFFER_SIZE)
	for {
		n, err := file_for_read.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			break
		}

		if _, err := w.Write(buf[:n]); err != nil {
			panic(err)
		}
	}
	fmt.Printf("Finished download file %s from %s\n", read_path, r.RemoteAddr)
}

func read_dir(w http.ResponseWriter, path_part string) {
	read_path := path.Join(root_path, path_part)
	dir, derr := ioutil.ReadDir(read_path)
	if derr != nil {
		fmt.Fprintf(w, "Oops reading %s", read_path)
		return
	}
	if path_part != "/" {
		fmt.Fprintf(w, get_table_string_up(path_part))
	}

	var dirs, files []string

	for _, sub_dir := range dir {
		name := sub_dir.Name()
		if name[0] == '.' {
			if !show_hidden {
				continue
			}
		}
		npath := path.Join(path_part, name)

		r_npath := path.Join(root_path, npath)
		stat, err := os.Stat(r_npath)
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
		fmt.Fprintf(w, get_table_string(dnpath))
	}
	for _, fnpath := range files {
		fmt.Fprintf(w, get_table_string(fnpath))
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	html_type := []string{"text/html"}
	w.Header()["Content-Type"] = html_type

	path_part := r.URL.Path
	read_path := path.Join(root_path, path_part)

	rp_stat, err := os.Stat(read_path)
	if err != nil {
		fmt.Fprintf(w, "No directory or file")
		return
	}
	if rp_stat.IsDir() {
		fmt.Fprintf(w, get_html_header())
		fmt.Fprintf(w, get_table_header())
		read_dir(w, path_part)
		fmt.Fprintf(w, get_table_footer())
		fmt.Fprintf(w, get_html_footer())
	} else {
		read_file(w, r, path_part)
	}
}

func get_real_path(rp string) string {
	if rp[0] != '~' {
		return rp
	}

	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	home_dir := usr.HomeDir
	new_rp := strings.Replace(rp, "~", home_dir, 1)
	return new_rp
}

func parse_env_vars() {
	if ev_address := os.Getenv("FWEGO_LISTEN"); ev_address != "" {
		listen_string = ev_address
	}
	if ev_block_size := os.Getenv("FWEGO_BLOCK_SIZE"); ev_block_size != "" {
		t_buffer_size, err := strconv.ParseUint(ev_block_size, 10, 64)
		if err == nil {
			buffer_size = t_buffer_size
		}
	}
	if ev_show_hidden := os.Getenv("FWEGO_SHOW_HIDDEN"); ev_show_hidden != "" {
		if ev_show_hidden == "true" || ev_show_hidden == "on" || ev_show_hidden == "1" {
			show_hidden = true
		} else if ev_show_hidden == "false" || ev_show_hidden == "off" || ev_show_hidden == "0" {
			show_hidden = false
		}
	}
	if ev_root_path := os.Getenv("FWEGO_ROOT_PATH"); ev_root_path != "" {
		root_path = get_real_path(ev_root_path)
	}
}

func init() {
	var ip = flag.String("address", "127.0.0.1", "Listen IP address")
	var port = flag.Uint("port", 4000, "Listen IP port")
	var bs = flag.Uint64("block-size", uint64(BUFFER_SIZE), "Block size for download files")
	var rp = flag.String("root-path", "", "Root path for browse")
	var sh = flag.Bool("show-hidden", false, "Show hidden files")

	flag.Parse()

	if *rp == "" && os.Getenv("FWEGO_ROOT_PATH") == "" {
		fmt.Printf("!!! Root path (-root-path) doesn't set.\n!!! Or environment variable FWEGO_ROOT_PATH is not set.\n!!! Please set it.\n")
		flag.PrintDefaults()
		os.Exit(1)
	} else if *rp != "" {
		root_path = get_real_path(*rp)
	}

	listen_string = fmt.Sprintf("%s:%d", *ip, *port)
	buffer_size = *bs
	show_hidden = *sh

	// MIME-types
	mime.AddExtensionType(".sh", "text/plain; charset=UTF-8")
	mime.AddExtensionType(".repo", "text/plain; charset=UTF-8")

	// Change variables if environment variables is set
	parse_env_vars()
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(listen_string, nil)
}
