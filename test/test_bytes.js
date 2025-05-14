import avro from 'k6/x/avro';
import encoding from 'k6/encoding';

const keyCodec = avro.newCodec(`{
  "type": "record",
  "name": "TestKeyRecord",
  "namespace": "org.test.schemas",
  "fields": [
    { "name": "uniqueID", "type": "long" },
    { "name": "dummyStr", "type": "string" }
  ]
}`);

const byteCodec = avro.newBytesCodec();

export const options = {
  vus: 1,
  iterations: 1,
};

export default function () {
  const dummyStr = "me-dummy";
  const keys = [
    { uniqueID: 247589500, dummyStr },
    { uniqueID: 247589501, dummyStr }
  ];

  const byteEncodedArr = [];

  for (const record of keys) {
    const raw = keyCodec.binaryFromTextual(JSON.stringify(record));
    const encodedBytes = byteCodec.encodeBytes(raw);
    byteEncodedArr.push(...encodedBytes);
  }

  const combined = new Uint8Array(byteEncodedArr);
  const base64Final = encoding.b64encode(combined, 'url');

  console.log(`base64:${base64Final}`);
}
