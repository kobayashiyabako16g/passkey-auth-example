import { useState } from "react";
import { signUp, signUpFinish } from "@/api/auth";
import { createPasskey } from "./authn";

function SignUp() {
  const [username, setUsername] = useState("");

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const data = await signUp(username);
    console.log("data", data);
    try {
      const res = await createPasskey(data.publicKey);
      console.log("res", res);

      const result = await signUpFinish(res);
      console.log("result", result);
    } catch (error) {
      console.error(error);
    }
  };

  return (
    <div className="flex items-center justify-center h-screen">
      <form
        className="bg-white shadow-md rounded px-8 pt-6 pb-8 mb-4"
        onSubmit={handleSubmit}
      >
        <div className="mb-4">
          <label className="block text-sm font-bold mb-2" htmlFor="username">
            Username
          </label>
          <input
            className="shadow appearance-none border rounded py-2 px-3 leading-tight"
            type="text"
            placeholder="Username"
            maxLength={30}
            required
            onChange={(e) => setUsername(e.target.value)}
            value={username}
          />
        </div>
        <div className="flex items-center justify-between">
          <button
            className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded border-white"
            type="submit"
          >
            Sign Up
          </button>
        </div>
      </form>
    </div>
  );
}

export default SignUp;
