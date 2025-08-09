#include <iostream>
#include <string>
#include <memory>

#include "grpc_server.h"

int main(int argc, char** argv) {
    std::string uds = "/var/run/mediacore.sock";
    if (const char* env = std::getenv("MEDIA_CORE_UDS")) {
        uds = env;
    }
    auto server = StartServerOnUDS(uds);
    if (!server) {
        std::cerr << "failed to start gRPC server on UDS: " << uds << std::endl;
        return 1;
    }
    std::cout << "media-core gRPC listening on UDS: " << uds << std::endl;
    server->Wait();
    return 0;
}