const express = require('express');
const atomicdocs = require('../../npm/atomicdocs');

const app = express();
app.use(express.json());
app.use(atomicdocs());

app.get('/users', (req, res) => {
  res.json([{ id: 1, name: 'John' }]);
});

app.post('/users', (req, res) => {
  res.json({ id: 2, name: req.body.name });
});

app.get('/users/:id', (req, res) => {
  res.json({ id: req.params.id, name: 'John' });
});

app.delete('/users/:id', (req, res) => {
  res.json({ deleted: true });
});

const PORT = 6767;
app.set('port', PORT);

app.listen(PORT, () => {
  console.log(`App: http://localhost:${PORT}`);
  console.log(`Docs: http://localhost:${PORT}/docs`);
});
