import { useEffect, useState } from "react";
import logo from "./logo.svg";
import "./App.css";
import { wsConnection } from "./service";
import { EntryForm } from "./component/EntryForm";

function App() {
  const [count, setCount] = useState(0);
  const socket = wsConnection("//localhost:8080/ws");

  const handleSubmit = (e) => {
    e.preventDefault();
    socket({
      content: e.target.querySelector("[name='message']").value,
    });

    e.target.querySelector("[name='message']").value = "";
  };
  return (
    <div className="bg-slate-50 w-screen min-h-screen flex justify-center">
      <EntryForm handleSubmit={handleSubmit} />
    </div>
  );
}

export default App;
