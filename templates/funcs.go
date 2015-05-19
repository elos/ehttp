package templates

import (
	"errors"
	"fmt"
	"html/template"
)

type Interpolator func(string) string
type HTMLInterpolator func(string) template.HTML

var (
	// Dict constructs a map out of the sequential key value pairs provided,
	// Used to construct custom context while in a template. i.e.,
	// if Dict was defined in the funcMap to be "dict"
	// {{template "CallThisTemplate" dict "user" $user "routes" .Data.Routes}}
	// now the CallThisTemplate template  gets a context with .user and .routes defined
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

	// Has checks whether the value is nil
	//
	// i.e., if Has was defined in the func map to be "has"
	// {{if has .Data.Flash}} <div class="flash"> {{.Data.Flash.Message}}  </div> {{end}}
	Has = func(v interface{}) bool {
		return v != nil
	}

	// CharSet constructs valid HTML using the charset specified by string, s
	CharSet = func(s string) template.HTML {
		return template.HTML(fmt.Sprintf("<meta charset=\"%s\">", s))
	}

	// UTF8 constrcuts valid HTML utf-8 charset definition
	UTF8 = func() template.HTML {
		return CharSet("utf-8")
	}

	// CSS constructs a valid stylesheet link html reference using the href string
	CSS = func(href string) template.HTML {
		return template.HTML(fmt.Sprintf("<link rel=\"stylesheet\" type=\"text/css\" href=\"%s\">", href))
	}

	// JS constructs a vlid script link html reference using the href string
	JS = func(href string) template.HTML {
		return template.HTML(fmt.Sprintf("<script src=\"%s\"></script>", href))
	}
)
