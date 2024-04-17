package response

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"github.com/golang-base-template/util/cors"
	"github.com/golang-base-template/util/state"
)

type (
	// Response is ..
	Response struct {
		Callback string  `json:"-"`
		Links    *Link   `json:"links,omitempty"`
		Meta     *Meta   `json:"meta,omitempty"`
		Errors   []Error `json:"errors,omitempty"`
	}
	// Meta is ...
	Meta struct {
		ProcessTime float64 `json:"-"`
		TotalData   int     `json:"total_data,omitempty"`
	}

	// Link is ...
	Link struct {
		Self  string `json:"self,omitempty"`
		First string `json:"first,omitempty"`
		Last  string `json:"last,omitempty"`
		Next  string `json:"next,omitempty"`
		Prev  string `json:"prev,omitempty"`
	}
	// Error is ...
	Error struct {
		Code   string `json:"code"`
		Title  string `json:"title"`
		Detail string `json:"detail"`
	}

	RespAPIStandard struct {
		Header HeaderData  `json:"header"`
		Data   interface{} `json:"data"`

		startTime       time.Time
		allowedOrigin   string
		allowCredential string
	}

	// HeaderData is struct for header
	HeaderData struct {
		ProccessTime float64     `json:"process_time"`
		Messages     interface{} `json:"messages"`
		Reason       string      `json:"reason"`
		ErrorCode    string      `json:"error_code"`
	}
	// EmptyData is ...
	EmptyData struct {
	}
)

const (
	// Version21 for v2.1 response
	Version21 = "2.1"
	// Version3 for version 3.0 response
	Version3 = "3.0"
	// VersionWS for version ws response
	VersionWS = "ws"
)

const (
	ErrDecodeRequestBody = "error when decode request body"
)

var corsChecker cors.IChecker

// New with origin param will create RespAPIStandard object with start time and origin
func New(origin string, allowCredential string, startTime ...time.Time) RespAPIStandard {
	start := time.Now()
	//used mainly in test
	if len(startTime) > 0 {
		start = startTime[0]
	}

	if origin == "" {
		origin = "*"
	}

	if corsChecker == nil {
		corsChecker = cors.New()
	}

	ok := corsChecker.Check(origin)

	if ok {
		return RespAPIStandard{
			startTime:       start,
			allowedOrigin:   origin,
			allowCredential: allowCredential,
		}
	}

	return RespAPIStandard{
		startTime:       start,
		allowedOrigin:   state.State.LocalUrl,
		allowCredential: allowCredential,
	}
}

// WriteError is wrapper function for write
func (res *RespAPIStandard) WriteError(w http.ResponseWriter, code int, message interface{}, reason string, v ...interface{}) {
	if len(v) > 0 {
		res.Write(w, v[0], code, message, reason)
	} else {
		res.Write(w, nil, code, message, reason)
	}
}

// WriteResponse is wrapper for write function with http status OK
func (res *RespAPIStandard) WriteResponse(w http.ResponseWriter, v interface{}) {
	res.Write(w, v, http.StatusOK, []string{}, "")
}

// WriteCSVFile is a wrapper for serving csv file over http response
func (res *RespAPIStandard) WriteCSVFile(w http.ResponseWriter, records [][]string, filename string, separator rune) {

	contentDisposition := fmt.Sprintf("attachment;filename=%s", filename)
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", contentDisposition)

	// expose content disposition to get filename on client
	w.Header().Set("Access-Control-Expose-Headers", "Content-Disposition")

	wr := csv.NewWriter(w)
	wr.Comma = separator
	err := wr.WriteAll(records)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to serve %s file", filename)
		err = errors.Wrap(err, errMsg)
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Write is function to write response with standard JSON format
func (res *RespAPIStandard) Write(w http.ResponseWriter, v interface{}, code int, message interface{}, reason string) {
	res.Header.ProccessTime = time.Since(res.startTime).Seconds()
	res.Header.Reason = reason
	res.Header.Messages = message
	if v == nil {
		res.Data = EmptyData{}
	} else {
		res.Data = v
	}

	e, _ := json.Marshal(res)
	if res.allowedOrigin != "" {
		w.Header().Set("Access-Control-Allow-Origin", res.allowedOrigin)
	}

	w.Header().Set("Access-Control-Allow-Credentials", res.allowCredential)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(e)
}

// GetStartTime is use to get start time when response created in the first time
func (res *RespAPIStandard) GetStartTime() time.Time {
	return res.startTime
}

// SetStartTime is use to override start time on response
func (res *RespAPIStandard) SetStartTime(startTime time.Time) {
	res.startTime = startTime
}

// ResponseAPIStandard is old function for write response header and body
// 		modified from ResponseV4
func ResponseAPIStandard(w http.ResponseWriter, v interface{}, code int, message []string, reason string) {
	startTimeStr := w.Header().Get("Date")
	startTime, err := time.Parse(time.RFC3339Nano, startTimeStr)
	if err != nil {
		startTime = time.Now()
	}
	response := RespAPIStandard{
		startTime:     startTime,
		allowedOrigin: "*",
	}

	response.Write(w, v, code, message, reason)
}
