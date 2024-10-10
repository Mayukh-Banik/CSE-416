import express from 'express';
import controllers from '../controllers';
import { body } from 'express-validator';
import { validateRequest } from '../middlewares';

const router = express.Router();

router.post(
  '/signup',
  [
    body('email').isEmail().withMessage('Please provide a valid email'),
    body('password').isLength({ min: 6 }).withMessage('Password must be at least 6 characters long'),
    validateRequest,
  ],
  controllers.createUser
);

export default router;
