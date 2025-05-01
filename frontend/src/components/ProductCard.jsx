import { Link } from 'react-router-dom';
import { likeProduct } from '../services/api';

export default function ProductCard({ product }) {
  const handleLike = async (e) => {
    e.preventDefault(); // Prevent navigation when clicking like button
    try {
      await likeProduct(product.id);
      // In a real app, you might want to update the UI or refetch the data
      console.log('Product liked successfully');
    } catch (error) {
      console.error('Failed to like product:', error);
    }
  };

  return (
    <Link to={`/products/${product.id}`} className="block">
      <div className="bg-white rounded-lg shadow-md overflow-hidden hover:shadow-lg transition-shadow">
        <div className="aspect-w-16 aspect-h-9">
          <img
            src={product.imageUrl || 'https://via.placeholder.com/300x200'}
            alt={product.name}
            className="w-full h-48 object-cover"
          />
        </div>
        <div className="p-4">
          <h3 className="text-lg font-semibold text-gray-800">{product.name}</h3>
          <p className="text-gray-600 mt-1">${product.price.toFixed(2)}</p>
          <div className="mt-2 flex items-center justify-between">
            <span className="text-sm text-gray-500">{product.likes} likes</span>
            <button
              onClick={handleLike}
              className="text-red-500 hover:text-red-600 focus:outline-none"
            >
              ❤️
            </button>
          </div>
        </div>
      </div>
    </Link>
  );
} 