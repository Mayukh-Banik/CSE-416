import { Schema, model } from 'mongoose';

// 1. Create an interface representing a document in MongoDB.
export interface IFile {
    fileName: string;
    hash: string;
    reputation: number;
    fileSize: number;// in byte
    createdAt: Date;
    //tbd
    // file_extension: string; 
    // extraneous: string; //for search bar other than that useless
    // fileType: string; 
}

// 2. Create a Schema corresponding to the document interface.
const fileSchema = new Schema<IFile>({
    fileName: { type: String, required: true, trim: true },
    hash: { type: String, required: true, unique: true },
    reputation: { type: Number, required: true },
    fileSize: { type: Number, required: true, min: 0 },
    createdAt: { type: Date, default: Date.now, immutable: true },
    // fileExtension: { type: String, trim: true },
    // extraneous: { type: String, trim: true },
    // fileType: { type: String, trim: true }
});

// 3. Create a Model.
const File = model<IFile>('File', fileSchema);

// Export the File model
export default File;