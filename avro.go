package k6avro

import (
	"reflect"

	"github.com/grafana/sobek"
	"github.com/linkedin/goavro/v2"
	"go.k6.io/k6/js/common"
	"go.k6.io/k6/js/modules"

	"fmt"
)

func init() {
	modules.Register("k6/x/avro", new(RootModule))
}

// RootModule is the global module instance that will create module
// instances for each VU.
type RootModule struct{}

type Avro struct {
	// vu provides methods for accessing internal k6 objects for a VU
	vu modules.VU

	exports modules.Exports
}

// New returns a pointer to a new RootModule instance.
func New() *RootModule {
	return &RootModule{}
}

// NewModuleInstance implements the modules.Module interface returning a new instance for each VU.
func (*RootModule) NewModuleInstance(vu modules.VU) modules.Instance {
	mod := &Avro{
		vu: vu,
	}
	mod.exports.Default = mod
	return mod
}

func (a *Avro) Exports() modules.Exports {
	return a.exports
}

// XNewCodec parses an avro schema and creates a new codec.
func (a *Avro) NewCodec(c sobek.ConstructorCall) *sobek.Object {
	rt := a.vu.Runtime()

	if len(c.Arguments) != 1 || c.Argument(0).ExportType().Kind() != reflect.String {
		common.Throw(rt, fmt.Errorf("Avro schema missing"))
	}
	schema := c.Argument(0).String()

	codec, err := goavro.NewCodec(schema)
	if err != nil {
		return nil
	}
	return rt.ToValue(&AvroCodec{codec: codec}).ToObject(rt)
}

func (a *Avro) NewBytesCodec() *sobek.Object {
	schema := `{"type":"bytes"}`
	codec, err := a.newCodec(schema)

	if err != nil {
		return nil
	}
	rt := a.vu.Runtime()
	return rt.ToValue(&AvroCodec{codec: codec, schema: schema}).ToObject(rt)
}

func (a *Avro) newCodec(schema string) (*goavro.Codec, error) {
	codec, err := goavro.NewCodec(schema)
	if err != nil {
		return nil, err
	}
	return codec, nil
}

type AvroCodec struct {
	codec  *goavro.Codec
	schema string
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

// EncodeBytes encodes native binary data as byte type
func (ac *AvroCodec) EncodeBytes(data []byte) ([]byte, error) {
	if ac.schema == "" || ac.schema != `{"type":"bytes"}` {
		return nil, fmt.Errorf("EncodeNative only supports 'bytes' schema")
	}

	if ac.codec == nil {
		return nil, fmt.Errorf("codec not initialized")
	}

	// Encode []byte as native Avro "bytes"
	return ac.codec.BinaryFromNative(nil, data)
}

// Ensure the interfaces are implemented correctly.
var (
	_ modules.Instance = &Avro{}
	_ modules.Module   = &RootModule{}
)
