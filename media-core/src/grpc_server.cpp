#include <grpcpp/grpcpp.h>
#include <sys/un.h>
#include <sys/socket.h>
#include <unistd.h>
#include <filesystem>
#include <memory>
#include <string>

#include "../../proto/gen/cpp/mediacore/session.grpc.pb.h"
#include "../../proto/gen/cpp/mediacore/session.pb.h"

using grpc::Server;
using grpc::ServerBuilder;
using grpc::ServerContext;
using grpc::Status;

class MediaCoreService : public mediacore::v1::MediaCore::Service {
public:
    Status CreateSession(ServerContext* context,
                         const mediacore::v1::CreateSessionRequest* request,
                         mediacore::v1::CreateSessionResponse* response) override {
        (void)context;
        response->set_session_id("sess-uds-0001");
        response->set_sdp_answer(std::string("v=0\n; gRPC answer for: ") + request->sdp_offer());
        return Status::OK;
    }
    Status UpdateSession(ServerContext* context,
                         const mediacore::v1::UpdateSessionRequest* request,
                         mediacore::v1::UpdateSessionResponse* response) override {
        (void)context; (void)request; (void)response;
        return Status::OK;
    }
    Status TerminateSession(ServerContext* context,
                            const mediacore::v1::TerminateSessionRequest* request,
                            mediacore::v1::TerminateSessionResponse* response) override {
        (void)context; (void)request; response->set_ok(true); return Status::OK;
    }
};

std::unique_ptr<Server> StartServerOnUDS(const std::string& uds_path) {
    std::error_code ec;
    std::filesystem::create_directories(std::filesystem::path(uds_path).parent_path(), ec);
    ::unlink(uds_path.c_str());

    MediaCoreService service;
    ServerBuilder builder;
    std::string addr = std::string("unix:") + uds_path;
    builder.AddListeningPort(addr, grpc::InsecureServerCredentials());
    builder.RegisterService(&service);
    std::unique_ptr<Server> server(builder.BuildAndStart());
    return server;
}