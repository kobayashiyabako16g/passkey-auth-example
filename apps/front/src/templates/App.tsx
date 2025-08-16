import { Link } from "@tanstack/react-router";

function App() {
  return (
    <div className="flex items-center justify-center h-screen">
      <div className="bg-white rounded-lg p-10 flex flex-col shadow-md">
        <Link
          to="/auth/signup"
          className="mt-3 text-indigo-500 inline-flex items-center"
        >
          Sign Up
        </Link>
        <Link
          to="/auth/signin"
          className="mt-3 text-indigo-500 inline-flex items-center"
        >
          Sign In
        </Link>
      </div>
    </div>
  );
}

export default App;
