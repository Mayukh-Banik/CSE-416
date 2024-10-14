// authController.ts (Controller)
import { Request, Response } from 'express';
import jwt from 'jsonwebtoken';

/**
 * Checks the authentication status of the user.
 * @param req - The request object.
 * @param res - The response object.
 */
export const checkAuthStatus = (req: Request, res: Response): void => {
    console.log("[authController] Request cookies:", req.cookies); // 쿠키 값 확인
    console.log("[authController] Authorization header:", req.headers.authorization); // Authorization 헤더 확인

    const token = req.cookies?.token;

    if (!token) {
        res.status(401).json({ isAuthenticated: false });
        return;
    }

    try {
        const decoded = jwt.verify(token, process.env.JWT_SECRET as string);
        res.status(200).json({ isAuthenticated: true, user: decoded });
    } catch (error) {
        res.status(401).json({ isAuthenticated: false });
    }
};
