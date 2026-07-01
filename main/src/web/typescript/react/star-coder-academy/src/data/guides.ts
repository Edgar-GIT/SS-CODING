import cppLogo from "@resources/img/languages/c++.svg";
import goLogo from "@resources/img/languages/golang.svg";
import rustLogo from "@resources/img/languages/rust.svg";
import typescriptLogo from "@resources/img/languages/typescript.svg";

export type Guide = {
  id: string;
  name: string;
  chapters: number;
  status: string;
  logo: string;
  aliases: string[];
  iconClassName: string;
  logoClassName?: string;
};

export const guides: Guide[] = [
  {
    id: "cpp",
    name: "C++",
    chapters: 50,
    status: "Coming Soon",
    logo: cppLogo,
    aliases: ["cpp", "c plus plus"],
    iconClassName: "bg-gradient-to-b from-gray-700 via-gray-900 to-black",
  },
  {
    id: "go",
    name: "Go",
    chapters: 50,
    status: "Coming Soon",
    logo: goLogo,
    aliases: ["golang"],
    iconClassName: "bg-slate-900 to-sky-600",
    logoClassName: "h-11 w-11",
  },
  {
    id: "typescript",
    name: "TypeScript",
    chapters: 50,
    status: "Coming Soon",
    logo: typescriptLogo,
    aliases: ["ts"],
    iconClassName: "bg-[conic-gradient(at_top,_#1e293b,_#e2e8f0,_#1e293b)]"
  },
  {
    id: "rust",
    name: "Rust",
    chapters: 50,
    status: "Coming Soon",
    logo: rustLogo,
    aliases: ["rs"],
    iconClassName: "bg-[radial-gradient(ellipse_at_bottom,_#500724,_#fb7185)]",
  },
];
