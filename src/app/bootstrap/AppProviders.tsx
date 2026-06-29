import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { WorkbenchRoot } from "@lwmacct/260627-antd-workbench";
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
      <WorkbenchRoot
        appearance={{ storageKey: "miniport-theme" }}
        withAntdApp
      >
        {children}
      </WorkbenchRoot>
    </QueryClientProvider>
  );
}
