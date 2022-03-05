import React, { useState } from "react";

export const WsLinkContext = React.createContext<{
  ws: WebSocket | null;
  msg: any;
}>({
  ws: null,
  msg: null,
});

export const WsLinkProvider = ({
  children,
  ws,
}: {
  children: React.ReactNode;
  ws: WebSocket;
}) => {
  const [msg, setMsg] = useState(null);
  ws.onmessage = (eventM) => {
    setMsg(eventM.data);
  };
  return (
    <WsLinkContext.Provider value={{ ws: ws, msg: msg }}>
      {children}
    </WsLinkContext.Provider>
  );
};
