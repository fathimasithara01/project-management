import { useRouter } from "next/router";

export const useAuth = () => {
  const router = useRouter();

  const login = (token) => {
    localStorage.setItem("token", token);
    router.push("/projects");
  };

  const logout = () => {
    localStorage.removeItem("token");
    router.push("/");
  };

  const isAuthenticated = () => {
    return !!localStorage.getItem("token");
  };

  return { login, logout, isAuthenticated };
};
