export type Meta = {
  name: string;
  version: string;
  listen: string;
  database: string;
  docsPath: string;
};

export type PortAllocation = {
  id: string;
  ownerSubject: string;
  ownerName: string;
  portStart: number;
  portEnd: number;
  status: string;
  notes: string;
};

export type RepositoryItem = {
  id?: string;
  ownerSubject?: string;
  name: string;
  url: string;
  kind: string;
  role?: string;
  notes: string;
};

export type DependencyItem = {
  id?: string;
  ownerSubject?: string;
  name: string;
  type: string;
  url: string;
  version: string;
  role?: string;
  notes: string;
};

export type ServiceItem = {
  id: string;
  ownerSubject: string;
  ownerName: string;
  portAllocationId: string;
  portAllocation?: PortAllocation;
  name: string;
  projectName: string;
  dindIp: string;
  dindContainer: string;
  status: string;
  owner: string;
  tags: string;
  notes: string;
  repositories: RepositoryItem[];
  dependencies: DependencyItem[];
};

export type ServiceForm = Partial<Omit<ServiceItem, "id" | "ownerName" | "portAllocation">> & {
  name?: string;
};

export type PortAllocationForm = Partial<PortAllocation> & {
  portStart?: number;
};

export type PortsvcSnapshot = {
  meta: Meta;
  ports: PortAllocation[];
  services: ServiceItem[];
};

export type PortsvcQuery = {
  projectName?: string;
  serviceQuery?: string;
  serviceSort?: string;
  status?: string;
  ownerSubject?: string;
};

export type AppStats = {
  services: number;
  ports: number;
  boundPorts: number;
  freePorts: number;
  repositories: number;
  dependencies: number;
};
