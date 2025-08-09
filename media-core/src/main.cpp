#include <iostream>
#include <filesystem>
#include <string>
#include <sys/socket.h>
#include <sys/un.h>
#include <unistd.h>

int main(int argc, char** argv) {
    std::string uds = "/var/run/mediacore.sock";
    if (const char* env = std::getenv("MEDIA_CORE_UDS")) {
        uds = env;
    }
    std::error_code ec;
    std::filesystem::create_directories(std::filesystem::path(uds).parent_path(), ec);

    // Create a dummy UDS listener placeholder (not gRPC yet)
    int fd = ::socket(AF_UNIX, SOCK_STREAM, 0);
    if (fd < 0) {
        std::perror("socket");
        return 1;
    }
    sockaddr_un addr{};
    addr.sun_family = AF_UNIX;
    std::snprintf(addr.sun_path, sizeof(addr.sun_path), "%s", uds.c_str());
    ::unlink(uds.c_str());
    if (::bind(fd, reinterpret_cast<sockaddr*>(&addr), sizeof(addr)) != 0) {
        std::perror("bind");
        ::close(fd);
        return 1;
    }
    if (::listen(fd, 4) != 0) {
        std::perror("listen");
        ::close(fd);
        return 1;
    }
    std::cout << "media-core placeholder listening on UDS: " << uds << std::endl;
    // Block forever
    for(;;) { ::pause(); }
    return 0;
}