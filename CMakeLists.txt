cmake_minimum_required(VERSION 3.7)
project(trinket)

set(CMAKE_CXX_STANDARD 14)
set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -std=c++11")

link_libraries(crypto)
set(SOURCE_FILES main.cpp)
add_executable(trinket ${SOURCE_FILES})