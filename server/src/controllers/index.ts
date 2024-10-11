import * as userController from './userController';
import * as transactionController from './transactionController';

export default {
  ...userController,
  ...transactionController,
};
