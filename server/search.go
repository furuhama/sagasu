package server

import (
	"log"
	"os"
	"unsafe"

	"github.com/furuhama/codesearch/index"
	"github.com/furuhama/codesearch/regexp"
)

type Result struct {
	Repo   string `json:"repo,string"`
	Path   string `json:"path,string"`
	Lineno int    `json:"lineno,int"`
	Text   string `json:"test,string"`
}

type Results struct {
	Count   int       `json:"count,int"`
	Results []*Result `json:"results,Result"`
}

func Search(query []string, iFlag bool, bruteFlag bool) (*Results, error) {
	searchResults := search(query, iFlag, bruteFlag)
	// Here's sample result
	return &Results{
		Count: 10,
		Results: []*Result{
			{
				Repo:   "sample_application",
				Path:   "app/controllers/sample_controller.rb",
				Lineno: 4,
				Text:   "before_action :sample_filter",
			},
		},
	}, nil
}

// search is almost the same as Main() in furuhama/codesearch/cmd/csearch/csearch.go
func search(query []string, iFlag bool, bruteFlag bool) []*regexp.SearchResult {
	g := regexp.Grep{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	g.AddFlags()

	pat := "(?m)" + query[0]
	if iFlag {
		pat = "(?i)" + pat
	}
	re, err := regexp.Compile(pat)
	// TODO: Need to handle errors
	if err != nil {
		log.Fatal(err)
	}
	g.Regexp = re
	var fre *regexp.Regexp
	q := index.RegexpQuery(re.Syntax)
	log.Printf("query: %s\n", q)

	ix := index.Open(index.File())
	ix.Verbose = true
	var post []uint32
	if bruteFlag {
		post = ix.PostingQuery(&index.Query{Op: index.QAll})
	} else {
		post = ix.PostingQuery(q)
	}
	log.Printf("post query identified %d possible files\n", len(post))

	if fre != nil {
		fnames := make([]uint32, 0, len(post))

		for _, fileid := range post {
			name := ix.Name(fileid)
			if fre.MatchString(name, true, true) < 0 {
				continue
			}
			fnames = append(fnames, fileid)
		}

		log.Printf("filename regexp matched %d files\n", len(fnames))
		post = fnames
	}

	var results []*regexp.SearchResult

	for _, fileid := range post {
		name := ix.Name(fileid)
		g.File(name)
		results, err = g.FileToSearchResult(name)
	}

	return results
}

func convert(searchResults []*regexp.SearchResult) *Results {
	var (
		count   = 0
		results *Results
	)

	for _, elem := range searchResults {
		count++
		result := Result{
			Repo:   "iv_web",
			Path:   elem.Path,
			Lineno: elem.Lineno,
			Text:   *(*string)(unsafe.Pointer(&elem.Line)),
		}
		results.Results = append(results.Results, &result)
	}
	results.Count = count

	return results
}
