import { Link } from 'react-router-dom';
import { ShoppingCart } from 'lucide-react';
import { useCart } from '../context/CartContext';

const ProductCard = ({ product }) => {
  const { addItem } = useCart();

  const handleAddToCart = async () => {
    try {
      await addItem(product.id, 1);
    } catch (error) {
      // Error already handled by CartContext
    }
  };

  return (
    <div className="card hover:shadow-lg transition-shadow">
      <Link to={`/products/${product.id}`}>
        <div className="aspect-w-1 aspect-h-1 mb-4 overflow-hidden rounded-lg">
          <img
            src={product.imageUrl || 'https://via.placeholder.com/300'}
            alt={product.name}
            className="w-full h-48 object-cover hover:scale-105 transition-transform duration-300"
          />
        </div>
      </Link>
      
      <div className="space-y-2">
        <Link to={`/products/${product.id}`}>
          <h3 className="font-semibold text-lg hover:text-primary-600 transition-colors">
            {product.name}
          </h3>
        </Link>
        
        <p className="text-gray-600 text-sm line-clamp-2">{product.description}</p>
        
        <div className="flex items-center justify-between">
          <span className="text-2xl font-bold text-primary-600">
            ${product.price.toFixed(2)}
          </span>
          
          {product.stock > 0 ? (
            <button
              onClick={handleAddToCart}
              className="btn btn-primary flex items-center space-x-2"
            >
              <ShoppingCart className="w-4 h-4" />
              <span>Add to Cart</span>
            </button>
          ) : (
            <span className="text-red-600 font-semibold">Out of Stock</span>
          )}
        </div>
        
        {product.category && (
          <span className="inline-block px-3 py-1 bg-gray-200 text-gray-700 text-xs rounded-full">
            {product.category}
          </span>
        )}
      </div>
    </div>
  );
};

export default ProductCard;
