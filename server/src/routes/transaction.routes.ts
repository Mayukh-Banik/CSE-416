// routes/transaction.routes.ts
import express, { Request, Response } from 'express';
import Transaction from '../models/transaction.model';

const router = express.Router();

// Create a new transaction
router.post('/create', async (req: Request, res: Response) => {
  try {
    const { sender, receiver, amount } = req.body;

    const newTransaction = new Transaction({
      sender,
      receiver,
      amount,
    });

    await newTransaction.save();

    res.status(201).json(newTransaction);
  } catch (error) {
    res.status(500).json({ message: 'Error creating transaction', error });
  }
});

// Get all transactions
router.get('/', async (req: Request, res: Response) => {
  try {
    const transactions = await Transaction.find().populate('sender receiver', 'username email');
    res.json(transactions);
  } catch (error) {
    res.status(500).json({ message: 'Error retrieving transactions', error });
  }
});

export default router;
