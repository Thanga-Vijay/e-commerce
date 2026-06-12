import { Link } from 'react-router-dom';
import { ShoppingCart, User, LogOut, LayoutDashboard, Package } from 'lucide-react';
import { useAuth } from '../context/AuthContext';
import { useCart } from '../context/CartContext';

const Header = () => {
  const { user, logout, isAdmin } = useAuth();
  const { itemCount } = useCart();

  return (
    <header className="bg-white shadow-md sticky top-0 z-50">
      <div className="container mx-auto px-4">
        <div className="flex items-center justify-between h-16">
          <Link to="/" className="text-2xl font-bold text-primary-600">
            E-Commerce
          </Link>

          <nav className="hidden md:flex items-center space-x-6">
            <Link to="/products" className="hover:text-primary-600 transition-colors">
              Products
            </Link>
            {user && (
              <>
                <Link to="/orders" className="hover:text-primary-600 transition-colors">
                  Orders
                </Link>
                {isAdmin && (
                  <Link to="/admin" className="hover:text-primary-600 transition-colors">
                    Admin
                  </Link>
                )}
              </>
            )}
          </nav>

          <div className="flex items-center space-x-4">
            {user ? (
              <>
                <Link to="/cart" className="relative">
                  <ShoppingCart className="w-6 h-6 hover:text-primary-600 transition-colors" />
                  {itemCount > 0 && (
                    <span className="absolute -top-2 -right-2 bg-primary-600 text-white text-xs rounded-full w-5 h-5 flex items-center justify-center">
                      {itemCount}
                    </span>
                  )}
                </Link>
                
                <div className="relative group">
                  <button className="flex items-center space-x-2 hover:text-primary-600 transition-colors">
                    <User className="w-6 h-6" />
                    <span className="hidden md:inline">{user.firstName}</span>
                  </button>
                  
                  <div className="absolute right-0 mt-2 w-48 bg-white rounded-lg shadow-lg py-2 hidden group-hover:block">
                    <Link
                      to="/profile"
                      className="flex items-center px-4 py-2 hover:bg-gray-100"
                    >
                      <User className="w-4 h-4 mr-2" />
                      Profile
                    </Link>
                    <Link
                      to="/orders"
                      className="flex items-center px-4 py-2 hover:bg-gray-100"
                    >
                      <Package className="w-4 h-4 mr-2" />
                      My Orders
                    </Link>
                    {isAdmin && (
                      <Link
                        to="/admin"
                        className="flex items-center px-4 py-2 hover:bg-gray-100"
                      >
                        <LayoutDashboard className="w-4 h-4 mr-2" />
                        Admin Dashboard
                      </Link>
                    )}
                    <button
                      onClick={logout}
                      className="flex items-center w-full px-4 py-2 hover:bg-gray-100 text-red-600"
                    >
                      <LogOut className="w-4 h-4 mr-2" />
                      Logout
                    </button>
                  </div>
                </div>
              </>
            ) : (
              <div className="flex items-center space-x-4">
                <Link to="/login" className="hover:text-primary-600 transition-colors">
                  Login
                </Link>
                <Link to="/register" className="btn btn-primary">
                  Sign Up
                </Link>
              </div>
            )}
          </div>
        </div>
      </div>
    </header>
  );
};

export default Header;
