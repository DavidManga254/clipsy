import { ReactNode } from "react";

type NavSidebarItemProps = {
  name: string;
  Icon: ReactNode;
  active?: boolean;
};

const NavSidebarItem = ({ name, Icon, active }: NavSidebarItemProps) => {
  const bgColor = active ? "bg-slate-800" : "bg-transparent";

  return (
    <div className={`h-12 flex items-center gap-2 pl-5 ${bgColor}`}>
      {Icon}
      <p className="text-[14px] text-[#E8E6E3] font-medium">{name}</p>
    </div>
  );
};

export default NavSidebarItem;
