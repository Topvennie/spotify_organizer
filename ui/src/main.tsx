import React from "react";
import ReactDOM from "react-dom/client";
import { QueryClientProvider } from "@tanstack/react-query";
import { queryClient } from "./lib/api/query";
import { theme } from "./lib/theme/theme";
import { MantineProvider } from "@mantine/core";

import "./index.css";
import "@mantine/core/styles.layer.css";
import "@mantine/notifications/styles.layer.css";
import "@mantine/dates/styles.css";
import 'mantine-datatable/styles.layer.css';
import { router } from "./router";
import { RouterProvider } from "@tanstack/react-router";
import { Notifications } from "@mantine/notifications";
import { AuthProvider } from "./lib/providers/AuthProvider";

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <QueryClientProvider client={queryClient}>
      <MantineProvider theme={theme}>
        <AuthProvider>
          <Notifications />
          <RouterProvider router={router} />
        </AuthProvider>
      </MantineProvider>
    </QueryClientProvider>
  </React.StrictMode>
);
