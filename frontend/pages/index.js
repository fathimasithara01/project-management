import { useState } from "react";
import API from "../utils/api";
import { useAuth } from "../hooks/useAuth";

export default function Login() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const { login } = useAuth();

 const handleLogin = async () => {
  try {
    const res = await API.post("/login", { email, password });

    console.log(res.data); 

    login(res.data.data.token);

  } catch (err) {
    alert("Login failed");
  }
};

  return (
    <div>
      <h2>Login</h2>
      <input onChange={(e) => setEmail(e.target.value)} placeholder="Email" />
      <input type="password" onChange={(e) => setPassword(e.target.value)} placeholder="Password" />
      <button onClick={handleLogin}>Login</button>
    </div>
  );
}