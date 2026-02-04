import * as React from "react";
import { Outlet, createRootRoute } from "@tanstack/react-router";
import Sidebar from "@/components/sidebar/sidebar";

export const Route = createRootRoute({
  component: RootComponent,
});

function RootComponent() {
  return (
    <React.Fragment>
      <div className="w-screen h-screen flex">
        <div className="flex-1 max-w-56 bg-red-500">
          <Sidebar />
        </div>
        <div>
          <Outlet />
        </div>
      </div>
    </React.Fragment>
  );
}
