import axios from 'axios';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL;

const api = axios.create({
  baseURL: API_BASE_URL,
  withCredentials: true, // Enable sending cookies if needed
});

export const getProducts = async () => {
  try {
    const response = await api.get('/api/products');
    return response.data;
  } catch (error) {
    throw new Error(error.response?.data?.error || 'Failed to fetch products');
  }
};

export const getProductById = async (id) => {
  try {
    const response = await api.get(`/api/products/${id}`);
    return response.data;
  } catch (error) {
    if (error.response?.status === 404) {
      throw new Error('Product not found');
    }
    throw new Error(error.response?.data?.error || 'Failed to fetch product');
  }
};

export const likeProduct = async (id) => {
  try {
    await api.post(`/api/products/${id}/like`);
    return true;
  } catch (error) {
    throw new Error(error.response?.data?.error || 'Failed to like product');
  }
}; 