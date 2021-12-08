# xk6-avro

![CI](https://github.com/xvzf/xk6-avro/actions/workflows/test.yaml/badge.svg)


This extension wraps the [goavro](https://github.com/linkedin/goavro`) library into a k6 extension.

You can build the extension using:
```bash
xk6 build --with github.com/xvzf/xk6-avro
```

## Example
```javascript
import avro from "k6/x/avro"

const codec = avro.NewCodec(`{
  "type": "record",
  "name": "LongList",
  "fields" : [
    {"name": "next", "type": ["null", "LongList"], "default": null}
  ]
}`)
let binary = codec.binaryFromTextual(`{"next":{"LongList":{}}}`)
console.log(binary)

let textual = codec.textualFromBinary(binary)
console.log(textual)

// When actually doing an HTTP request using the k6 http client, data has to be passed as ArrayBuffer:
const binaryAB = new Uint8Array(binary).buffer
http.post(`https://example.com/avro-receiver`, binaryAB, {headers: {"Content-Type": "binary/avro"}})
```
