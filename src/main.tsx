import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { AppRoot } from "./app/bootstrap/AppRoot";
import "antd/dist/reset.css";
import "@lwmacct/260627-antd-workbench/styles.css";
import "./styles/global.css";

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <AppRoot />
  </StrictMode>,
);
