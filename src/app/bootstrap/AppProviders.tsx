import { App as AntdApp, ConfigProvider } from "antd";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { useState, type PropsWithChildren } from "react";
import { appTheme, type ThemeMode } from "../../shared/theme/theme";

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
      <ConfigProvider theme={appTheme(themeMode)}>
        <AntdApp component="div">{children}</AntdApp>
      </ConfigProvider>
    </QueryClientProvider>
  );
}
