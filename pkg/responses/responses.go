package responses

type Homepage struct {
	NumWords          int64
	NumDefinitions    int64
	NumUndefinedWords int64
	NumPhrases        int64
	NumTopics         int64
	NumNotes          int64
	UndefinedWords    []string
	UntaggedWords     []string
}
