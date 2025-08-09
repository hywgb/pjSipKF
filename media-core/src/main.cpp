#include <iostream>
#include <string>
#include <memory>
#include <cstdlib>

#include "grpc_server.h"

extern std::unique_ptr<grpc::Server> StartServerOnTCP(const std::string& hostport);

int main(int argc, char** argv) {
    const char* tcp = std::getenv("MEDIA_CORE_TCP");
    if (tcp && std::string(tcp) == std::string("1")) {
        std::string hp = "127.0.0.1:50051";
        auto server = StartServerOnTCP(hp);
        if (!server) { std::cerr << "failed to start gRPC TCP server on " << hp << std::endl; return 1; }
        std::cout << "media-core gRPC listening on TCP: " << hp << std::endl;
        server->Wait();
        return 0;
    }

    std::string uds = "/tmp/mediacore.sock";
    if (const char* env = std::getenv("MEDIA_CORE_UDS")) { uds = env; }
    auto server = StartServerOnUDS(uds);
    if (!server) {
        std::cerr << "failed to start gRPC server on UDS: " << uds << std::endl;
        return 1;
    }
    std::cout << "media-core gRPC listening on UDS: " << uds << std::endl;
    server->Wait();
    return 0;
}