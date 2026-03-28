import { useEffect, useState } from "react";
import API from "../utils/api";
import TaskModal from "../components/TaskModal";

export default function Tasks() {
  const [tasks, setTasks] = useState([]);
  const [open, setOpen] = useState(false);

  const fetchTasks = async () => {
    const res = await API.get("/tasks?page=1&limit=10");
    console.log(res.data)
    setTasks(res.data.data.data);
  };

  useEffect(() => {
  const token = localStorage.getItem("token");
  if (token) {
    fetchTasks();
  }
}, []);

  return (
    <div>
      <h2>Tasks</h2>
      <button onClick={() => setOpen(true)}>Create Task</button>

      {tasks.map((t) => (
        <div key={t.id}>
          {t.title} - {t.status}
        </div>
      ))}

      {open && (
        <TaskModal
          onClose={() => setOpen(false)}
          refresh={fetchTasks}
        />
      )}
    </div>
  );
}