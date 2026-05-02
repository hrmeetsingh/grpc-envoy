import * as grpcWeb from "grpc-web";
const {
  SayHelloRequest,
  SayHelloResponse,
  SayGoodbyeRequest,
  SayGoodbyeResponse,
} = require("./greeter_pb");

export {
  SayHelloRequest,
  SayHelloResponse,
  SayGoodbyeRequest,
  SayGoodbyeResponse,
};

export class GreeterServiceClient {
  private client: grpcWeb.GrpcWebClientBase;
  private hostname: string;

  constructor(hostname: string) {
    this.hostname = hostname;
    this.client = new grpcWeb.GrpcWebClientBase({ format: "binary" });
  }

  sayHello(
    request: InstanceType<typeof SayHelloRequest>
  ): Promise<InstanceType<typeof SayHelloResponse>> {
    const methodDescriptor = new grpcWeb.MethodDescriptor(
      "/greeter.v1.GreeterService/SayHello",
      grpcWeb.MethodType.UNARY,
      SayHelloRequest,
      SayHelloResponse,
      (req: any) => req.serializeBinary(),
      SayHelloResponse.deserializeBinary
    );
    return new Promise((resolve, reject) => {
      this.client.rpcCall(
        this.hostname + "/greeter.v1.GreeterService/SayHello",
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

  sayGoodbye(
    request: InstanceType<typeof SayGoodbyeRequest>
  ): Promise<InstanceType<typeof SayGoodbyeResponse>> {
    const methodDescriptor = new grpcWeb.MethodDescriptor(
      "/greeter.v1.GreeterService/SayGoodbye",
      grpcWeb.MethodType.UNARY,
      SayGoodbyeRequest,
      SayGoodbyeResponse,
      (req: any) => req.serializeBinary(),
      SayGoodbyeResponse.deserializeBinary
    );
    return new Promise((resolve, reject) => {
      this.client.rpcCall(
        this.hostname + "/greeter.v1.GreeterService/SayGoodbye",
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
