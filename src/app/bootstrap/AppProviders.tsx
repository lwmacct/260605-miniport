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
        storageKey="miniport-theme"
        withAntdApp
      >
        {children}
      </WorkbenchProvider>
    </QueryClientProvider>
  );
}
