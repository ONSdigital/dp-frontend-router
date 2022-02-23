package helpers

import (
	"fmt"
	"strings"
)

// ReturnSecondSegmentFromPath returns the second segment of a path and assumes the path is formed /firstSegment/secondSegment
func ReturnSecondSegmentFromPath(path string) (secondSegment string, err error) {
	subs := strings.Split(path, "/")
	if len(subs) < 3 {
		err = fmt.Errorf("unable to extract secondSegment from path: %s", path)
		return
	}
	return subs[2], nil
}
