package models

// Imgur ...
type Imgur struct {
	Data struct {
		ID          string      `json:"id"`
		Title       string      `json:"title"`
		Description string      `json:"description"`
		Datetime    int         `json:"datetime"`
		Type        string      `json:"type"`
		Animated    bool        `json:"animated"`
		Width       int         `json:"width"`
		Height      int         `json:"height"`
		Views       int         `json:"views"`
		Bandwidth   int         `json:"bandwidth"`
		Vote        interface{} `json:"vote"`
		Favorite    bool        `json:"favorite"`
		Nsfw        interface{} `json:"nsfw"`
		Section     interface{} `json:"section"`
		AccountURL  interface{} `json:"account_url"`
		AccountID   interface{} `json:"account_id"`
		IsAd        bool        `json:"is_ad"`
		InMostViral bool        `json:"in_most_viral"`
		Tags        []string    `json:"tags"`
		AdType      int         `json:"ad_type"`
		AdURL       string      `json:"ad_url"`
		InGallery   bool        `json:"in_gallery"`
		DeleteHash  string      `json:"delete_hash"`
		Name        string      `json:"name"`
		Link        string      `json:"link"`
	} `json:"data"`
	Success bool `json:"success"`
	Status  int  `json:"status"`
}
