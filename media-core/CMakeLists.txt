cmake_minimum_required(VERSION 3.20)
project(media_core LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 20)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

add_executable(media_core src/main.cpp)

# TODO: add gRPC/protobuf, pjsip, pjnath, libsrtp dependencies
# target_link_libraries(media_core PRIVATE ...)