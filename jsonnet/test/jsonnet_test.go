package test

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/google/go-jsonnet"
	"github.com/stretchr/testify/require"
	"github.com/xeipuuv/gojsonschema"
)

func TestJSONNet(t *testing.T) {
	vm := jsonnet.MakeVM()
	fjn, err := os.Open("root.jsonnet")
	require.NoError(t, err)

	defer fjn.Close()

	b, err := io.ReadAll(fjn)
	require.NoError(t, err)

	j, err := vm.EvaluateAnonymousSnippet("root.jsonnet", string(b))
	require.NoError(t, err)

	fschema, err := os.Open("../../dzhttpjsonschema/schema.json")
	require.NoError(t, err)

	defer fschema.Close()

	b, err = io.ReadAll(fschema)
	require.NoError(t, err)

	schema := gojsonschema.NewBytesLoader(b)
	results, err := gojsonschema.Validate(schema, gojsonschema.NewStringLoader(j))

	require.NoError(t, err)
	require.Empty(t, results.Errors())

	fmt.Println(j)
}
