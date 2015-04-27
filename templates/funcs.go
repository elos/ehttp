package templates

import (
	"errors"
	"fmt"
	"html/template"
)

type Interpolator func(string) string
type HTMLInterpolator func(string) template.HTML

var (
	Dict = func(vals ...interface{}) (map[string]interface{}, error) {
		if len(vals)%2 != 0 {
			return nil, errors.New("Must pass element pairs")
		}

		dict := make(map[string]interface{}, len(vals)/2)

		for i := 0; i < len(vals); i += 2 {
			key, ok := vals[i].(string)
			if !ok {
				return nil, errors.New("Keys must be strings")
			}

			dict[key] = vals[i+1]
		}

		return dict, nil
	}

	CharSet = func(s string) template.HTML {
		return template.HTML(fmt.Sprintf("<meta charset=\"%s\">", s))
	}

	UTF8 = func() template.HTML {
		return CharSet("utf-8")
	}

	CSS = func(href string) template.HTML {
		return template.HTML(fmt.Sprintf("<link rel=\"stylesheet\" type=\"text/css\" href=\"%s\">", href))
	}

	JS = func(href string) template.HTML {
		return template.HTML(fmt.Sprintf("<script src=\"%s\"></script>", href))
	}
)
