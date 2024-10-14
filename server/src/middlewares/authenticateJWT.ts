import { Request, Response, NextFunction } from 'express';
import jwt from 'jsonwebtoken';

export interface UserRequest extends Request {
    user?: { userId: string; };
}

// JWT 인증 미들웨어
export const authenticateJWT = (req: UserRequest, res: Response, next: NextFunction) => {
    const authHeader = req.headers.authorization;
    console.log("[authenticateJWT] Authorization header:", authHeader);

    if (authHeader) {
        const token = authHeader.split(' ')[1]; // 'Bearer <token>'에서 토큰 부분만 추출

        jwt.verify(token, process.env.JWT_SECRET as string, (err, decoded) => {
            if (err) {
                res.sendStatus(403); // 토큰이 유효하지 않을 경우 접근 금지
                return;
            }

            // 토큰에서 사용자 정보를 추출하여 req.user에 저장
            req.user = {
                userId: (decoded as any).userId,
            };
            next(); // 다음 미들웨어 또는 라우트 핸들러로 이동
        });
    } else {
        res.sendStatus(401); // 인증 헤더가 없을 경우 접근 금지
        return;
    }
};

export default authenticateJWT;
