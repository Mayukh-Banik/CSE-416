import { RequestHandler } from 'express';
import Transaction from '../models/transaction';
import mongoose from 'mongoose';

export const getTransactionHistory: RequestHandler = async (req, res, next) => {
    try {
        const { userId } = req.params;

        if (!mongoose.Types.ObjectId.isValid(userId)) {
            res.status(400).json({ error: 'Invalid user ID' });
            return;
        }

        const transactions = await Transaction.find({
            $or: [{ sender: userId }, { receiver: userId }],
        }).sort({ timestamp: -1 }); 

        res.status(200).json(transactions);
        return;
    } catch (error) {
        console.error('Error fetching transaction history:', error);
        next(error);
    }
};
