cmake_minimum_required (VERSION 3.20)
project (ctpgo C CXX)

set(MY_LIB_PATH ${CMAKE_CURRENT_SOURCE_DIR}/../linux64)
set(DIR_LIB_SRCS ${CMAKE_CURRENT_SOURCE_DIR}/../linux64/ctpgo_wrap.cxx)
set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -std=c++11")
include_directories(${CMAKE_CURRENT_SOURCE_DIR}/../linux64)

add_library(ctpgo SHARED ${DIR_LIB_SRCS})

TARGET_LINK_LIBRARIES(ctpgo ${MY_LIB_PATH}/libthosttraderapi_se.so ${MY_LIB_PATH}/libthostmduserapi_se.so)
