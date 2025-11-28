import { Hono } from 'hono';
import { serve } from '@hono/node-server';

const app = new Hono();

app.get('/__atomicdocs_routes', (c) => {
  const routes = app.routes.map((route: any) => ({
    method: route.method.toUpperCase(),
    path: route.path,
    summary: `${route.method.toUpperCase()} ${route.path}`
  }));
  return c.json(routes);
});

app.get('/products', (c) => c.json([{ id: 1, name: 'Widget' }]));
app.post('/products', (c) => c.json({ id: 2, name: 'New Product' }));
app.get('/products/:id', (c) => c.json({ id: c.req.param('id'), name: 'Widget' }));
app.put('/products/:id', (c) => c.json({ updated: true }));

serve({ fetch: app.fetch, port: 3000 }, () => {
  console.log('Hono app on http://localhost:3000');
});
