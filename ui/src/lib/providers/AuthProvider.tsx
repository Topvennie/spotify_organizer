import { useState, useEffect, useCallback, useMemo, PropsWithChildren } from "react";
import { notifications } from "@mantine/notifications";
import { isResponseNot200Error } from "../api/query";
import { useUser, useUserLogin, useUserLogout } from "../api/user";
import { AuthContext } from "../contexts/authContext";
import { User } from "../types/user";

export const AuthProvider = ({ children }: PropsWithChildren) => {
  const [user, setUser] = useState<User | null>(null);
  const [forbidden, setForbidden] = useState(false);

  const { data, isLoading, error } = useUser();
  const { mutate: logoutMutation } = useUserLogout();

  useEffect(() => {
    if (data) {
      setUser(data);
      setForbidden(false);
    }
  }, [data]);

  useEffect(() => {
    if (error && isResponseNot200Error(error)) {
      if (error.response.status === 403) {
        setForbidden(true);
        return;
      }
    }

    setForbidden(false);
  }, [error]);

  const logout = useCallback(() => {
    logoutMutation(undefined, {
      onSuccess: () => notifications.show({ variant: "success", message: "Logged out" }),
      onError: (err) => { throw new Error(`Logout failed ${err}`) },
      onSettled: () => setUser(null),
    });
  }, [logoutMutation]);

  const value = useMemo(() => ({ user, isLoading, forbidden, login: useUserLogin, logout }), [user, isLoading, forbidden, logout]);

  return <AuthContext value={value}>{children}</AuthContext>;
}
