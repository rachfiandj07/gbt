package cors

import (
	"net"
	"net/url"
	"regexp"
)

type (
	// IChecker is the the interface that will be used by outside package to interact with this package
	IChecker interface {
		Check(string) bool
	}

	// Checker is the main object for csrf checker package
	checker struct {
		regex *regexp.Regexp
	}
)

var (
	checkerObj      checker
	whitelistOrigin map[string]string
)

// New function will create new csrf checker object
func New() IChecker {

	// if checkerObj.regex == nil {
	// 	reg, _ := regexp.Compile("^(https://)([a-zA-Z-]*)(.)$|(https://)([a-zA-Z-]*)(.t)([/?])|^(http://)([a-z.-]*)(.ndvl)|^(https://)([a-zA-Z0-9-]*)(staging.com)$|(https://)([a-zA-Z0-9-]*)(staging.com)([/?])|^(http://)(localhost-intools)")
	// 	checkerObj = checker{reg}
	// }

	return &checkerObj
}

// Check for origin and referer header content to validate csrf attack
// will return false if suspected from untrusted source
func (chk *checker) Check(origin string) bool {
	// if chk.regex.MatchString(origin) == true {
	// 	return true
	// }
	return isURLCors(origin)

}

// isURLCors handles common . URL except for the one without subdomain.
// This will only handle a URL with protocol.
func isURLCors(u string) bool {
	if u == "" {
		return false
	}

	uObject, err := url.Parse(u)
	if err != nil {
		return false
	}

	host, _, err := net.SplitHostPort(uObject.Host)
	if err != nil {
		host = uObject.Host
	}

	// if !strings.HasSuffix(host, ".com") &&
	// 	!strings.HasSuffix(host, ".net") &&
	// 	!strings.HasSuffix(host, ".tkpd") &&
	// 	!strings.HasSuffix(host, ".ndvl") &&
	// 	!strings.HasSuffix(host, ".id") {
	// 	return false
	// }

	return true
}

// we can add more origin client that we trusted here
func initWhitelistOrigin() {
	if whitelistOrigin == nil {
		whitelistOrigin = make(map[string]string)
	}
	// staging
	whitelistOrigin["https://.com"] = "https://.com"
}

// CheckWhitelistOrigin filter unknown origin
func CheckWhitelistOrigin(origin string) (trustedOrigin string, trusted bool) {
	if whitelistOrigin == nil {
		initWhitelistOrigin()
	}

	if trustedOrigin, trusted = whitelistOrigin[origin]; trusted {
		return
	}

	// empty string and false (untrusted)
	return
}
