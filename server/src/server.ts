import express from 'express';
import mongoose from 'mongoose';
import dotenv from 'dotenv';
import router from './routes';
import {createMonkUsers, createMonkTransactions} from './monkData';

// Load environment variables
dotenv.config();

const app = express();
const PORT = process.env.PORT || 5000;

// MongoDB connection string from environment variables
const MONGODB_URI = process.env.MONGODB_URI || 'mongodb://localhost:27017/test';

// Middleware for parsing JSON request bodies
app.use(express.json());

// Use routes
app.use('/api', router);

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

const startServer = async (): Promise<void> => {
  await connectDB();

  app.listen(PORT, () => {
    console.log(`Server is running on port ${PORT}`);
  });

  //monk data
  await createMonkUsers();
  await createMonkTransactions();

};

startServer();

