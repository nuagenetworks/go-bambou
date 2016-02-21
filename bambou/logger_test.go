package bambou

import (
	"testing"

	"github.com/ccding/go-logging/logging"
)

func TestLogger_Logger(t *testing.T) {

	l := Logger()

	if w := l.Level(); w != logging.ERROR {
		t.Error("Log level should be '%d' but is '%d'", logging.ERROR, w)
	}
}
