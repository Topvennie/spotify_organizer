import { createContext } from "react";
import { User } from "../types/user";

interface AuthContextType {
  user: User | null;
  isLoading: boolean;
  forbidden: boolean;
  login: () => void;
  logout: () => void;
}

export const AuthContext = createContext<AuthContextType | undefined>(undefined);
