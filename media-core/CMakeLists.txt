cmake_minimum_required(VERSION 3.20)
project(media_core LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 20)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

find_package(Protobuf REQUIRED)
find_package(gRPC REQUIRED)

add_executable(media_core
    src/main.cpp
    src/grpc_server.cpp
    ../proto/gen/cpp/mediacore/session.pb.cc
    ../proto/gen/cpp/mediacore/session.grpc.pb.cc
)

target_include_directories(media_core PRIVATE ../proto/gen/cpp)

target_link_libraries(media_core PRIVATE gRPC::grpc++ protobuf::libprotobuf)