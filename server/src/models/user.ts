import { Schema, model} from 'mongoose';

export interface IUser {
  name: string;
  email: string;
  password: string;
  publicKey: string;
  balance: number;
  reputation: number;
  createdAt: Date;
}

const userSchema = new Schema<IUser>({
  name: { type: String, required: true, unique: true, trim: true, minlength: 3 },
  email: { type: String, required: true, unique: true, trim: true, lowercase: true },
  password: { type: String, required: true, minlength: 6 },
  publicKey: { type: String, required: true},
  balance: { type: Number, default: 0, min: 0 },
  createdAt: { type: Date, default: Date.now, immutable: true }
});

const User = model<IUser>('User', userSchema);

export default User;
