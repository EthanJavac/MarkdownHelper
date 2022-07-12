package params

import "flag"

var (
	Path     = flag.String("path", "", "relative path to the markdown folder")
	Original = flag.Bool("place", false, "false means generate new md file here; and true means new md file will be generated in the original markdown folder")
)

func init() {
	flag.Parse()
}
