package dto

type SearchKeywordDto struct {
    Keyword string `json:"keyword" default:""`
}

type SearchKeyword2Dto struct {
    Keyword  string `json:"keyword" default:""`
    Keyword2 string `json:"keyword2" default:""`
}

type SearchKeyword3Dto struct {
    Keyword  string `json:"keyword" default:""`
    Keyword2 string `json:"keyword2" default:""`
    Keyword3 string `json:"keyword3" default:""`
}
