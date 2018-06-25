/*
	The package services provides the main implementation of the Comment Parsing service. The object of
	this service is to go through all source files of a provided package (eg: "fmt") and going through
	all the present comments to find any of the provided tokens/words. All such matches of tokens in comments
	will be returned
*/
package services

import (
	"commentparser/logging"
	"commentparser/models"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// Import a package from a dir and return it if it is valid (not binary or a command)
func importPkg(path, dir string, logging logging.Logging) (*build.Package, error) {

	p, err := build.Import(path, dir, build.ImportComment)
	if err != nil {
		return nil, err
	}

	// we can tell if the package is binary only alongside the rest of the
	// comment parsing

	if p.IsCommand() {
		logging.Debug("The package %s is a command", p.Name)
		return nil, nil
	}

	return p, nil
}

// Go through all the sources at filename and if there are any comments containing
// the terms in search terms, return the file name, line number and the comment itself
func extractCommentsWithTerms(
	searchTerms []string,
	fileName string,
	logging logging.Logging) (map[string][]models.MatchedComment, bool) {

	logging.Debug("Beginning extraction of %s", fileName)
	fileSet := token.NewFileSet()
	f, err := parser.ParseFile(fileSet, fileName, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	var resultMap map[string][]models.MatchedComment
	resultMap = make(map[string][]models.MatchedComment)

	commentMap := ast.NewCommentMap(fileSet, f, f.Comments)
	for commentGroupIdx, commentGroups := range commentMap {
		fileSet := fileSet.File(commentGroupIdx.Pos())
		for _, commentGroup := range commentGroups {
			commentGroupText := commentGroup.Text()
			for _, searchTerm := range searchTerms {
				if strings.Contains(commentGroupText, "go:binary-only-package") {
					logging.Info("Found binary-only flag in %s", fileName)
					return nil, true // this is a binary only package
				} else if strings.Contains(commentGroupText, searchTerm) {
					childItems := &[]models.MatchedComment{}
					if matchedTokens, found := resultMap[searchTerm]; found {
						childItems = &matchedTokens
					}
					resultMap[searchTerm] = append(*childItems, models.MatchedComment{
						FileName:    fileName,
						LineNumber:  fileSet.Position(commentGroup.Pos()).Line,
						LineContent: commentGroupText,
					})
				}
			}
		}
	}

	return resultMap, false
}

// Go through all the sources belonging to the provided package name and if there are any comments containing
// the terms in search terms, return the file name, line number and the comment itself
func ExtractRelevantComments(
	request models.CommentParsingRequest,
	logging logging.Logging) (models.CommentParsingResult, error) {

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	logging.Debug("Beginning extraction of package %s", request.PackageName)
	p, err := importPkg(request.PackageName, dir, logging)

	if err != nil {
		return models.CommentParsingResult{}, err
	}

	resultMap := make(map[string][]models.MatchedComment)

	result := models.CommentParsingResult{
		PackageName: p.Name,
		BinaryOnly:  false,
	}

	if p == nil {
		return result, nil
	}

	for _, goFile := range p.GoFiles {
		matchesForTokens, binaryOnly := extractCommentsWithTerms(request.Tokens, filepath.Join(p.Dir, goFile), logging)
		if binaryOnly {
			result.Matches = nil
			result.BinaryOnly = true
			return result, nil
			break
		} else if len(matchesForTokens) > 0 {
			for key, val := range matchesForTokens {
				if currentMatches, found := resultMap[key]; found {
					resultMap[key] = append(currentMatches, val...)
				} else {
					resultMap[key] = val
				}
			}
		}
	}
	result.Matches = resultMap
	return result, nil
}
