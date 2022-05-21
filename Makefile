ctpgo:
	go build -v -x -linkshared -o ./bin/ctpgo ./src/*.go

all:
	if [ -d "/home/godev" ]; then \
  			cd /home/godev; \
            git clone https://github.com/pengzuyun/ctpgo; \
        fi \
    cd ctpgo
	go get -u github.com/axgle/mahonia
	swig -go -cgo -intgosize 64 -module ctpgo -I./lib/linux64 -c++ -outdir ./ -o ./lib/linux64/ctpgo_wrap.cxx -oh ./lib/linux64/ctpgo_wrap.h ./lib/ctpgo.swigcxx
	cd lib/cmake
	cmake .
	make
	mv libctpgo.so ../linux64/
	cd ../..
	go build -v -x -linkshared -o ./bin/ctpgo ./src/*.go
