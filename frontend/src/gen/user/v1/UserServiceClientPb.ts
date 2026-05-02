import * as grpcWeb from "grpc-web";
const {
  CreateUserRequest,
  CreateUserResponse,
  GetUserRequest,
  GetUserResponse,
} = require("./user_pb");

export { CreateUserRequest, CreateUserResponse, GetUserRequest, GetUserResponse };

export class UserServiceClient {
  private client: grpcWeb.GrpcWebClientBase;
  private hostname: string;

  constructor(hostname: string) {
    this.hostname = hostname;
    this.client = new grpcWeb.GrpcWebClientBase({ format: "binary" });
  }

  createUser(
    request: InstanceType<typeof CreateUserRequest>
  ): Promise<InstanceType<typeof CreateUserResponse>> {
    const methodDescriptor = new grpcWeb.MethodDescriptor(
      "/user.v1.UserService/CreateUser",
      grpcWeb.MethodType.UNARY,
      CreateUserRequest,
      CreateUserResponse,
      (req: any) => req.serializeBinary(),
      CreateUserResponse.deserializeBinary
    );
    return new Promise((resolve, reject) => {
      this.client.rpcCall(
        this.hostname + "/user.v1.UserService/CreateUser",
        request,
        {},
        methodDescriptor,
        (err: grpcWeb.RpcError, response: any) => {
          if (err) reject(err);
          else resolve(response);
        }
      );
    });
  }

  getUser(
    request: InstanceType<typeof GetUserRequest>
  ): Promise<InstanceType<typeof GetUserResponse>> {
    const methodDescriptor = new grpcWeb.MethodDescriptor(
      "/user.v1.UserService/GetUser",
      grpcWeb.MethodType.UNARY,
      GetUserRequest,
      GetUserResponse,
      (req: any) => req.serializeBinary(),
      GetUserResponse.deserializeBinary
    );
    return new Promise((resolve, reject) => {
      this.client.rpcCall(
        this.hostname + "/user.v1.UserService/GetUser",
        request,
        {},
        methodDescriptor,
        (err: grpcWeb.RpcError, response: any) => {
          if (err) reject(err);
          else resolve(response);
        }
      );
    });
  }
}
