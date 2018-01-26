package placeholder

import (
	"testing"
)

func BenchmarkImage(b *testing.B) {
	for n := 0; n < b.N; n++ {
		GenerateImage(DefaultImageOptions)
	}
}
