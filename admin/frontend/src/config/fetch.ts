import { ApplyAuth } from "../auth/auth";
import type { ConfigHistoryType, ConfigPostResponseBody, ConfigType } from "./types";

export async function ConfigHistory(): Promise<ConfigHistoryType> {
  // return [
  //   {
  //     redirect: true,
  //     unix: 1,
  //     gateways: [
  //       {
  //         address: "0.0.0.0:443",
  //         services: [
  //           {
  //             hostName: "google.org",
  //             balancer: {
  //               algorithm: "RoundRobin",
  //               servers: [
  //                 {
  //                   address: "localhost:8080",
  //                 },
  //               ],
  //             },
  //           },
  //           {
  //             hostName: "google.org",
  //             balancer: {
  //               algorithm: "RoundRobin",
  //               servers: [],
  //             },
  //           },
  //         ],
  //       },
  //     ],
  //   },
  // ];
  const res = await fetch("config", { method: "GET", headers: ApplyAuth({}) });
  const json = await res.json();
  return json as ConfigHistoryType;
}

export async function ApplyConfig(config: ConfigType): Promise<ConfigPostResponseBody> {
  const res = await fetch("config", {
    method: "POST",
    headers: ApplyAuth({}),
    body: JSON.stringify(config),
  });
  const json = await res.json();

  return json as ConfigPostResponseBody;
}
