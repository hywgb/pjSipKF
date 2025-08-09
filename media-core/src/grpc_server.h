#pragma once

#include <memory>
#include <string>

#include <grpcpp/grpcpp.h>

std::unique_ptr<grpc::Server> StartServerOnUDS(const std::string& uds_path);