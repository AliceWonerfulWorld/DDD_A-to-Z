import { createConnectTransport } from "@connectrpc/connect-web";
import { createClient } from "@connectrpc/connect";
import type { DescService } from "@bufbuild/protobuf";
import type { Client } from "@connectrpc/connect";

const transport = createConnectTransport({ baseUrl: "" });

export function connectClient<S extends DescService>(service: S): Client<S> {
  return createClient(service, transport);
}
