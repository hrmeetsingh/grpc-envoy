// @generated hand-written grpc-web stubs matching proto/user/v1/user.proto
const jspb = require("google-protobuf");

class CreateUserRequest extends jspb.Message {
  constructor(opt_data) {
    super();
    jspb.Message.initialize(this, opt_data, 0, -1, null, null);
  }
  getName() { return jspb.Message.getFieldWithDefault(this, 1, ""); }
  setName(value) { jspb.Message.setField(this, 1, value); }
  serializeBinary() {
    const writer = new jspb.BinaryWriter();
    if (this.getName().length > 0) writer.writeString(1, this.getName());
    return writer.getResultBuffer();
  }
  static deserializeBinary(bytes) {
    const msg = new CreateUserRequest();
    const reader = new jspb.BinaryReader(bytes);
    while (reader.nextField()) {
      if (reader.isEndGroup()) break;
      switch (reader.getFieldNumber()) {
        case 1: msg.setName(reader.readString()); break;
        default: reader.skipField(); break;
      }
    }
    return msg;
  }
}
CreateUserRequest.displayName = "CreateUserRequest";

class CreateUserResponse extends jspb.Message {
  constructor(opt_data) {
    super();
    jspb.Message.initialize(this, opt_data, 0, -1, null, null);
  }
  getId() { return jspb.Message.getFieldWithDefault(this, 1, ""); }
  setId(value) { jspb.Message.setField(this, 1, value); }
  getName() { return jspb.Message.getFieldWithDefault(this, 2, ""); }
  setName(value) { jspb.Message.setField(this, 2, value); }
  serializeBinary() {
    const writer = new jspb.BinaryWriter();
    if (this.getId().length > 0) writer.writeString(1, this.getId());
    if (this.getName().length > 0) writer.writeString(2, this.getName());
    return writer.getResultBuffer();
  }
  static deserializeBinary(bytes) {
    const msg = new CreateUserResponse();
    const reader = new jspb.BinaryReader(bytes);
    while (reader.nextField()) {
      if (reader.isEndGroup()) break;
      switch (reader.getFieldNumber()) {
        case 1: msg.setId(reader.readString()); break;
        case 2: msg.setName(reader.readString()); break;
        default: reader.skipField(); break;
      }
    }
    return msg;
  }
}
CreateUserResponse.displayName = "CreateUserResponse";

class GetUserRequest extends jspb.Message {
  constructor(opt_data) {
    super();
    jspb.Message.initialize(this, opt_data, 0, -1, null, null);
  }
  getId() { return jspb.Message.getFieldWithDefault(this, 1, ""); }
  setId(value) { jspb.Message.setField(this, 1, value); }
  serializeBinary() {
    const writer = new jspb.BinaryWriter();
    if (this.getId().length > 0) writer.writeString(1, this.getId());
    return writer.getResultBuffer();
  }
  static deserializeBinary(bytes) {
    const msg = new GetUserRequest();
    const reader = new jspb.BinaryReader(bytes);
    while (reader.nextField()) {
      if (reader.isEndGroup()) break;
      switch (reader.getFieldNumber()) {
        case 1: msg.setId(reader.readString()); break;
        default: reader.skipField(); break;
      }
    }
    return msg;
  }
}
GetUserRequest.displayName = "GetUserRequest";

class GetUserResponse extends jspb.Message {
  constructor(opt_data) {
    super();
    jspb.Message.initialize(this, opt_data, 0, -1, null, null);
  }
  getId() { return jspb.Message.getFieldWithDefault(this, 1, ""); }
  setId(value) { jspb.Message.setField(this, 1, value); }
  getName() { return jspb.Message.getFieldWithDefault(this, 2, ""); }
  setName(value) { jspb.Message.setField(this, 2, value); }
  serializeBinary() {
    const writer = new jspb.BinaryWriter();
    if (this.getId().length > 0) writer.writeString(1, this.getId());
    if (this.getName().length > 0) writer.writeString(2, this.getName());
    return writer.getResultBuffer();
  }
  static deserializeBinary(bytes) {
    const msg = new GetUserResponse();
    const reader = new jspb.BinaryReader(bytes);
    while (reader.nextField()) {
      if (reader.isEndGroup()) break;
      switch (reader.getFieldNumber()) {
        case 1: msg.setId(reader.readString()); break;
        case 2: msg.setName(reader.readString()); break;
        default: reader.skipField(); break;
      }
    }
    return msg;
  }
}
GetUserResponse.displayName = "GetUserResponse";

module.exports = {
  CreateUserRequest,
  CreateUserResponse,
  GetUserRequest,
  GetUserResponse,
};
