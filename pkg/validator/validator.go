package validator

import (
	"log"
	"net/url"
	"regexp"
)

type TokenMap map[int]string

type PreprocessFunc func(string) (string, error)

type Validator interface {
	ValidateLine(string) []string
	ValidateTokens([]string) (TokenMap, bool)
	GetToken(int, TokenMap) string
}

var (
	regexLine                        = `(\S+)\s+[-]\s+(\S+)\s+\[(.*)\]\s+"(\S+)\s+(\S+)\s+(\S+)"\s+([0-9]+)\s+([0-9]+)\s+"(.*)"`
	regexLineCompiled *regexp.Regexp = nil
	regexMap                         = map[int]string{
		// ip addr
		0: `^(([0-9]{1,3})\.([0-9]{1,3})\.([0-9]{1,3})\.([0-9]{1,3})|(([0-9a-fA-F]{0,4}):){1,7}[0-9a-fA-F]{0,4})$`,
		// username
		1: `^([0-9A-Za-z\._-]+)$`,
		// date
		2: `^(([0-9]{1,2})\/([A-Za-z]{3})\/([1-9][0-9]{3}):([0-9]+):([0-9]+):([0-9]+)\s+([-+][0-9]{4}))$`,
		// http method
		3: `^([A-Za-z]+)$`,
		// uri
		4: `^(\/?([0-9a-zA-Z\._][~0-9a-zA-Z\#\+\%@\.\/_-]*))?(\?[0-9a-zA-Z\+\%@\/&\[\],;=_-]+)?$`,
		// http version
		5: `^(HTTP/[0-9]+\.[0-9]+)$`,
		// http status
		6: `^([1-5][0-9]{2})$`,
		// response time (ms)
		7: `^([1-9][0-9]+)$`,
		// user agent
		8: `^(.*)$`,
	}
	regexMapCompiled                       = make(map[int]*regexp.Regexp)
	regexMapDefault                        = `(.*)`
	regexMapDefaultCompiled *regexp.Regexp = nil

	preprocessMap = map[int]PreprocessFunc{
		4: UrlDecode,
	}
)

func init() {
	regexLineCompiled = regexp.MustCompile(regexLine)
	for k, v := range regexMap {
		regexMapCompiled[k] = regexp.MustCompile(v)
	}
	regexMapDefaultCompiled = regexp.MustCompile(regexMapDefault)
}

type RegexValidator struct {
	l *log.Logger
}

func UrlDecode(path string) (string, error) {
	decodedPath, e := url.PathUnescape(path)
	return decodedPath, e
}

func NewRegexValidator(logger *log.Logger) *RegexValidator {
	return &RegexValidator{
		l: logger,
	}
}

func (rv *RegexValidator) ValidateLine(in string) []string {
	match := regexLineCompiled.FindStringSubmatch(in)
	if len(match) < 1 {
		//rv.l.Printf("error validating line '%s' against regex '%s'", in, regexLine)
		return nil
	}
	return match[1:]
}

func (rv *RegexValidator) ValidateTokens(tokens []string) (TokenMap, bool) {
	returnMap := make(TokenMap)
	retval := true
	var processedToken string
	var e error
	for i, token := range tokens {
		// preprocess - if preprocess func available
		f, found := preprocessMap[i]
		processedToken = token
		if found {
			processedToken, e = f(token)
			if e != nil {
				retval = false
				break
			}
		}
		r, found := regexMapCompiled[i]
		var match []string
		if !found {
			r = regexMapDefaultCompiled
		}
		match = r.FindStringSubmatch(processedToken)
		if len(match) < 2 {
			retval = false
			//rv.l.Printf("error validating token[%v] '%s' against regex '%s'", i, r, processedToken)
			break
		}
		returnMap[i] = match[1]
	}
	return returnMap, retval
}

func (rv *RegexValidator) GetToken(key int, tm TokenMap) string {
	return tm[key]
}
