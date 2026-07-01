export type Meta = {
  name: string;
  version: string;
  listen: string;
  database: string;
  docsPath: string;
};

export type HostItem = {
  id: string;
  name: string;
  ip: string;
  spec: string;
  status: string;
  notes: string;
};

export type PortSlotItem = {
  id?: string;
  portGroupId?: string;
  port: number;
  name: string;
  kind: string;
  protocol: string;
  containerName: string;
  status: string;
  notes: string;
};

export type RepositoryItem = {
  id?: string;
  ownerSubject?: string;
  portGroupId?: string;
  name: string;
  url: string;
  kind: string;
  notes: string;
};

export type DependencyItem = {
  id?: string;
  ownerSubject?: string;
  portGroupId?: string;
  name: string;
  type: string;
  url: string;
  version: string;
  notes: string;
};

export type PortGroupItem = {
  id: string;
  ownerSubject: string;
  ownerName: string;
  hostId: string;
  host?: HostItem;
  portStart: number;
  portEnd: number;
  projectName: string;
  projectOwner: string;
  runtimeMode: string;
  runtimeName: string;
  serviceIp: string;
  status: string;
  tags: string;
  notes: string;
  slots: PortSlotItem[];
  repositories: RepositoryItem[];
  dependencies: DependencyItem[];
};

export type PortGroupForm = Partial<Omit<PortGroupItem, "id" | "ownerName" | "host">> & {
  projectName?: string;
};

export type HostForm = Partial<HostItem> & {
  name?: string;
};

export type PortsvcSnapshot = {
  meta: Meta;
  hosts: HostItem[];
  portGroups: PortGroupItem[];
};

export type PortsvcQuery = {
  query?: string;
  sort?: string;
  status?: string;
  ownerSubject?: string;
};

export type AppStats = {
  groups: number;
  hosts: number;
  runningGroups: number;
  freeGroups: number;
  slots: number;
  repositories: number;
  dependencies: number;
};
