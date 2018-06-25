package services

import (
	"bufio"
	"bytes"
	"commentparser/logging"
	"commentparser/models"
	"github.com/stretchr/testify/assert"
	"github.com/tcnksm/go-binary-only-package"
	"strings"
	"testing"
)

func TestMainWithFmt(t *testing.T) {

	req := models.CommentParsingRequest{
		Tokens:      []string{"TODO"},
		PackageName: "fmt",
	}
	res, _ := ExtractRelevantComments(req, logging.NewMockLogging())

	assert.Equal(t, "fmt", res.PackageName)
	assert.Equal(t, false, res.BinaryOnly)
	assert.NotNil(t, res.Matches["TODO"])
	assert.Equal(t, 2, len(res.Matches["TODO"]))
	assert.True(t, strings.Contains(res.Matches["TODO"][0].FileName, "fmt/format.go"))
	assert.Equal(t, "TODO: accept N and Ni independently?\n", res.Matches["TODO"][1].LineContent)
	assert.Equal(t, 747, res.Matches["TODO"][1].LineNumber)
}

func TestMainWithFmt_BinaryOnly(t *testing.T) {

	hello.Hello("")
	req := models.CommentParsingRequest{
		Tokens:      []string{"TODO"},
		PackageName: "github.com/tcnksm/go-binary-only-package/src/github.com/tcnksm/hello",
	}
	res, _ := ExtractRelevantComments(req, logging.NewMockLogging())

	assert.Equal(t, true, res.BinaryOnly)
	assert.Equal(t, "hello", res.PackageName)
	assert.Nil(t, res.Matches)
}

func TestMainWithFmtNotFound(t *testing.T) {

	req := models.CommentParsingRequest{
		Tokens:      []string{"EXAMPLE"},
		PackageName: "fmt",
	}
	res, _ := ExtractRelevantComments(req, logging.NewMockLogging())
	assert.Equal(t, 0, len(res.Matches))
}

func TestMainWithUnknownPackage(t *testing.T) {

	buf := bufio.NewWriter(bytes.NewBufferString(""))
	req := models.CommentParsingRequest{
		Tokens:      []string{"TODO"},
		PackageName: "voodoo1231",
	}
	_, err := ExtractRelevantComments(req, logging.NewWriterLogging(buf))
	errStr := err.Error()
	assert.True(t,
		strings.Contains(errStr, "cannot find package \"voodoo1231\" in any of"),
		"The output should be empty")
}
