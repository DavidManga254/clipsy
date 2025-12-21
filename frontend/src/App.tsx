import "./App.css";
import React, { useState } from "react";
import {
  Copy,
  Check,
  Download,
  File,
  FileText,
  Monitor,
  Wifi,
  WifiOff,
  Clock,
  X,
  Minimize2,
  Maximize2,
  Settings,
  Search,
  Zap,
} from "lucide-react";

const mockPeers = [
  {
    id: "1",
    name: "David's MacBook",
    status: "connected",
    ip: "192.168.1.100",
    lastSeen: "2m",
  },
  {
    id: "2",
    name: "Work PC",
    status: "connected",
    ip: "192.168.1.101",
    lastSeen: "now",
  },
  {
    id: "3",
    name: "iPad Pro",
    status: "discovered",
    ip: "192.168.1.102",
    lastSeen: "5m",
  },
];

const mockHistory = {
  "1": [
    {
      id: "h1",
      type: "text",
      content: "https://github.com/wailsapp/wails",
      timestamp: "14:32",
      copied: false,
    },
    {
      id: "h2",
      type: "text",
      content: "npm install react-router-dom",
      timestamp: "14:27",
      copied: false,
    },
    {
      id: "h3",
      type: "image",
      url: "https://via.placeholder.com/300x200",
      name: "screenshot.png",
      size: "245 KB",
      timestamp: "14:22",
    },
    {
      id: "h4",
      type: "file",
      name: "project-docs.pdf",
      size: "2.4 MB",
      timestamp: "14:17",
    },
    {
      id: "h5",
      type: "text",
      content: "const handleSubmit = async () => {}",
      timestamp: "14:10",
      copied: false,
    },
  ],
  "2": [
    {
      id: "h6",
      type: "text",
      content: "Meeting at 3pm in conference room B",
      timestamp: "13:45",
      copied: false,
    },
    {
      id: "h7",
      type: "image",
      url: "https://via.placeholder.com/300x180",
      name: "design-mockup.jpg",
      size: "512 KB",
      timestamp: "12:30",
    },
  ],
  "3": [],
};

const mockIncomingTransfer = {
  id: "t1",
  from: "David's MacBook",
  fileName: "presentation.pptx",
  fileSize: "15.2 MB",
  type: "file",
};

function ClipsyApp() {
  const [selectedPeer, setSelectedPeer] = useState(mockPeers[0]);
  const [history, setHistory] = useState(mockHistory);
  const [showTransferPrompt, setShowTransferPrompt] = useState(true);
  const [activeTransfer, setActiveTransfer] = useState(null);
  const [copiedId, setCopiedId] = useState(null);
  const [searchQuery, setSearchQuery] = useState("");

  const handleCopy = (itemId, content) => {
    setCopiedId(itemId);
    setTimeout(() => setCopiedId(null), 1500);
  };

  const handleAcceptTransfer = () => {
    setShowTransferPrompt(false);
    setActiveTransfer({
      ...mockIncomingTransfer,
      progress: 0,
      speed: "2.3 MB/s",
      eta: "6s",
    });

    let progress = 0;
    const interval = setInterval(() => {
      progress += 10;
      if (progress > 100) {
        clearInterval(interval);
        setActiveTransfer(null);
        const newItem = {
          id: "h-new",
          type: "file",
          name: mockIncomingTransfer.fileName,
          size: mockIncomingTransfer.fileSize,
          timestamp: new Date().toLocaleTimeString("en-US", {
            hour: "2-digit",
            minute: "2-digit",
            hour12: false,
          }),
        };
        setHistory((prev) => ({
          ...prev,
          [selectedPeer.id]: [newItem, ...prev[selectedPeer.id]],
        }));
      } else {
        setActiveTransfer((prev) => ({ ...prev, progress }));
      }
    }, 600);
  };

  const handleRejectTransfer = () => {
    setShowTransferPrompt(false);
  };

  const getTypeIcon = (type) => {
    switch (type) {
      case "image":
        return (
          <div className="w-9 h-9 rounded-xl bg-gradient-to-br from-blue-500/20 to-cyan-500/20 border border-blue-500/30 flex items-center justify-center">
            <div className="w-4 h-4 bg-gradient-to-br from-blue-400 to-cyan-400 rounded" />
          </div>
        );
      case "file":
        return (
          <div className="w-9 h-9 rounded-xl bg-gradient-to-br from-purple-500/20 to-pink-500/20 border border-purple-500/30 flex items-center justify-center">
            <File className="w-4 h-4 text-purple-300" />
          </div>
        );
      default:
        return (
          <div className="w-9 h-9 rounded-xl bg-gradient-to-br from-slate-600/20 to-slate-500/20 border border-slate-500/30 flex items-center justify-center">
            <FileText className="w-4 h-4 text-slate-300" />
          </div>
        );
    }
  };

  return (
    <div className="flex flex-col h-screen bg-gradient-to-br from-slate-950 via-slate-900 to-slate-950 text-gray-100 select-none">
      {/* Modern Title Bar with Gradient */}
      <div
        className="h-11 bg-gradient-to-r from-slate-900/50 via-slate-800/50 to-slate-900/50 backdrop-blur-xl border-b border-white/5 flex items-center justify-between px-4"
        style={{ WebkitAppRegion: "drag" }}
      >
        <div className="flex items-center space-x-3 ml-16">
          <div className="flex items-center justify-center w-6 h-6 rounded-lg bg-gradient-to-br from-blue-500 to-cyan-500 shadow-lg shadow-blue-500/20">
            <Zap className="w-3.5 h-3.5 text-white" />
          </div>
          <span className="text-sm font-semibold bg-gradient-to-r from-blue-400 to-cyan-400 bg-clip-text text-transparent">
            Clipsy
          </span>
        </div>
        <div
          className="flex items-center space-x-1.5"
          style={{ WebkitAppRegion: "no-drag" }}
        >
          <button className="w-9 h-7 flex items-center justify-center hover:bg-white/5 rounded-lg transition-all">
            <Minimize2 className="w-3.5 h-3.5 text-gray-400 hover:text-gray-200" />
          </button>
          <button className="w-9 h-7 flex items-center justify-center hover:bg-white/5 rounded-lg transition-all">
            <Maximize2 className="w-3.5 h-3.5 text-gray-400 hover:text-gray-200" />
          </button>
          <button className="w-9 h-7 flex items-center justify-center hover:bg-red-500/90 rounded-lg transition-all group">
            <X className="w-3.5 h-3.5 text-gray-400 group-hover:text-white" />
          </button>
        </div>
      </div>

      <div className="flex flex-1 overflow-hidden">
        {/* Modern Sidebar with Glassmorphism */}
        <div className="w-64 bg-slate-900/30 backdrop-blur-xl border-r border-white/5 flex flex-col">
          <div className="p-4 space-y-3">
            <div className="flex items-center justify-between">
              <span className="text-xs font-semibold text-gray-400 uppercase tracking-wider">
                Network Peers
              </span>
              <button className="p-1.5 hover:bg-white/5 rounded-lg transition-all">
                <Settings className="w-4 h-4 text-gray-400 hover:text-gray-200" />
              </button>
            </div>
            <div className="relative">
              <Search className="absolute left-3 top-2.5 w-4 h-4 text-gray-500" />
              <input
                type="text"
                placeholder="Search peers..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="w-full bg-slate-800/50 border border-slate-700/50 rounded-xl text-sm px-9 py-2 text-gray-200 placeholder-gray-500 focus:outline-none focus:border-blue-500/50 focus:bg-slate-800/70 transition-all"
              />
            </div>
          </div>

          <div className="flex-1 overflow-y-auto px-2 pb-2">
            <div className="space-y-1.5">
              {mockPeers.map((peer) => (
                <button
                  key={peer.id}
                  onClick={() => setSelectedPeer(peer)}
                  className={`w-full text-left px-3 py-3 rounded-xl transition-all ${
                    selectedPeer.id === peer.id
                      ? "bg-gradient-to-r from-blue-500/10 to-cyan-500/10 border border-blue-500/20 shadow-lg shadow-blue-500/5"
                      : "hover:bg-white/5 border border-transparent"
                  }`}
                >
                  <div className="flex items-center space-x-3">
                    <div className="relative">
                      <div
                        className={`w-10 h-10 rounded-xl flex items-center justify-center ${
                          selectedPeer.id === peer.id
                            ? "bg-gradient-to-br from-blue-500/20 to-cyan-500/20 border border-blue-500/30"
                            : "bg-slate-800/50 border border-slate-700/50"
                        }`}
                      >
                        <Monitor className="w-5 h-5 text-gray-300" />
                      </div>
                      <div
                        className={`absolute -bottom-0.5 -right-0.5 w-3 h-3 rounded-full border-2 border-slate-900 ${
                          peer.status === "connected"
                            ? "bg-emerald-400 shadow-lg shadow-emerald-400/50"
                            : "bg-slate-600"
                        }`}
                      />
                    </div>
                    <div className="flex-1 min-w-0">
                      <p className="text-sm font-medium text-gray-200 truncate">
                        {peer.name}
                      </p>
                      <div className="flex items-center space-x-1.5 mt-0.5">
                        <span className="text-xs text-gray-500">{peer.ip}</span>
                        <span className="text-gray-700">•</span>
                        <span className="text-xs text-gray-500">
                          {peer.lastSeen}
                        </span>
                      </div>
                    </div>
                  </div>
                </button>
              ))}
            </div>
          </div>

          <div className="p-3 border-t border-white/5 bg-slate-900/30">
            <div className="flex items-center justify-between text-xs">
              <span className="text-gray-500">Connected</span>
              <div className="flex items-center space-x-1.5">
                <div className="w-2 h-2 rounded-full bg-emerald-400 animate-pulse" />
                <span className="font-medium text-emerald-400">
                  {mockPeers.filter((p) => p.status === "connected").length}
                </span>
              </div>
            </div>
          </div>
        </div>

        {/* Main Content with Modern Styling */}
        <div className="flex-1 flex flex-col">
          {/* Header with Gradient */}
          <div className="h-14 px-5 flex items-center justify-between border-b border-white/5 bg-slate-900/20 backdrop-blur-xl">
            <div className="flex items-center space-x-4">
              <div className="flex items-center space-x-2.5">
                <Monitor className="w-5 h-5 text-blue-400" />
                <span className="text-base font-semibold text-gray-100">
                  {selectedPeer.name}
                </span>
              </div>
              <div className="h-5 w-px bg-white/10" />
              <div className="flex items-center space-x-2 px-2.5 py-1 rounded-lg bg-emerald-500/10 border border-emerald-500/20">
                <div className="w-1.5 h-1.5 rounded-full bg-emerald-400 animate-pulse" />
                <span className="text-xs font-medium text-emerald-400">
                  Live
                </span>
              </div>
            </div>
            <div className="text-sm text-gray-400">
              {history[selectedPeer.id]?.length || 0} items
            </div>
          </div>

          {/* Transfer Prompt with Modern Alert */}
          {showTransferPrompt && (
            <div className="bg-gradient-to-r from-amber-500/10 via-orange-500/10 to-amber-500/10 border-b border-amber-500/20 px-5 py-3 backdrop-blur-xl">
              <div className="flex items-center justify-between">
                <div className="flex items-center space-x-3">
                  <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-amber-500/20 to-orange-500/20 border border-amber-500/30 flex items-center justify-center">
                    <File className="w-5 h-5 text-amber-400" />
                  </div>
                  <div>
                    <p className="text-sm font-semibold text-gray-100">
                      {mockIncomingTransfer.fileName}
                    </p>
                    <p className="text-xs text-gray-400">
                      {mockIncomingTransfer.fileSize} from{" "}
                      {mockIncomingTransfer.from}
                    </p>
                  </div>
                </div>
                <div className="flex items-center space-x-2">
                  <button
                    onClick={handleRejectTransfer}
                    className="px-4 py-2 text-sm font-medium text-gray-300 bg-slate-800/50 hover:bg-slate-700/50 rounded-xl border border-slate-700/50 transition-all"
                  >
                    Decline
                  </button>
                  <button
                    onClick={handleAcceptTransfer}
                    className="px-4 py-2 text-sm font-medium text-white bg-gradient-to-r from-blue-500 to-cyan-500 hover:from-blue-600 hover:to-cyan-600 rounded-xl shadow-lg shadow-blue-500/20 transition-all"
                  >
                    Accept
                  </button>
                </div>
              </div>
            </div>
          )}

          {/* Active Transfer with Modern Progress */}
          {activeTransfer && (
            <div className="bg-gradient-to-r from-blue-500/10 via-cyan-500/10 to-blue-500/10 border-b border-blue-500/20 px-5 py-3 backdrop-blur-xl">
              <div className="flex items-center justify-between mb-2.5">
                <div className="flex items-center space-x-3">
                  <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-blue-500/20 to-cyan-500/20 border border-blue-500/30 flex items-center justify-center">
                    <Download className="w-5 h-5 text-blue-400 animate-pulse" />
                  </div>
                  <div>
                    <p className="text-sm font-semibold text-gray-100">
                      {activeTransfer.fileName}
                    </p>
                    <p className="text-xs text-gray-400">
                      {activeTransfer.speed} • {activeTransfer.eta} remaining
                    </p>
                  </div>
                </div>
                <span className="text-sm font-bold text-blue-400">
                  {activeTransfer.progress}%
                </span>
              </div>
              <div className="relative w-full bg-slate-800/50 rounded-full h-2 overflow-hidden">
                <div
                  className="absolute top-0 left-0 h-full bg-gradient-to-r from-blue-500 to-cyan-500 rounded-full transition-all duration-300 shadow-lg shadow-blue-500/50"
                  style={{ width: `${activeTransfer.progress}%` }}
                />
              </div>
            </div>
          )}

          {/* History List with Modern Cards */}
          <div className="flex-1 overflow-y-auto p-4">
            {history[selectedPeer.id]?.length > 0 ? (
              <div className="space-y-2 max-w-4xl">
                {history[selectedPeer.id].map((item) => (
                  <div
                    key={item.id}
                    className="group bg-slate-800/30 backdrop-blur-xl hover:bg-slate-800/50 border border-slate-700/30 hover:border-slate-600/50 rounded-2xl p-4 transition-all hover:shadow-xl hover:shadow-black/20"
                  >
                    <div className="flex items-start space-x-3.5">
                      <div className="flex-shrink-0">
                        {getTypeIcon(item.type)}
                      </div>
                      <div className="flex-1 min-w-0">
                        {item.type === "text" && (
                          <p className="text-sm text-gray-200 break-words font-mono leading-relaxed bg-slate-900/30 rounded-lg p-3 border border-slate-700/30">
                            {item.content}
                          </p>
                        )}
                        {item.type === "image" && (
                          <div>
                            <img
                              src={item.url}
                              alt={item.name}
                              className="rounded-xl max-w-md mb-3 border border-slate-700/50 shadow-lg"
                            />
                            <p className="text-xs text-gray-400">
                              {item.name} • {item.size}
                            </p>
                          </div>
                        )}
                        {item.type === "file" && (
                          <div>
                            <p className="text-sm font-semibold text-gray-200">
                              {item.name}
                            </p>
                            <p className="text-xs text-gray-400 mt-1">
                              {item.size}
                            </p>
                          </div>
                        )}
                        <div className="flex items-center justify-between mt-3">
                          <span className="text-xs text-gray-500">
                            {item.timestamp}
                          </span>
                        </div>
                      </div>
                      <button
                        onClick={() => handleCopy(item.id, item.content)}
                        className="opacity-0 group-hover:opacity-100 p-2 text-gray-400 hover:text-white hover:bg-slate-700/50 rounded-xl transition-all border border-transparent hover:border-slate-600/50"
                        title="Copy"
                      >
                        {copiedId === item.id ? (
                          <Check className="w-4 h-4 text-emerald-400" />
                        ) : (
                          <Copy className="w-4 h-4" />
                        )}
                      </button>
                    </div>
                  </div>
                ))}
              </div>
            ) : (
              <div className="flex flex-col items-center justify-center h-full text-center">
                <div className="w-20 h-20 rounded-2xl bg-gradient-to-br from-slate-800/50 to-slate-700/50 border border-slate-700/50 flex items-center justify-center mb-4 shadow-xl">
                  <FileText className="w-10 h-10 text-gray-500" />
                </div>
                <h3 className="text-base font-semibold text-gray-300 mb-2">
                  No clipboard history yet
                </h3>
                <p className="text-sm text-gray-500 max-w-sm leading-relaxed">
                  When {selectedPeer.name} copies something, it will appear here
                  in real-time
                </p>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}

export default ClipsyApp;
