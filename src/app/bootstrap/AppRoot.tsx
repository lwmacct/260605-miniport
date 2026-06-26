import { RouterProvider } from "react-router-dom";
import { ThemeModeProvider } from "@/shared/theme/ThemeModeProvider";
import { useThemeModeContext } from "@/shared/theme/ThemeModeContext";
import { AppProviders } from "./AppProviders";
import { router } from "../router";

function RoutedApplication() {
  const { themeMode } = useThemeModeContext();

  return (
    <AppProviders themeMode={themeMode}>
      <RouterProvider router={router} />
    </AppProviders>
  );
}

export function AppRoot() {
  return (
    <ThemeModeProvider>
      <RoutedApplication />
    </ThemeModeProvider>
  );
}
