import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { WorkbenchProvider } from "@lwmacct/260627-antd-workbench";
import { useState, type PropsWithChildren } from "react";
import type { ThemeMode } from "@/shared/theme/theme";

interface AppProvidersProps {
  themeMode: ThemeMode;
}

export function AppProviders({
  children,
  themeMode,
}: PropsWithChildren<AppProvidersProps>) {
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
        storageKey={false}
        themeMode={themeMode}
        withAntdApp
      >
        {children}
      </WorkbenchProvider>
    </QueryClientProvider>
  );
}
