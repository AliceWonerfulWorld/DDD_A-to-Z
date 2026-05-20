import { connectClient } from "../../lib/api/connect";
import { HomeService } from "@lang-war/proto-ts/langwar/home/v1/home_service_pb";

const client = connectClient(HomeService);

export async function fetchHomeViaConnect() {
  return client.getHome({});
}
