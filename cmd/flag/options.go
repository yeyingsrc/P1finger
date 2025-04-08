package flag

type Options struct {
	Url     string
	UrlFile string

	// proxy usage
	SocksProxy string
	HttpProxy  string

	FingerDir  string
	FingerOnly string

	Debug bool

	// fingers file
	P1fingerFile string
	CollectMode  string

	Rate int

	Update bool
	Output string
}
