import { Request, Response } from 'express';
import User from '../models/user';
import crypto from 'crypto';
import bcrypt from 'bcryptjs';
import jwt from 'jsonwebtoken';

export const createUser = async (req: Request, res: Response) => {
    try {
        const { name, email, password } = req.body;
        
        const existingUser = await User.findOne({ email });
        if (existingUser) {
            res.status(400).json({ message: 'User with this email already exists' });
            return;
        }
        
        if (password.length < 8) {
            res.status(400).json({ message: 'Password must be at least 8 characters long' });
            return;
        }
        
        const hashedPassword = await bcrypt.hash(password, 10);
        
        const { publicKey, privateKey } = crypto.generateKeyPairSync('rsa', {
            modulusLength: 2048,
            publicKeyEncoding: {
                type: 'spki',
                format: 'pem'
            },
            privateKeyEncoding: {
                type: 'pkcs8',
                format: 'pem'
            }
        });
        
        const newUser = new User({
            name,
            email,
            password: hashedPassword,
            publicKey, 
        });

        await newUser.save();
        
        res.status(201).json({ 
            message: 'User created successfully', 
            user: newUser, 
            secretKey: privateKey 
        });
    } catch (error) {
        res.status(400).json({ message: 'Error creating user', error });
    }
};

export const loginUser = async (req: Request, res: Response) => {
    try {
        const { email, password } = req.body;
        const user = await User.findOne({ email });

        if (!user) {
            res.status(404).json({ message: 'Invalid email or password' });
            return;
        }

        const isPasswordValid = await bcrypt.compare(password, user.password);
        if (!isPasswordValid) {
            res.status(404).json({ message: 'Invalid email or password' });
            return;
        }

        if (!process.env.JWT_SECRET) {
            throw new Error('JWT_SECRET is not defined in environment variables');
        }

        const token = jwt.sign({ userId: user._id }, process.env.JWT_SECRET, { expiresIn: '5m' });

        res.cookie('token', token, {
            httpOnly: true,
            sameSite: 'strict', 
            maxAge: 5 * 60 * 1000
        });

        res.status(200).json({ message: 'Login successful'});
    } catch (error) {
        res.status(400).json({ message: 'Error logging in', error });
    }
};
