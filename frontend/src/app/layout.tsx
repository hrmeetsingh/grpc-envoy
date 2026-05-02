import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "gRPC-Envoy Demo",
  description: "Browser → Envoy → gRPC services",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body className="bg-gray-950 text-gray-100 min-h-screen antialiased">
        <header className="border-b border-gray-800 px-6 py-4">
          <h1 className="text-xl font-semibold tracking-tight">
            gRPC-Envoy Demo
          </h1>
          <p className="text-sm text-gray-400">
            Browser → Envoy grpc-web proxy → Go gRPC services
          </p>
        </header>
        <main className="max-w-4xl mx-auto px-6 py-8">{children}</main>
      </body>
    </html>
  );
}
