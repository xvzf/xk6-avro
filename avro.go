package k6avro

import (
	"context"
	"github.com/linkedin/goavro/v2"
	"go.k6.io/k6/js/common"
	"go.k6.io/k6/js/modules"

  "fmt"
)

func init() {
	modules.Register("k6/x/avro", new(Avro))
}

type Avro struct{}

type AvroCodec struct {
	codec *goavro.Codec
}

// XNewCodec parses an avro schema and creates a new codec.
func (*Avro) XNewCodec(ctxPtr *context.Context, schema string) interface{} {
	rt := common.GetRuntime(*ctxPtr)
	c, err := goavro.NewCodec(schema)
	if err != nil {
		return nil
	}
	return common.Bind(rt, &AvroCodec{codec: c}, ctxPtr)
}

// BinaryFromTextual converts an avro json encoded document into its binary representation.
func (ac *AvroCodec) BinaryFromTextual(textual string) []byte {
	native, _, err := ac.codec.NativeFromTextual([]byte(textual))
	if err != nil {
		return []byte(fmt.Sprintf("%v", err))
	}
	binary, err := ac.codec.BinaryFromNative(nil, native)
	if err != nil {
		return []byte(fmt.Sprintf("%v", err))
	}
	return binary
}

// TextualFromBinary converts an avro binary encoded document into its textual representation.
func (ac *AvroCodec) TextualFromBinary(binary []byte) string {
	native, _, err := ac.codec.NativeFromBinary(binary)
	if err != nil {
		return fmt.Sprintf("%v", err)
	}
	textual, err := ac.codec.TextualFromNative(nil, native)
	if err != nil {
		return fmt.Sprintf("%v", err)
	}
	return string(textual)
}
