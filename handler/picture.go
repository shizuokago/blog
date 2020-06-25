package handler

import (
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/tools/present"
)

func init() {
	present.Register("picture", parsePicture)
}

func parsePicture(ctx *present.Context, fileName string, lineno int, text string) (present.Elem, error) {

	args := strings.Fields(text)
	img := present.Image{URL: "/file/data/" + args[1]}
	a, err := parseArgs(fileName, lineno, args[2:])
	if err != nil {
		return nil, err
	}
	switch len(a) {
	case 0:
	case 2:
		if v, ok := a[0].(int); ok {
			img.Height = v
		}
		if v, ok := a[1].(int); ok {
			img.Width = v
		}
	default:
		return nil, fmt.Errorf("incorrect image invocation: %q", text)
	}
	return img, nil

}

func parseArgs(name string, line int, args []string) (res []interface{}, err error) {
	res = make([]interface{}, len(args))
	for i, v := range args {
		if len(v) == 0 {
			return nil, fmt.Errorf("%s:%d bad code argument %q", name, line, v)
		}
		switch v[0] {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			n, err := strconv.Atoi(v)
			if err != nil {
				return nil, fmt.Errorf("%s:%d bad code argument %q", name, line, v)
			}
			res[i] = n
		case '/':
			if len(v) < 2 || v[len(v)-1] != '/' {
				return nil, fmt.Errorf("%s:%d bad code argument %q", name, line, v)
			}
			res[i] = v
		case '$':
			res[i] = "$"
		case '_':
			if len(v) == 1 {
				// Do nothing; "_" indicates an intentionally empty parameter.
				break
			}
			fallthrough
		default:
			return nil, fmt.Errorf("%s:%d bad code argument %q", name, line, v)
		}
	}
	return
}
