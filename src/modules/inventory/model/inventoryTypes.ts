export type Meta = {
  name: string;
  version: string;
  listen: string;
  database: string;
  docsPath: string;
};

export type Host = {
  id: number;
  ip: string;
  name: string;
  network: string;
  environment: string;
  notes: string;
};

export type PortSlot = {
  id?: number;
  port: number;
  name: string;
  protocol: string;
  purpose: string;
  status: string;
  notes: string;
};

export type ComponentItem = {
  id?: number;
  name: string;
  type: string;
  url: string;
  version: string;
  notes: string;
};

export type RepositoryItem = {
  id?: number;
  name: string;
  url: string;
  kind: string;
  notes: string;
};

export type PortGroup = {
  id: number;
  hostId: number;
  host?: Host;
  portStart: number;
  portEnd: number;
  serviceName: string;
  containerName: string;
  dindHost: string;
  status: string;
  owner: string;
  tags: string;
  notes: string;
  slots: PortSlot[];
  components: ComponentItem[];
  repositories: RepositoryItem[];
};

export type GroupForm = Partial<Omit<PortGroup, "id" | "host">> & {
  hostId?: number;
  portStart?: number;
  portEnd?: number;
  serviceName?: string;
};

export type HostForm = Partial<Omit<Host, "id">> & {
  ip?: string;
};

export type InventorySnapshot = {
  meta: Meta;
  hosts: Host[];
  groups: PortGroup[];
};

export type InventoryQuery = {
  environment?: string;
  hostSort?: string;
  hostQuery?: string;
  hostId?: number;
  portGroupQuery?: string;
  portGroupSort?: string;
  status?: string;
};

export type BatchPortGroupUpdate = {
  owner?: string;
  status?: string;
  tags?: string;
};

export type AppStats = {
  hosts: number;
  groups: number;
  usedSlots: number;
  emptySlots: number;
  components: number;
  repositories: number;
};
