import { createFileRoute } from "@tanstack/react-router";
import SignIn from "@/templates/auth/SignIn";

export const Route = createFileRoute("/auth/signin")({
  component: RouteComponent,
});

function RouteComponent() {
  return <SignIn />;
}
