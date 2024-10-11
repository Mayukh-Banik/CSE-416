import mongoose, { Schema, Document } from 'mongoose';
import { ObjectId } from 'mongodb';

export interface ITransaction extends Document {
  transactionId: string;
  sender: ObjectId;
  receiver: ObjectId;
  amount: number;
  timestamp: Date;
  fileName: string;
  fileId: string;
  fileSize: string;
  fee: number;
  status: 'pending' | 'completed' | 'failed';
}

const TransactionSchema: Schema<ITransaction> = new Schema({
  transactionId: { type: String, required: true, unique: true, trim: true, minlength: 3 },
  sender: { type: Schema.Types.ObjectId, ref: 'User', required: true },
  receiver: { type: Schema.Types.ObjectId, ref: 'User', required: true },
  amount: { type: Number, required: true, min: 0 },
  timestamp: { type: Date, default: Date.now },
  fileName: { type: String, required: true, trim: true },
  fileId: { type: String, required: true, trim: true },
  fileSize: { type: String, required: true, trim: true },
  fee: { type: Number, required: true, min: 0 },
  status: { type: String, enum: ['pending', 'completed', 'failed'], default: 'pending' },
}, { timestamps: true });

TransactionSchema.pre('save', function (next) {
  if (this.sender.equals(this.receiver)) {
    return next(new Error('Sender and receiver cannot be the same user.'));
  }
  next();
});

const Transaction = mongoose.model<ITransaction>('Transaction', TransactionSchema);

export default Transaction;
