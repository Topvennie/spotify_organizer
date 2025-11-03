import { API } from "./api";

export interface User {
  id: number;
  name: string;
  email: string;
  uid: string;
}

// Converters

export const convertUserToModel = (user: API.User): User => {
  return user
}
