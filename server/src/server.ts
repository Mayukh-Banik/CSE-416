import express, { Request, Response, NextFunction } from 'express';
import mongoose from 'mongoose';
import dotenv from 'dotenv';
import routes from './routes';
// Import your routes (uncomment when available)
// import userRoutes from './routes/user.routes';
// import transactionRoutes from './routes/transaction.routes';

// Load environment variables
dotenv.config();

const app = express();
const PORT = process.env.PORT || 5000;

// MongoDB connection string from environment variables
const MONGODB_URI = process.env.MONGODB_URI || 'mongodb://localhost:27017/test';

// Middleware for parsing JSON request bodies
app.use(express.json());

// Use routes
app.use('/api', routes);

// Function to connect to MongoDB
const connectDB = async (): Promise<void> => {
  try {
    // Connect to MongoDB using mongoose
    await mongoose.connect(MONGODB_URI);
    console.log('MongoDB connected successfully');
  } catch (error) {
    console.error('Error connecting to MongoDB:', error);
    process.exit(1); // Exit the process if the connection fails
  }
};
connectDB();

// Default route
app.get('/', (req: Request, res: Response) => {
  res.send('Hello from TypeScript Express server!');
});

// // Routes (uncomment these lines when you have the routes ready)
// app.use('/api/users', userRoutes);
// app.use('/api/transactions', transactionRoutes);

// Global error handler middleware
app.use((err: Error, req: Request, res: Response, next: NextFunction) => {
  console.error(err.stack);
  res.status(500).send({ error: 'Something went wrong!' });
});

// Start the server
app.listen(PORT, () => {
  console.log(`Server is running on port ${PORT}`);
});
