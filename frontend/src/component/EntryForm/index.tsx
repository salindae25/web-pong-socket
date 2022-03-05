import React from "react";
export function EntryForm({ handleSubmit }) {
  return (
    <div className="row-start-7 px-4">
      <form
        onSubmit={handleSubmit}
        className="w-full flex items-center gap-3 py-4 "
      >
        <div className="w-full grid items-center">
          {/* <label
            htmlFor="message"
            className="block mb-2 text-sm font-medium text-gray-900 dark:text-gray-300"
          >
            Your messgae
          </label> */}
          <input
            type="text"
            name="message"
            className="shadow-sm bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500 dark:shadow-sm-light"
            placeholder="Enter your message"
            required
          />
        </div>
        <button
          type="submit"
          className="text-white w-40 bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800"
        >
          Send Message
        </button>
      </form>
    </div>
  );
}
