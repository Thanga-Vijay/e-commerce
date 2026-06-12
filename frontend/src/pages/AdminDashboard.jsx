import { useState, useEffect } from 'react';
import { DollarSign, Package, Users, TrendingUp, Download } from 'lucide-react';
import { reportingService } from '../api/reporting';
import LoadingSpinner from '../components/LoadingSpinner';
import toast from 'react-hot-toast';

const AdminDashboard = () => {
  const [metrics, setMetrics] = useState(null);
  const [topProducts, setTopProducts] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchDashboard();
  }, []);

  const fetchDashboard = async () => {
    try {
      setLoading(true);
      const [dashboardResponse, productsResponse] = await Promise.all([
        reportingService.getDashboard(),
        reportingService.getTopProducts(5),
      ]);
      setMetrics(dashboardResponse.data);
      setTopProducts(productsResponse.data);
    } catch (error) {
      toast.error('Failed to fetch dashboard data');
    } finally {
      setLoading(false);
    }
  };

  const handleExport = async (format) => {
    try {
      const blob = await reportingService.exportDashboard(format);
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `dashboard_metrics.${format}`;
      a.click();
      window.URL.revokeObjectURL(url);
      toast.success(`Report exported as ${format.toUpperCase()}`);
    } catch (error) {
      toast.error('Failed to export report');
    }
  };

  if (loading) {
    return <LoadingSpinner />;
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="flex items-center justify-between mb-8">
        <h1 className="text-4xl font-bold">Admin Dashboard</h1>
        <div className="flex space-x-2">
          <button onClick={() => handleExport('csv')} className="btn btn-secondary flex items-center">
            <Download className="w-4 h-4 mr-2" />
            Export CSV
          </button>
          <button onClick={() => handleExport('pdf')} className="btn btn-secondary flex items-center">
            <Download className="w-4 h-4 mr-2" />
            Export PDF
          </button>
        </div>
      </div>

      {/* Metrics Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        <div className="card">
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-gray-600 font-medium">Total Revenue</h3>
            <DollarSign className="w-8 h-8 text-green-600" />
          </div>
          <p className="text-3xl font-bold">${metrics?.totalRevenue?.toFixed(2) || '0.00'}</p>
          <p className="text-sm text-gray-600 mt-2">
            Today: ${metrics?.todayRevenue?.toFixed(2) || '0.00'}
          </p>
        </div>

        <div className="card">
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-gray-600 font-medium">Total Orders</h3>
            <Package className="w-8 h-8 text-blue-600" />
          </div>
          <p className="text-3xl font-bold">{metrics?.totalOrders || 0}</p>
          <p className="text-sm text-gray-600 mt-2">Today: {metrics?.todayOrders || 0}</p>
        </div>

        <div className="card">
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-gray-600 font-medium">Total Customers</h3>
            <Users className="w-8 h-8 text-purple-600" />
          </div>
          <p className="text-3xl font-bold">{metrics?.totalCustomers || 0}</p>
        </div>

        <div className="card">
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-gray-600 font-medium">Avg Order Value</h3>
            <TrendingUp className="w-8 h-8 text-orange-600" />
          </div>
          <p className="text-3xl font-bold">${metrics?.averageOrderValue?.toFixed(2) || '0.00'}</p>
        </div>
      </div>

      {/* Monthly Overview */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="card">
          <h2 className="text-2xl font-bold mb-6">This Month</h2>
          <div className="space-y-4">
            <div className="flex justify-between items-center">
              <span className="text-gray-600">Revenue</span>
              <span className="text-2xl font-bold text-primary-600">
                ${metrics?.monthRevenue?.toFixed(2) || '0.00'}
              </span>
            </div>
            <div className="flex justify-between items-center">
              <span className="text-gray-600">Orders</span>
              <span className="text-2xl font-bold text-primary-600">
                {metrics?.monthOrders || 0}
              </span>
            </div>
          </div>
        </div>

        {/* Top Products */}
        <div className="card">
          <h2 className="text-2xl font-bold mb-6">Top Products</h2>
          <div className="space-y-4">
            {topProducts.map((product, index) => (
              <div key={product.productId} className="flex items-center justify-between">
                <div className="flex items-center space-x-3">
                  <span className="text-2xl font-bold text-gray-400">#{index + 1}</span>
                  <div>
                    <p className="font-semibold">{product.productName}</p>
                    <p className="text-sm text-gray-600">{product.totalSold} sold</p>
                  </div>
                </div>
                <span className="font-bold text-primary-600">
                  ${product.totalRevenue.toFixed(2)}
                </span>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
};

export default AdminDashboard;
