import { RouterProvider } from "react-router-dom";
import { AuthBoundary } from "../auth";
import type { SessionState } from "../session";
import { AppProviders } from "./AppProviders";
import { router } from "../router";

export function AppRoot({ initialSession }: { initialSession: SessionState }) {
  return (
    <AppProviders>
      <AuthBoundary initialSession={initialSession}>
        <RouterProvider router={router} />
      </AuthBoundary>
    </AppProviders>
  );
}
