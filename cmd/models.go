package cmd

type promptContent struct {
	errorMsg string
	label    string
}

type bbbContent struct {
	folder        string
	rawUrl        string
	jsonName      string
	meetingId     string
	downloadLinks []downloadLink
}

type downloadLink struct {
	name string
	ext  string
	link string
}
