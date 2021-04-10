package datastore

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"io"

	"io/ioutil"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net/http"
	"net/mail"
	"strings"

	"github.com/pborman/uuid"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	"google.golang.org/appengine/log"
)

func init() {
}

type MailData struct {
	subject string
	msg     bytes.Buffer
	file    bytes.Buffer
	mime    string
}

const (
	SUBJ_PREFIX_ISO2022JP_B = "=?iso-2022-jp?b?"
	SUBJ_PREFIX_ISO2022JP_Q = "=?iso-2022-jp?q?"
	SUBJ_PREFIX_UTF8_B      = "=?utf-8?b?"
	SUBJ_PREFIX_UTF8_Q      = "=?utf-8?q?"
	CHARSET_ISO2022JP       = "iso-2022-jp"
	ENC_QUOTED_PRINTABLE    = "quoted-printable"
	ENC_BASE64              = "base64"
	MEDIATYPE_TEXT          = "text/"
	MEDIATYPE_MULTI         = "multipart/"
	MEDIATYPE_MULTI_REL     = "multipart/related"
	MEDIATYPE_MULTI_ALT     = "multipart/alternative"
)

func incomingMail(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	defer r.Body.Close()

	//var b bytes.Buffer
	//_, err := b.ReadFrom(r.Body)
	//if err != nil {
	//log.Errorf(c, "ReadFrom Error[%s]", err.Error())
	//return
	//}

	msg, err := mail.ReadMessage(r.Body)
	if err != nil {
		log.Errorf(ctx, "ReadMessage Error[%s]", err.Error())
		return
	}

	contentType := msg.Header.Get("content-type")

	data := MailData{
		subject: decSubject(msg.Header.Get("Subject")),
	}
	log.Infof(ctx, "Subject = [%s]", msg.Header.Get("Subject"))

	parseBody(r, msg.Body, contentType, &data)

	err = CreateHtmlFromMail(ctx, &data)
	if err != nil {
		log.Errorf(ctx, "ReadMessage Error[%s]", err.Error())
		return
	}
}

func parseBody(r *http.Request, body io.Reader, contentType string, data *MailData) {

	c := r.Context()
	log.Infof(c, "Parse ContentType [%s]", contentType)

	var mediatype string
	var params map[string]string
	var err error

	if mediatype, params, err = mime.ParseMediaType(contentType); err != nil {
		log.Errorf(c, "Parse Error[%s]", err.Error())
	}

	switch strings.Split(mediatype, "/")[0] {
	case "text", "html":
		lines, err := encode(body)
		if err != nil {
			log.Errorf(c, "Parse Error[%s]", err.Error())
		}

		for _, line := range lines {
			data.msg.Write([]byte(line + "\n"))
		}

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
			parseBody(r, part, contentType, data)
			part.Close()
		}

	default:
		data.mime = mediatype
		io.Copy(&data.file, transform.NewReader(body, japanese.ISO2022JP.NewDecoder()))
	}
}

func encode(inStream io.Reader) ([]string, error) {

	//read from stream (Shift-JIS to UTF-8)
	scanner := bufio.NewScanner(transform.NewReader(inStream, japanese.ISO2022JP.NewDecoder()))
	list := make([]string, 0)
	for scanner.Scan() {
		list = append(list, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func decSubject(subject string) string {
	splitsubj := strings.Fields(subject)
	var bufSubj bytes.Buffer
	for seq, parts := range splitsubj {
		switch {
		case !strings.HasPrefix(parts, "=?"):
			// エンコードなし
			if seq > 0 {
				// 先頭以外はSpaceで区切りなおし
				bufSubj.WriteByte(' ')
			}
			bufSubj.WriteString(parts)

		case len(parts) > len(SUBJ_PREFIX_ISO2022JP_B) && strings.HasPrefix(strings.ToLower(parts[0:len(SUBJ_PREFIX_ISO2022JP_B)]), SUBJ_PREFIX_ISO2022JP_B):
			// iso-2022-jp / base64
			beforeDecode := parts[len(SUBJ_PREFIX_ISO2022JP_B):strings.LastIndex(parts, "?=")]
			afterDecode := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(beforeDecode))
			subj_bytes, _ := ioutil.ReadAll(transform.NewReader(afterDecode, japanese.ISO2022JP.NewDecoder()))
			bufSubj.Write(subj_bytes)

		case len(parts) > len(SUBJ_PREFIX_ISO2022JP_Q) && strings.HasPrefix(strings.ToLower(parts[0:len(SUBJ_PREFIX_ISO2022JP_Q)]), SUBJ_PREFIX_ISO2022JP_Q):
			// iso-2022-jp / quoted-printable
			beforeDecode := parts[len(SUBJ_PREFIX_ISO2022JP_Q):strings.LastIndex(parts, "?=")]
			afterDecode := quotedprintable.NewReader(bytes.NewBufferString(beforeDecode))
			subj_bytes, _ := ioutil.ReadAll(transform.NewReader(afterDecode, japanese.ISO2022JP.NewDecoder()))
			bufSubj.Write(subj_bytes)

		case len(parts) > len(SUBJ_PREFIX_UTF8_B) && strings.HasPrefix(strings.ToLower(parts[0:len(SUBJ_PREFIX_UTF8_B)]), SUBJ_PREFIX_UTF8_B):
			// utf-8 / base64
			beforeDecode := parts[len(SUBJ_PREFIX_UTF8_B):strings.LastIndex(parts, "?=")]
			subj_bytes, _ := base64.StdEncoding.DecodeString(beforeDecode)
			bufSubj.Write(subj_bytes)

		case len(parts) > len(SUBJ_PREFIX_UTF8_Q) && strings.HasPrefix(strings.ToLower(parts[0:len(SUBJ_PREFIX_UTF8_Q)]), SUBJ_PREFIX_UTF8_Q):
			// utf-8 / quoted-printable
			beforeDecode := parts[len(SUBJ_PREFIX_UTF8_Q):strings.LastIndex(parts, "?=")]
			afterDecode := quotedprintable.NewReader(bytes.NewBufferString(beforeDecode))
			subj_bytes, _ := ioutil.ReadAll(afterDecode)
			bufSubj.Write(subj_bytes)
		}
	}
	return bufSubj.String()

}

func CreateHtmlFromMail(ctx context.Context, d *MailData) error {

	id := uuid.New()

	bgd := GetBlog(ctx)
	article := &Article{
		Title:    d.subject,
		Tags:     bgd.Tags,
		Markdown: d.msg.Bytes(),
	}

	article.Key = getArticleKey(id)
	client, err := createClient(ctx)
	_, err = client.Put(ctx, article.Key, article)
	if err != nil {
		return err
	}

	fid := "bg/" + id
	fb := d.file.Bytes()
	file := &File{
		Size: int64(len(fb)),
		Type: FileTypeBG,
	}

	file.Key = getFileKey(fid)

	_, err = client.Put(ctx, file.Key, file)
	if err != nil {
		return err
	}
	fileData := &FileData{
		Content: fb,
		Mime:    d.mime,
	}

	fdk := getFileDataKey(fid)
	fileData.SetKey(fdk)
	_, err = client.Put(ctx, fileData.GetKey(), fileData)
	if err != nil {
		return err
	}
	return nil

}
