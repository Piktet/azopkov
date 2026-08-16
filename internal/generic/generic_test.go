package generic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type resetTest struct{}

func (p resetTest) Reset() {

}
func TestPool(t *testing.T) {

	tests := []struct {
		name string
	}{
		{
			name: "positive",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := New[resetTest]()
			assert.NotEmpty(t, got)

			gotVal := resetTest{}
			got.Put(gotVal)
			wantVal := got.Get()

			assert.Equal(t, gotVal, wantVal)
		})
	}
}
