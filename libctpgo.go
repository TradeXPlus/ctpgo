package ctpgo

/*
#cgo linux LDFLAGS: -fPIC -L${SRCDIR}/lib -L${SRCDIR}/lib/linux64 -Wl,-rpath,${SRCDIR}/lib/linux64 -lctpgo -lthostmduserapi_se -lthosttraderapi_se -lstdc++
#cgo linux CPPFLAGS: -fPIC -I${SRCDIR}/lib/linux64

// windows 不可用，go 部分功能不支持

*/
import "C"
