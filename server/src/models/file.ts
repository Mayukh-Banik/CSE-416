import { Schema, model } from 'mongoose';

export interface IFile {
    fileName: string;
    hash: string;
    reputation: string;
    fileSize: number;// in byte
    createdAt: Date;
}

const fileSchema = new Schema<IFile>({
    fileName: { type: String, required: true, trim: true },
    hash: { type: String, required: true, unique: true },
    reputation: { type: String, required: true },
    fileSize: { type: Number, required: true, min: 0 },
    createdAt: { type: Date, default: Date.now, immutable: true }
});

const File = model<IFile>('File', fileSchema);

export default File;