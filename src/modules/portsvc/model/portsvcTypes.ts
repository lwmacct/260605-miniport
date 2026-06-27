export type Meta = {
  name: string;
  version: string;
  listen: string;
  database: string;
  docsPath: string;
};

export type PortAllocation = {
  id: number;
  userId: number;
  username: string;
  portStart: number;
  portEnd: number;
  status: string;
  notes: string;
};

export type RepositoryItem = {
  id?: number;
  userId?: number;
  name: string;
  url: string;
  kind: string;
  role?: string;
  notes: string;
};

export type DependencyItem = {
  id?: number;
  userId?: number;
  name: string;
  type: string;
  url: string;
  version: string;
  role?: string;
  notes: string;
};

export type ServiceItem = {
  id: number;
  userId: number;
  username: string;
  portAllocationId: number;
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

export type ServiceForm = Partial<Omit<ServiceItem, "id" | "username" | "portAllocation">> & {
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
  userId?: number;
};

export type AppStats = {
  services: number;
  ports: number;
  boundPorts: number;
  freePorts: number;
  repositories: number;
  dependencies: number;
};
