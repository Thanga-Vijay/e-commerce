# E-Commerce Frontend

React.js application built with Vite, TypeScript, and Material UI.

## Technology Stack

- **Framework:** React 18
- **Build Tool:** Vite
- **Language:** TypeScript
- **State Management:** Redux Toolkit
- **Routing:** React Router v6
- **HTTP Client:** Axios
- **UI Library:** Material UI (MUI)
- **Form Handling:** React Hook Form
- **Validation:** Yup
- **Authentication:** JWT stored in httpOnly cookies

## Features

### Customer Features

#### Authentication
- Login / Signup
- Email verification
- Password reset
- Profile management

#### Shopping
- Browse products
- Search & filter products
- Product details & reviews
- Add to cart
- Wishlist
- Checkout process
- Order tracking

#### Account
- Order history
- Saved addresses
- Profile settings

### Admin Features

#### Dashboard
- Revenue analytics
- Order statistics
- Customer metrics
- Product performance

#### Management
- Product CRUD
- Category management
- Inventory tracking
- Order management
- Customer management

#### Reports
- Sales reports
- Export CSV/PDF

## Project Structure

```
frontend/
├── public/
├── src/
│   ├── assets/
│   ├── components/
│   │   ├── common/
│   │   ├── layout/
│   │   ├── product/
│   │   ├── cart/
│   │   └── order/
│   ├── features/
│   │   ├── auth/
│   │   ├── products/
│   │   ├── cart/
│   │   ├── orders/
│   │   └── admin/
│   ├── hooks/
│   ├── pages/
│   │   ├── Home.tsx
│   │   ├── Products.tsx
│   │   ├── ProductDetail.tsx
│   │   ├── Cart.tsx
│   │   ├── Checkout.tsx
│   │   ├── Orders.tsx
│   │   ├── Login.tsx
│   │   ├── Register.tsx
│   │   ├── Profile.tsx
│   │   └── admin/
│   ├── services/
│   │   ├── api.ts
│   │   ├── auth.service.ts
│   │   ├── product.service.ts
│   │   ├── cart.service.ts
│   │   └── order.service.ts
│   ├── store/
│   │   ├── store.ts
│   │   └── slices/
│   ├── types/
│   ├── utils/
│   ├── App.tsx
│   └── main.tsx
├── .env.example
├── .eslintrc.js
├── index.html
├── package.json
├── tsconfig.json
├── vite.config.ts
└── README.md
```

## API Integration

Base URL: `http://localhost:8080/api/v1`

Axios instance with interceptors for:
- JWT token injection
- Token refresh on 401
- Error handling
- Loading states

## State Management

Redux Toolkit slices:
- `authSlice` - User authentication
- `productSlice` - Products & categories
- `cartSlice` - Shopping cart
- `orderSlice` - Orders
- `wishlistSlice` - Wishlist
- `uiSlice` - UI state (loading, notifications)

## Routing

```tsx
<Routes>
  {/* Public routes */}
  <Route path="/" element={<Home />} />
  <Route path="/products" element={<Products />} />
  <Route path="/products/:id" element={<ProductDetail />} />
  <Route path="/login" element={<Login />} />
  <Route path="/register" element={<Register />} />
  
  {/* Protected routes */}
  <Route element={<ProtectedRoute />}>
    <Route path="/cart" element={<Cart />} />
    <Route path="/checkout" element={<Checkout />} />
    <Route path="/orders" element={<Orders />} />
    <Route path="/profile" element={<Profile />} />
    <Route path="/wishlist" element={<Wishlist />} />
  </Route>
  
  {/* Admin routes */}
  <Route element={<AdminRoute />}>
    <Route path="/admin" element={<AdminDashboard />} />
    <Route path="/admin/products" element={<AdminProducts />} />
    <Route path="/admin/orders" element={<AdminOrders />} />
    <Route path="/admin/customers" element={<AdminCustomers />} />
    <Route path="/admin/reports" element={<AdminReports />} />
  </Route>
</Routes>
```

## Environment Variables

```env
VITE_API_BASE_URL=http://localhost:8080/api/v1
VITE_STRIPE_PUBLIC_KEY=pk_test_...
```

## Development

```bash
npm install
npm run dev
```

## Build

```bash
npm run build
npm run preview
```

## Docker

```bash
docker build -t frontend:latest .
docker run -p 3000:80 frontend:latest
```


Test