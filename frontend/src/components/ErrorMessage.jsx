export default function ErrorMessage({ message }) {
  return (
    <div className="p-4 text-red-600 bg-red-50 border border-red-200 rounded-lg">
      <p>{message}</p>
    </div>
  );
} 