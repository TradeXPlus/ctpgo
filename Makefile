ctpgo:
	go build -o ./bin/ctpgo main.go

all:
	go get -u github.com/axgle/mahonia && go get -u github.com/tidwall/gjson

	swig -go -cgo -intgosize 64 -module lib -I./lib/linux64 -c++ -outdir ./lib -o ./lib/linux64/ctpgo_wrap.cxx -oh ./lib/linux64/ctpgo_wrap.h ./lib/ctpgo_swigcxx.i

	cd lib/cmake && cmake . && make && mv libctpgo.so ../linux64/

	go build -o ./bin/ctpgo main.go

