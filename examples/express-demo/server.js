const express = require('express');
const atomicdocs = require('atomicdocs');

const app = express();
app.use(express.json());
app.use(atomicdocs());

// In-memory data
let users = [
  { id: 1, name: 'John Doe', email: 'john@example.com', age: 30, role: 'admin' },
  { id: 2, name: 'Jane Smith', email: 'jane@example.com', age: 25, role: 'user' },
  { id: 3, name: 'Bob Johnson', email: 'bob@example.com', age: 35, role: 'user' }
];

let posts = [
  { id: 1, userId: 1, title: 'First Post', content: 'Hello World', likes: 10 },
  { id: 2, userId: 2, title: 'Second Post', content: 'Another post', likes: 5 }
];

let comments = [
  { id: 1, postId: 1, userId: 2, text: 'Great post!', createdAt: '2025-01-01' },
  { id: 2, postId: 1, userId: 3, text: 'Thanks for sharing', createdAt: '2025-01-02' }
];

let nextUserId = 4;
let nextPostId = 3;
let nextCommentId = 3;

// ============ AUTH ROUTES ============

// POST login
app.post('/auth/login', (req, res) => {
  const { username, password } = req.body;
  
  if (username === 'admin' && password === 'password') {
    return res.json({ 
      token: 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9',
      user: { username: 'admin', role: 'admin' }
    });
  }
  
  res.status(401).json({ error: 'Invalid credentials' });
});

// POST register
app.post('/auth/register', (req, res) => {
  const { username, email, password, age } = req.body;
  
  if (!username || !email || !password) {
    return res.status(400).json({ error: 'Missing required fields' });
  }
  
  const newUser = {
    id: nextUserId++,
    name: username,
    email,
    age: age || 18,
    role: 'user'
  };
  
  users.push(newUser);
  res.status(201).json({ message: 'User registered', userId: newUser.id });
});

// ============ USER ROUTES ============

// GET all users
app.get('/users', (req, res) => {
  res.json(users);
});

// GET user by ID
app.get('/users/:id', (req, res) => {
  const user = users.find(u => u.id === parseInt(req.params.id));
  if (!user) {
    return res.status(404).json({ error: 'User not found' });
  }
  res.json(user);
});

// POST create user
app.post('/users', (req, res) => {
  const { name, email, age, role } = req.body;
  
  if (!name || !email) {
    return res.status(400).json({ error: 'Name and email are required' });
  }
  
  const newUser = {
    id: nextUserId++,
    name,
    email,
    age: age || 0,
    role: role || 'user'
  };
  
  users.push(newUser);
  res.status(201).json(newUser);
});

// PUT update user
app.put('/users/:id', (req, res) => {
  const user = users.find(u => u.id === parseInt(req.params.id));
  if (!user) {
    return res.status(404).json({ error: 'User not found' });
  }
  
  const { name, email, age, role } = req.body;
  if (name) user.name = name;
  if (email) user.email = email;
  if (age !== undefined) user.age = age;
  if (role) user.role = role;
  
  res.json(user);
});

// DELETE user
app.delete('/users/:id', (req, res) => {
  const index = users.findIndex(u => u.id === parseInt(req.params.id));
  if (index === -1) {
    return res.status(404).json({ error: 'User not found' });
  }
  
  users.splice(index, 1);
  res.json({ message: 'User deleted successfully' });
});

// ============ POST ROUTES ============

// GET all posts
app.get('/posts', (req, res) => {
  res.json(posts);
});

// GET post by ID
app.get('/posts/:id', (req, res) => {
  const post = posts.find(p => p.id === parseInt(req.params.id));
  if (!post) {
    return res.status(404).json({ error: 'Post not found' });
  }
  res.json(post);
});

// POST create post
app.post('/posts', (req, res) => {
  const { userId, title, content } = req.body;
  
  if (!userId || !title || !content) {
    return res.status(400).json({ error: 'userId, title, and content are required' });
  }
  
  const newPost = {
    id: nextPostId++,
    userId: parseInt(userId),
    title,
    content,
    likes: 0
  };
  
  posts.push(newPost);
  res.status(201).json(newPost);
});

// PUT update post
app.put('/posts/:id', (req, res) => {
  const post = posts.find(p => p.id === parseInt(req.params.id));
  if (!post) {
    return res.status(404).json({ error: 'Post not found' });
  }
  
  const { title, content } = req.body;
  if (title) post.title = title;
  if (content) post.content = content;
  
  res.json(post);
});

// POST like a post
app.post('/posts/:id/like', (req, res) => {
  const post = posts.find(p => p.id === parseInt(req.params.id));
  if (!post) {
    return res.status(404).json({ error: 'Post not found' });
  }
  
  post.likes++;
  res.json({ likes: post.likes });
});

// DELETE post
app.delete('/posts/:id', (req, res) => {
  const index = posts.findIndex(p => p.id === parseInt(req.params.id));
  if (index === -1) {
    return res.status(404).json({ error: 'Post not found' });
  }
  
  posts.splice(index, 1);
  res.json({ message: 'Post deleted successfully' });
});

// ============ COMMENT ROUTES ============

// GET comments for a post
app.get('/posts/:postId/comments', (req, res) => {
  const postComments = comments.filter(c => c.postId === parseInt(req.params.postId));
  res.json(postComments);
});

// POST create comment
app.post('/posts/:postId/comments', (req, res) => {
  const { userId, text } = req.body;
  
  if (!userId || !text) {
    return res.status(400).json({ error: 'userId and text are required' });
  }
  
  const newComment = {
    id: nextCommentId++,
    postId: parseInt(req.params.postId),
    userId: parseInt(userId),
    text,
    createdAt: new Date().toISOString().split('T')[0]
  };
  
  comments.push(newComment);
  res.status(201).json(newComment);
});

// DELETE comment
app.delete('/comments/:id', (req, res) => {
  const index = comments.findIndex(c => c.id === parseInt(req.params.id));
  if (index === -1) {
    return res.status(404).json({ error: 'Comment not found' });
  }
  
  comments.splice(index, 1);
  res.json({ message: 'Comment deleted successfully' });
});

// ============ SEARCH ROUTES ============

// GET search users
app.get('/search/users', (req, res) => {
  const { query, role } = req.query;
  let results = users;
  
  if (query) {
    results = results.filter(u => 
      u.name.toLowerCase().includes(query.toLowerCase()) ||
      u.email.toLowerCase().includes(query.toLowerCase())
    );
  }
  
  if (role) {
    results = results.filter(u => u.role === role);
  }
  
  res.json(results);
});

const PORT = 6767;
app.set('port', PORT);

app.listen(PORT, () => {
  console.log(`App: http://localhost:${PORT}`);
  console.log(`Docs: http://localhost:${PORT}/docs`);
  
  atomicdocs.register(app, PORT);
});
