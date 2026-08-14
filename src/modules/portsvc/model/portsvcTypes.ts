import type { components, operations } from "@/shared/api/schema.gen";

type Schema = components["schemas"];

export type Meta = Schema["MetaDTO"];
export type HostItem = Schema["HostDTO"];
export type PortSlotItem = Schema["PortSlotDTO"];
export type DependencyAssetItem = Schema["DependencyAssetDTO"];
export type PortGroupAssetLinkItem = Schema["PortGroupAssetLinkDTO"];
export type GitHubRepositoryItem = Schema["GithubRepositoryDTO"];
export type PortGroupRepositoryLinkItem = Schema["PortGroupRepositoryLinkDTO"];
export type ServiceGroupPortGroupItem = Schema["ServiceGroupPortGroupDTO"];
export type ServiceGroupItem = Schema["ServiceGroupDTO"];
export type PortGroupItem = Schema["PortGroupDTO"];

export type PortGroupForm = Schema["PortGroupCreateDTO"];
export type HostForm = Schema["HostCreateDTO"];
export type ServiceGroupForm = Schema["ServiceGroupCreateDTO"];
export type PortsvcQuery = NonNullable<operations["console-list-port-groups"]["parameters"]["query"]>;

export type PortsvcSnapshot = {
  meta: Meta;
  hosts: HostItem[];
  portGroups: PortGroupItem[];
  dependencyAssets: DependencyAssetItem[];
  repositories: GitHubRepositoryItem[];
  serviceGroups: ServiceGroupItem[];
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
  repositories: number;
  repositoryLinks: number;
};
