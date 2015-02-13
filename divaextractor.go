package divaextractor

import (
	"encoding/json"
	"io"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	gq "github.com/PuerkitoBio/goquery"
	"github.com/ikeikeikeike/gopkg/convert"
	behavior "github.com/ikeikeikeike/gopkg/net/http"
	"github.com/ikeikeikeike/gopkg/str"
)

const EndPoint = "http://ja.wikipedia.org/w/api.php"

func tee(r io.Reader, debug bool) io.Reader {
	if !debug {
		return r
	}
	return io.TeeReader(r, os.Stdout)
}

type Wikipedia struct {
	*behavior.UserBehavior
	doc *gq.Document

	Unit  string
	Debug bool
}

func NewWikipedia() *Wikipedia {
	return &Wikipedia{
		UserBehavior: behavior.NewUserBehavior(),
		Unit:         "cm",
		Debug:        false,
	}
}

func (w *Wikipedia) Doc(page string) (*gq.Document, error) {
	api := EndPoint + "?action=parse&format=json&prop=text&uselang=ja&page="

	resp, err := w.Behave(api + url.QueryEscape(page))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var r struct {
		Parse struct {
			Title string
			Text  struct {
				Source string `json:"*"`
			}
		}
	}
	err = json.NewDecoder(tee(resp.Body, w.Debug)).Decode(&r)
	if err != nil {
		return nil, err
	}

	reader := strings.NewReader(r.Parse.Text.Source)
	return gq.NewDocumentFromReader(reader)
}

func (w *Wikipedia) Do(page string) error {
	doc, err := w.Doc(page)

	if err != nil {
		return err
	}

	w.doc = doc
	return nil
}

func (w *Wikipedia) Birthday() (r time.Time) {
	w.doc.Find(`table th:contains(生年月日)`).Each(func(i int, s *gq.Selection) {
		text := str.Clean(s.Next().Text())
		r, _ = time.Parse("2006年1月2日", text)
	})
	return
}

func (w *Wikipedia) Blood() (r string) {
	w.doc.Find(`table th:contains(血液型)`).Each(func(i int, s *gq.Selection) {
		r = str.Clean(strings.Replace(s.Next().Text(), "型", "", -1))
	})
	return
}

func (w *Wikipedia) HW() (r string) {
	w.doc.Find(`table th:contains(身長), table th:contains(体重)`).Each(func(i int, s *gq.Selection) {
		text := s.Next().Text()
		if strings.Contains(text, w.Unit) {
			r = str.Clean(strings.Replace(str.Clean(strings.Replace(text, w.Unit, "", -1)), "kg", "", -1))
		}
	})
	return
}

func (w *Wikipedia) Height() (r int) {
	hw := strings.Split(w.HW(), "/")
	if len(hw) > 0 {
		r, _ = convert.StrTo(str.Clean(strings.Split(w.HW(), "/")[0])).Int()
	}
	return
}

func (w *Wikipedia) Weight() (r int) {
	hw := strings.Split(w.HW(), "/")
	if len(hw) > 1 {
		r, _ = convert.StrTo(str.Clean(hw[1])).Int()
	}
	return
}

func (w *Wikipedia) BWH() (r string) {
	w.doc.Find(`table th:contains(スリーサイズ)`).Each(func(i int, s *gq.Selection) {
		text := s.Next().Text()
		if strings.Contains(text, w.Unit) {
			r = str.Clean(strings.Replace(text, w.Unit, "", -1))
		}
	})
	return
}

func (w *Wikipedia) Bust() (r int) {
	bhw := strings.Split(w.BWH(), "-")
	if len(bhw) > 0 {
		r, _ = convert.StrTo(str.Clean(bhw[0])).Int()
	}
	return
}

func (w *Wikipedia) Waste() (r int) {
	bhw := strings.Split(w.BWH(), "-")
	if len(bhw) > 1 {
		r, _ = convert.StrTo(str.Clean(strings.Split(w.BWH(), "-")[1])).Int()
	}
	return
}

func (w *Wikipedia) Hip() (r int) {
	bhw := strings.Split(w.BWH(), "-")
	if len(bhw) > 2 {
		r, _ = convert.StrTo(str.Clean(strings.Split(w.BWH(), "-")[2])).Int()
	}
	return
}

func (w *Wikipedia) Bracup() (r string) {
	var re = regexp.MustCompile("(?:[a-z]|[A-Z]){1}")
	w.doc.Find(`table th:contains(ブラのサイズ)`).Each(func(i int, s *gq.Selection) {
		text := s.Next().Text()
		if re.MatchString(text) {
			r = str.Clean(strings.ToUpper(string([]rune(text)[0])))
		}
	})
	return
}
