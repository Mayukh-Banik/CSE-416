import { Schema, model} from 'mongoose';

// 1. Create an interface representing a document in MongoDB.
export interface IUser {
  name: string;
  email: string;
  password: string;
  publicKey: string;
  balance: number;
  reputation: number;
  createdAt: Date;
}

// 2. Create a Schema corresponding to the document interface.
const userSchema = new Schema<IUser>({
  name: { type: String, required: true, unique: true, trim: true, minlength: 3 },
  email: { type: String, required: true, unique: true, trim: true, lowercase: true },
  password: { type: String, required: true, minlength: 6 },
  publicKey: { type: String, required: true},
  balance: { type: Number, default: 0, min: 0 },
  createdAt: { type: Date, default: Date.now, immutable: true }
});

// 3. Create a Model.
const User = model<IUser>('User', userSchema);

// Export the User model
export default User;
