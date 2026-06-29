import { Alert, Result, Spin } from "antd";
import { CenterState as WorkbenchCenterState } from "@lwmacct/260627-antd-workbench";
import { Navigate, Outlet } from "react-router-dom";
import { useAuthStateQuery } from "@/modules/auth";
import { appPaths } from "./navigation";

function CenterState({ children }: { children: React.ReactNode }) {
  return <WorkbenchCenterState>{children}</WorkbenchCenterState>;
}

export function GuestOnlyBoundary() {
  const authState = useAuthStateQuery();

  if (authState.isPending) {
    return (
      <CenterState>
        <Spin />
      </CenterState>
    );
  }

  if (authState.isError) {
    return (
      <CenterState>
        <Alert showIcon message="应用初始化失败" description={authState.error.message} type="error" />
      </CenterState>
    );
  }

  if (authState.data.session.authenticated) {
    return <Navigate to={appPaths.overview} replace />;
  }

  return <Outlet />;
}

export function ProtectedBoundary() {
  const authState = useAuthStateQuery();

  if (authState.isPending) {
    return (
      <CenterState>
        <Spin />
      </CenterState>
    );
  }

  if (authState.isError) {
    return (
      <CenterState>
        <Alert showIcon message="应用初始化失败" description={authState.error.message} type="error" />
      </CenterState>
    );
  }

  if (!authState.data.session.authenticated) {
    return <Navigate to={appPaths.login} replace />;
  }

  return <Outlet />;
}

export function AdminBoundary() {
  const authState = useAuthStateQuery();

  if (authState.isPending) {
    return (
      <CenterState>
        <Spin />
      </CenterState>
    );
  }

  if (authState.isError) {
    return (
      <CenterState>
        <Alert showIcon message="应用初始化失败" description={authState.error.message} type="error" />
      </CenterState>
    );
  }

  if (!authState.data.session.user?.admin) {
    return (
      <CenterState>
        <Result status="403" title="403" subTitle="当前账号没有管理员权限。" />
      </CenterState>
    );
  }

  return <Outlet />;
}
