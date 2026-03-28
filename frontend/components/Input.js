export default function Input({ placeholder, type = "text", onChange }) {
  return <input type={type} placeholder={placeholder} onChange={onChange} />;
}