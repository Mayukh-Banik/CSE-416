import mongoose, { Schema, Document } from 'mongoose';
import { ObjectId } from 'mongodb';

// Define User interface extending Mongoose Document
export interface IUser extends Document {
  username: string;
  email: string;
  password: string;
  balance: number;
  createdAt: Date;
}

// Define User schema
const UserSchema: Schema<IUser> = new Schema({
  username: { type: String, required: true, unique: true, trim: true, minlength: 3 },
  email: { type: String, required: true, unique: true, trim: true, lowercase: true },
  password: { type: String, required: true, minlength: 6 },
  publicKey: { type: String, required: true},
  balance: { type: Number, default: 0, min: 0 },
  createdAt: { type: Date, default: Date.now, immutable: true },
}, { timestamps: true });

// Define indexes for performance
UserSchema.index({ email: 1, username: 1 });

// Define User model based on the schema
const User = mongoose.model<IUser>('User', UserSchema);

// Export the User model
export default User;
