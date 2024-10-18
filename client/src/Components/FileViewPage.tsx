import React from "react";
import { useParams } from "react-router-dom";
import { Box, Typography } from "@mui/material";
import Sidebar from "./Sidebar";

interface File {
  name: string;
  size: number; // size in bytes
  uploadedAt: string;
  rating: number;
}

const files: File[] = [
  { name: "file001.pdf", size: 102400, uploadedAt: "2024-10-01", rating: 4 },
  { name: "file009.docx", size: 204800, uploadedAt: "2024-10-05", rating: 5 }
  // Add other files...
];

const FileViewPage: React.FC = () => {
  const { fileName } = useParams<{ fileName: string }>();
  const file = files.find((file) => file.name === fileName);

  if (!file) {
    return <Typography variant="h5">File not found</Typography>;
  }

  return (
    <Box sx={{ display: "flex" }}>
      <Sidebar />
      <Box sx={{ flexGrow: 1, padding: 3 }}>
        <Typography variant="h4" gutterBottom>
          File Details
        </Typography>
        <Typography variant="body1">File Name: {file.name}</Typography>
        <Typography variant="body1">
          Size: {(file.size / 1024).toFixed(2)} KB
        </Typography>
        <Typography variant="body1">Uploaded At: {file.uploadedAt}</Typography>
        <Typography variant="body1">Rating: {file.rating}/5</Typography>

        {/* Account's Hosting section */}
        <Box sx={{ mt: 4 }}>
          <Typography variant="h5">Account's Hosting</Typography>
          {/* Display account hosting details */}
        </Box>
      </Box>
    </Box>
  );
};

export default FileViewPage;
