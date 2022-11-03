export type ConfigHistoryType = ConfigType[];

export type ConfigType = {
  unix: number;
  gateways: GatewayType[];
};

export type GatewayType = {
  services: ServiceType[];
  address: string;
};

export type ServiceType = {
  hostName: string;
  balancer: BalancerType;
};

export type BalancerType = {
  servers: ServerType[];
  algorithm: string;
};

export type ServerType = {
  address: string;
};

export type ConfigPostResponseBody = {
  ok: boolean;
  exceptions: string[];
};
