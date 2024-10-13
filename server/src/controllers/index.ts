// controllers/index.ts (Index for Controllers)
import * as userController from './userController';
import * as transactionController from './transactionController';
import { checkAuthStatus } from './authController';

export default {
  ...userController,
  ...transactionController,
  checkAuthStatus,
};