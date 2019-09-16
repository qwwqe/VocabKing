package requests

const (
	ResultOK    = "ok"
	ResultError = "error"
)

type ErrorResponse struct {
	Result string `json:"result"`
	Data   error  `json:"data"`
}

func NewErrorResponse(err error) ErrorResponse {
	return ErrorResponse{
		Result: ResultError,
		Data:   err,
	}
}

type LoginForm struct {
	Data struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	} `json:"data" binding:"required"`
}

type LoginResponseData struct {
	Expiry int64  `json:"expiry"`
	Token  string `json:"token"`
}

type LoginResponse struct {
	Result string            `json:"result"`
	Data   LoginResponseData `json:"data" binding:"required"`
}

func NewLoginResponse(expiry int64, token string) LoginResponse {
	return LoginResponse{
		Result: ResultOK,
		Data: LoginResponseData{
			Expiry: expiry,
			Token:  token,
		},
	}
}

type StatsForm struct {
	Data struct{} `json:"data" binding:"required"`
}

type StatsResponseData struct {
	NumWords          int64    `json:"num_words"`
	NumDefinitions    int64    `json:"num_definitions"`
	NumUndefinedWords int64    `json:"num_undefined_words"`
	NumPhrases        int64    `json:"num_phrases"`
	NumTopics         int64    `json:"num_topics"`
	NumNotes          int64    `json:"num_notes"`
	UndefinedWords    []string `json:"undefined_words"`
	UntaggedWords     []string `json:"untagged_words"`
}

type StatsResponse struct {
	Result string            `json:"result"`
	Data   StatsResponseData `json:"data"`
}

func NewStatsResponse() StatsResponse {
	return StatsResponse{
		Result: ResultOK,
		Data:   StatsResponseData{},
	}
}

type SaveWordForm struct {
	Data struct{} `json:"data" binding:"required"`
}

type SaveWordResponseData struct {
}

type SaveWordResponse struct {
	Result string               `json:"result"`
	Data   SaveWordResponseData `json:"data"`
}

func NewSaveWordResponse() SaveWordResponse {
	return SaveWordResponse{
		Result: ResultOK,
		Data:   SaveWordResponseData{},
	}
}

type SavePictureForm struct {
	Data struct{} `json:"data" binding:"required"`
}

type SavePictureResponseData struct {
}

type SavePictureResponse struct {
	Result string                  `json:"result"`
	Data   SavePictureResponseData `json:"data"`
}

func NewSavePictureResponse() SavePictureResponse {
	return SavePictureResponse{
		Result: ResultOK,
		Data:   SavePictureResponseData{},
	}
}
