package email

type Attachment struct {
	Name       string `json:"name"`
	Content    string `json:"content_type"`
	Base64Data string `json:"data"`
}
