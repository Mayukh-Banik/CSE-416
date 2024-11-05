const mongoose = require("mongoose");
var Schema = mongoose.Schema;

export interface FileMetadata {
    name: string;
    type: string;
    size: number;
    description: string;
    hash: string;
    createdAt?: Date; // Use Date type
    reputation?: number;
    isPublished?: boolean;
    fee?: number;
  }

/** 
  let metadata : FileMetadata = {
    id: `${file.name}-${file.size}-${Date.now()}`, // Unique ID for the uploaded file
    name: file.name,
    type: file.type,
    size: file.size,
    // file_data: base64FileData, // Encode file data as Base64 if required
    description: descriptions[file.name] || "",
    hash: fileHashes[file.name], // not needed - computed on backend
    isPublished: false, // Initially not published
};
*/

  
const fileSchema = new Schema({
    id:{type:String, required:true},
    name:{type:String, required:true},
    type:{type:String,required:true},
    size:{type:Number, required:true},
    description:{type: String, required:true},
    hash:{type:String},
    // createdAt:{type:Date, default: Date.now}
    isPublished:{type:Boolean},
    fee: {type:Number}
})

const File_Data = mongoose.model("File_Data", fileSchema)

module.exports = File_Data;