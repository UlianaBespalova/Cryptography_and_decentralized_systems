cmake_minimum_required(VERSION 3.10)
project(Shamirs_secret_sharing)

set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -std=c++11")

link_libraries(crypto)
add_executable(Shamir main.cpp)