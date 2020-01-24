package uuid

import (
	"regexp"
	"testing"
)

func TestCreateUuidV4(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	uuid := CreateUuidV4()

	if len(uuid) != 36 {
		t.Errorf("Length of '%s' is %d, expected 36", uuid, len(uuid))
	}

	r, _ := regexp.Compile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
	matched := r.MatchString(uuid)

	if !matched {
		t.Errorf("Pattern of '%s' is not matching UUID v4.", uuid)
	}

}
