import React from "react";
import ReactDOM from "react-dom";
import "./index.css";
import App from "./App";
import { WsLinkProvider } from "./GlobalContext";

ReactDOM.render(
  <React.StrictMode>
    <WsLinkProvider ws={new WebSocket("ws://localhost:8080/ws")}>
      <App />
    </WsLinkProvider>
  </React.StrictMode>,
  document.getElementById("root")
);
