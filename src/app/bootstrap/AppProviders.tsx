import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { WorkbenchProvider } from "@lwmacct/260627-antd-workbench";
import enUS from "antd/es/locale/en_US";
import zhCN from "antd/es/locale/zh_CN";
import { useState, type PropsWithChildren } from "react";

export function AppProviders({ children }: PropsWithChildren) {
  const [queryClient] = useState(
    () =>
      new QueryClient({
        defaultOptions: {
          queries: {
            refetchOnWindowFocus: false,
            retry: 1,
          },
        },
      }),
  );

  return (
    <QueryClientProvider client={queryClient}>
      <WorkbenchProvider
        appearance={{ storageKey: "miniport-theme" }}
        locale={{
          defaultValue: "zh-CN",
          documentLang: (locale) => locale,
          options: [
            {
              antdLocale: zhCN,
              documentLang: "zh-CN",
              label: "简体中文",
              shortLabel: "中",
              value: "zh-CN",
            },
            {
              antdLocale: enUS,
              documentLang: "en-US",
              label: "English",
              shortLabel: "EN",
              value: "en-US",
            },
          ],
          storageKey: "miniport-locale",
        }}
        withAntdApp
      >
        {children}
      </WorkbenchProvider>
    </QueryClientProvider>
  );
}
