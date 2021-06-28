// This file is part of the ego-cli distribution.
// Copyright (c) Next.e.GO Mobile SE, Aachen, Germany (https://e-go-mobile.com/)
//
// ego-cli is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as
// published by the Free Software Foundation, version 3.
//
// ego-cli is distributed in the hope that it will be useful, but
// WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU
// Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"fmt"
	"html"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/thatisuday/commando"
)

type HttpServerItem struct {
	File     fs.FileInfo
	PathName string
	Name     string
}

func serve_createRequestHandler(cwd string) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		var fullPath = filepath.Join(cwd, req.URL.Path)

		if strings.HasPrefix(fullPath, cwd) {
			file, err := os.Stat(fullPath)

			if err != nil {
				if os.IsNotExist(err) {
					serve_sendNotFound(w, req)
				} else {
					serve_sendError(err, w, req)
				}
			} else {
				if file.IsDir() {
					serve_listDirectory(fullPath, w, req)
				} else {
					serve_sendFile(fullPath, w, req)
				}
			}
		} else {
			serve_sendNotFound(w, req)
		}
	}
}

func serve_listDirectory(dir string, w http.ResponseWriter, req *http.Request) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		serve_sendError(err, w, req)
		return
	}

	var items []HttpServerItem
	for _, f := range files {
		if !strings.HasPrefix(f.Name(), ".") {
			var newItem HttpServerItem
			newItem.Name = f.Name()
			newItem.PathName = newItem.Name
			newItem.File = f

			items = append(items, newItem)
		}
	}

	sort.SliceStable(items, func(i, j int) bool {
		var x = items[i]
		var y = items[j]

		// first check if X is a directory
		var isXDir = 1
		if x.File.IsDir() {
			isXDir = 0
		}
		// then check if Y is a directory
		var isYDir = 1
		if y.File.IsDir() {
			isYDir = 0
		}

		var isDirDiff = isXDir - isYDir
		if isDirDiff != 0 {
			return isDirDiff < 0 // sort if directory first
		}

		// ... now by name

		var nameX = strings.TrimSpace(strings.ToLower(x.Name))
		var nameY = strings.TrimSpace(strings.ToLower(y.Name))

		return strings.Compare(nameX, nameY) < 0
	})

	var h = "<html>"
	h += "<head>"
	h += "<title>e.GO CLI</title>"
	h += "</head>"

	h += "<body>"
	h += "<h1>Index of " + html.EscapeString(req.URL.Path) + "</h1>"

	if req.URL.Path != "/" {
		var parentPath = filepath.Join(dir, "..")
		if parentPath != dir {
			parentDir, err := os.Stat(parentPath)
			if err != nil {
				serve_sendError(err, w, req)
				return
			}

			if parentDir.IsDir() {
				var newItem HttpServerItem
				newItem.Name = ".."
				newItem.PathName = ".."
				newItem.File = parentDir

				items = append([]HttpServerItem{newItem}, items...)
			}
		}
	}

	if len(items) > 0 {
		h += "<table>"

		h += "<thead>"
		h += "<tr>"
		h += "<th>&nbsp;&nbsp;&nbsp;Name&nbsp;&nbsp;&nbsp;</th>"
		h += "<th>&nbsp;&nbsp;&nbsp;Last modified&nbsp;&nbsp;&nbsp;</th>"
		h += "<th align=\"right\">&nbsp;&nbsp;&nbsp;Size&nbsp;&nbsp;&nbsp;</th>"
		h += "</tr>"
		h += "</thead>"

		h += "<tbody>"

		for _, i := range items {
			var urlPath = filepath.Join(req.URL.Path, i.PathName)

			var name = html.EscapeString(i.Name)
			var href = html.EscapeString(urlPath)
			var size = html.EscapeString("-")
			var lastModified = html.EscapeString(i.File.ModTime().Format(time.RFC3339))
			if !i.File.IsDir() {
				size = html.EscapeString(strconv.FormatInt(i.File.Size(), 10))
			}

			h += "<tr>"
			h += "<td>&nbsp;&nbsp;&nbsp;<a href=\"" + href + "\">" + name + "</a>&nbsp;&nbsp;&nbsp;</td>"
			h += "<td>&nbsp;&nbsp;&nbsp;" + lastModified + "&nbsp;&nbsp;&nbsp;</td>"
			h += "<td align=\"right\">&nbsp;&nbsp;&nbsp;" + size + "&nbsp;&nbsp;&nbsp;</td>"
			h += "</tr>"
		}

		h += "</tbody>"

		h += "</table>"
	}

	h += "<hr>"

	var footerText = fmt.Sprintf("<a href=\"https://github.com/egomobile/ego-cli\" target=\"_blank\">ego-cli/%v</a>", html.EscapeString(commando.DefaultCommandRegistry.Version))
	footerText += fmt.Sprintf(" Server at %v", html.EscapeString(req.Host))

	h += "<div>" + footerText + "</div>"

	h += "</body>"

	h += "</html>"

	w.WriteHeader(200)
	w.Write([]byte(h))
}

func serve_sendError(err error, w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(500)
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")

	w.Write([]byte(err.Error()))
}

func serve_sendFile(file string, w http.ResponseWriter, req *http.Request) {
	fileName := filepath.Base(file)

	sendAsDownload, contentType := serve_shouldSendAsDownload(file)

	if sendAsDownload {
		w.Header().Set("Content-Disposition", "attachment; filename="+strconv.Quote(fileName))
		w.Header().Set("Content-Type", contentType)
	}

	http.ServeFile(w, req, file)
}

func serve_sendNotFound(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(404)
}

func serve_shouldSendAsDownload(file string) (bool, string) {
	var contentType = ""

	mime, err := mimetype.DetectFile(file)
	if err == nil {
		contentType = strings.TrimSpace(strings.ToLower(mime.String()))
	}

	if contentType == "" {
		contentType = "application/octet-stream"
	}

	var sendAsDownload = true

	if strings.HasPrefix(contentType, "text/") {
		sendAsDownload = false
	}
	if strings.Contains(contentType, "/json") {
		sendAsDownload = false
	}
	if strings.Contains(contentType, "/pdf") {
		sendAsDownload = false
	}
	if strings.HasPrefix(contentType, "image/") {
		sendAsDownload = false
	}

	return sendAsDownload, contentType
}

func serve_run(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
	var port = 8080

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	var handler = http.HandlerFunc(serve_createRequestHandler(cwd))

	fmt.Println("Running server on port", port)
	err = http.ListenAndServe(":"+strconv.Itoa(port), handler)
	if err != nil {
		log.Fatal(err)
	}
}

func Setup_serve_Command() {
	commando.
		Register("serve").
		SetShortDescription("serve files via HTTP").
		SetDescription("Starts a HTTP server, serving the files in the current directory").
		SetAction(serve_run)
}
