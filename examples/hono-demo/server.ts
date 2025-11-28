import { Hono } from 'hono';
import { serve } from '@hono/node-server';
import atomicdocs from '../../npm/atomicdocs/index.js';

const app = new Hono();

const PORT = 8080;
// atomicdocs auto-detects Hono and returns the appropriate middleware
app.use('*', atomicdocs(app, PORT));

// In-memory data
let users = [
  { id: 1, name: 'Alice', email: 'alice@example.com', age: 28, role: 'admin' },
  { id: 2, name: 'Bob', email: 'bob@example.com', age: 32, role: 'user' }
];

let products = [
  { id: 1, name: 'Laptop', price: 999, stock: 10, category: 'electronics' },
  { id: 2, name: 'Mouse', price: 25, stock: 50, category: 'electronics' },
  { id: 3, name: 'Desk', price: 299, stock: 5, category: 'furniture' }
];

let orders = [
  { id: 1, userId: 1, productId: 1, quantity: 1, status: 'pending' },
  { id: 2, userId: 2, productId: 2, quantity: 2, status: 'shipped' }
];

let nextUserId = 3;
let nextProductId = 4;
let nextOrderId = 3;

// ============ AUTH ROUTES ============

app.post('/auth/login', async (c) => {
  const { username, password } = await c.req.json();
  
  if (username === 'admin' && password === 'secret') {
    return c.json({ 
      token: 'hono-jwt-token-xyz',
      user: { username, role: 'admin' }
    });
  }
  
  return c.json({ error: 'Invalid credentials' }, 401);
});

app.post('/auth/signup', async (c) => {
  const { name, email, password, age } = await c.req.json();
  
  if (!name || !email || !password) {
    return c.json({ error: 'Missing required fields' }, 400);
  }
  
  const newUser = {
    id: nextUserId++,
    name,
    email,
    age: age || 18,
    role: 'user'
  };
  
  users.push(newUser);
  return c.json({ message: 'User created', userId: newUser.id }, 201);
});

// ============ USER ROUTES ============

app.get('/users', (c) => c.json(users));

app.get('/users/:id', (c) => {
  const user = users.find(u => u.id === parseInt(c.req.param('id')));
  if (!user) return c.json({ error: 'User not found' }, 404);
  return c.json(user);
});

app.post('/users', async (c) => {
  const { name, email, age, role } = await c.req.json();
  
  if (!name || !email) {
    return c.json({ error: 'Name and email required' }, 400);
  }
  
  const newUser = {
    id: nextUserId++,
    name,
    email,
    age: age || 0,
    role: role || 'user'
  };
  
  users.push(newUser);
  return c.json(newUser, 201);
});

app.put('/users/:id', async (c) => {
  const user = users.find(u => u.id === parseInt(c.req.param('id')));
  if (!user) return c.json({ error: 'User not found' }, 404);
  
  const { name, email, age, role } = await c.req.json();
  if (name) user.name = name;
  if (email) user.email = email;
  if (age !== undefined) user.age = age;
  if (role) user.role = role;
  
  return c.json(user);
});

app.delete('/users/:id', (c) => {
  const index = users.findIndex(u => u.id === parseInt(c.req.param('id')));
  if (index === -1) return c.json({ error: 'User not found' }, 404);
  
  users.splice(index, 1);
  return c.json({ message: 'User deleted' });
});

// ============ PRODUCT ROUTES ============

app.get('/products', (c) => c.json(products));

app.get('/products/:id', (c) => {
  const product = products.find(p => p.id === parseInt(c.req.param('id')));
  if (!product) return c.json({ error: 'Product not found' }, 404);
  return c.json(product);
});

app.post('/products', async (c) => {
  const { name, price, stock, category } = await c.req.json();
  
  if (!name || !price) {
    return c.json({ error: 'Name and price required' }, 400);
  }
  
  const newProduct = {
    id: nextProductId++,
    name,
    price: parseFloat(price),
    stock: stock || 0,
    category: category || 'general'
  };
  
  products.push(newProduct);
  return c.json(newProduct, 201);
});

app.put('/products/:id', async (c) => {
  const product = products.find(p => p.id === parseInt(c.req.param('id')));
  if (!product) return c.json({ error: 'Product not found' }, 404);
  
  const { name, price, stock, category } = await c.req.json();
  if (name) product.name = name;
  if (price !== undefined) product.price = parseFloat(price);
  if (stock !== undefined) product.stock = stock;
  if (category) product.category = category;
  
  return c.json(product);
});

app.delete('/products/:id', (c) => {
  const index = products.findIndex(p => p.id === parseInt(c.req.param('id')));
  if (index === -1) return c.json({ error: 'Product not found' }, 404);
  
  products.splice(index, 1);
  return c.json({ message: 'Product deleted' });
});

// ============ ORDER ROUTES ============

app.get('/orders', (c) => c.json(orders));

app.post('/orders', async (c) => {
  const { userId, productId, quantity } = await c.req.json();
  
  if (!userId || !productId || !quantity) {
    return c.json({ error: 'userId, productId, and quantity required' }, 400);
  }
  
  const newOrder = {
    id: nextOrderId++,
    userId: parseInt(userId),
    productId: parseInt(productId),
    quantity: parseInt(quantity),
    status: 'pending'
  };
  
  orders.push(newOrder);
  return c.json(newOrder, 201);
});

app.get('/stats/summary', (c) => {
  return c.json({
    totalUsers: users.length,
    totalProducts: products.length,
    totalOrders: orders.length
  });
});

serve({ fetch: app.fetch, port: PORT }, () => {
  console.log(`Hono app: http://localhost:${PORT}`);
  console.log(`Docs: http://localhost:${PORT}/docs`);
});
