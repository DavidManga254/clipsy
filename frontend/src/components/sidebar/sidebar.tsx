import {
  ArrowLeftRight,
  Clipboard,
  History,
  LaptopMinimal,
  Settings,
} from "lucide-react";
import NavSidebarItem from "./NavSidebarItem/NavSidebarItem";
import NavDeviceItem from "./NavDeviceItem/NavDeviceItem";

const Sidebar = () => {
  return (
    <div className="bg-slate-900 h-full relative">
      <div className="h-14 flex items-center pl-5 border-b border-b-white">
        <Clipboard fill="white" />
        <p className="text-[#E1DFDB] text-[16px] font-bold ml-3">ClipSync</p>
      </div>
      <NavSidebarItem
        name="Devices"
        Icon={<LaptopMinimal color="white" />}
        active
      />
      <NavSidebarItem name="History" Icon={<History color="white" />} />
      <NavSidebarItem name="Activty" Icon={<ArrowLeftRight color="white" />} />

      <div className="mt-5">
        <h2 className="ml-5 text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2">
          Connected
        </h2>
        <NavDeviceItem
          deviceName="work-pc"
          deviceType="windows"
          ipAddress="192.08.05.5"
          isOnline
        />
      </div>
      <div className="w-full border-gray-800 p-3 absolute bottom-0 border-t-[2px]">
        <button className="w-full px-3 py-2 text-left flex items-center space-x-3 text-gray-400 hover:bg-gray-800 hover:text-white rounded">
          <Settings color="#9CA3AF" />
          <span className="text-sm">Settings</span>
        </button>
      </div>
    </div>
  );
};

export default Sidebar;
