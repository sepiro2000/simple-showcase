import { useState, useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';
import { getProductById, likeProduct } from '../services/api';
import LoadingSpinner from '../components/LoadingSpinner';
import ErrorMessage from '../components/ErrorMessage';

export default function ProductDetailPage() {
  const { id } = useParams();
  const [product, setProduct] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchProduct = async () => {
      try {
        const data = await getProductById(id);
        setProduct(data);
        setError(null);
      } catch (err) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };

    fetchProduct();
  }, [id]);

  const handleLike = async () => {
    try {
      await likeProduct(id);
      setProduct(prev => ({ ...prev, likes: prev.likes + 1 }));
    } catch (err) {
      console.error('Failed to like product:', err);
    }
  };

  if (loading) return <LoadingSpinner />;
  if (error) return <ErrorMessage message={error} />;
  if (!product) return null;

  return (
    <div className="container mx-auto px-4 py-8">
      <Link
        to="/"
        className="inline-block mb-6 text-blue-500 hover:text-blue-600"
      >
        ← Back to List
      </Link>
      
      <div className="bg-white rounded-lg shadow-lg overflow-hidden">
        <img
          src={product.imageUrl || 'https://via.placeholder.com/600x400'}
          alt={product.name}
          className="w-full h-96 object-cover"
        />
        
        <div className="p-6">
          <h1 className="text-3xl font-bold text-gray-800 mb-2">{product.name}</h1>
          <p className="text-2xl text-blue-600 mb-4">${product.price.toFixed(2)}</p>
          <p className="text-gray-600 mb-6">{product.description}</p>
          
          <div className="flex items-center justify-between">
            <span className="text-gray-500">{product.likes} likes</span>
            <button
              onClick={handleLike}
              className="px-4 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600 focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-opacity-50"
            >
              ❤️ Like
            </button>
          </div>
        </div>
      </div>
    </div>
  );
} 