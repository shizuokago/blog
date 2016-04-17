package blog

import (
	"bufio"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/mail"
	"strings"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func init() {
}

func incomingMail(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)
	defer r.Body.Close()

	//var b bytes.Buffer
	//_, err := b.ReadFrom(r.Body)
	//if err != nil {
	//log.Errorf(c, "ReadFrom Error[%s]", err.Error())
	//return
	//}

	msg, err := mail.ReadMessage(r.Body)
	if err != nil {
		log.Errorf(c, "ReadMessage Error[%s]", err.Error())
		return
	}

	contentType := msg.Header.Get("content-type")

	parseBody(r, msg.Body, contentType)

	// Article
	// File
	// Html
}

func parseBody(r *http.Request, body io.Reader, contentType string) {

	c := appengine.NewContext(r)
	log.Infof(c, "Parse ContentType [%s]", contentType)

	var mediatype string
	var params map[string]string
	var err error

	if mediatype, params, err = mime.ParseMediaType(contentType); err != nil {
		log.Errorf(c, "Parse Error[%s]", err.Error())
	}

	switch strings.Split(mediatype, "/")[0] {
	case "text", "html":
		log.Infof(c, "Text")
	case "message":
		log.Infof(c, "Message")
	case "multipart":
		boundary := params["boundary"]
		reader := multipart.NewReader(body, boundary)
		var part *multipart.Part

		for {
			if part, _ = reader.NextPart(); part == nil {
				break
			}

			contentType := part.Header.Get("content-type")
			parseBody(r, part, contentType)
			part.Close()
		}

	default:
		log.Infof(c, "%v,name=%v", mediatype, params["name"])
	}
}

func conversion(inStream io.Reader, outStream io.Writer) error {
	//read from stream (Shift-JIS to UTF-8)
	scanner := bufio.NewScanner(transform.NewReader(inStream, japanese.ISO2022JP.NewDecoder()))
	list := make([]string, 0)
	for scanner.Scan() {
		list = append(list, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	for _, line := range list {
		_, err := fmt.Fprintln(outStream, line)
		if err != nil {
			return err
		}
	}
	return outStream.Flush()
}
