package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/streadway/amqp"
)

var xchgName string = "cache-hit"

func main() {

	dir := os.Getenv("WS_DIR")
	port := os.Getenv("WS_PORT")

	// connect to local memcached
	mc := memcache.New("127.0.0.1:11211")

	// set up rabbitmq
	rmqHost := os.Getenv("RMQ_HOST")
	conn, err := amqp.Dial("amqp://guest:guest@" + rmqHost + ":5672/")
	if err != nil {
		log.Println("Failed to connect to rabbitmq: " + err.Error())
	}
	defer conn.Close()

	ch, err := conn.Channel()
	defer ch.Close()

	// rabbitmq declare exchange for cache-hit log
	ch.ExchangeDeclare(
		xchgName, // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)

	// serve from the file, but check cache first
	fileserver := cacheFileServer(http.Dir(dir), mc, ch)

	log.Println("Starting up on port " + port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, fileserver))
}

type cacheFileHandler struct {
	root http.FileSystem
	mc   *memcache.Client
	ch   *amqp.Channel
}

func cacheFileServer(root http.FileSystem, mc *memcache.Client, ch *amqp.Channel) http.Handler {
	return &cacheFileHandler{root, mc, ch}
}

func (f *cacheFileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.Path, string(filepath.Separator)) {
		log.Println("Directory: no cache")
		f.ServeDir(w, r)
	} else {
		dir, file := filepath.Split(r.URL.Path)
		log.Println("Dir: " + dir)
		log.Println("File: " + file)

		hashString := hashUrl(r.URL.Path)

		if it, err := f.mc.Get(hashString); err != nil {
			log.Println("Cache miss: " + hashString)
			f.reportCacheHit(false)
			f.ServeFile(w, r)
		} else {
			log.Println("Cache hit: " + hashString)
			f.reportCacheHit(true)
			w.Write(it.Value)
		}
	}
}

func (f *cacheFileHandler) ServeFile(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path
	fs := f.root
	file, err := fs.Open(name)
	if err != nil {
		msg, code := toHTTPError(err)
		http.Error(w, msg, code)
		return
	}
	defer file.Close()

	_, err = file.Stat()
	if err != nil {
		msg, code := toHTTPError(err)
		http.Error(w, msg, code)
		return
	}

	// copy object to http
	var reader io.Reader = file

	var buf bytes.Buffer
	mw := io.MultiWriter(w, &buf)

	io.Copy(mw, reader)

	// cache
	hashString := hashUrl(r.URL.Path)
	f.mc.Set(&memcache.Item{Key: hashString, Value: buf.Bytes()})
}

func (f *cacheFileHandler) ServeDir(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path
	fs := f.root
	file, err := fs.Open(name)
	if err != nil {
		msg, code := toHTTPError(err)
		http.Error(w, msg, code)
		return
	}
	defer file.Close()

	_, err = file.Stat()
	if err != nil {
		msg, code := toHTTPError(err)
		http.Error(w, msg, code)
		return
	}

	dirs, err := file.Readdir(-1)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<pre>\n")
	for _, d := range dirs {
		name := d.Name()
		if d.IsDir() {
			name += "/"
		}
		url := url.URL{Path: name}
		fmt.Fprintf(w, "<a href=\"%s\">%s</a>\n", url.String(), html.EscapeString(name))
	}
	fmt.Fprintf(w, "</pre>\n")
}

func (f *cacheFileHandler) reportCacheHit(hit bool) {
	var body string
	if hit {
		body = "h"
	} else {
		body = "m"
	}

	err := f.ch.Publish(
		xchgName, // exchange
		"",       // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})

	if err == nil {
		log.Printf("report cache hit '%s'", body)
	} else {
		log.Printf(err.Error())
	}
}

func hashUrl(url string) string {
	hash := md5.Sum([]byte(url))
	return hex.EncodeToString(hash[:])
}

func toHTTPError(err error) (msg string, httpStatus int) {
	if os.IsNotExist(err) {
		return "404 page not found", 404
	}
	if os.IsPermission(err) {
		return "403 Forbidden", 403
	}
	// Default:
	return "500 Internal Server Error", 500
}
