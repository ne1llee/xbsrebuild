package booksourcecheck

import (
	"encoding/json"
	"fmt"
)

type BookDetail struct {
	ActionID string `json:"actionID"`
	ParserID string `json:"parserID"`
}

type RequestFilters struct {
	NewTime string `json:"new_time"`
	Channel string `json:"channel"`
	Name    string `json:"name"`
	Time    string `json:"time"`
	Filter  string `json:"filter"`
}

type MoreKeys struct {
	RequestFilters map[string]RequestFilters `json:"requestFilters"`
}

type BookWorld struct {
	ActionID           string   `json:"actionID"`
	ValidConfig        string   `json:"validConfig"`
	ResponseJavascript string   `json:"responseJavascript"`
	ResponseFunction   string   `json:"responseFunction"`
	MoreKeys           MoreKeys `json:"moreKeys"`
	Host               string   `json:"host"`
	SIndex             int      `json:"_sIndex"`
	RequestFunction    string   `json:"requestFunction"`
	RequestJavascript  string   `json:"requestJavascript"`
	ParserID           string   `json:"parserID"`
	ResponseFormatType string   `json:"responseFormatType"`
}

type ChapterList struct {
	ValidConfig        string `json:"validConfig"`
	ActionID           string `json:"actionID"`
	ResponseJavascript string `json:"responseJavascript"`
	ResponseFunction   string `json:"responseFunction"`
	RequestFunction    string `json:"requestFunction"`
	Host               string `json:"host"`
	ResponseFormatType string `json:"responseFormatType"`
	ParserID           string `json:"parserID"`
	RequestJavascript  string `json:"requestJavascript"`
}

type SearchShudan struct {
	ActionID string `json:"actionID"`
	ParserID string `json:"parserID"`
}

type RelatedWord struct {
	ActionID string `json:"actionID"`
	ParserID string `json:"parserID"`
}

type ShudanDetail struct {
	ActionID string `json:"actionID"`
	ParserID string `json:"parserID"`
}

type SearchBook struct {
	ResponseJavascript string `json:"responseJavascript"`
	ActionID           string `json:"actionID"`
	ValidConfig        string `json:"validConfig"`
	ResponseFunction   string `json:"responseFunction"`
	RequestFunction    string `json:"requestFunction"`
	Host               string `json:"host"`
	ResponseFormatType string `json:"responseFormatType"`
	ParserID           string `json:"parserID"`
	RequestJavascript  string `json:"requestJavascript"`
}

type ChapterContent struct {
	ValidConfig        string `json:"validConfig"`
	ActionID           string `json:"actionID"`
	ResponseJavascript string `json:"responseJavascript"`
	ResponseFunction   string `json:"responseFunction"`
	RequestFunction    string `json:"requestFunction"`
	Host               string `json:"host"`
	RequestJavascript  string `json:"requestJavascript"`
	ParserID           string `json:"parserID"`
}

type ShupingList struct {
	ActionID string `json:"actionID"`
	ParserID string `json:"parserID"`
}

type ShudanList struct{}

type ShupingHome struct {
	ActionID string `json:"actionID"`
	ParserID string `json:"parserID"`
}

type BookSource struct {
	BookDetail     BookDetail           `json:"bookDetail"`
	BookWorld      map[string]BookWorld `json:"bookWorld"`
	Weight         string               `json:"weight"`
	MiniAppVersion string               `json:"miniAppVersion"`
	ChapterList    ChapterList          `json:"chapterList"`
	SearchShudan   SearchShudan         `json:"searchShudan"`
	RelatedWord    RelatedWord          `json:"relatedWord"`
	Enable         bool                 `json:"enable"`
	SourceName     string               `json:"sourceName"`
	SourceUrl      string               `json:"sourceUrl"`
	ShudanDetail   ShudanDetail         `json:"shudanDetail"`
	LastModifyTime string               `json:"lastModifyTime"`
	SearchBook     SearchBook           `json:"searchBook"`
	ChapterContent ChapterContent       `json:"chapterContent"`
	ShupingList    ShupingList          `json:"shupingList"`
	Password       string               `json:"password"`
	ShudanList     ShudanList           `json:"shudanList"`
	Desc           string               `json:"desc"`
	ShupingHome    ShupingHome          `json:"shupingHome"`
	AuthorId       string               `json:"authorId"`
}

type BookSourceCheck struct {
	source map[string]*BookSource
}

func NewBookSourceCheck(bookSource []byte) (*BookSourceCheck, error) {
	bsc := &BookSourceCheck{
		source: map[string]*BookSource{},
	}
	err := json.Unmarshal(bookSource, &bsc.source)
	if err != nil {
		return nil, err
	}
	return bsc, nil
}

func (bsc *BookSourceCheck) Check() {
	for _, v := range bsc.source {
		sbc, err := NewSearchBookCheck(&v.SearchBook)
		if err != nil {
			fmt.Println(err)
			continue
		}
		err = sbc.Check()
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
}
