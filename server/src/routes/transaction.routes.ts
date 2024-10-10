import express from 'express';
import controllers from '../controllers';

const router = express.Router();

router.get('/user/:userId/transactions', controllers.getTransactionHistory);

export default router;
