// routes/user.routes.ts
import express, { Request, Response } from 'express';
import User from '../models/user.model';

const router = express.Router();

// Create a new user
router.post('/register', async (req: Request, res: Response) => {
  try {
    const { username, email, password } = req.body;

    const newUser = new User({
      username,
      email,
      password,
    });

    await newUser.save();

    res.status(201).json(newUser);
  } catch (error) {
    res.status(500).json({ message: 'Error creating user', error });
  }
});

// Get all users
router.get('/', async (req: Request, res: Response) => {
  try {
    const users = await User.find();
    res.json(users);
  } catch (error) {
    res.status(500).json({ message: 'Error retrieving users', error });
  }
});

export default router;
