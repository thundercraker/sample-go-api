package models

// the request model for Comment Parsing
type CommentParsingRequest struct {
	PackageName string   // the package name to search for comments
	Tokens      []string // the tokens/words to search for
}

// the result model for Comment Parsing
type CommentParsingResult struct {
	PackageName string                      // the package name in which the matches were made
	BinaryOnly  bool                        // true if the package was binary only
	Matches     map[string][]MatchedComment // the matched comments
}

// result model for a single matched comment
type MatchedComment struct {
	FileName    string // the file name where the comment was found
	LineNumber  int    // the line number where the comment was found
	LineContent string // the content of the comment itself
}
