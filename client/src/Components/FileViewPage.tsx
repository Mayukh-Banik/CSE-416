import React from "react";
import { useParams } from "react-router-dom";
import { Box, Paper, Typography } from "@mui/material";
import Sidebar from "./Sidebar";
import { useTheme } from '@mui/material/styles';

interface File {
  id: string;
  name: string;
  size: number; // size in bytes
  uploadedAt: string;
  rating: number;
}

const drawerWidth = 300;
const collapsedDrawerWidth = 100;

const files: File[] = [
  { id: "tx001", name: "file001.pdf", size: 102400, uploadedAt: "2024-10-01", rating: 4 },
  { id: "tx009", name: "file009.docx", size: 204800, uploadedAt: "2024-10-05", rating: 5 },
  { id: "file1.txt", name: "file1.txt", size: 51200, uploadedAt: "2024-10-10", rating: 3 },
  { id: "file2.txt", name: "file2.txt", size: 76800, uploadedAt: "2024-10-12", rating: 4 },
  // Add more files as needed...
];

const FileViewPage: React.FC = () => {
  const theme = useTheme();
  const { fileId } = useParams<{ fileId: string }>();
  const file = files.find((file) => file.id === fileId);

  if (!fileId || !file) {
    return (
      <Box sx={{ padding: 3 }}>
        <Typography variant="h5">File not found</Typography>
      </Box>
    );
  }

  return (
    <Box
      sx={{
        padding: 2,
        marginTop: '70px',
        marginLeft: `${drawerWidth}px`,
        transition: 'margin-left 0.3s ease',
        [theme.breakpoints.down('sm')]: {
          marginLeft: `${collapsedDrawerWidth}px`,
        },
      }}
    >
      <Sidebar />
      <Box sx={{ flexGrow: 1, padding: 3 }}>
        <Paper elevation={3} sx={{ padding: 3 }}>
          <Typography variant="h4" gutterBottom>
            File Details
          </Typography>
          <Typography variant="body1" sx={{ fontWeight: 'bold' }}>File Name:</Typography>
          <Typography variant="body1">{file.name}</Typography>
          
          <Typography variant="body1" sx={{ fontWeight: 'bold' }}>Size:</Typography>
          <Typography variant="body1">{(file.size / 1024).toFixed(2)} KB</Typography>
          
          <Typography variant="body1" sx={{ fontWeight: 'bold' }}>Uploaded At:</Typography>
          <Typography variant="body1">{file.uploadedAt}</Typography>
          
          <Typography variant="body1" sx={{ fontWeight: 'bold' }}>Hash:</Typography>
          <Typography variant="body1" sx={{ wordWrap: "break-word" }}>
            41c90bc99a00b80a7a8c032f302d4d994c7209e9f911146a15478348b7793bd27fb9075866e3c6f7f6d1c39830c389030d36142ab75e0cfdd4f749282d9fbc69...
          </Typography>
        </Paper>
      </Box>
    </Box>
  );
};

export default FileViewPage;