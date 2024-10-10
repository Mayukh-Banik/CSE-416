import { Request, Response } from 'express';
import User from '../models/user';
import crypto from 'crypto';
import bcrypt from 'bcryptjs';


// Create a new user
export const createUser = async (req: Request, res: Response) => {
    console.log('Request received at /api/users/signup');  // 로그 추가
    try {
        const { name, email, password } = req.body;
        const publicKey = crypto.randomBytes(32).toString('hex');
        const newUser = new User({
            name,
            email,
            password,
            publicKey, 
        });

        await newUser.save();
        res.status(201).json({ message: 'User created successfully', user: newUser });
    } catch (error) {
        res.status(400).json({ message: 'Error creating user', error });
    }
}; // ToDo: return sk

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

        // const token = jwt.sign({ userId: user._id }, process.env.JWT_SECRET as string, { expiresIn: '1h' });
        // res.status(200).json({ message: 'Login successful', token });
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
