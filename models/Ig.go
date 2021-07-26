package models

type IgResponse struct {
	EntryData struct {
		PostPage []PostPageIg `json:"PostPage"`
	} `json:"entry_data"`
}

type PostPageIg struct {
	GraphQL struct {
		ShortcodeMedia struct {
			VideoUrl   string `json:"video_url"`
			DisplayUrl string `json:"display_url"`
		} `json:"shortcode_media"`
	} `json:"graphql"`
}
