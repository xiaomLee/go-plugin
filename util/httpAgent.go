package common

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"golang.org/x/net/publicsuffix"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

type HttpMethod string

const (
	POST    HttpMethod = "POST"
	GET     HttpMethod = "GET"
	HEAD    HttpMethod = "HEAD"
	PUT     HttpMethod = "PUT"
	DELETE  HttpMethod = "DELETE"
	PATCH   HttpMethod = "PATCH"
	OPTIONS HttpMethod = "OPTIONS"
)

type HttpContentType string

const (
	TypeJSON              HttpContentType = "application/json"
	TypeXML               HttpContentType = "application/xml"
	TypeFormUrlencoded    HttpContentType = "application/x-www-form-urlencoded"
	TypeHTML              HttpContentType = "text/html"
	TypeText              HttpContentType = "text/plain"
	TypeMultipartFormData HttpContentType = "multipart/form-data"
)

var TypeMap = map[HttpContentType]string{
	TypeJSON:           string(TypeJSON),
	TypeXML:            string(TypeXML),
	TypeFormUrlencoded: string(TypeFormUrlencoded),
	TypeHTML:           string(TypeHTML),
	TypeText:           string(TypeText),
}

type File struct {
	FileName  string
	FieldName string
	Data      []byte
}

// HttpAgent is a package of a http util, it comply with chain call rules
// For Example
//
// to send a GET request, you can write like this:
// 		resp, body, err := NewHttpAgent().Get(url).End()
//
// with data:
// 		resp, body, err := NewHttpAgent().Get(url).Query(map[string]string{"hello": "world"}).End()
//
// and set a timeout:
// 		resp, body, err := NewHttpAgent().Get(url).Query(map[string]string{"hello": "world"}).SetTimeout(time.Second*3).End()
//
//
//
//
// the other request such as POST, PUT, DELETE are very similar to GET,
// but use the SendForm, SendData, SendFile function to send data.
// For Example
//
// to send a POST request with application/x-www-form-urlencoded contentType:
// 		resp, body, err := NewHttpAgent().Post(url).ContentType(TypeFormUrlencoded).SendForm(map...).End()
//
// send a POST request with application/json:
// 		resp, body, err := NewHttpAgent().Post(url).ContentType(TypeJson).SendData([]byte...).End()
//
// send a POST request with multipart/form-data:
// 		request := NewHttpAgent().Post(url).ContentType(TypeMultipartFormData)
// 		request = request.SendFile(f1)
// 		request = request.SendFile(f2)
// 		resp, body, err := request.End()
//
type HttpAgent struct {
	Url       string
	Method    HttpMethod
	Type      HttpContentType
	Header    http.Header
	Data      []byte     // (POST) TypeJSON, TypeXML, TypeHtml, TypeText
	FileData  []File     // (POST) TypeMultipart
	FormData  url.Values // (POST) TypeForm, TypeMultipart
	QueryData url.Values // (GET)
	Client    *http.Client
	Transport *http.Transport
	Cookies   []*http.Cookie
	Errors    []error
}

func NewHttpAgent() *HttpAgent {
	cookieJarOptions := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	jar, _ := cookiejar.New(&cookieJarOptions)

	s := &HttpAgent{
		Type:      TypeFormUrlencoded,
		Data:      make([]byte, 0),
		Header:    http.Header{},
		FormData:  url.Values{},
		QueryData: url.Values{},
		FileData:  make([]File, 0),
		Client:    &http.Client{Jar: jar},
		Transport: &http.Transport{},
		Cookies:   make([]*http.Cookie, 0),
		Errors:    nil,
	}

	return s
}

// --------------------------------------------------------------
//                         METHOD
// --------------------------------------------------------------

func (s *HttpAgent) Get(targetUrl string) *HttpAgent {
	s.Method = GET
	s.Url = targetUrl
	s.Errors = nil
	return s
}

func (s *HttpAgent) Post(targetUrl string) *HttpAgent {
	s.Method = POST
	s.Url = targetUrl
	s.Errors = nil
	return s
}

func (s *HttpAgent) Head(targetUrl string) *HttpAgent {
	s.Method = HEAD
	s.Url = targetUrl
	s.Errors = nil
	return s
}

func (s *HttpAgent) Put(targetUrl string) *HttpAgent {
	s.Method = PUT
	s.Url = targetUrl
	s.Errors = nil
	return s
}

func (s *HttpAgent) Delete(targetUrl string) *HttpAgent {
	s.Method = DELETE
	s.Url = targetUrl
	s.Errors = nil
	return s
}

func (s *HttpAgent) Patch(targetUrl string) *HttpAgent {
	s.Method = PATCH
	s.Url = targetUrl
	s.Errors = nil
	return s
}

func (s *HttpAgent) Options(targetUrl string) *HttpAgent {
	s.Method = OPTIONS
	s.Url = targetUrl
	s.Errors = nil
	return s
}

// --------------------------------------------------------------

// ContentType specify the content type to send
func (s *HttpAgent) ContentType(t HttpContentType) *HttpAgent {
	s.Type = t
	return s
}

// AddCookie adds a cookie to the request. The behavior is the same as AddCookie on Request from net/http
func (s *HttpAgent) AddCookie(c *http.Cookie) *HttpAgent {
	s.Cookies = append(s.Cookies, c)
	return s
}

// AddCookies is a convenient method to add multiple cookies
func (s *HttpAgent) AddCookies(cookies []*http.Cookie) *HttpAgent {
	s.Cookies = append(s.Cookies, cookies...)
	return s
}

// SetHeader setting header fields
func (s *HttpAgent) SetHeader(key, val string) *HttpAgent {
	s.Header.Set(key, val)
	return s
}

// AddHeader setting header fields with multiple val
func (s *HttpAgent) AddHeader(key, val string) *HttpAgent {
	s.Header.Add(key, val)
	return s
}

// Timeout config the timeout to request
func (s *HttpAgent) Timeout(d time.Duration) *HttpAgent {
	s.Transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(network, addr, d)
		if err != nil {
			s.Errors = append(s.Errors, err)
			return nil, err
		}
		if err = conn.SetDeadline(time.Now().Add(d)); err != nil {
			return nil, err
		}
		return conn, nil
	}
	s.Client.Transport = s.Transport
	return s
}

// TLSClientConfig setting the tls config
func (s *HttpAgent) TLSClientConfig(config *tls.Config) *HttpAgent {
	s.Transport.TLSClientConfig = config
	s.Client.Transport = s.Transport
	return s
}

// Query accept a map, and will form a query-string in url of GET method
func (s *HttpAgent) Query(queryData map[string]string) *HttpAgent {
	for key, val := range queryData {
		s.QueryData.Add(key, val)
	}
	return s
}

// SendForm accept a map which is usually used to assign data to POST or PUT method
// and it usually used to TypeFormUrlencoded and TypeMultipartFormData
func (s *HttpAgent) SendForm(formData map[string]string) *HttpAgent {
	for key, val := range formData {
		s.FormData.Add(key, val)
	}
	return s
}

// SendData specify data to send of POST or PUT method
// and it usually used to TypeJSON, TypeXML, TypeHTML, TypeText
func (s *HttpAgent) SendData(data []byte) *HttpAgent {
	// TODO 校验json、xml合法性??
	s.Data = data
	return s
}

// SendFile add a file used to send of POST method and TypeMultipartFormData
func (s *HttpAgent) SendFile(f File) *HttpAgent {
	s.FileData = append(s.FileData, f)
	return s
}

// ResetAllDate reset the Data, QueryData, FormData, FileData but reserve headers, cookies
func (s *HttpAgent) ResetAllDate() *HttpAgent {
	s.Data = make([]byte, 0)
	s.QueryData = url.Values{}
	s.FormData = url.Values{}
	s.FileData = make([]File, 0)
	return s
}

// MakeRequest convert the entire HttpAgent to http.Request
func (s *HttpAgent) MakeRequest() (*http.Request, error) {
	var (
		contentType = ""
		urlStr      = s.Url
		body        io.Reader
	)

	switch s.Method {
	case GET:
		urlStr = urlStr + "?" + s.QueryData.Encode()

	case POST, HEAD, PUT, DELETE, PATCH, OPTIONS:
		switch s.Type {
		case TypeFormUrlencoded:
			contentType = TypeMap[s.Type]
			body = strings.NewReader(s.FormData.Encode())

		case TypeMultipartFormData:
			if len(s.FileData) > 0 {
				buf := &bytes.Buffer{}
				writer := multipart.NewWriter(buf)

				for key, value := range s.FormData {
					writer.WriteField(key, value[0])
				}

				for _, f := range s.FileData {
					fw, _ := writer.CreateFormFile(f.FieldName, f.FileName)
					fw.Write(f.Data)
				}

				body = buf
				writer.Close()
				contentType = writer.FormDataContentType()
			}

		case TypeJSON, TypeXML, TypeHTML, TypeText:
			contentType = TypeMap[s.Type]
			body = bytes.NewReader(s.Data)

		default:
			e := errors.New("unknow content type")
			s.Errors = append(s.Errors, e)
			return nil, e
		}

	default:
		return nil, errors.New("unknow method")
	}

	request, err := http.NewRequest(string(s.Method), urlStr, body)
	if err != nil {
		s.Errors = append(s.Errors, err)
		return nil, err
	}

	if contentType != "" {
		request.Header.Set("Content-Type", contentType)
	}
	for key, vals := range s.Header {
		if key == "Content-Type" {
			continue
		}
		for _, v := range vals {
			request.Header.Add(key, v)
		}
	}

	return request, nil
}

// End is final function of the entire call chain, and body will be closed
func (s *HttpAgent) End() (*http.Response, []byte, error) {
	var (
		resp *http.Response
		req  *http.Request
		body []byte
		err  error
	)

	req, err = s.MakeRequest()
	if err != nil {
		s.Errors = append(s.Errors, err)
		return nil, nil, err
	}

	resp, err = s.Client.Do(req)
	if err != nil {
		s.Errors = append(s.Errors, err)
		return nil, nil, err
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		s.Errors = append(s.Errors, err)
		return nil, nil, err
	}

	return resp, body, nil
}
