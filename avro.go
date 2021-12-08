package k6avro

import (
	"context"
	"github.com/linkedin/goavro/v2"
	"go.k6.io/k6/js/common"
	"go.k6.io/k6/js/modules"
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
func (ac *AvroCodec) BinaryFromTextual(textual string) string {
	native, _, err := ac.codec.NativeFromTextual([]byte(textual))
	if err != nil {
		return ""
	}
	binary, err := ac.codec.BinaryFromNative(nil, native)
	if err != nil {
		return ""
	}
	return string(binary)
}

// TextualFromBinary converts an avro binary encoded document into its textual representation.
func (ac *AvroCodec) TextualFromBinary(binary []byte) string {
	native, _, err := ac.codec.NativeFromBinary(binary)
	if err != nil {
		return ""
	}
	textual, err := ac.codec.TextualFromNative(nil, native)
	if err != nil {
		return ""
	}
	return string(textual)
}
