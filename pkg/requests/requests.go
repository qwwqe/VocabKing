package requests

type Action string

const (
	ActionAddWord    Action = "add_word"
	ActionAddPicture Action = "add_picture"
)

type StatsForm struct {
}

type StatsResponse struct {
	NumWords          int64    `json:"num_words"`
	NumDefinitions    int64    `json:"num_definitions"`
	NumUndefinedWords int64    `json:"num_undefined_words"`
	NumPhrases        int64    `json:"num_phrases"`
	NumTopics         int64    `json:"num_topics"`
	NumNotes          int64    `json:"num_notes"`
	UndefinedWords    []string `json:"undefined_words"`
	UntaggedWords     []string `json:"untagged_words"`
}

type AddWordForm struct {
	Action Action `json:"action"`
	Data   struct {
		Word          string   `json:"word"`
		Pronunciation string   `json:"pronunciation"`
		Definition    string   `json:"definition"`
		Tags          []string `json:"tags"`
		Note          string   `json:"note"`
	} `json:"data"`
}

type AddWordResponse struct {
	Success Status `json:"success"`
	Data    struct {
		DefinitionID int64 `json:"definition_id"`
	} `json:"data"`
}

type AddPictureForm struct {
	Action Action `json:"action"`
	Data   struct {
		DefinitionID int64  `json:"definition_id"`
		Picture      []byte `json:"picture"`
	} `json:"data"`
}

type AddPictureResponse struct {
	Success Status `json:"success"`
	Data    struct {
		PictureID int64 `json:"definition_id"`
	} `json:"data"`
}
