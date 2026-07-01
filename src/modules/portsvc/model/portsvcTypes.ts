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

export type DependencyAssetItem = {
  id?: string;
  ownerSubject?: string;
  ownerName?: string;
  name: string;
  assetKind: string;
  assetType: string;
  provider: string;
  url: string;
  fullName: string;
  externalId: string;
  visibility: string;
  controllability: string;
  status: string;
  description: string;
  metadata: string;
  notes: string;
};

export type PortGroupAssetLinkItem = {
  id?: string;
  portGroupId?: string;
  portSlotId?: string;
  assetId: string;
  asset?: DependencyAssetItem;
  relationType: string;
  required: boolean;
  notes: string;
};

export type ServiceGroupPortGroupItem = {
  id?: string;
  serviceGroupId?: string;
  portGroupId: string;
  portGroup?: PortGroupItem;
  role: string;
  notes: string;
};

export type ServiceGroupItem = {
  id?: string;
  ownerSubject?: string;
  ownerName?: string;
  name: string;
  kind: string;
  status: string;
  description: string;
  notes: string;
  portGroups: ServiceGroupPortGroupItem[];
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
  assetLinks: PortGroupAssetLinkItem[];
};

export type PortGroupForm = Partial<Omit<PortGroupItem, "id" | "ownerName" | "host">> & {
  projectName?: string;
};

export type HostForm = Partial<HostItem> & {
  name?: string;
};

export type ServiceGroupForm = Partial<ServiceGroupItem> & {
  name?: string;
};

export type PortsvcSnapshot = {
  meta: Meta;
  hosts: HostItem[];
  portGroups: PortGroupItem[];
  dependencyAssets: DependencyAssetItem[];
  serviceGroups: ServiceGroupItem[];
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
  dependencyAssets: number;
  assetLinks: number;
  serviceGroups: number;
};
