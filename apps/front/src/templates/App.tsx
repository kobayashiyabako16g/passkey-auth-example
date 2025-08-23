import { Link } from "@tanstack/react-router";

function App() {
  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-slate-50 to-slate-100">
      <div className="text-center space-y-8 max-w-lg mx-auto px-6">
        <div className="space-y-4">
          <h1 className="text-5xl font-black text-slate-900 tracking-tight">
            Welcome Back
          </h1>
          <p className="text-slate-600 text-xl leading-relaxed">
            Passkey Authentication is a secure and convenient way to
            authenticate users without the need for passwords.
          </p>
        </div>

        <div className="flex flex-col sm:flex-row gap-4 justify-center pt-4">
          <Link
            to="/auth/signin"
            className="inline-flex items-center justify-center px-8 py-3 text-lg font-medium text-white bg-emerald-600 hover:bg-emerald-700 rounded-lg transition-colors duration-200 focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:ring-offset-2"
          >
            Sign In
          </Link>
          <Link
            to="/auth/signup"
            className="inline-flex items-center justify-center px-8 py-3 text-lg font-medium text-emerald-600 bg-transparent border-2 border-emerald-600 hover:bg-emerald-50 rounded-lg transition-colors duration-200 focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:ring-offset-2"
          >
            Create Account
          </Link>
        </div>

        <div className="pt-8 text-sm text-slate-500">
          <p>Passkey Authentication Example</p>
        </div>
      </div>
    </div>
  );
}

export default App;
