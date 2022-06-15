package lib

/*
#cgo linux LDFLAGS: -fPIC -L${SRCDIR} -L${SRCDIR}/linux64 -Wl,-rpath,${SRCDIR}/linux64 -lctpgo -lthostmduserapi_se -lthosttraderapi_se -lstdc++
#cgo linux CPPFLAGS: -fPIC -I${SRCDIR}/linux64

// windows 不可用，go 部分功能不支持
*/
import "C"
