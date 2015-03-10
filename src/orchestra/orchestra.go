package orchestra

import "strconv"
import "net/http"
import "regexp"
import "strings"

type Orchestra struct {
	Address string
	Port int
	Handles map[string]HandlerFunc
	Server http.Server
}

type HandlerFunc func(http.ResponseWriter, *http.Request, map[string]string)
func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request, p map[string]string) {
	f(w, r, p)
}

func (o Orchestra) String() string {
	return o.Address + ":" + strconv.FormatInt(int64(o.Port), 10)
}

func (o Orchestra) HandleFunc(pattern string, fn HandlerFunc) {
	o.Handles[pattern] = fn
}

func (o Orchestra) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var found = false
	for pattern, fn := range o.Handles {
		if URLMatchesPattern(r.URL.Path, pattern) {
			found = true
			fn(w, r, URLToParameters(r.URL.Path, pattern))
			break
		}
	}
	if !found {
		w.Write([]byte("4-oh-4"))
	}
}

func URLMatchesPattern(url, pattern string) bool {
	var reParamNames *regexp.Regexp
	var _ error
	var rePattern *regexp.Regexp
	var parsedPattern string = "^" + pattern + "$"
	var parameterNames []string
	reParamNames, _ = regexp.Compile("(:[a-zA-Z]+[a-zA-Z0-9_\\-]*)")
	parameterNames = reParamNames.FindAllString(pattern, -1)
	for _, parameter := range parameterNames {
		parsedPattern = strings.Replace(
			parsedPattern,
			parameter,
			"(?P<" + parameter[1:] + ">[a-z0-9]+)",
			-1,
		)
	}
	rePattern, _ = regexp.Compile(parsedPattern)
	return rePattern.MatchString(url)
}

func URLToParameters(url, pattern string) map[string]string {
	var out map[string]string = make(map[string]string)
	var reParamNames *regexp.Regexp
	var _ error
	var rePattern *regexp.Regexp
	var parsedPattern string = "^" + pattern + "$"
	var parameterNames []string
	var match []string
	reParamNames, _ = regexp.Compile("(:[a-zA-Z]+[a-zA-Z0-9_\\-]*)")
	parameterNames = reParamNames.FindAllString(pattern, -1)
	for idx, parameter := range parameterNames {
		parsedPattern = strings.Replace(
			parsedPattern,
			parameter,
			"(?P<" + parameter[1:] + ">[a-zA-Z0-9_\\-]+)",
			-1 + (0 * idx),
		)
	}
	rePattern, _ = regexp.Compile(parsedPattern)
	match = rePattern.FindStringSubmatch(url)
	if match != nil {
		for i, name := range rePattern.SubexpNames() {
			if i == 0 || name == "" {
				continue
			}
			out[name] = match[i]
		}
	}
	return out
}

func NewOrchestra(address string, port int) *Orchestra {
	var o Orchestra
	o.Address = address
	o.Port = port
	o.Handles = make(map[string]HandlerFunc)
	return &o
}

func (o Orchestra) ListenAndServe() {
	o.Server = http.Server{Addr: o.String(), Handler: o}
	o.Server.ListenAndServe()
}
