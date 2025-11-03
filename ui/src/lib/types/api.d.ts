export namespace API {
  interface Base extends JSON {
    id: number;
  }

  export interface User extends Base {
    name: string;
    email: string;
    uid: string;
  }
}
