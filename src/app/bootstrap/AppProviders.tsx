import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { WorkbenchProvider } from "@lwmacct/260627-antd-workbench";
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
        defaultLocale="zh-CN"
        localeStorageKey="miniport-locale"
        withAntdApp
      >
        {children}
      </WorkbenchProvider>
    </QueryClientProvider>
  );
}
