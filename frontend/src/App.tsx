import { useContext, useEffect, useState } from "react";
import { wsConnection } from "./service";
import { EntryForm } from "./component/EntryForm";
import { WsLinkContext } from "./GlobalContext";
import { useEventListener } from "./hooks/useEvents";
import { GameBoard } from "./component/GameBoard";
import { ChatArea } from "./component/ChatArea";
import { UserList as UserArea } from "./component/UserList";
function randomColor(str: string) {
  var hash = 0;
  for (var i = 0; i < str.length; i++) {
    hash = str.charCodeAt(i) + ((hash << 5) - hash);
  }
  var colour = "#";
  for (var i = 0; i < 3; i++) {
    var value = (hash >> (i * 8)) & 0xff;
    colour += ("00" + value.toString(16)).substr(-2);
  }
  return colour;
}
const ESCAPE_KEYS = ["KeyW", "KeyS", "ArrowUp", "ArrowDown", "Space"];
function App() {
  const [messageArray, setMessageArray] = useState<any[]>([]);
  const [userList, setUserList] = useState<any>({});
  const [view, setView] = useState({
    Ball: { x: 200, y: 200, radius: 10 },
    Player1: { x: 2, y: 150 },
    Player2: { x: 378, y: 150 },
  });
  const { ws: socket, msg } = useContext(WsLinkContext);
  useEventListener("keydown", handler);

  function handler(e: { key: any; code: any }) {
    const { key, code } = e;
    if (ESCAPE_KEYS.includes(String(code))) {
      socket?.send(
        JSON.stringify({
          content: code.toUpperCase(),
          type: "GAME_CMD",
        })
      );
    }
  }

  useEffect(() => {
    handleSocketMessage(msg);
  }, [msg]);

  const handleSocketMessage = (msg: string) => {
    const data = JSON.parse(msg);
    switch (data?.type) {
      case "USER_NAME":
        addUser(data);
        // add to user list  and decide random color
        break;
      case "MESSAGE":
        addNewMessage(data);
        // add message to message list
        break;
      case "NEW_USER":
        addUser(data);
        break;
      case "EXIT":
        removeUser(data);
        break;
      case "GAME_STATE":
        drawElement(data);
        break;
      default:
        break;
    }
  };

  const removeUser = (data: { content: string }) => {
    setUserList((prev: any) => {
      const users = { ...prev };
      delete users[data.content];
      return users;
    });
  };

  const addUser = (data: { content: string }) => {
    const color = randomColor(data.content);
    setUserList((prev: any) => {
      const users = { ...prev };
      users[data.content] = color;
      return users;
    });
  };

  const addNewMessage = (data: any) => {
    setMessageArray((prev) => [...prev, data]);
  };

  const handleSubmit = (e: {
    preventDefault: () => void;
    target: {
      querySelector: (arg0: string) => { (): any; new (): any; value: string };
    };
  }) => {
    e.preventDefault();
    socket?.send(
      JSON.stringify({
        content: e.target.querySelector("[name='message']").value,
        type: "MESSAGE",
      })
    );
    e.target.querySelector("[name='message']").value = "";
  };

  const drawElement = (data: { state: any }) => {
    const element = data?.state;
    if (!element) return;
    setView({ ...element });
  };

  return (
    <div className="bg-slate-50 h-screen grid grid-cols-12 py-4 px-20">
      <div className="col-start-1 col-span-5 h-full">
        <GameBoard height={400} width={400} state={view} />
      </div>
      <ChatArea
        handleSubmit={handleSubmit}
        messageArray={messageArray}
        userList={userList}
      />
      <UserArea userMap={userList} />
    </div>
  );
}

export default App;
