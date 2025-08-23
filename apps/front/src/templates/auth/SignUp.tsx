import { Link } from "@tanstack/react-router";
import { useState } from "react";
import { signUp, signUpFinish } from "@/api/auth";
import { createPasskey } from "./authn";

function SignUp() {
  const [username, setUsername] = useState("");
  const [error, setError] = useState("");
  const [isLoading, setIsLoading] = useState(false);

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
    <div className="min-h-screen flex items-center justify-center bg-slate-50 p-4">
      <div className="w-full max-w-md bg-white rounded-lg shadow-lg border border-slate-200">
        <div className="p-6 text-center border-b border-slate-200">
          <h1 className="text-2xl font-black text-slate-900">Create Account</h1>
          <p className="text-slate-600 mt-2">
            Sign up to get started with your new account
          </p>
        </div>
        <div className="p-6">
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="space-y-2">
              <label
                htmlFor="name"
                className="block text-sm font-medium text-slate-700"
              >
                User Name
              </label>
              <input
                name="username"
                type="text"
                maxLength={30}
                placeholder="Enter your name"
                onChange={(e) => setUsername(e.target.value)}
                value={username}
                className="w-full px-3 py-2 border border-slate-300 rounded-md focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500 transition-colors"
              />
            </div>
            {error && (
              <div className="text-red-600 text-sm text-center">{error}</div>
            )}
            <button
              type="submit"
              disabled={isLoading}
              className="w-full px-4 py-2 text-white bg-emerald-600 hover:bg-emerald-700 disabled:bg-emerald-400 rounded-md font-medium transition-colors duration-200 focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:ring-offset-2"
            >
              {isLoading ? "Creating Account..." : "Sign Up"}
            </button>
          </form>
          <div className="mt-6 text-center">
            <p className="text-slate-600 text-sm">
              Already have an account?{" "}
              <Link
                to="/auth/signin"
                className="text-emerald-600 hover:text-emerald-700 font-medium"
              >
                Sign in
              </Link>
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}

export default SignUp;
