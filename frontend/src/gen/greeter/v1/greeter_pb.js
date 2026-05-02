// @generated hand-written grpc-web stubs matching proto/greeter/v1/greeter.proto
const jspb = require("google-protobuf");

class SayHelloRequest extends jspb.Message {
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
    const msg = new SayHelloRequest();
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
SayHelloRequest.displayName = "SayHelloRequest";

class SayHelloResponse extends jspb.Message {
  constructor(opt_data) {
    super();
    jspb.Message.initialize(this, opt_data, 0, -1, null, null);
  }
  getMessage() { return jspb.Message.getFieldWithDefault(this, 1, ""); }
  setMessage(value) { jspb.Message.setField(this, 1, value); }
  serializeBinary() {
    const writer = new jspb.BinaryWriter();
    if (this.getMessage().length > 0) writer.writeString(1, this.getMessage());
    return writer.getResultBuffer();
  }
  static deserializeBinary(bytes) {
    const msg = new SayHelloResponse();
    const reader = new jspb.BinaryReader(bytes);
    while (reader.nextField()) {
      if (reader.isEndGroup()) break;
      switch (reader.getFieldNumber()) {
        case 1: msg.setMessage(reader.readString()); break;
        default: reader.skipField(); break;
      }
    }
    return msg;
  }
}
SayHelloResponse.displayName = "SayHelloResponse";

class SayGoodbyeRequest extends jspb.Message {
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
    const msg = new SayGoodbyeRequest();
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
SayGoodbyeRequest.displayName = "SayGoodbyeRequest";

class SayGoodbyeResponse extends jspb.Message {
  constructor(opt_data) {
    super();
    jspb.Message.initialize(this, opt_data, 0, -1, null, null);
  }
  getMessage() { return jspb.Message.getFieldWithDefault(this, 1, ""); }
  setMessage(value) { jspb.Message.setField(this, 1, value); }
  serializeBinary() {
    const writer = new jspb.BinaryWriter();
    if (this.getMessage().length > 0) writer.writeString(1, this.getMessage());
    return writer.getResultBuffer();
  }
  static deserializeBinary(bytes) {
    const msg = new SayGoodbyeResponse();
    const reader = new jspb.BinaryReader(bytes);
    while (reader.nextField()) {
      if (reader.isEndGroup()) break;
      switch (reader.getFieldNumber()) {
        case 1: msg.setMessage(reader.readString()); break;
        default: reader.skipField(); break;
      }
    }
    return msg;
  }
}
SayGoodbyeResponse.displayName = "SayGoodbyeResponse";

module.exports = {
  SayHelloRequest,
  SayHelloResponse,
  SayGoodbyeRequest,
  SayGoodbyeResponse,
};
