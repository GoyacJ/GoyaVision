package dto

type StreamListQuery struct {
	Enabled *bool `query:"enabled"`
}

type StreamCreateReq struct {
	URL     string `json:"url"`
	Name    string `json:"name"`
	Enabled *bool  `json:"enabled"`
}

type StreamUpdateReq struct {
	URL     *string `json:"url"`
	Name    *string `json:"name"`
	Enabled *bool   `json:"enabled"`
}
