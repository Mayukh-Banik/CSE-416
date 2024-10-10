import mongoose, { Schema, Document } from 'mongoose';
import { ObjectId } from 'mongodb';

// Define Transaction interface extending Mongoose Document
export interface ITransaction extends Document {
  sender: ObjectId;
  receiver: ObjectId;
  amount: number;
  timestamp: Date;
  status: 'pending' | 'completed' | 'failed';
}

// Define Transaction schema
const TransactionSchema: Schema<ITransaction> = new Schema({
  sender: { type: Schema.Types.ObjectId, ref: 'User', required: true },
  receiver: { type: Schema.Types.ObjectId, ref: 'User', required: true },
  amount: { type: Number, required: true, min: 0 },
  timestamp: { type: Date, default: Date.now },
  status: { type: String, enum: ['pending', 'completed', 'failed'], default: 'pending' },
}, { timestamps: true });

// Ensure sender and receiver are not the same
TransactionSchema.pre('save', function (next) {
  if (this.sender.equals(this.receiver)) {
    return next(new Error('Sender and receiver cannot be the same user.'));
  }
  next();
});

// Define Transaction model
const Transaction = mongoose.model<ITransaction>('Transaction', TransactionSchema);

// Export the Transaction model
export default Transaction;
