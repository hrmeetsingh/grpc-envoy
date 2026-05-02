"use client";

import { useState } from "react";
import { userClient, greeterClient } from "@/lib/clients";
import {
  CreateUserRequest,
  GetUserRequest,
} from "@/gen/user/v1/UserServiceClientPb";
import {
  SayHelloRequest,
  SayGoodbyeRequest,
} from "@/gen/greeter/v1/GreeterServiceClientPb";

export default function Home() {
  const [userName, setUserName] = useState("");
  const [userId, setUserId] = useState("");
  const [greeterName, setGreeterName] = useState("");
  const [results, setResults] = useState<string[]>([]);

  function log(msg: string) {
    setResults((prev) => [`[${new Date().toLocaleTimeString()}] ${msg}`, ...prev]);
  }

  async function handleCreateUser() {
    try {
      const req = new CreateUserRequest();
      req.setName(userName);
      const res = await userClient.createUser(req);
      log(`CreateUser → id=${res.getId()}, name=${res.getName()}`);
    } catch (e: any) {
      log(`CreateUser error: ${e.message}`);
    }
  }

  async function handleGetUser() {
    try {
      const req = new GetUserRequest();
      req.setId(userId);
      const res = await userClient.getUser(req);
      log(`GetUser → id=${res.getId()}, name=${res.getName()}`);
    } catch (e: any) {
      log(`GetUser error: ${e.message}`);
    }
  }

  async function handleSayHello() {
    try {
      const req = new SayHelloRequest();
      req.setName(greeterName);
      const res = await greeterClient.sayHello(req);
      log(`SayHello → ${res.getMessage()}`);
    } catch (e: any) {
      log(`SayHello error: ${e.message}`);
    }
  }

  async function handleSayGoodbye() {
    try {
      const req = new SayGoodbyeRequest();
      req.setName(greeterName);
      const res = await greeterClient.sayGoodbye(req);
      log(`SayGoodbye → ${res.getMessage()}`);
    } catch (e: any) {
      log(`SayGoodbye error: ${e.message}`);
    }
  }

  return (
    <div className="space-y-10">
      {/* User Service */}
      <section className="rounded-lg border border-gray-800 p-6 space-y-4">
        <h2 className="text-lg font-medium">UserService</h2>
        <div className="flex gap-3 items-end">
          <div className="flex-1">
            <label className="block text-sm text-gray-400 mb-1">Name</label>
            <input
              value={userName}
              onChange={(e) => setUserName(e.target.value)}
              className="w-full rounded bg-gray-900 border border-gray-700 px-3 py-2 text-sm focus:outline-none focus:ring-1 focus:ring-blue-500"
              placeholder="Alice"
            />
          </div>
          <button
            onClick={handleCreateUser}
            className="rounded bg-blue-600 hover:bg-blue-500 px-4 py-2 text-sm font-medium transition"
          >
            CreateUser
          </button>
        </div>
        <div className="flex gap-3 items-end">
          <div className="flex-1">
            <label className="block text-sm text-gray-400 mb-1">User ID</label>
            <input
              value={userId}
              onChange={(e) => setUserId(e.target.value)}
              className="w-full rounded bg-gray-900 border border-gray-700 px-3 py-2 text-sm focus:outline-none focus:ring-1 focus:ring-blue-500"
              placeholder="paste id here"
            />
          </div>
          <button
            onClick={handleGetUser}
            className="rounded bg-gray-700 hover:bg-gray-600 px-4 py-2 text-sm font-medium transition"
          >
            GetUser
          </button>
        </div>
      </section>

      {/* Greeter Service */}
      <section className="rounded-lg border border-gray-800 p-6 space-y-4">
        <h2 className="text-lg font-medium">GreeterService</h2>
        <div className="flex gap-3 items-end">
          <div className="flex-1">
            <label className="block text-sm text-gray-400 mb-1">Name</label>
            <input
              value={greeterName}
              onChange={(e) => setGreeterName(e.target.value)}
              className="w-full rounded bg-gray-900 border border-gray-700 px-3 py-2 text-sm focus:outline-none focus:ring-1 focus:ring-blue-500"
              placeholder="Bob"
            />
          </div>
          <button
            onClick={handleSayHello}
            className="rounded bg-green-600 hover:bg-green-500 px-4 py-2 text-sm font-medium transition"
          >
            SayHello
          </button>
          <button
            onClick={handleSayGoodbye}
            className="rounded bg-amber-600 hover:bg-amber-500 px-4 py-2 text-sm font-medium transition"
          >
            SayGoodbye
          </button>
        </div>
      </section>

      {/* Response Log */}
      <section className="rounded-lg border border-gray-800 p-6">
        <h2 className="text-lg font-medium mb-3">Response Log</h2>
        {results.length === 0 ? (
          <p className="text-sm text-gray-500">No responses yet</p>
        ) : (
          <ul className="space-y-1 font-mono text-sm max-h-64 overflow-y-auto">
            {results.map((r, i) => (
              <li key={i} className="text-gray-300">
                {r}
              </li>
            ))}
          </ul>
        )}
      </section>
    </div>
  );
}
