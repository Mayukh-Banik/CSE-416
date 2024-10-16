import { Request, Response } from 'express';
import User from '../models/user';
import crypto from 'crypto';
import bcrypt from 'bcryptjs';
import jwt from 'jsonwebtoken';

// Create a new user
export const createUser = async (req: Request, res: Response) => {
    try {
        const { name, email, password } = req.body;
        
        // Check if user already exists
        const existingUser = await User.findOne({ email });
        if (existingUser) {
            res.status(400).json({ message: 'User with this email already exists' });
            return;
        }
        
        // Validate password length
        if (password.length < 8) {
            res.status(400).json({ message: 'Password must be at least 8 characters long' });
            return;
        }
        
        // Hash the password
        const hashedPassword = await bcrypt.hash(password, 10);
        
        // Generate key pair
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
            secretKey: privateKey // Returning private key for download
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

        // Ensure JWT_SECRET is properly set
        if (!process.env.JWT_SECRET) {
            throw new Error('JWT_SECRET is not defined in environment variables');
        }

        const token = jwt.sign({ userId: user._id }, process.env.JWT_SECRET, { expiresIn: '5m' });

        res.cookie('token', token, {
            httpOnly: true,
            // secure: process.env.NODE_ENV === 'production', //Send cookies only on HTTPS connections (for use in production environments)
            sameSite: 'strict', 
            maxAge: 5 * 60 * 1000
        });

        res.status(200).json({ message: 'Login successful'});
    } catch (error) {
        res.status(400).json({ message: 'Error logging in', error });
    }
};


// // Get a user by ID
// export const getUserById = async (req: Request, res: Response) => {
//   try {
//     const userId = req.params.id;
//     const user = await User.findById(userId);

//     if (!user) {
//       return res.status(404).json({ message: 'User not found' });
//     }

//     res.status(200).json(user);
//   } catch (error) {
//     res.status(400).json({ message: 'Error fetching user', error });
//   }
// };

// // Update a user's balance
// export const updateUserBalance = async (req: Request, res: Response) => {
//   try {
//     const userId = req.params.id;
//     const { balance } = req.body;

//     const updatedUser = await User.findByIdAndUpdate(
//       userId,
//       { balance },
//       { new: true, runValidators: true }
//     );

//     if (!updatedUser) {
//       return res.status(404).json({ message: 'User not found' });
//     }

//     res.status(200).json({ message: 'User balance updated', user: updatedUser });
//   } catch (error) {
//     res.status(400).json({ message: 'Error updating user balance', error });
//   }
// };

// // Delete a user by ID
// export const deleteUser = async (req: Request, res: Response) => {
//   try {
//     const userId = req.params.id;
//     const deletedUser = await User.findByIdAndDelete(userId);

//     if (!deletedUser) {
//       return res.status(404).json({ message: 'User not found' });
//     }

//     res.status(200).json({ message: 'User deleted successfully' });
//   } catch (error) {
//     res.status(400).json({ message: 'Error deleting user', error });
//   }
// };
