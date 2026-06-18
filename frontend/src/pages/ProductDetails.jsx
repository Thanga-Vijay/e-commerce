import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { ShoppingCart, ArrowLeft } from 'lucide-react';
import { productService } from '../api/product';
import { useCart } from '../context/CartContext';
import LoadingSpinner from '../components/LoadingSpinner';
import toast from 'react-hot-toast';

const ProductDetails = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const [product, setProduct] = useState(null);
  const [loading, setLoading] = useState(true);
  const [quantity, setQuantity] = useState(1);
  const { addItem } = useCart();

  useEffect(() => {
    fetchProduct();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [id]);

  const fetchProduct = async () => {
    try {
      setLoading(true);
      const response = await productService.getProductById(id);
      setProduct(response.data);
    } catch (error) {
      toast.error('Failed to fetch product');
      navigate('/products');
    } finally {
      setLoading(false);
    }
  };

  const handleAddToCart = async () => {
    try {
      await addItem(product.id, quantity);
    } catch (error) {
      // Error handled by CartContext
    }
  };

  if (loading) {
    return <LoadingSpinner />;
  }

  if (!product) {
    return null;
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <button
        onClick={() => navigate('/products')}
        className="flex items-center text-primary-600 hover:text-primary-700 mb-6"
      >
        <ArrowLeft className="w-5 h-5 mr-2" />
        Back to Products
      </button>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-12">
        {/* Product Image */}
        <div className="rounded-lg overflow-hidden">
          <img
            src={product.imageUrl || 'https://via.placeholder.com/600'}
            alt={product.name}
            className="w-full h-auto"
          />
        </div>

        {/* Product Info */}
        <div className="space-y-6">
          <div>
            <h1 className="text-4xl font-bold mb-2">{product.name}</h1>
            {product.category && (
              <span className="inline-block px-3 py-1 bg-gray-200 text-gray-700 text-sm rounded-full">
                {product.category}
              </span>
            )}
          </div>

          <div className="text-4xl font-bold text-primary-600">
            ${product.price.toFixed(2)}
          </div>

          <p className="text-gray-600 leading-relaxed">{product.description}</p>

          {product.stock > 0 ? (
            <div className="space-y-4">
              <div className="flex items-center space-x-4">
                <label className="font-medium">Quantity:</label>
                <div className="flex items-center border rounded-lg">
                  <button
                    onClick={() => setQuantity(Math.max(1, quantity - 1))}
                    className="px-4 py-2 hover:bg-gray-100"
                  >
                    -
                  </button>
                  <span className="px-6 py-2 border-x">{quantity}</span>
                  <button
                    onClick={() => setQuantity(Math.min(product.stock, quantity + 1))}
                    className="px-4 py-2 hover:bg-gray-100"
                  >
                    +
                  </button>
                </div>
              </div>

              <button
                onClick={handleAddToCart}
                className="btn btn-primary flex items-center justify-center space-x-2 w-full md:w-auto px-8 py-3 text-lg"
              >
                <ShoppingCart className="w-5 h-5" />
                <span>Add to Cart</span>
              </button>

              <p className="text-sm text-gray-600">
                {product.stock} items in stock
              </p>
            </div>
          ) : (
            <div className="text-red-600 font-semibold text-lg">Out of Stock</div>
          )}

          {/* Product Details */}
          <div className="border-t pt-6 space-y-2">
            <h3 className="font-semibold text-lg mb-4">Product Details</h3>
            <div className="text-sm text-gray-600 space-y-1">
              <p><span className="font-medium">SKU:</span> {product.sku}</p>
              <p><span className="font-medium">Category:</span> {product.category || 'N/A'}</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default ProductDetails;
