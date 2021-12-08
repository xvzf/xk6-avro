import { fail } from "k6"
import avro from "k6/x/avro"


const testPayload = `{"next":{"LongList":{"next":null}}}`
const codec = avro.NewCodec(`{
  "type": "record",
  "name": "LongList",
  "fields" : [
    {"name": "next", "type": ["null", "LongList"], "default": null}
  ]
}`)

export function setup() {
  const binary = codec.binaryFromTextual(testPayload)
  const textual = codec.textualFromBinary(binary)

  if(testPayload !== textual) {
    console.log(`${textPayload} !== ${textual}`)
    fail("idempotence failed")
  }

}

export default function() {
  // noop for test
}
