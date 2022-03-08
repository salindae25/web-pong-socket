import { useContext, useEffect, useRef, useState } from "react";
import { wsConnection } from "./service";
import { EntryForm } from "./component/EntryForm";
import { WsLinkContext } from "./GlobalContext";
import { useEventListener } from "./hooks/useEvents";
function randomColor(str) {
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
const ESCAPE_KEYS = ["a", "A", "d", "D", "40", "ArrowUp", "38", "ArrowDown"];
function App() {
  const [messageArray, setMessageArry] = useState<any[]>([]);
  const [userList, setUserList] = useState<any>({});
  const [ball, setBall] = useState({ x: 100, y: 75, radius: 50 });
  const { ws: socket, msg } = useContext(WsLinkContext);
  function handler({ key }) {
    console.log(key);

    if (ESCAPE_KEYS.includes(String(key))) {
      socket?.send(
        JSON.stringify({
          content: key.toUpperCase(),
          type: "GAME_CMD",
        })
      );
    }
  }

  useEventListener("keydown", handler);
  useEffect(() => {
    handleSocketMessage(msg);
  }, [msg]);
  const handleSocketMessage = (msg) => {
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
  const removeUser = (data) => {
    setUserList((prev) => {
      const users = { ...prev };
      console.log(data?.content + " removed");

      delete users[data.content];
      return users;
    });
  };
  const addUser = (data) => {
    const color = randomColor(data.content);
    setUserList((prev) => {
      const users = { ...prev };
      users[data.content] = color;
      return users;
    });
  };
  const addNewMessage = (data) => {
    setMessageArry((prev) => [...prev, data]);
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    socket?.send(
      JSON.stringify({
        content: e.target.querySelector("[name='message']").value,
        type: "MESSAGE",
      })
    );
    e.target.querySelector("[name='message']").value = "";
  };

  const drawElement = (data) => {
    const element = data?.state?.element;
    if (!element) return;
    setBall({
      x: element.position.x,
      y: element.position.y,
      radius: element.radius,
    });
  };

  return (
    <div className="bg-slate-50 h-screen grid grid-cols-12 py-4 px-20">
      <div className="col-start-1 col-span-5 h-full">
        <GameBoard height={400} width={400} ball={ball} />
      </div>
      <div className=" grid grid-rows-6  col-start-6 col-span-4 bg-slate-300 rounded-lg ">
        <EntryForm handleSubmit={handleSubmit} />
        <div className="row-start-1 row-span-6  h-5/6 flex flex-col gap-4 w-full overflow-y-scroll scroll-m-1 pt-2">
          {messageArray
            //.filter((y) => y?.type === "msg")
            .map((x) => {
              return (
                <div
                  key={x?.content}
                  className="px-2 py-1  shadow-md w-4/6 rounded-md mx-2 font-medium text-yellow-200 capitalize"
                  style={{ background: userList[x?.userName] }}
                >
                  {x?.content}
                </div>
              );
            })}
        </div>
      </div>
      <div className="col-start-10 col-span-2 flex flex-col gap-4 bg-red-200 px-4 py-4">
        {Object.keys(userList).map((x) => {
          return (
            <div
              className="font-medium text-yellow-200 rounded-md p-3 flex w-full items-center"
              style={{ background: userList[x] }}
            >
              <span className="w-full">{x}</span>
              <span className="w-2 h-2 bg-green-400 rounded-full border border-yellow-200"></span>
            </div>
          );
        })}
      </div>
    </div>
  );
}

export default App;

function GameBoard({
  height,
  width,
  ball,
}: {
  height: number;
  width: number;
  ball?: any;
}) {
  const canvasRef = useRef<any>(null);
  useEffect(() => {
    const canvas = canvasRef.current;
    const ctx = canvas.getContext("2d");
    ctx.fillStyle = "gray";
    ctx.fillRect(0, 0, 400, 400);
  }, []);
  useEffect(() => {
    const canvas = canvasRef.current;
    const ctx = canvas.getContext("2d");
    if (ball) {
      ctx.clearRect(0, 0, 400, 400);
      ctx.fillStyle = "gray";
      ctx.fillRect(0, 0, 400, 400);
      ctx.beginPath();
      ctx.arc(ball.x, ball.y, ball.radius, 0, 2 * Math.PI);
      ctx.stroke();
    }
  }, [ball]);
  return <canvas ref={canvasRef} id="canvas" height={height} width={width} />;
}
