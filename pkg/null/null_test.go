package null

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type MyInt int

func Test_NullScan(t *testing.T) {
	i := Null[MyInt]{}

	require.NoError(
		t,
		i.Scan("123432"),
	)

	require.Equal(
		t,
		i.Val,
		MyInt(123432),
	)
}
