import { useEffect, useState } from "react";
import API from "../utils/api";
import { useRouter } from "next/router";

export default function Projects() {
  const [projects, setProjects] = useState([]);
  const router = useRouter();

  const fetchProjects = async () => {
    const res = await API.get("/projects");
    const projects = res.data?.data?.data || [];
  setProjects(projects);
  };

  useEffect(() => {
  const token = localStorage.getItem("token");
    if (!token) {
    router.push("/"); 
  } else {
    fetchProjects();
  }
}, []);

  return (
    <div>
      <h2>Projects</h2>
      <button onClick={() => router.push("/tasks")}>Tasks</button>
      {projects.map((p) => (
        <div key={p.id}>{p.name}</div>
      ))}
    </div>
  );
}
