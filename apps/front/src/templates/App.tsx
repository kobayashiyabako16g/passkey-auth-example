import { Link } from "@tanstack/react-router";

function App() {
  return (
    <div className="flex items-center justify-center h-screen">
      <Link to="/auth/signup">Sign Up</Link>
    </div>
  );
}

export default App;
