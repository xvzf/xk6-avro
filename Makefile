build:
	xk6 build --with github.com/xvzf/xk6-avro=.

test: build
	./k6 run -i 1 -u 1 test/test.js

