const express = require('express');
const mongoose = require('mongoose');
const cors = require('cors');
const bodyParser = require('body-parser');

//Model
const File_Data = require("./models/file");


const app = express();
const mongoDB = "mongodb://localhost:27017/files"
const PORT = 3000; // Set your desired port

// Middleware
app.use(cors());
app.use(bodyParser.json());
app.use(express.json());
// MongoDB connection
mongoose.connect(mongoDB, {
  useNewUrlParser: true,
  useUnifiedTopology: true,
});

var db = mongoose.connection()
db.on("error", console.error.bind(console, "MongoDB connection error"));

// testing ruote 
app.get("/",function (req,res){
  res.send("Hello world!");
})


app.post('/upload',async function(req,res){
  const {name,type,size,description,hash,isPublished,fee} = req.body;
  if(!name||!type||!size||!description||!hash||!isPublished||!fee)
  {
    return res.status(400).json({message:'Missing required fields'});
  }

  const newFile = new File_Data({
    name,
    type,
    size,
    description,
    hash,
    isPublished,
  })
  
  try{
    const existingFile = await File_Data.findOne(newFile);
    if(existingFile)
    {
      return res.status(400).send("File already published");
    }
    const savedFile = await newFile.save();
    res.status(201).json(savedFile);
  } catch(error){
    console.log("Error saving file:",error);
    res.status(500).json({message: "Failed uploading"});
  }
});

app.get('/fetchUploadedFiles', async (req, res) => {
  try {
    const files = await File_Data.find();
    res.json(files);
  } catch (error) {
  res.status(500).json({error: "Failed to fetch all files"})
  }
});

app.delete('/delete/:hash', async(req,res) => {
  try{
    const {hash} = req.params;
    const deletedFile = await File_Data.findOneAndDelete({hash});
    if (!deletedFile) {
      return res.status(404).json({error: 'File not found'});
    }
    res.json({message: 'File deleted'});
  } catch (error) {
    res.status(500).json({error: 'Failed to delete'})
  }
})

app.put('/update/:hash'), async(req, res)=> {
  try {
    const { hash } = req.params;
    const updatedFile = await File_Data.findOneAndUpdate({ hash }, req.body, { new: true });
    if (!updatedFile) {
        return res.status(404).json({ error: 'File not found' });
    }
    res.json({ message: 'File updated' });
} catch (error) {
    res.status(500).json({ error: 'Failed to update file' });
}
}


app.listen(PORT,()=>{
  console.log("Server running on port 3000");
})