import { createFileRoute } from "@tanstack/react-router";
import SignUp from "@/templates/auth/SignUp";

export const Route = createFileRoute("/auth/signup")({
  component: RouteComponent,
});

function RouteComponent() {
  return <SignUp />;
}
