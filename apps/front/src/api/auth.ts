import { API_URL } from "@/configs/env";

export const signUp = async (username: string) => {
  console.log("Signing up user:", username);
  const response = await fetch(`${API_URL}/posts`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ username }),
  });

  if (!response.ok) {
    throw new Error("Failed to sign in");
  }

  const data = await response.json();
  return data;
};
