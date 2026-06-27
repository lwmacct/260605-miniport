import path from "node:path";
import { defineConfig, type Alias, type UserConfig } from "vite";

interface LocalPackage {
  dedupe?: string[];
  exports: Record<string, string>;
  name: string;
  root: string;
}

const localPackages: LocalPackage[] = [
  {
    name: "@lwmacct/260627-antd-workbench",
    root: "/data/project/260627-antd-workbench/workspace",
    exports: {
      ".": "src/index.ts",
      "styles.css": "src/styles.css",
    },
    dedupe: ["@ant-design/icons", "antd", "react", "react-dom"],
  },
];

function createLocalWorkspaceConfig(packages: LocalPackage[]): UserConfig {
  const alias: Alias[] = [];
  const allow = new Set<string>([import.meta.dirname]);
  const dedupe = new Set<string>();
  const exclude = new Set<string>();

  for (const pkg of packages) {
    allow.add(pkg.root);
    exclude.add(pkg.name);

    for (const dep of pkg.dedupe ?? []) {
      dedupe.add(dep);
    }

    for (const [subpath, target] of Object.entries(pkg.exports)) {
      const specifier = subpath === "." ? pkg.name : `${pkg.name}/${subpath}`;

      alias.push({
        find: new RegExp(`^${escapeRegExp(specifier)}$`),
        replacement: path.resolve(pkg.root, target),
      });
    }
  }

  return {
    resolve: {
      alias,
      dedupe: [...dedupe],
    },
    server: {
      fs: {
        allow: [...allow],
      },
    },
    optimizeDeps: {
      exclude: [...exclude],
    },
  };
}

function escapeRegExp(value: string): string {
  return value.replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
}

export default defineConfig(createLocalWorkspaceConfig(localPackages));
