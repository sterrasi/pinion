package pinion

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNormalize(t *testing.T) {

	assert.Equal(t, Normalize(""), "")
	assert.Equal(t, Normalize("nothing to normalize"), "nothing to normalize")
	assert.Equal(t, Normalize("To LoWeRCasE"), "to lowercase")
	assert.Equal(t, Normalize("  trim whitespace   "), "trim whitespace")
	assert.Equal(t, Normalize("		trim Tabs	"), "trim tabs")
	assert.Equal(t, Normalize("		trim NewLines\n\r	"), "trim newlines")

}
