import { UserServiceClient } from "@/gen/user/v1/UserServiceClientPb";
import { GreeterServiceClient } from "@/gen/greeter/v1/GreeterServiceClientPb";

const ENVOY_URL = process.env.NEXT_PUBLIC_ENVOY_URL || "http://localhost:8080";

export const userClient = new UserServiceClient(ENVOY_URL);
export const greeterClient = new GreeterServiceClient(ENVOY_URL);
