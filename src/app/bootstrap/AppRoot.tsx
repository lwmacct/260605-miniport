import { RouterProvider } from "react-router-dom";
import { AppProviders } from "./AppProviders";
import { router } from "../router";
import { ThemeModeProvider } from "../../shared/theme/ThemeModeProvider";
import { useThemeModeContext } from "../../shared/theme/ThemeModeContext";

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
