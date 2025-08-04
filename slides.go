package main

type slideshow struct {
	NoAutoPlay bool
	PrimaryFooter,
	SecondaryFooter,
	Script,
	Styles opt[string]
	CustomDuration opt[float64]

	Slides []slide
}

type slide struct {
	IsTitleSlide   bool
	CustomDuration opt[float64]
	Image, Id      opt[string]

	HeaderLines []string
	Bullets     []string
	Lines       []string
}
