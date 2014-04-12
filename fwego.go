package main

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
)

var root_path string

func get_table_header() string {
	result := `    <table width="100%%">
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
		"<td></td>" +
		"<td colspan=\"3\">" +
		"<a href=\"%s\">%s</a></td>" +
		"</tr><tr><td colspan=\"4\">&nbsp;</td></tr>"
	result = fmt.Sprintf(result, parent_path(target), "Up to")
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
		t_type = "DIR"
		t_class = "dir"
	} else {
		t_type = "[*]"
		t_size = humanize.Bytes(uint64(stat.Size()))
		t_class = "file"
	}

	t_date = humanize.Time(stat.ModTime())
	t_name = stat.Name()

	result = "<tr class=\"%s\">" +
		"<td align=\"center\">%s</td>" +
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

	oct_length := fmt.Sprintf("%d", stat.Size())
	download_type := []string{"application/download"}
	length := []string{oct_length}
	w.Header()["Content-Type"] = download_type
	w.Header()["Content-Length"] = length

	file_for_read, err := os.Open(read_path)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := file_for_read.Close(); err != nil {
			panic(err)
		}
	}()

	buf := make([]byte, 1024)
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

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <filesystem root path>\n", os.Args[0])
		os.Exit(1)
	}
	root_path = os.Args[1]
	rp_stat, err := os.Stat(root_path)
	if err != nil || !rp_stat.IsDir() {
		fmt.Printf("Target %s doesn't exist or not a directory!", root_path)
		os.Exit(2)
	}

	http.HandleFunc("/", handler)
	http.ListenAndServe(":3000", nil)
}