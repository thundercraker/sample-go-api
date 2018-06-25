/*
	The package server provides the HTTP server that makes the Comment Parsing service as a RESTful service
*/
package server

import (
	"net/http"

	"commentparser/logging"
	"commentparser/models"
	"commentparser/services"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/url"
	"strings"
	"time"
)

// this is the model representing the configuration options that the server uses
type Configuration struct {
	Development          bool   // true if the application is in development mode, false in production
	Address              string // the address to bind the server to
	LogName              string // path to log to
	GoogleCloudProjectID string // the google cloud project ID
	GoogleCloudCredFile  string // google cloud API credentials file
}

// POST "/parse"
// Extract the comments where comments contains the specified tokens in the
// provided package name, the body should be a models.CommentParsingRequest
func ParseAction(
	writer http.ResponseWriter,
	body []byte,
	logging logging.Logging) ErrorPkg {

	var err error
	var request models.CommentParsingRequest
	err = json.Unmarshal(body, &request)

	if err != nil {
		return ErrorWithCodeSantized(400, err)
	}

	if len(request.PackageName) < 1 {
		return ErrorWithCodeSantized(
			400,
			errors.New("The parameter `PackageName` cannot be empty"))
	}

	if len(request.Tokens) < 1 {
		return ErrorWithCodeSantized(
			400,
			errors.New("The parameter `Tokens` cannot be empty"))
	}

	resObj, err := services.ExtractRelevantComments(request, logging)

	if err != nil {
		return Error(err)
	}

	res, err := json.Marshal(resObj)

	if err != nil {
		return Error(err)
	}

	writer.Write(res)

	return ErrorPkg{}
}

// GET "/"
// Extract the comments where comments contains the specified tokens in the
// provided package name
func IndexAction(
	writer http.ResponseWriter,
	values url.Values,
	logging logging.Logging) ErrorPkg {

	qPackage := values.Get("package")
	qTokens := values.Get("tokens")
	if len(qPackage) < 1 {
		return ErrorWithCodeSantized(400, errors.New("the query must contain the parameter `package`"))
	}
	if len(qTokens) < 1 {
		return ErrorWithCodeSantized(400, errors.New("the query must contain the parameter `tokens`"))
	}
	request := models.CommentParsingRequest{
		PackageName: qPackage,
		Tokens:      strings.Split(qTokens, ","),
	}

	resObj, err := services.ExtractRelevantComments(request, logging)

	if err != nil {
		return Error(err)
	}

	res, err := json.Marshal(resObj)

	if err != nil {
		return Error(err)
	}

	writer.Write(res)

	return ErrorPkg{}
}

// Represents a POST action that handles a request body
type apiPostAction func(w http.ResponseWriter, body []byte, logging logging.Logging) ErrorPkg

// Represents a GET action that handles a request body
type apiGetAction func(w http.ResponseWriter, values url.Values, logging logging.Logging) ErrorPkg

// Mask errors and log them at the top level
func (config *Configuration) errorHandle(
	err error,
	writer http.ResponseWriter,
	logging logging.Logging) bool {
	if err != nil {
		errorString := err.Error()
		logging.Error(errorString)
		if !config.Development {
			errorString = "An internal server has occurred, please contact support@corporate.biz"
		}
		http.Error(writer, errorString, 500)
		return true
	}
	return false
}

// Mask errors in an ErrorPkg and log them at the top level
func (config *Configuration) errorPkgHandle(
	err ErrorPkg,
	writer http.ResponseWriter,
	logging logging.Logging) bool {
	if err.Error() {
		errorString := err.innerError.Error()
		logging.Error(errorString)
		if !config.Development &&
			err.httpStatus == 500 &&
			!err.isSanitized {
			errorString = "An internal server has occurred, please contact support@corporate.biz"
		}
		http.Error(writer, errorString, err.httpStatus)
		return true
	}
	return false
}

// Mask errors and log them at the top level
// also the central point to measure Http performance
func baseGetHandler(
	hander apiGetAction,
	config Configuration,
	logging logging.Logging,
	measurement Measurement) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		if request.Method == "GET" {
			start := time.Now()
			err := hander(writer, request.URL.Query(), logging)
			measurement.Log(request.URL.Path, time.Since(start).Nanoseconds()/1000000)

			if config.errorPkgHandle(err, writer, logging) {
				return
			}
		} else {
			http.Error(writer, "Unsupported HTTP method", 422)
		}
	}
}

// basic handling for all actions will withhold the actual error message
// if not in Development mode
func basePostHandler(
	handler apiPostAction,
	config Configuration,
	logging logging.Logging,
	measurement Measurement) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == "POST" {

			requestBody, err := ioutil.ReadAll(request.Body)
			defer request.Body.Close()

			if config.errorHandle(err, writer, logging) {
				return
			}

			start := time.Now()
			errPkg := handler(writer, requestBody, logging)
			measurement.Log(request.URL.Path, time.Since(start).Nanoseconds()/1000000)

			if config.errorPkgHandle(errPkg, writer, logging) {
				return
			}
		} else {
			http.Error(writer, "Unsupported HTTP method", 422)
		}
	}
}

// common settings for all Post routes
func commonPostRouteSetup(routes ...*mux.Route) {
	for _, route := range routes {
		route.
			Methods("POST").
			Headers("Content-Type", "application/json")
	}
}

// common settings for all Get routes
func commonGetRouteSetup(routes ...*mux.Route) {
	for _, route := range routes {
		route.
			Methods("GET")
	}
}

// This is the entry point for the server application, will start a server that provides comment parsing
// as a REST-ful service
func CommentParserHttpServer(
	config Configuration,
	logging logging.Logging,
	measurement Measurement) error {

	router := mux.NewRouter().StrictSlash(true)
	commonPostRouteSetup(
		router.HandleFunc("/parse", basePostHandler(ParseAction, config, logging, measurement)),
	)
	commonGetRouteSetup(
		router.HandleFunc("/", baseGetHandler(IndexAction, config, logging, measurement)),
	)
	http.Handle("/", router)
	srv := &http.Server{
		Handler:      router,
		Addr:         config.Address,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	srvError := srv.ListenAndServe()
	logging.Critical(srvError.Error())
	return nil
}
