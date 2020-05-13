package getter

import (
	"fmt"
	"github.com/hashicorp/go-getter/helper/url"
)

// Detect is a method used to detect if a Getter is a candidate for downloading an artifact
// by calling the Getter.Detect(*Request) method
func Detect(req *Request, getter Getter) (bool, error) {
	originalSrc := req.Src

	getForce, getSrc := getForcedGetter(req.Src)
	if getForce != "" {
		req.Forced = getForce
	}

	req.Src = getSrc
	ok, err := getter.Detect(req)
	if err != nil {
		return true, err
	}
	if !ok {
		// Write back the original source
		req.Src = originalSrc
		return ok, nil
	}

	result, detectSubdir := SourceDirSubdir(req.Src)

	// If we have a subdir from the detection, then prepend it to our
	// requested subdir.
	if detectSubdir != "" {
		u, err := url.Parse(result)
		if err != nil {
			return true, fmt.Errorf("Error parsing URL: %s", err)
		}
		u.Path += "//" + detectSubdir

		// a subdir may contain wildcards, but in order to support them we
		// have to ensure the path isn't escaped.
		u.RawPath = u.Path

		result = u.String()
	}

	req.Src = result
	return true, nil
}
