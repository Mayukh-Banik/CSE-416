import express from 'express';
import controllers from '../controllers';

const router = express.Router();

router.get('/status', controllers.checkAuthStatus);

export default router;