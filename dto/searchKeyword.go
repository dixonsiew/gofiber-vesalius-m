package dto

type SearchKeywordDto struct {
    Keyword string `json:"keyword" default:""`
}

type SearchKeyword2Dto struct {
    Keyword  string `json:"keyword" default:""`
    Keyword2 string `json:"keyword2" default:""`
}
