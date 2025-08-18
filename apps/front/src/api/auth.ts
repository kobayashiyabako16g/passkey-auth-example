import { API_URL } from "@/configs/env";
import type {
  RegistrationResponse,
  AuthenticationResponse,
} from "@/templates/auth/authn";

export const signUp = async (username: string) => {
  const response = await fetch(`${API_URL}/passkey/register/start`, {
    method: "POST",
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ username }),
  });

  if (!response.ok) {
    throw new Error("Failed to sign up");
  }

  const data = await response.json();

  return data;
};

export const signUpFinish = async (dto: RegistrationResponse) => {
  const response = await fetch(`${API_URL}/passkey/register/finish`, {
    method: "POST",
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(dto),
  });

  if (!response.ok) {
    throw new Error("Failed to sign up");
  }

  const data = await response.json();

  return data;
};

export const loginStart = async (username: string) => {
  const response = await fetch(`${API_URL}/passkey/login/start`, {
    method: "POST",
    credentials: "include",
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

export const loginFinish = async (dto: AuthenticationResponse) => {
  const response = await fetch(`${API_URL}/passkey/login/finish`, {
    method: "POST",
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(dto),
  });

  if (!response.ok) {
    throw new Error("Failed to sign in");
  }

  const data = await response.json();

  return data;
};
