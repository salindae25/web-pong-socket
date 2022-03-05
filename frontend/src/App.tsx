import { useContext, useEffect, useState } from "react";
import { wsConnection } from "./service";
import { EntryForm } from "./component/EntryForm";
import { WsLinkContext } from "./GlobalContext";
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

function App() {
  const [messageArray, setMessageArry] = useState<any[]>([]);
  const [userList, setUserList] = useState<any>({});
  const { ws: socket, msg } = useContext(WsLinkContext);

  useEffect(() => {
    handleSocketMessage(msg);
  }, [msg]);
  const handleSocketMessage = (msg) => {
    const data = JSON.parse(msg);
    switch (data?.type) {
      case "newJoin":
        addUser(data);
        // add to user list  and decide random color
        break;
      case "msg":
        addNewMessage(data);
        // add message to message list
        break;
      default:
        break;
    }
  };

  const addUser = (data) => {
    const color = randomColor(data.userName);
    setUserList((prev) => {
      const users = { ...prev };
      users[data.userName] = color;
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
        type: "msg",
      })
    );
    e.target.querySelector("[name='message']").value = "";
  };
  console.log(userList);

  return (
    <div className="bg-slate-50 h-screen grid grid-cols-12 py-4 px-20">
      <div className=" grid grid-rows-6  col-start-1 col-span-8 bg-slate-300 rounded-lg ">
        <EntryForm handleSubmit={handleSubmit} />
        <div className="row-start-1 row-span-6  h-5/6 flex flex-col gap-4 w-full overflow-y-scroll scroll-m-1 pt-2">
          {messageArray
            .filter((y) => y?.type === "msg")
            .map((x) => {
              return (
                <div
                  key={x?.content}
                  className="px-2 py-1  shadow-md w-4/6 rounded-md mx-2 font-medium text-yellow-200 capitalize"
                  style={{ background: userList[x.userName] }}
                >
                  {x?.content}
                </div>
              );
            })}
        </div>
      </div>
      <div className="col-start-9 col-span-4 flex flex-col gap-4 bg-red-200 px-4 py-4">
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
