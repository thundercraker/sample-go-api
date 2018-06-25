package server

import (
	"bytes"
	"commentparser/logging"
	"commentparser/models"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestServer_PostParse_Success(t *testing.T) {

	reqBody := models.CommentParsingRequest{
		Tokens:      []string{"TODO"},
		PackageName: "fmt",
	}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/", bytes.NewReader(body))

	rrec := httptest.NewRecorder()

	config := Configuration{Development: false}

	handlerFunc := basePostHandler(ParseAction, config, logging.NewMockLogging(), NewBlankMeasurementTool())
	handler := http.HandlerFunc(handlerFunc)

	handler.ServeHTTP(rrec, req)

	str := fmt.Sprintf("%s", rrec.Body)

	assert.Equal(t, http.StatusOK, rrec.Code)
	assert.NotEqual(t, 0, len(str))

	var res models.CommentParsingResult
	json.Unmarshal(rrec.Body.Bytes(), &res)

	assert.Equal(t, "fmt", res.PackageName)
	assert.Equal(t, false, res.BinaryOnly)
	assert.NotNil(t, res.Matches["TODO"])
	assert.Equal(t, 2, len(res.Matches["TODO"]))
	assert.True(t, strings.Contains(res.Matches["TODO"][0].FileName, "fmt/format.go"))
	assert.Equal(t, "TODO: accept N and Ni independently?\n", res.Matches["TODO"][1].LineContent)
	assert.Equal(t, 747, res.Matches["TODO"][1].LineNumber)
}

func TestServer_PostParse_BadRequestJson(t *testing.T) {
	config := Configuration{Development: false}

	handlerFunc := basePostHandler(ParseAction, config, logging.NewMockLogging(), NewBlankMeasurementTool())
	handler := http.HandlerFunc(handlerFunc)

	{
		reqStr := "{{\"Package\":\"fmt\",\"Tokens\":\"TODO,voodoo\"}"
		req, _ := http.NewRequest("POST", "/", strings.NewReader(reqStr))

		rrec := httptest.NewRecorder()
		handler.ServeHTTP(rrec, req)

		assert.Equal(t, http.StatusBadRequest, rrec.Code)

		str := fmt.Sprintf("%s", rrec.Body)
		assert.Equal(t, "invalid character '{' looking for beginning of object key string\n", str)
	}
	{
		// missing parameter packages
		rrec := httptest.NewRecorder()
		reqStr := "{\"Tokens\":[\"TODO\",\"voodoo\"]}"
		req, _ := http.NewRequest("POST", "/", strings.NewReader(reqStr))
		handler.ServeHTTP(rrec, req)

		assert.Equal(t, http.StatusBadRequest, rrec.Code)

		str := fmt.Sprintf("%s", rrec.Body)
		assert.Equal(t, "The parameter `PackageName` cannot be empty\n", str)
	}
	{
		// missing parameter tokens
		rrec := httptest.NewRecorder()
		reqStr := "{\"Tokens\":[\"TODO\",\"voodoo\"]}"
		req, _ := http.NewRequest("POST", "/", strings.NewReader(reqStr))
		handler.ServeHTTP(rrec, req)

		assert.Equal(t, http.StatusBadRequest, rrec.Code)

		str := fmt.Sprintf("%s", rrec.Body)
		assert.Equal(t, "The parameter `PackageName` cannot be empty\n", str)
	}
}

func TestServer_GetIndex_Success(t *testing.T) {

	config := Configuration{Development: false}
	handlerFunc := baseGetHandler(IndexAction, config, logging.NewMockLogging(), NewBlankMeasurementTool())
	handler := http.HandlerFunc(handlerFunc)

	req, _ := http.NewRequest("GET", "/?package=fmt&tokens=TODO%2Cvoodo", nil)
	rrec := httptest.NewRecorder()

	handler.ServeHTTP(rrec, req)

	assert.Equal(t, http.StatusOK, rrec.Code)

	var res models.CommentParsingResult
	json.Unmarshal(rrec.Body.Bytes(), &res)

	assert.Equal(t, "fmt", res.PackageName)
	assert.Equal(t, false, res.BinaryOnly)
	assert.NotNil(t, res.Matches["TODO"])
	assert.Equal(t, 2, len(res.Matches["TODO"]))
	assert.True(t, strings.Contains(res.Matches["TODO"][0].FileName, "fmt/format.go"))
	assert.Equal(t, "TODO: accept N and Ni independently?\n", res.Matches["TODO"][1].LineContent)
	assert.Equal(t, 747, res.Matches["TODO"][1].LineNumber)
}

func TestServer_GetIndex_MissingQueryParams(t *testing.T) {

	config := Configuration{Development: false}
	handlerFunc := baseGetHandler(IndexAction, config, logging.NewMockLogging(), NewBlankMeasurementTool())
	handler := http.HandlerFunc(handlerFunc)

	{
		req, _ := http.NewRequest("GET", "/?tokens=TODO%2Cvoodo", nil)

		rrec := httptest.NewRecorder()

		handler.ServeHTTP(rrec, req)

		assert.Equal(t, http.StatusBadRequest, rrec.Code)
		resStr := fmt.Sprintf("%s", rrec.Body)
		assert.Equal(t, "the query must contain the parameter `package`\n", resStr)
	}
	{
		req, _ := http.NewRequest("GET", "/?package=fmt", nil)
		rrec := httptest.NewRecorder()

		handler.ServeHTTP(rrec, req)

		assert.Equal(t, http.StatusBadRequest, rrec.Code)
		resStr2 := fmt.Sprintf("%s", rrec.Body)
		assert.Equal(t, "the query must contain the parameter `tokens`\n", resStr2)
	}
}
