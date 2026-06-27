export type Meta = {
  name: string;
  version: string;
  listen: string;
  database: string;
  docsPath: string;
};

export type PortSlot = {
  id?: number;
  allocationId?: number;
  port: number;
  name: string;
  protocol: string;
  purpose: string;
  status: string;
  notes: string;
};

export type ProjectItem = {
  id?: number;
  allocationId?: number;
  name: string;
  description: string;
  notes: string;
};

export type ComponentItem = {
  id?: number;
  allocationId?: number;
  name: string;
  type: string;
  url: string;
  version: string;
  notes: string;
};

export type RepositoryItem = {
  id?: number;
  allocationId?: number;
  projectId?: number;
  name: string;
  url: string;
  kind: string;
  notes: string;
};

export type PortGroup = {
  id: number;
  userId: number;
  username: string;
  portStart: number;
  portEnd: number;
  name: string;
  dindIp: string;
  dindContainer: string;
  status: string;
  owner: string;
  tags: string;
  notes: string;
  slots: PortSlot[];
  projects: ProjectItem[];
  components: ComponentItem[];
  repositories: RepositoryItem[];
};

export type GroupForm = Partial<Omit<PortGroup, "id" | "username">> & {
  name?: string;
  portStart?: number;
};

export type InventorySnapshot = {
  meta: Meta;
  groups: PortGroup[];
};

export type InventoryQuery = {
  dindIp?: string;
  portGroupQuery?: string;
  portGroupSort?: string;
  projectName?: string;
  status?: string;
  userId?: number;
};

export type BatchPortGroupUpdate = {
  owner?: string;
  status?: string;
  tags?: string;
};

export type AppStats = {
  allocations: number;
  users: number;
  usedSlots: number;
  emptySlots: number;
  projects: number;
  components: number;
  repositories: number;
};
