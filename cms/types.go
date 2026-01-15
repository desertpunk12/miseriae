package cms

// BlogPost represents a blog post fetched from Google Drive
type BlogPost struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Date        string   `json:"date"` // ISO 8601 YYYY-MM-DD
	Tags        []string `json:"tags"`
	HTMLContent string   `json:"html_content"`
	ImageURL    string   `json:"image_url"` // Optional cover image from the doc? or metadata
	Type        string   `json:"type"`      // Tutorial, Life Update, Vlog
	Summary     string   `json:"summary"`
}

// CosplayAlbum represents a cosplay album from Google Photos
type CosplayAlbum struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`        // From "Title | Series"
	Series       string   `json:"series"`       // From "Title | Series"
	CoverImage   string   `json:"cover_image"`  // First image in album
	Images       []string `json:"images"`       // List of all image URLs
	Photographer string   `json:"photographer"` // Parsed from Description
	Assistant    string   `json:"assistant"`    // Parsed from Description
	Location     string   `json:"location"`     // Parsed from Description
	Description  string   `json:"description"`  // Parsed from Description
}
