package mopinion

// Meta is a struct that takes place in responses from Mopinion API, with the json tag `_meta`.
// It holds basic information about the resources.
type Meta struct {
	Code     int
	Count    int
	HasMore  bool `json:"has_more"`
	Message  string
	Next     interface{} // boolean false, or a string with a url
	Previous interface{}
	Total    int
}

// Token contains a string which is retrieved by Token service.
// It is used to generate a hmac signature.
type Token struct {
	Token string
}

// Account is a struct which reflects to mopinion Account resource.
type Account struct {
	Meta          Meta `json:"_meta,omitempty"`
	Name          string
	Package       string
	EndDate       string
	NumberUsers   int `json:"number_users"`
	NumberCharts  int `json:"number_charts"`
	NumberForms   int `json:"number_forms"`
	NumberReports int `json:"number_reports"`
	Reports       []Report
}

// Deployments is a struct which reflects to mopinion Deployments resource.
type Deployments struct {
	Meta        Meta `json:"_meta"`
	Deployments []Deployment
}

// Deployment is a struct which reflects to mopinion Deployment resource.
type Deployment struct {
	ID   int
	Key  string
	Name string
}

// Report is a struct which reflects to mopinion Report resource.
type Report struct {
	Meta        *Meta  `json:"_meta,omitempty"`
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Language    string `json:"language,omitempty"`
	Created     string `json:"created,omitempty"`
	Datasets    []Dataset
}

// Dataset is a struct which reflects to mopinion Dataset resource.
type Dataset struct {
	Meta        *Meta  `json:"_meta,omitempty"`
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	ReportID    int    `json:"report_id,omitempty"`
	Description string `json:"description,omitempty"`
	DataSource  string `json:"data_source,omitempty"`
}

// Fields is a struct which reflects to mopinion Fields resource.
type Fields struct {
	Meta Meta `json:"_meta,omitempty"`
	Data []FieldData
}

// FieldData is a struct which reflects to mopinion FieldData resource.
type FieldData struct {
	AnswerOptions *AnswerOptions `json:"answer_options"`
	AnswerValues  []string       `json:"answer_values"`
	DatasetID     int
	Key           string
	Label         string
	ReportID      int    `json:"report_id"`
	ShortLabel    string `json:"short_label"`
	Type          string
}

// AnswerOptions is a struct which reflects to mopinion AnswerOptions resource.
type AnswerOptions struct {
	Scale       int
	StartAtZero bool `json:"start_at_zero"`
	Type        string
}

// Feedback is a struct which reflects to mopinion Feedback resource.
type Feedback struct {
	Meta Meta `json:"_meta"`
	Data []FeedbackData
}

// FeedbackData is a struct which reflects to mopinion FeedbackData resource.
type FeedbackData struct {
	Created   string
	DatasetID int `json:"dataset_id"`
	ID        int
	ReportID  int `json:"report_id"`
	Tags      []string
	Fields    []FeedbackField
}

// FeedbackField is a struct which reflects to mopinion FeedbackField resource.
type FeedbackField struct {
	Key   string
	Label string
	Value interface{}
}

// DeleteResponse is a struct which reflects to mopinion DeleteResponse resource.
type DeleteResponse struct {
	Executed          bool
	ResourcesAffected map[string]interface{} `json:"resources_affected"`
}
