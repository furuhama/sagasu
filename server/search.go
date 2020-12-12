package server

type Result struct {
	Repo   string `json:"repo,string"`
	Path   string `json:"path,string"`
	Lineno uint32 `json:"lineno,int"`
	Text   string `json:"test,string"`
}

type Results struct {
	Count   uint32    `json:"count,int"`
	Results []*Result `json:"results,Result"`
}

func Search(query []string) (*Results, error) {
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
