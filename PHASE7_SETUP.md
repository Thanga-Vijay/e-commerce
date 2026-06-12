# E-Commerce Platform - Phase 7 Implementation

This document provides instructions for running the React Frontend (Phase 7).

## Service Overview

**Phase 7** adds the React Frontend:
- **Frontend React App** (Port 3000): Complete e-commerce user interface with admin dashboard

## Quick Start with Docker Compose

Start all services including frontend:

```bash
# Build and start all services
docker-compose up --build

# Run in detached mode
docker-compose up -d --build

# View logs for frontend
docker-compose logs -f frontend

# Stop all services
docker-compose down
```

Services available:
- **Frontend**: http://localhost:3000
- **Auth Service**: http://localhost:8081
- **Product Service**: http://localhost:8082
- **Cart Service**: http://localhost:8083
- **Wishlist Service**: http://localhost:8084
- **Order Service**: http://localhost:8085
- **Payment Service**: http://localhost:8086
- **Inventory Service**: http://localhost:8087
- **Notification Service**: http://localhost:8088
- **Reporting Service**: http://localhost:8089

## Manual Setup

### 1. Install Dependencies

```bash
cd frontend
npm install
```

### 2. Configure Environment

```bash
cp .env.example .env
# Edit .env with your API URLs
```

### 3. Run Development Server

```bash
npm run dev
```

The app will be available at http://localhost:3000

### 4. Build for Production

```bash
npm run build
npm run preview
```

## Features

### Customer Features

#### 1. User Authentication
- **Registration**: Create new account with email, password, name, phone
- **Login**: Secure authentication with JWT tokens
- **Token Refresh**: Automatic token refresh on expiry
- **Profile Management**: Update user information
- **Logout**: Secure logout with token invalidation

#### 2. Product Catalog
- **Product Listing**: Grid view with pagination
- **Search**: Real-time product search
- **Category Filter**: Filter by product categories
- **Product Details**: Detailed product information
- **Product Images**: High-quality product images
- **Stock Status**: Real-time stock availability

#### 3. Shopping Cart
- **Add to Cart**: Add products with quantity selection
- **Update Quantity**: Increase/decrease item quantities
- **Remove Items**: Delete items from cart
- **Cart Summary**: Real-time price calculations
- **Persistent Cart**: Cart saved per user

#### 4. Checkout Process
- **Shipping Address**: Enter shipping details
- **Billing Address**: Separate or same as shipping
- **Order Summary**: Review before placing order
- **Order Confirmation**: Order number and details
- **Tax Calculation**: 8% tax applied
- **Shipping Cost**: $10 flat rate

#### 5. Order Management
- **Order History**: View all past orders
- **Order Details**: Complete order information
- **Order Tracking**: Track order status
- **Order Status**: pending, confirmed, processing, shipped, delivered, cancelled
- **Cancel Order**: Cancel pending/confirmed orders

### Admin Features

#### 1. Dashboard Metrics
- **Total Revenue**: All-time revenue
- **Total Orders**: Order count
- **Total Customers**: Unique customer count
- **Average Order Value**: Revenue per order
- **Today's Metrics**: Today's revenue and orders
- **Monthly Metrics**: Current month statistics

#### 2. Top Products
- **Product Ranking**: Top 5 best-selling products
- **Sales Data**: Total quantity sold
- **Revenue Data**: Total revenue per product

#### 3. Report Export
- **CSV Export**: Download metrics as CSV
- **PDF Export**: Download formatted PDF report
- **Dashboard Export**: Export current dashboard data

## Technology Stack

### Frontend Framework
- **React 18**: Latest React with hooks
- **Vite**: Fast build tool and dev server
- **React Router v6**: Client-side routing

### Styling
- **Tailwind CSS**: Utility-first CSS framework
- **Custom Components**: Reusable styled components
- **Responsive Design**: Mobile-first approach
- **Lucide Icons**: Beautiful icon library

### State Management
- **Context API**: Global state management
- **AuthContext**: Authentication state
- **CartContext**: Shopping cart state

### API Integration
- **Axios**: HTTP client with interceptors
- **Token Management**: Automatic token refresh
- **Error Handling**: Centralized error handling
- **React Hot Toast**: Toast notifications

### Development Tools
- **ESLint**: Code linting
- **PostCSS**: CSS processing
- **Autoprefixer**: CSS vendor prefixes

## Pages

### Public Pages
- **Home** (`/`): Landing page with features
- **Login** (`/login`): User login
- **Register** (`/register`): User registration
- **Products** (`/products`): Product catalog with search and filters
- **Product Details** (`/products/:id`): Individual product page

### Protected Pages (Require Login)
- **Cart** (`/cart`): Shopping cart
- **Checkout** (`/checkout`): Checkout process
- **Profile** (`/profile`): User profile management
- **Orders** (`/orders`): Order history
- **Order Details** (`/orders/:id`): Individual order details

### Admin Pages (Require Admin Role)
- **Admin Dashboard** (`/admin`): Analytics dashboard with reports

## Components

### Layout Components
- **Header**: Navigation with cart counter and user menu
- **Footer**: Footer with links and contact info

### Route Guards
- **PrivateRoute**: Protects authenticated routes
- **AdminRoute**: Protects admin-only routes

### UI Components
- **ProductCard**: Product display card
- **LoadingSpinner**: Loading indicator

## API Services

### Authentication Service
- `register(userData)`: Create new account
- `login(credentials)`: User login
- `logout(refreshToken)`: User logout
- `refreshToken(refreshToken)`: Refresh access token
- `getProfile(token)`: Get user profile
- `updateProfile(token, userData)`: Update profile

### Product Service
- `getProducts(page, limit, search, category)`: Get products with filters
- `getProductById(id)`: Get single product
- `getCategories()`: Get all categories

### Cart Service
- `getCart()`: Get user's cart
- `addItem(productId, quantity)`: Add item to cart
- `updateItem(itemId, quantity)`: Update item quantity
- `removeItem(itemId)`: Remove item from cart
- `clearCart()`: Clear entire cart

### Order Service
- `createOrder(orderData)`: Create new order
- `getOrders(page, limit)`: Get user's orders
- `getOrderById(id)`: Get single order
- `trackOrder(id)`: Track order status
- `cancelOrder(id)`: Cancel order
- `updateOrderStatus(id, status)`: Update order status (admin)

### Reporting Service
- `getDashboard()`: Get dashboard metrics
- `exportDashboard(format)`: Export dashboard as CSV/PDF
- `getRevenueReport(startDate, endDate, period)`: Get revenue report
- `getTopProducts(limit)`: Get top products
- `getCustomerReport(limit)`: Get customer analytics

## User Flows

### Customer Journey
1. **Browse Products**
   - Visit home page
   - Navigate to products page
   - Search or filter products
   - View product details

2. **Add to Cart**
   - Select quantity
   - Add to cart
   - View cart summary
   - Update quantities or remove items

3. **Checkout**
   - Proceed to checkout
   - Enter shipping address
   - Enter billing address
   - Review order summary
   - Place order

4. **Track Order**
   - View order confirmation
   - Check order history
   - Track order status
   - View order details

### Admin Workflow
1. **Access Dashboard**
   - Login with admin account
   - Navigate to admin dashboard

2. **View Metrics**
   - Check total revenue
   - Review order statistics
   - Monitor customer count
   - Analyze average order value

3. **Analyze Performance**
   - View top products
   - Check monthly trends
   - Review today's performance

4. **Export Reports**
   - Export dashboard as CSV
   - Download PDF report
   - Share with stakeholders

## Environment Variables

```bash
# API Service URLs
VITE_API_URL=http://localhost:8081
VITE_AUTH_SERVICE_URL=http://localhost:8081
VITE_PRODUCT_SERVICE_URL=http://localhost:8082
VITE_CART_SERVICE_URL=http://localhost:8083
VITE_WISHLIST_SERVICE_URL=http://localhost:8084
VITE_ORDER_SERVICE_URL=http://localhost:8085
VITE_PAYMENT_SERVICE_URL=http://localhost:8086
VITE_INVENTORY_SERVICE_URL=http://localhost:8087
VITE_NOTIFICATION_SERVICE_URL=http://localhost:8088
VITE_REPORTING_SERVICE_URL=http://localhost:8089
```

## Responsive Design

### Breakpoints
- **Mobile**: < 768px
- **Tablet**: 768px - 1024px
- **Desktop**: > 1024px

### Mobile Features
- Hamburger menu for navigation
- Touch-friendly buttons
- Responsive grid layouts
- Optimized images
- Mobile-first CSS

## Security Features

### Authentication
- JWT token-based authentication
- Secure token storage in localStorage
- Automatic token refresh
- Protected routes with guards
- HTTPS support (production)

### API Security
- CORS configuration
- Request/response interceptors
- Token validation
- Error handling
- XSS protection (React escaping)

## Performance Optimizations

### Code Splitting
- Lazy loading with React.lazy()
- Route-based code splitting
- Dynamic imports

### Caching
- Browser caching for static assets
- API response caching
- Image optimization

### Build Optimization
- Minification (JS, CSS)
- Tree shaking
- Asset compression (gzip)
- CDN-ready builds

## Testing

### Unit Tests
```bash
npm run test
```

### Linting
```bash
npm run lint
```

### Build Check
```bash
npm run build
```

## Deployment

### Docker Deployment
```bash
# Build image
docker build -t ecommerce-frontend .

# Run container
docker run -p 3000:80 ecommerce-frontend
```

### Production Build
```bash
npm run build
# Deploy dist/ folder to web server
```

## Troubleshooting

### Common Issues

**1. API Connection Failed:**
- Check service URLs in .env
- Verify backend services are running
- Check CORS configuration

**2. Login Not Working:**
- Verify auth service is running on port 8081
- Check JWT token in localStorage
- Clear browser cache and localStorage

**3. Cart Not Updating:**
- Check cart service connection
- Verify user is authenticated
- Check browser console for errors

**4. Build Errors:**
- Delete node_modules and package-lock.json
- Run `npm install` again
- Check Node.js version (requires 18+)

**5. Styling Issues:**
- Verify Tailwind CSS is configured
- Check PostCSS configuration
- Clear browser cache

### Debug Mode
```bash
# Run with verbose logging
npm run dev -- --debug
```

### Clear Cache
```bash
# Clear npm cache
npm cache clean --force

# Clear browser data
# Use browser dev tools > Application > Clear storage
```

## Browser Support

- Chrome/Edge: Latest 2 versions
- Firefox: Latest 2 versions
- Safari: Latest 2 versions
- Mobile browsers: iOS Safari, Chrome Android

## Accessibility

- Semantic HTML elements
- ARIA labels where needed
- Keyboard navigation support
- Focus visible states
- Screen reader friendly

## Future Enhancements

### Phase 8: Additional Features
- Product reviews and ratings
- Wishlist functionality
- Product recommendations
- Advanced search with filters
- Order returns and refunds

### Phase 9: Performance
- Server-side rendering (SSR)
- Progressive Web App (PWA)
- Offline support
- Push notifications

### Phase 10: Analytics
- Google Analytics integration
- User behavior tracking
- Conversion tracking
- A/B testing

## Resources

- [React Documentation](https://react.dev/)
- [Vite Documentation](https://vitejs.dev/)
- [Tailwind CSS](https://tailwindcss.com/)
- [React Router](https://reactrouter.com/)
- [Axios](https://axios-http.com/)
