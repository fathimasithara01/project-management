import { useState } from "react";
import API from "../utils/api";

export default function TaskModal({ onClose, refresh }) {
  const [title, setTitle] = useState("");

  const createTask = async () => {
    try {
      const payload = {
        title,
        description: "Task",
        project_id: 1,
        assigned_to: 1,
        status: "todo",
        due_date: new Date().toISOString(),
      };

      console.log("Payload:", payload); // 👈 debug

      await API.post("/tasks", payload);

      refresh();
      onClose();
    } catch (err) {
      console.error(err.response?.data || err);
      alert("Error creating task");
    }
  };

  return (
    <div>
      <h3>Create Task</h3>
      <input onChange={(e) => setTitle(e.target.value)} placeholder="Title" />
      <button onClick={createTask}>Create</button>
      <button onClick={onClose}>Close</button>
    </div>
  );
}