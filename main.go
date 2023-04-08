package main

/**
 *	iServer - Basic file server written in Go Lang
 *
 *	@author: Monimoy Saha
 *  @created 08 Apr, 2023
 */

import (
	_ "embed"
	"flag"
	"fmt"
	"html/template"
	"mime"
	"net"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

// get version string during build
var versionStr string = "0.0.0.0000"

//go:embed index.html
var indexHTML string

// defined flags
var (
	directory string
	port      string
	version   bool
	verbose   bool
	help      bool
)

func main() {
	flag.StringVar(&directory, "d", ".", "the directory to serve")
	flag.StringVar(&port, "p", "8080", "the port to listen on")
	flag.BoolVar(&verbose, "v", false, "verbose output")
	flag.BoolVar(&version, "V", false, "show version")
	flag.BoolVar(&help, "h", false, "show help message")
	flag.Parse()

	// help message
	if help {
		fmt.Printf("iServer v%s - Basic fileserver using golang. \n", versionStr)
		println()
		println("Usage: iServer [options]")
		println()
		println("Options:")
		flag.PrintDefaults()
		return
	}

	if version {
		fmt.Printf("iServer v%s\n", versionStr)
		return
	}

	// start the server
	http.HandleFunc("/", fileServer)
	fmt.Printf("Serving HTTP on 0.0.0.0 port %s (http://0.0.0.0:%s/) ...\n", port, port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

type templateData struct {
	Version string
	Dir
}

type Dir struct {
	BreadCrumbLinks []BreadCrumbLink
	Path            string
	Files           []File
}

type File struct {
	Name  string
	Mode  string
	Size  string
	Link  string
	IsDir bool
}

type BreadCrumbLink struct {
	Name string
	Link string
}

func fileServer(w http.ResponseWriter, r *http.Request) {

	urlPath := r.URL.Path
	filePath := path.Join(directory, urlPath)

	// check if the file exists
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		http.NotFound(w, r)
		log(r, http.StatusNotFound)
		return
	}

	// Handle directory
	if fileInfo.IsDir() {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		// get directory contents
		files, err := os.ReadDir(filePath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log(r, http.StatusInternalServerError)
			return
		}

		// sort.Slice(files, func(i, j int) bool {
		// 	return strings.ToLower(files[i].Name()) < strings.ToLower(files[j].Name())
		// })

		dir := Dir{
			Path: urlPath,
		}

		directories := []File{
			{
				Name:  "..",
				Mode:  "drwxr-xr-x",
				Size:  "0 B",
				Link:  path.Join(urlPath, ".."),
				IsDir: true,
			},
		}

		// build breadcrumb
		pathParts := strings.Split(r.URL.Path, "/")
		var currentPath string = "/"
		for _, part := range pathParts {
			if part == "" {
				continue
			}
			currentPath = path.Join(currentPath, part)
			dir.BreadCrumbLinks = append(dir.BreadCrumbLinks, BreadCrumbLink{
				Name: part,
				Link: currentPath,
			})
		}

		// build file array
		for _, i := range files {
			fileInfo, _ := os.Stat(filePath + "/" + i.Name())
			file := File{
				Name:  fileInfo.Name(),
				Mode:  fileInfo.Mode().String(),
				Size:  ByteCountIEC(fileInfo.Size()),
				Link:  path.Join(urlPath, fileInfo.Name()),
				IsDir: fileInfo.IsDir(),
			}
			if fileInfo.IsDir() {
				directories = append(directories, file)
			} else {
				dir.Files = append(dir.Files, file)
			}
		}

		dir.Files = append(directories, dir.Files...)

		// lead and execute template
		templateData := templateData{
			Version: versionStr,
			Dir:     dir,
		}
		tmpl := template.Must(template.New("fileServer").Parse(indexHTML))
		// tmpl := template.Must(template.ParseFS(indexHTML))
		if err := tmpl.Execute(w, templateData); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log(r, http.StatusInternalServerError)
			return
		}

		log(r, http.StatusOK)
		return
	}

	// Handle file
	{
		file, err := os.Open(filePath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log(r, http.StatusInternalServerError)
			return
		}
		defer file.Close()

		// get file info
		fileInfo, err = file.Stat()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log(r, http.StatusInternalServerError)
			return
		}

		// set headers
		// w.Header().Set("Content-Disposition", "attachment; filename="+fileInfo.Name())
		mimeType := mime.TypeByExtension(path.Ext(fileInfo.Name()))
		w.Header().Set("Content-Type", mimeType)
		w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))

		// serve the file
		http.ServeContent(w, r, fileInfo.Name(), fileInfo.ModTime(), file)
		log(r, http.StatusOK)
	}
}

// convert size to human readable string
func ByteCountIEC(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "KMGTPE"[exp])
}

// log handling function
func log(r *http.Request, status int) {
	if verbose {
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		fmt.Printf("%s - - [%s] \"%s %s %s\" %d -\n", ip, time.Now().Format("02/Jan/2006 15:04:05"), r.Method, r.URL.Path, r.Proto, status)
	}
}
