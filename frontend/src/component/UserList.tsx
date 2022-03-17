import React from "react";

export function UserList({ userMap }: { userMap: any }) {
  return (
    <div className="col-start-10 col-span-2 flex flex-col gap-4 bg-red-200 px-4 py-4">
      {Object.keys(userMap).map((x) => {
        return (
          <div
            className="font-medium text-yellow-200 rounded-md p-3 flex w-full items-center"
            style={{
              background: userMap[x],
            }}
          >
            <span className="w-full">{x}</span>
            <span className="w-2 h-2 bg-green-400 rounded-full border border-yellow-200"></span>
          </div>
        );
      })}
    </div>
  );
}
