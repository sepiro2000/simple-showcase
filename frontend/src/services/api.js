import axios from 'axios';

const api = axios.create({
  baseURL: '/api',
});

export const getProducts = async () => {
  try {
    const response = await api.get('/products');
    return response.data;
  } catch (error) {
    throw new Error(error.response?.data?.error || 'Failed to fetch products');
  }
};

export const getProductById = async (id) => {
  try {
    const response = await api.get(`/products/${id}`);
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
    await api.post(`/products/${id}/like`);
    return true;
  } catch (error) {
    throw new Error(error.response?.data?.error || 'Failed to like product');
  }
}; 