package main

import (
	"fmt"
	"log"
	"os"

	humanize "github.com/dustin/go-humanize"
)

const (
	tableHeader = `    <table>
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

func getTableStringUp(target string) string {
	result := "<tr>" +
		"<td align=\"center\">" +
		"<a class=\"icon\" href=\"%s\">%s</a></td>" +
		"<td colspan=\"4\"></td>" +
		"</tr>"
	result = fmt.Sprintf(result, parentPath(getPathWithoutRoot(target)), "")
	return result
}

func getTableString(target string) string {
	var (
		tClass, tType, tName, tSize, tDate, result, tPath string
	)

	stat, err := os.Stat(target)

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
	tPath = getPathWithoutRoot(target)

	result = "<tr class=\"%s\">" +
		"<td class=\"icon\" align=\"center\">%s</td>" +
		"<td><a href=\"%s\">%s</a></td>" +
		"<td align=\"right\">%s</td>" +
		"<td align=\"center\">%s</td></tr>"
	result = fmt.Sprintf(result, tClass, tType, tPath, tName, tSize, tDate)
	return result
}
