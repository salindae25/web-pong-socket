import React from "react";
import { EntryForm } from "./EntryForm";

export function ChatArea({
  handleSubmit,
  messageArray,
  userList,
}: {
  handleSubmit: Function;
  messageArray: any[];
  userList: {};
}): JSX.Element {
  return (
    <div className=" grid grid-rows-6  col-start-6 col-span-4 bg-slate-300 rounded-lg ">
      <EntryForm handleSubmit={handleSubmit} />
      <div className="row-start-1 row-span-6  h-5/6 flex flex-col gap-4 w-full overflow-y-scroll scroll-m-1 pt-2">
        {messageArray //.filter((y) => y?.type === "msg")
          .map((x: { content: string; userName: string | number }) => {
            return (
              x && (
                <div
                  key={x.content}
                  className="px-2 py-1  shadow-md w-4/6 rounded-md mx-2 font-medium text-yellow-200 capitalize"
                  style={{
                    background: userList[x.userName],
                  }}
                >
                  {x?.content}
                </div>
              )
            );
          })}
      </div>
    </div>
  );
}
