import React from "react";
import { useParams } from "react-router-dom";
import { Box, Typography } from "@mui/material";
import Sidebar from "./Sidebar";
import { useTheme } from '@mui/material/styles';


// gotta make this entire file removed later
interface File {
  id: string;
  name: string;
  size: number; 
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
        <Typography variant="h4" gutterBottom>
          File Details
        </Typography>
        <Typography variant="body1">File Name: {file.name}</Typography>
        <Typography variant="body1">
          Size: {(file.size / 1024).toFixed(2)} KB {/* Optionally, handle larger files */}
        </Typography>
        <Typography variant="body1">Uploaded At: {file.uploadedAt}</Typography>
        <Typography variant="body1">Hash: 41c90bc99a00b80a7a8c032f302d4d994c7209e9f911146a15478348b7793bd27fb9075866e3c6f7f6d1c39830c389030d36142ab75e0cfdd4f749282d9fbc69
ef30f5d6d54627ba983f11e430b5852641806d162dd7f06b914798a09d79f9b9f369cf916f16dcf61416127b0c5fade1e3eedad21ac1b4c6032143b77d019c8d
66da2bc420688356c1a0ae59d0f40285e9cb84755376252934f70281b9dc7953c0b2a5bd79a6a4021d39d63041d82191a7b8c1a884f2ae8706758bc260b0658f
df206149a89ed09011b7c9f7e12ef258988819eecb72c94f0e5cf319c4ceec26f4e684988e443cbe321e8f9f6f4f1e254445fec27d14b21afb83a378db309479
d5e9e5bea308799a5e891d5eb27f6cb2322fe3e44cfb9fb9eaa17926df3d337358feb283e7d9402fa53d733eaccbc29953cbc95e8886cd2c992d95378343f296</Typography>
      </Box>
    </Box>
  );
};

export default FileViewPage;
