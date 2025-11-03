import { useAuth } from "@/lib/hooks/useAuth";
import { Forbidden } from "@/pages/Forbidden";
import { Login } from "@/pages/Login";
import { PropsWithChildren } from "react";

export const AuthLayout = ({ children }: PropsWithChildren) => {
  const { user, isLoading, forbidden } = useAuth();

  if (isLoading) {
    // Avoid a brief flickering of the loging view when the user is already logged in
    return null
  }

  if (forbidden) {
    return <Forbidden />
  }

  if (!user) {
    return <Login />
  }

  return children;
}
