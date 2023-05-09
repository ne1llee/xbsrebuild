package booksourcecheck

import (
	"encoding/json"
	"fmt"
	"time"
	"xbsrebuild/common/bookquery"

	"github.com/dop251/goja"
)

type SearchBookCheck struct {
	session    *bookquery.CachedHTTPClient
	searchBook *SearchBook
}

type RequestParams struct {
	PageIndex int    `json:"pageIndex"`
	KeyWord   string `json:"keyWord"`
	Offset    int    `json:"offset"`
}

func NewSearchBookCheck(searchBook *SearchBook) (*SearchBookCheck, error) {
	sbc := &SearchBookCheck{
		searchBook: searchBook,
	}

	session, err := bookquery.NewCachedHTTPClient(true, "")
	if err != nil {
		return nil, err
	}
	sbc.session = session

	return sbc, nil
}

func (sbc *SearchBookCheck) generateUrl() (string, error) {
	var url string
	vm := goja.New()

	params := map[string]string{
		"pageIndex": "1",
		"keyWord":   "都市",
		"offset":    "0",
	}

	_, err := vm.RunString(sbc.searchBook.RequestJavascript)
	if err != nil {
		return "", err
	}

	fn, ok := goja.AssertFunction(vm.Get(sbc.searchBook.RequestFunction))
	if !ok {
		return "", fmt.Errorf("%s is not a function", sbc.searchBook.RequestFunction)
	}

	out, err := fn(nil, vm.ToValue(nil), vm.ToValue(params))
	if err != nil {
		return "", fmt.Errorf("failed call %s function, err: %v", sbc.searchBook.RequestFunction, err)
	}
	url = out.String()
	return url, nil
}

func (sbc *SearchBookCheck) responseCheck(res []byte) error {
	vm := goja.New()
	_, err := vm.RunString(sbc.searchBook.ResponseJavascript)
	if err != nil {
		return err
	}
	fn, ok := goja.AssertFunction(vm.Get(sbc.searchBook.ResponseFunction))
	if !ok {
		return fmt.Errorf("%s is not a function", sbc.searchBook.ResponseFunction)
	}
	var resObj interface{}
	if sbc.searchBook.ResponseFormatType == "json" {
		err = json.Unmarshal(res, &resObj)
		if err != nil {
			return fmt.Errorf("failed call unmarshal respose, err: %v", err)
		}
	} else {
		resObj = res
	}
	response, err := fn(nil, vm.ToValue(nil), vm.ToValue(nil), vm.ToValue(resObj))
	if err != nil {
		return fmt.Errorf("failed call %s function, err: %v", sbc.searchBook.ResponseFunction, err)
	}
	fmt.Println(response.Export())
	return nil
}

func (sbc *SearchBookCheck) Check() error {
	url, err := sbc.generateUrl()
	if err != nil {
		return fmt.Errorf("failed call generateUrl, err: %v", err)
	}
	fmt.Printf("request url: %s\n", url)

	res, err := sbc.session.CachedHttpGet(url, nil, 60*time.Second)

	if err != nil {
		return fmt.Errorf("failed call CachedHttpGet, err: %v", err)
	}
	err = sbc.responseCheck(res)
	if err != nil {
		return fmt.Errorf("failed call responseCheck, err: %v", err)
	}
	return nil
}
