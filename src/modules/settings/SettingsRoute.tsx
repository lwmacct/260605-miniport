import {
  WorkbenchAppearanceSettings,
  WorkbenchPage,
  WorkbenchPanel,
} from "@lwmacct/260627-antd-workbench";

export function SettingsRoute() {
  return (
    <WorkbenchPage title="Settings">
      <WorkbenchPanel title="外观设置">
        <WorkbenchAppearanceSettings />
      </WorkbenchPanel>
    </WorkbenchPage>
  );
}
