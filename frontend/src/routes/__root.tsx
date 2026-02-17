import * as React from "react";
import { Outlet, createRootRoute } from "@tanstack/react-router";
import Sidebar from "@/components/sidebar/sidebar";
import Header from "@/components/header/header";

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
        <div className="relative flex-1 bg-[#222527]">
          <div className="sticky top-0">
            <Header />
          </div>
          <main className="">
            <Outlet />
          </main>
        </div>
      </div>
    </React.Fragment>
  );
}
